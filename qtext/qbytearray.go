package qtext

/*
QByteArray — Go-обёртка над Qt 6.11 QtCore/QByteArray.

# Владение

Каждый *QByteArray — handle на C++ QByteArray в куче. Копирование (Clone) дёшево
(COW, refcount++), мутирующие методы внутри Qt выполняют detach. C++ объект
освобождается либо Close() (немедленно), либо runtime.AddCleanup (при сборке
мусора). Close() идемпотентен и безопасен для nil-получателя.

# Use after Close

После Close() любой метод, кроме Close(), паникует. Nil-safe-методы
(IsSharedWith, ToPercentEncoding) используют tryHandle().

# Варианты API (Go-native)

Каждая операция, у которой в Qt на входе или выходе есть QByteArray или
QByteArrayView, имеет два варианта (исключения — IsSharedWith и Swap: для них
[]byte-вариант не имеет смысла):

  - handle-вариант — *QByteArray (владение, COW, без копирования);
  - Go-вариант — []byte: методы с суффиксом Bytes возвращают копию содержимого,
    а функции, принимающие *QByteArray, имеют пару, принимающую []byte.

Для параметров const char* добавлены удобные ...String-варианты
(NewQByteArrayFromString, AppendString, PrependString, InsertString,
AssignString, ReplaceAllString).

QString, QStringView, QLatin1StringView и QStringList в публичном API QByteArray
(Qt 6.11) не участвуют — сравнения с QString удалены в Qt 6.8, — поэтому
отдельных string/[]string-перегрузок для них не требуется.

# Соглашения

Конструкторы начинаются с NewQByteArray*, статические фабрики — с
NewQByteArrayNumber*. При выходе индекса за диапазон методы паникуют: в Qt это
Q_ASSERT, а в release-сборке — неопределённое поведение. Left/Right/Mid
повторяют зажимающую семантику Qt (без паники).

# Что сознательно не вошло

  - fromRawData/setRawData — требуют runtime.Pinner и ручного владения;
    используйте NewQByteArrayFromBytes (копирует данные).
  - data() (mutable), operator[], итераторы — не идиоматичны в Go;
    используйте At/SetAt/Bytes.
  - removeIf/erase_if с C++-предикатом — есть Go-версия RemoveIf.
  - QDataStream и qHash — соответствующие пакеты ещё не реализованы.
  - Платформенные CFData/NSData/emscripten-конверсии.
  - QByteArray::FromBase64Result — заменён двухзначным возвратом
    NewQByteArrayFromBase64Encoding.
*/

/*
#cgo pkg-config: Qt6Core
#cgo CXXFLAGS: -std=c++17 -fPIC
#cgo LDFLAGS: -lstdc++
#include <stdint.h>
#include <stdbool.h>
#include <stddef.h>
#include "qbytearray_c.h"
*/
import "C"

import (
	"runtime"
	"sync/atomic"
	"unsafe"
)

// qbytearrayState хранит указатель на C++ QByteArray в атомарном слоте, чтобы
// Close() и cleanup могли безопасно соревноваться: кто первым обнулил слот,
// тот владеет объектом и удаляет его.
//
// Тип *byte, а не *C.QByteArray, выбран потому, что cgo объявляет C.QByteArray
// неполным типом, а atomic.Pointer[T] требует полного типа. Значение указателя
// не меняется — только статический тип на стороне Go; на границе C-вызовов
// выполняется обратное приведение к *C.QByteArray.
//
// State — отдельная аллокация (не поле QByteArray): cleanup-функция не должна
// захватывать *QByteArray, иначе q останется достижимым навсегда и cleanup
// никогда не выполнится.
type qbytearrayState struct {
	p atomic.Pointer[byte]
}

// cpp возвращает нижележащий *C.QByteArray или nil, если слот уже захвачен
// (объект закрыт). Только чтение, владение не забирает.
// Qt: аналога нет — внутренний помощник.
func (s *qbytearrayState) cpp() *C.QByteArray {
	if bp := s.p.Load(); bp != nil {
		return (*C.QByteArray)(unsafe.Pointer(bp))
	}
	return nil
}

// claim атомарно обнуляет слот и возвращает лежавший там *C.QByteArray (если
// был). Ровно один вызывающий получает не-nil для данного указателя; остальные
// видят nil. Именно это делает Close() и cleanup безопасными против double free.
// Qt: аналога нет — внутренний помощник.
func (s *qbytearrayState) claim() *C.QByteArray {
	if bp := s.p.Swap(nil); bp != nil {
		return (*C.QByteArray)(unsafe.Pointer(bp))
	}
	return nil
}

// QByteArray — обёртка над Qt6 QByteArray.
// Нулевое значение непригодно: используйте NewQByteArray / NewQByteArrayFromBytes
// и другие конструкторы.
type QByteArray struct {
	state   *qbytearrayState
	cleanup runtime.Cleanup
}

// wrapQByteArray оборачивает непрозрачный C++ QByteArray* в *QByteArray.
// Принимает unsafe.Pointer для симметрии с wrapQString и чтобы C-ABI мог
// возвращать как *C.QByteArray, так и go_qt6core_qbytearray_t*. Забирает
// владение: cleanup (или Close) освободит C++ объект.
// Qt: аналога нет — внутренний помощник.
func wrapQByteArray(p unsafe.Pointer) *QByteArray {
	if p == nil {
		return nil
	}
	st := &qbytearrayState{}
	st.p.Store((*byte)(p))
	q := &QByteArray{state: st}
	q.cleanup = runtime.AddCleanup(q, func(st *qbytearrayState) {
		if cpp := st.claim(); cpp != nil {
			C.qbytearray_delete(cpp)
		}
	}, st)
	return q
}

// handle возвращает нижележащий C++ указатель, паникуя, если q == nil или
// объект уже закрыт. Все методы, кроме Close() и nil-safe-методов, идут через
// этот guard.
// Qt: аналога нет — внутренний помощник.
func (q *QByteArray) handle() *C.QByteArray {
	if cpp, ok := q.tryHandle(); ok {
		return cpp
	}
	panic("qbytearray: use after Close (or nil receiver)")
}

// tryHandle — непаникующий вариант handle. Возвращает nil, false, если q == nil,
// без state или уже закрыт.
// Qt: аналога нет — внутренний помощник.
func (q *QByteArray) tryHandle() (*C.QByteArray, bool) {
	if q == nil || q.state == nil {
		return nil, false
	}
	if cpp := q.state.cpp(); cpp != nil {
		return cpp, true
	}
	return nil, false
}

// ---- construction --------------------------------------------------------

// NewQByteArrayFromBytes создаёт QByteArray копированием b.
// nil даёт null QByteArray (Qt: QByteArray(nullptr, 0)); не-nil пустой срез —
// пустой, но не null массив.
// Qt: QByteArray(const char *, qsizetype) — с копированием данных.
func NewQByteArrayFromBytes(b []byte) *QByteArray {
	if b == nil {
		return wrapQByteArray(unsafe.Pointer(C.qbytearray_new_empty()))
	}
	if len(b) == 0 {
		return wrapQByteArray(unsafe.Pointer(C.qbytearray_new_empty_nonnull()))
	}
	p := (*C.char)(unsafe.Pointer(&b[0]))
	q := wrapQByteArray(unsafe.Pointer(C.qbytearray_new_from_data(p, C.ptrdiff_t(len(b)))))
	runtime.KeepAlive(b)
	return q
}

// NewQByteArrayFromString создаёт QByteArray из строки s: байты UTF-8, без
// завершающего NUL. Встроенные NUL в s сохраняются.
// Qt: QByteArray(const char *, qsizetype) с size = len(s).
func NewQByteArrayFromString(s string) *QByteArray {
	return NewQByteArrayFromBytes([]byte(s))
}

// NewQByteArrayFromFill возвращает QByteArray размера n, заполненный c.
// Qt: QByteArray(qsizetype, char).
func NewQByteArrayFromFill(n int, c byte) *QByteArray {
	return wrapQByteArray(unsafe.Pointer(
		C.qbytearray_new_fill(C.ptrdiff_t(n), C.char(c))))
}

// NewQByteArrayFromQCharDecomposition возвращает Unicode-декомпозицию c в
// UTF-8 (QChar::decomposition().toUtf8()).
// Qt: QChar::decomposition() const, результат сконвертирован в QByteArray.
func NewQByteArrayFromQCharDecomposition(c QChar) *QByteArray {
	s := c.Decompose()
	if s == nil {
		return NewQByteArrayFromBytes(nil)
	}
	defer s.Close()
	return s.Utf8()
}

// NewQByteArrayFromQCharDecompositionLatin1 — то же, что
// NewQByteArrayFromQCharDecomposition, но в Latin-1.
// Qt: QChar::decomposition() const, результат сконвертирован в Latin-1 QByteArray.
func NewQByteArrayFromQCharDecompositionLatin1(c QChar) *QByteArray {
	s := c.Decompose()
	if s == nil {
		return NewQByteArrayFromBytes(nil)
	}
	defer s.Close()
	return s.Latin1()
}

// NewQByteArray диспетчеризует конструктор по конкретным Go-типам аргументов.
// Qt: единого аналога нет — объединяет перегрузки конструкторов.
//
//	NewQByteArray()                 → пустой (null) QByteArray
//	NewQByteArray("hello")          → из строки (UTF-8)
//	NewQByteArray([]byte{...})      → копия среза
//	NewQByteArray(n int, c byte)    → n байт, заполненных c
//	NewQByteArray(QChar)            → QChar::decomposition().toUtf8()
//	NewQByteArray(*QByteArray)      → COW-копия (Clone)
func NewQByteArray(args ...any) *QByteArray {
	switch len(args) {
	case 0:
		return wrapQByteArray(unsafe.Pointer(C.qbytearray_new_empty()))
	case 1:
		switch v := args[0].(type) {
		case []byte:
			return NewQByteArrayFromBytes(v)
		case string:
			return NewQByteArrayFromString(v)
		case QChar:
			return NewQByteArrayFromQCharDecomposition(v)
		case *QByteArray:
			return v.Clone()
		}
	case 2:
		n, ok1 := args[0].(int)
		c, ok2 := args[1].(byte)
		if ok1 && ok2 {
			return NewQByteArrayFromFill(n, c)
		}
	}
	panic("NewQByteArray: unsupported argument combination")
}

// NewQByteArrayNumberLL — Qt: QByteArray::number(qlonglong, int).
func NewQByteArrayNumberLL(n int64, base int) *QByteArray {
	return wrapQByteArray(unsafe.Pointer(C.qbytearray_number_ll(C.int64_t(n), C.int(base))))
}

// NewQByteArrayNumberULL — Qt: QByteArray::number(qulonglong, int).
func NewQByteArrayNumberULL(n uint64, base int) *QByteArray {
	return wrapQByteArray(unsafe.Pointer(C.qbytearray_number_ull(C.uint64_t(n), C.int(base))))
}

// NewQByteArrayNumberInt — Qt: QByteArray::number(int, int).
func NewQByteArrayNumberInt(n int32, base int) *QByteArray {
	return wrapQByteArray(unsafe.Pointer(C.qbytearray_number_int(C.int32_t(n), C.int(base))))
}

// NewQByteArrayNumberUInt — Qt: QByteArray::number(uint, int).
func NewQByteArrayNumberUInt(n uint32, base int) *QByteArray {
	return wrapQByteArray(unsafe.Pointer(C.qbytearray_number_uint(C.uint32_t(n), C.int(base))))
}

// NewQByteArrayNumberDouble — Qt: QByteArray::number(double, char, int).
func NewQByteArrayNumberDouble(v float64, format byte, precision int) *QByteArray {
	return wrapQByteArray(unsafe.Pointer(
		C.qbytearray_number_double(C.double(v), C.char(format), C.int(precision))))
}

// NewQByteArrayFromBase64 — Qt: QByteArray::fromBase64(const QByteArray &, Base64Options).
func NewQByteArrayFromBase64(b []byte, options Base64Options) *QByteArray {
	src := NewQByteArrayFromBytes(b)
	defer src.Close()
	return wrapQByteArray(unsafe.Pointer(C.qbytearray_from_base64(src.handle(), C.int(options))))
}

// NewQByteArrayFromBase64Encoding повторяет QByteArray::fromBase64Encoding и
// возвращает декодированный массив вместе со статусом.
// Qt: QByteArray::fromBase64Encoding(QByteArray &&, Base64Options) и его
// const-перегрузка; вместо FromBase64Result — два значения.
func NewQByteArrayFromBase64Encoding(b []byte, options Base64Options) (*QByteArray, Base64DecodingStatus) {
	var p *C.char
	if len(b) > 0 {
		p = (*C.char)(unsafe.Pointer(&b[0]))
	}
	var status C.int
	r := wrapQByteArray(unsafe.Pointer(
		C.qbytearray_from_base64_encoding(p, C.ptrdiff_t(len(b)), C.int(options), &status)))
	runtime.KeepAlive(b)
	return r, Base64DecodingStatus(status)
}

// NewQByteArrayFromHex — Qt: QByteArray::fromHex(const QByteArray &).
func NewQByteArrayFromHex(b []byte) *QByteArray {
	src := NewQByteArrayFromBytes(b)
	defer src.Close()
	return wrapQByteArray(unsafe.Pointer(C.qbytearray_from_hex(src.handle())))
}

// NewQByteArrayFromPercentEncoding — Qt: QByteArray::fromPercentEncoding(const QByteArray &, char).
func NewQByteArrayFromPercentEncoding(b []byte, percent byte) *QByteArray {
	src := NewQByteArrayFromBytes(b)
	defer src.Close()
	return wrapQByteArray(unsafe.Pointer(C.qbytearray_from_percent(src.handle(), C.char(percent))))
}

// ---- lifecycle -----------------------------------------------------------

// Close немедленно освобождает C++ объект и отменяет cleanup. Close опционален,
// идемпотентен и безопасен для nil. Если Close гоняется с cleanup, ровно один
// из них забирает объект; двойного free нет.
// Qt: QByteArray::~QByteArray().
func (q *QByteArray) Close() {
	if q == nil || q.state == nil {
		return
	}
	q.cleanup.Stop()
	if cpp := q.state.claim(); cpp != nil {
		C.qbytearray_delete(cpp)
	}
}

// Clone возвращает поверхностную COW-копию (refcount++; дёшево).
// Qt: QByteArray(const QByteArray &) — конструктор копирования.
func (q *QByteArray) Clone() *QByteArray {
	r := wrapQByteArray(unsafe.Pointer(C.qbytearray_new_copy(q.handle())))
	runtime.KeepAlive(q)
	return r
}

// CloneBytes возвращает []byte-копию содержимого (Go-вариант Clone).
// Qt: QByteArray(const QByteArray &), результат скопирован в []byte.
func (q *QByteArray) CloneBytes() []byte {
	r := q.Clone()
	defer r.Close()
	return r.Bytes()
}

// Swap обменивает содержимое q и other без копирования (COW-указатели).
// Паникует, если other == nil или закрыт.
// Qt: QByteArray::swap(QByteArray &).
func (q *QByteArray) Swap(other *QByteArray) {
	op, ok := other.tryHandle()
	if !ok {
		panic("qbytearray: Swap with nil or closed QByteArray")
	}
	C.qbytearray_swap(q.handle(), op)
	runtime.KeepAlive(q)
	runtime.KeepAlive(other)
}

// AssignBytes заменяет содержимое q копией b. nil делает массив null, не-nil
// пустой срез — пустым, но не null.
// Qt: QByteArray &QByteArray::assign(QByteArrayView).
func (q *QByteArray) AssignBytes(b []byte) {
	switch {
	case b == nil:
		q.Clear()
		return
	case len(b) == 0:
		C.qbytearray_assign_empty(q.handle())
		runtime.KeepAlive(q)
		return
	}
	p := (*C.char)(unsafe.Pointer(&b[0]))
	C.qbytearray_assign_data(q.handle(), p, C.ptrdiff_t(len(b)))
	runtime.KeepAlive(q)
	runtime.KeepAlive(b)
}

// AssignString заменяет содержимое q байтами строки s.
// Qt: QByteArray &QByteArray::operator=(const char *).
func (q *QByteArray) AssignString(s string) {
	q.AssignBytes([]byte(s))
}

// AssignFill заменяет содержимое q массивом из n байт c.
// Паникует при n < 0.
// Qt: QByteArray &QByteArray::assign(qsizetype, char).
func (q *QByteArray) AssignFill(n int, c byte) {
	if n < 0 {
		panic("qbytearray: AssignFill negative size")
	}
	C.qbytearray_assign_fill(q.handle(), C.ptrdiff_t(n), C.char(c))
	runtime.KeepAlive(q)
}

// AssignQByteArray заменяет содержимое q копией other (COW).
// Qt: QByteArray &QByteArray::operator=(const QByteArray &).
func (q *QByteArray) AssignQByteArray(other *QByteArray) {
	C.qbytearray_assign_ba(q.handle(), other.handle())
	runtime.KeepAlive(q)
	runtime.KeepAlive(other)
}

// ---- basic accessors -----------------------------------------------------

// Size возвращает число байт (без завершающего NUL).
// Qt: QByteArray::size() const.
func (q *QByteArray) Size() int {
	n := C.qbytearray_size(q.handle())
	runtime.KeepAlive(q)
	return int(n)
}

// Len — синоним Size.
// Qt: QByteArray::length() const.
func (q *QByteArray) Len() int { return q.Size() }

// Capacity возвращает вместимость внутреннего буфера.
// Qt: QByteArray::capacity() const.
func (q *QByteArray) Capacity() int {
	n := C.qbytearray_capacity(q.handle())
	runtime.KeepAlive(q)
	return int(n)
}

// IsEmpty сообщает, что размер равен нулю.
// Qt: QByteArray::isEmpty() const.
func (q *QByteArray) IsEmpty() bool {
	v := C.qbytearray_is_empty(q.handle())
	runtime.KeepAlive(q)
	return bool(v)
}

// IsNull сообщает, что массив null (не имеет данных); null-массив всегда пуст,
// но пустой массив не всегда null.
// Qt: QByteArray::isNull() const.
func (q *QByteArray) IsNull() bool {
	v := C.qbytearray_is_null(q.handle())
	runtime.KeepAlive(q)
	return bool(v)
}

// IsDetached сообщает, что refcount == 1 (буфер не разделяется).
// Qt: QByteArray::isDetached() const.
func (q *QByteArray) IsDetached() bool {
	v := C.qbytearray_is_detached(q.handle())
	runtime.KeepAlive(q)
	return bool(v)
}

// Detach принудительно отсоединяет буфер (O(n) при refcount > 1).
// Qt: QByteArray::detach().
func (q *QByteArray) Detach() {
	C.qbytearray_detach(q.handle())
	runtime.KeepAlive(q)
}

// IsSharedWith сообщает, делят ли q и other один буфер (указатель и размер).
// Возвращает false, если other == nil или закрыт. []byte-варианта нет: у среза
// нет COW-буфера Qt, сравнение указателей не имело бы смысла.
// Qt: QByteArray::isSharedWith(const QByteArray &) const.
func (q *QByteArray) IsSharedWith(other *QByteArray) bool {
	p := q.handle()
	op, ok := other.tryHandle()
	if !ok {
		runtime.KeepAlive(q)
		return false
	}
	v := C.qbytearray_is_shared_with(p, op)
	runtime.KeepAlive(q)
	runtime.KeepAlive(other)
	return bool(v)
}

// QByteArrayMaxSize возвращает максимально допустимый размер массива.
// Qt: QByteArray::maxSize() static.
func QByteArrayMaxSize() int { return int(C.qbytearray_max_size()) }

// Bytes возвращает копию байт q ([]byte). Для пустого массива — nil.
// Qt: аналога нет — constData() + копирование в Go.
func (q *QByteArray) Bytes() []byte {
	n := C.qbytearray_size(q.handle())
	if n == 0 {
		runtime.KeepAlive(q)
		return nil
	}
	data := C.qbytearray_data(q.handle())
	b := C.GoBytes(unsafe.Pointer(data), C.int(n))
	runtime.KeepAlive(q)
	return b
}

// String возвращает копию байт q как Go-строку.
// Qt: аналога нет — constData() + копирование в Go (std::string).
func (q *QByteArray) String() string {
	n := C.qbytearray_size(q.handle())
	if n == 0 {
		runtime.KeepAlive(q)
		return ""
	}
	data := C.qbytearray_data(q.handle())
	s := C.GoStringN(data, C.int(n))
	runtime.KeepAlive(q)
	return s
}

// ---- modifiers -----------------------------------------------------------

// Resize изменяет размер до n; новые байты не инициализированы.
// Qt: QByteArray::resize(qsizetype).
func (q *QByteArray) Resize(n int) {
	C.qbytearray_resize(q.handle(), C.ptrdiff_t(n))
	runtime.KeepAlive(q)
}

// ResizeFill изменяет размер до n, заполняя новые байты c.
// Qt: QByteArray::resize(qsizetype, char).
func (q *QByteArray) ResizeFill(n int, c byte) {
	C.qbytearray_resize_fill(q.handle(), C.ptrdiff_t(n), C.char(c))
	runtime.KeepAlive(q)
}

// ResizeForOverwrite изменяет размер до n, не инициализируя память (для
// последующей записи). Qt: QByteArray::resizeForOverwrite(qsizetype).
func (q *QByteArray) ResizeForOverwrite(n int) {
	C.qbytearray_resize_for_overwrite(q.handle(), C.ptrdiff_t(n))
	runtime.KeepAlive(q)
}

// Fill заполняет массив c. При n < 0 размер сохраняется (Qt-умолчание -1),
// иначе массив сначала расширяется/усекается до n.
// Qt: QByteArray::fill(char, qsizetype).
func (q *QByteArray) Fill(c byte, n int) {
	C.qbytearray_fill(q.handle(), C.char(c), C.ptrdiff_t(n))
	runtime.KeepAlive(q)
}

// Clear очищает массив и делает его null.
// Qt: QByteArray::clear().
func (q *QByteArray) Clear() {
	C.qbytearray_clear(q.handle())
	runtime.KeepAlive(q)
}

// Reserve резервирует память как минимум под n байт.
// Qt: QByteArray::reserve(qsizetype).
func (q *QByteArray) Reserve(n int) {
	C.qbytearray_reserve(q.handle(), C.ptrdiff_t(n))
	runtime.KeepAlive(q)
}

// Squeeze освобождает неиспользуемую вместимость.
// Qt: QByteArray::squeeze().
func (q *QByteArray) Squeeze() {
	C.qbytearray_squeeze(q.handle())
	runtime.KeepAlive(q)
}

// Truncate усекает массив до pos байт (если pos >= Size, ничего не делает).
// Qt: QByteArray::truncate(qsizetype).
func (q *QByteArray) Truncate(pos int) {
	C.qbytearray_truncate(q.handle(), C.ptrdiff_t(pos))
	runtime.KeepAlive(q)
}

// Chop удаляет n последних байт (если n > Size, массив становится пустым).
// Qt: QByteArray::chop(qsizetype).
func (q *QByteArray) Chop(n int) {
	C.qbytearray_chop(q.handle(), C.ptrdiff_t(n))
	runtime.KeepAlive(q)
}

// NullTerminate гарантирует наличие завершающего '\0' по адресу data()[Size()].
// В размер массива терминатор не входит (как и в Qt), поэтому Size() и Bytes()
// его не показывают.
// Qt: QByteArray::nullTerminate().
func (q *QByteArray) NullTerminate() {
	C.qbytearray_null_terminate(q.handle())
	runtime.KeepAlive(q)
}

// NullTerminated возвращает копию, у которой гарантирован завершающий '\0'
// (вне Size()).
// Qt: QByteArray::nullTerminated() const &.
func (q *QByteArray) NullTerminated() *QByteArray {
	r := wrapQByteArray(unsafe.Pointer(C.qbytearray_null_terminated(q.handle())))
	runtime.KeepAlive(q)
	return r
}

// NullTerminatedBytes — Go-вариант NullTerminated: возвращает копию байт.
// Завершающий '\0' лежит вне Size() и в []byte не попадает.
// Qt: QByteArray::nullTerminated() const &, результат скопирован в []byte.
func (q *QByteArray) NullTerminatedBytes() []byte {
	r := q.NullTerminated()
	defer r.Close()
	return r.Bytes()
}

// ---- element access ------------------------------------------------------

// At возвращает i-й байт. Паникует, если i вне [0, Size()).
// Qt: QByteArray::at(qsizetype) const.
func (q *QByteArray) At(i int) byte {
	n := q.Size()
	if i < 0 || i >= n {
		panic("qbytearray: At out of range")
	}
	c := C.qbytearray_at(q.handle(), C.ptrdiff_t(i))
	runtime.KeepAlive(q)
	return byte(c)
}

// SetAt записывает i-й байт. Паникует, если i вне [0, Size()).
// Qt: QByteArray::operator[](qsizetype) — mutable-перегрузка.
func (q *QByteArray) SetAt(i int, c byte) {
	n := q.Size()
	if i < 0 || i >= n {
		panic("qbytearray: SetAt out of range")
	}
	C.qbytearray_set_at(q.handle(), C.ptrdiff_t(i), C.char(c))
	runtime.KeepAlive(q)
}

// Front возвращает первый байт. Паникует, если массив пуст.
// Qt: QByteArray::front() const.
func (q *QByteArray) Front() byte {
	if q.IsEmpty() {
		panic("qbytearray: Front on empty QByteArray")
	}
	c := C.qbytearray_front(q.handle())
	runtime.KeepAlive(q)
	return byte(c)
}

// Back возвращает последний байт. Паникует, если массив пуст.
// Qt: QByteArray::back() const.
func (q *QByteArray) Back() byte {
	if q.IsEmpty() {
		panic("qbytearray: Back on empty QByteArray")
	}
	c := C.qbytearray_back(q.handle())
	runtime.KeepAlive(q)
	return byte(c)
}

// SetFront записывает первый байт. Паникует, если массив пуст.
// Qt: QByteArray::front() — mutable-перегрузка (char &).
func (q *QByteArray) SetFront(c byte) {
	if q.IsEmpty() {
		panic("qbytearray: SetFront on empty QByteArray")
	}
	C.qbytearray_set_front(q.handle(), C.char(c))
	runtime.KeepAlive(q)
}

// SetBack записывает последний байт. Паникует, если массив пуст.
// Qt: QByteArray::back() — mutable-перегрузка (char &).
func (q *QByteArray) SetBack(c byte) {
	if q.IsEmpty() {
		panic("qbytearray: SetBack on empty QByteArray")
	}
	C.qbytearray_set_back(q.handle(), C.char(c))
	runtime.KeepAlive(q)
}

// ---- search --------------------------------------------------------------

// IndexOf возвращает позицию первого вхождения c, начиная с from, или -1.
// Qt: QByteArray::indexOf(char, qsizetype) const.
func (q *QByteArray) IndexOf(c byte, from int) int {
	n := C.qbytearray_index_of_char(q.handle(), C.char(c), C.ptrdiff_t(from))
	runtime.KeepAlive(q)
	return int(n)
}

// IndexOfBytes возвращает позицию первого вхождения sub, начиная с from, или -1.
// Qt: QByteArray::indexOf(QByteArrayView, qsizetype) const.
func (q *QByteArray) IndexOfBytes(sub []byte, from int) int {
	var p *C.char
	if len(sub) > 0 {
		p = (*C.char)(unsafe.Pointer(&sub[0]))
	}
	n := C.qbytearray_index_of_data(q.handle(), p, C.ptrdiff_t(len(sub)), C.ptrdiff_t(from))
	runtime.KeepAlive(q)
	runtime.KeepAlive(sub)
	return int(n)
}

// LastIndexOf возвращает позицию последнего вхождения c не позже from, или -1.
// from < 0 означает поиск с конца (Qt-умолчание -1).
// Qt: QByteArray::lastIndexOf(char, qsizetype) const.
func (q *QByteArray) LastIndexOf(c byte, from int) int {
	n := C.qbytearray_last_index_of_char(q.handle(), C.char(c), C.ptrdiff_t(from))
	runtime.KeepAlive(q)
	return int(n)
}

// LastIndexOfBytes возвращает позицию последнего вхождения sub не позже from,
// или -1. from < 0 означает поиск с конца.
// Qt: QByteArray::lastIndexOf(QByteArrayView, qsizetype) const.
func (q *QByteArray) LastIndexOfBytes(sub []byte, from int) int {
	var p *C.char
	if len(sub) > 0 {
		p = (*C.char)(unsafe.Pointer(&sub[0]))
	}
	n := C.qbytearray_last_index_of_data(q.handle(), p, C.ptrdiff_t(len(sub)), C.ptrdiff_t(from))
	runtime.KeepAlive(q)
	runtime.KeepAlive(sub)
	return int(n)
}

// Count возвращает число вхождений c.
// Qt: QByteArray::count(char) const.
func (q *QByteArray) Count(c byte) int {
	n := C.qbytearray_count_char(q.handle(), C.char(c))
	runtime.KeepAlive(q)
	return int(n)
}

// CountBytes возвращает число (возможно, перекрывающихся) вхождений sub.
// Qt: QByteArray::count(QByteArrayView) const.
func (q *QByteArray) CountBytes(sub []byte) int {
	var p *C.char
	if len(sub) > 0 {
		p = (*C.char)(unsafe.Pointer(&sub[0]))
	}
	n := C.qbytearray_count_data(q.handle(), p, C.ptrdiff_t(len(sub)))
	runtime.KeepAlive(q)
	runtime.KeepAlive(sub)
	return int(n)
}

// Contains сообщает, есть ли байт c.
// Qt: QByteArray::contains(char) const.
func (q *QByteArray) Contains(c byte) bool {
	v := C.qbytearray_contains_char(q.handle(), C.char(c))
	runtime.KeepAlive(q)
	return bool(v)
}

// ContainsBytes сообщает, есть ли подпоследовательность sub.
// Qt: QByteArray::contains(QByteArrayView) const.
func (q *QByteArray) ContainsBytes(sub []byte) bool {
	var p *C.char
	if len(sub) > 0 {
		p = (*C.char)(unsafe.Pointer(&sub[0]))
	}
	v := C.qbytearray_contains_data(q.handle(), p, C.ptrdiff_t(len(sub)))
	runtime.KeepAlive(q)
	runtime.KeepAlive(sub)
	return bool(v)
}

// Compare возвращает результат сравнения как в Qt: <0, 0, >0. cs == true —
// с учётом регистра.
// Qt: QByteArray::compare(QByteArrayView, Qt::CaseSensitivity) const.
func (q *QByteArray) Compare(other []byte, cs bool) int {
	var p *C.char
	if len(other) > 0 {
		p = (*C.char)(unsafe.Pointer(&other[0]))
	}
	csInt := C.int(0)
	if !cs {
		csInt = 1
	}
	r := C.qbytearray_compare(q.handle(), p, C.ptrdiff_t(len(other)), csInt)
	runtime.KeepAlive(q)
	runtime.KeepAlive(other)
	return int(r)
}

// ---- slicing -------------------------------------------------------------

// Left возвращает первые n байт (n зажимается в [0, Size()]).
// Qt: QByteArray::left(qsizetype) const &.
func (q *QByteArray) Left(n int) *QByteArray {
	r := wrapQByteArray(unsafe.Pointer(C.qbytearray_left(q.handle(), C.ptrdiff_t(n))))
	runtime.KeepAlive(q)
	return r
}

// LeftBytes — Go-вариант Left: первые n байт копией в []byte.
// Qt: QByteArray::left(qsizetype) const &, результат скопирован в []byte.
func (q *QByteArray) LeftBytes(n int) []byte {
	r := q.Left(n)
	defer r.Close()
	return r.Bytes()
}

// Right возвращает последние n байт (n зажимается в [0, Size()]).
// Qt: QByteArray::right(qsizetype) const &.
func (q *QByteArray) Right(n int) *QByteArray {
	r := wrapQByteArray(unsafe.Pointer(C.qbytearray_right(q.handle(), C.ptrdiff_t(n))))
	runtime.KeepAlive(q)
	return r
}

// RightBytes — Go-вариант Right: последние n байт копией в []byte.
// Qt: QByteArray::right(qsizetype) const &, результат скопирован в []byte.
func (q *QByteArray) RightBytes(n int) []byte {
	r := q.Right(n)
	defer r.Close()
	return r.Bytes()
}

// Mid возвращает подмассив с позиции i длиной length. length < 0 означает «до
// конца» (Qt-умолчание -1); границы зажимаются.
// Qt: QByteArray::mid(qsizetype, qsizetype) const &.
func (q *QByteArray) Mid(i, length int) *QByteArray {
	r := wrapQByteArray(unsafe.Pointer(C.qbytearray_mid(q.handle(), C.ptrdiff_t(i), C.ptrdiff_t(length))))
	runtime.KeepAlive(q)
	return r
}

// MidBytes — Go-вариант Mid: подмассив с позиции i длиной length в []byte.
// Qt: QByteArray::mid(qsizetype, qsizetype) const &, результат скопирован в []byte.
func (q *QByteArray) MidBytes(i, length int) []byte {
	r := q.Mid(i, length)
	defer r.Close()
	return r.Bytes()
}

// First возвращает первые n байт; паникует, если n вне [0, Size()].
// Qt: QByteArray::first(qsizetype) const &.
func (q *QByteArray) First(n int) *QByteArray {
	if n < 0 || n > q.Size() {
		panic("qbytearray: First out of range")
	}
	r := wrapQByteArray(unsafe.Pointer(C.qbytearray_first(q.handle(), C.ptrdiff_t(n))))
	runtime.KeepAlive(q)
	return r
}

// FirstBytes — Go-вариант First: первые n байт копией в []byte.
// Qt: QByteArray::first(qsizetype) const &, результат скопирован в []byte.
func (q *QByteArray) FirstBytes(n int) []byte {
	r := q.First(n)
	defer r.Close()
	return r.Bytes()
}

// Last возвращает последние n байт; паникует, если n вне [0, Size()].
// Qt: QByteArray::last(qsizetype) const &.
func (q *QByteArray) Last(n int) *QByteArray {
	if n < 0 || n > q.Size() {
		panic("qbytearray: Last out of range")
	}
	r := wrapQByteArray(unsafe.Pointer(C.qbytearray_last(q.handle(), C.ptrdiff_t(n))))
	runtime.KeepAlive(q)
	return r
}

// LastBytes — Go-вариант Last: последние n байт копией в []byte.
// Qt: QByteArray::last(qsizetype) const &, результат скопирован в []byte.
func (q *QByteArray) LastBytes(n int) []byte {
	r := q.Last(n)
	defer r.Close()
	return r.Bytes()
}

// Sliced возвращает подмассив с позиции pos до конца; паникует, если pos вне
// [0, Size()].
// Qt: QByteArray::sliced(qsizetype) const &.
func (q *QByteArray) Sliced(pos int) *QByteArray {
	if pos < 0 || pos > q.Size() {
		panic("qbytearray: Sliced out of range")
	}
	r := wrapQByteArray(unsafe.Pointer(C.qbytearray_sliced(q.handle(), C.ptrdiff_t(pos))))
	runtime.KeepAlive(q)
	return r
}

// SlicedBytes — Go-вариант Sliced: подмассив с позиции pos копией в []byte.
// Qt: QByteArray::sliced(qsizetype) const &, результат скопирован в []byte.
func (q *QByteArray) SlicedBytes(pos int) []byte {
	r := q.Sliced(pos)
	defer r.Close()
	return r.Bytes()
}

// SlicedN возвращает n байт с позиции pos; паникует, если диапазон выходит за
// пределы массива.
// Qt: QByteArray::sliced(qsizetype, qsizetype) const &.
func (q *QByteArray) SlicedN(pos, n int) *QByteArray {
	sz := q.Size()
	if pos < 0 || n < 0 || pos+n > sz {
		panic("qbytearray: SlicedN out of range")
	}
	r := wrapQByteArray(unsafe.Pointer(
		C.qbytearray_sliced_n(q.handle(), C.ptrdiff_t(pos), C.ptrdiff_t(n))))
	runtime.KeepAlive(q)
	return r
}

// SlicedNBytes — Go-вариант SlicedN: n байт с позиции pos копией в []byte.
// Qt: QByteArray::sliced(qsizetype, qsizetype) const &, результат скопирован в []byte.
func (q *QByteArray) SlicedNBytes(pos, n int) []byte {
	r := q.SlicedN(pos, n)
	defer r.Close()
	return r.Bytes()
}

// Chopped возвращает массив без n последних байт; паникует, если n вне
// [0, Size()].
// Qt: QByteArray::chopped(qsizetype) const &.
func (q *QByteArray) Chopped(n int) *QByteArray {
	if n < 0 || n > q.Size() {
		panic("qbytearray: Chopped out of range")
	}
	r := wrapQByteArray(unsafe.Pointer(C.qbytearray_chopped(q.handle(), C.ptrdiff_t(n))))
	runtime.KeepAlive(q)
	return r
}

// ChoppedBytes — Go-вариант Chopped: массив без n последних байт в []byte.
// Qt: QByteArray::chopped(qsizetype) const &, результат скопирован в []byte.
func (q *QByteArray) ChoppedBytes(n int) []byte {
	r := q.Chopped(n)
	defer r.Close()
	return r.Bytes()
}

// Slice обрезает q на месте, оставляя байты с позиции pos (удаляет префикс).
// Паникует, если pos вне [0, Size()].
// Qt: QByteArray &QByteArray::slice(qsizetype).
func (q *QByteArray) Slice(pos int) {
	if pos < 0 || pos > q.Size() {
		panic("qbytearray: Slice out of range")
	}
	C.qbytearray_slice(q.handle(), C.ptrdiff_t(pos))
	runtime.KeepAlive(q)
}

// SliceN обрезает q на месте до n байт, начиная с позиции pos.
// Паникует, если диапазон выходит за пределы массива.
// Qt: QByteArray &QByteArray::slice(qsizetype, qsizetype).
func (q *QByteArray) SliceN(pos, n int) {
	sz := q.Size()
	if pos < 0 || n < 0 || pos+n > sz {
		panic("qbytearray: SliceN out of range")
	}
	C.qbytearray_slice_n(q.handle(), C.ptrdiff_t(pos), C.ptrdiff_t(n))
	runtime.KeepAlive(q)
}

// ---- predicates ----------------------------------------------------------

// StartsWith сообщает, начинается ли массив с байта c.
// Qt: QByteArray::startsWith(char) const.
func (q *QByteArray) StartsWith(c byte) bool {
	v := C.qbytearray_starts_with_char(q.handle(), C.char(c))
	runtime.KeepAlive(q)
	return bool(v)
}

// StartsWithBytes сообщает, начинается ли массив с prefix.
// Qt: QByteArray::startsWith(QByteArrayView) const.
func (q *QByteArray) StartsWithBytes(prefix []byte) bool {
	var p *C.char
	if len(prefix) > 0 {
		p = (*C.char)(unsafe.Pointer(&prefix[0]))
	}
	v := C.qbytearray_starts_with_data(q.handle(), p, C.ptrdiff_t(len(prefix)))
	runtime.KeepAlive(q)
	runtime.KeepAlive(prefix)
	return bool(v)
}

// EndsWith сообщает, заканчивается ли массив байтом c.
// Qt: QByteArray::endsWith(char) const.
func (q *QByteArray) EndsWith(c byte) bool {
	v := C.qbytearray_ends_with_char(q.handle(), C.char(c))
	runtime.KeepAlive(q)
	return bool(v)
}

// EndsWithBytes сообщает, заканчивается ли массив на suffix.
// Qt: QByteArray::endsWith(QByteArrayView) const.
func (q *QByteArray) EndsWithBytes(suffix []byte) bool {
	var p *C.char
	if len(suffix) > 0 {
		p = (*C.char)(unsafe.Pointer(&suffix[0]))
	}
	v := C.qbytearray_ends_with_data(q.handle(), p, C.ptrdiff_t(len(suffix)))
	runtime.KeepAlive(q)
	runtime.KeepAlive(suffix)
	return bool(v)
}

// IsUpper сообщает, что все буквы — заглавные (не-буквы игнорируются).
// Qt: QByteArray::isUpper() const.
func (q *QByteArray) IsUpper() bool {
	v := C.qbytearray_is_upper(q.handle())
	runtime.KeepAlive(q)
	return bool(v)
}

// IsLower сообщает, что все буквы — строчные (не-буквы игнорируются).
// Qt: QByteArray::isLower() const.
func (q *QByteArray) IsLower() bool {
	v := C.qbytearray_is_lower(q.handle())
	runtime.KeepAlive(q)
	return bool(v)
}

// IsValidUtf8 сообщает, является ли содержимое корректным UTF-8.
// Qt: QByteArray::isValidUtf8() const.
func (q *QByteArray) IsValidUtf8() bool {
	v := C.qbytearray_is_valid_utf8(q.handle())
	runtime.KeepAlive(q)
	return bool(v)
}

// ---- transforms ----------------------------------------------------------

// ToLower возвращает копию в нижнем регистре (ASCII).
// Qt: QByteArray::toLower() const &.
func (q *QByteArray) ToLower() *QByteArray {
	r := wrapQByteArray(unsafe.Pointer(C.qbytearray_to_lower(q.handle())))
	runtime.KeepAlive(q)
	return r
}

// ToLowerBytes — Go-вариант ToLower: копия в нижнем регистре в []byte.
// Qt: QByteArray::toLower() const &, результат скопирован в []byte.
func (q *QByteArray) ToLowerBytes() []byte {
	r := q.ToLower()
	defer r.Close()
	return r.Bytes()
}

// ToUpper возвращает копию в верхнем регистре (ASCII).
// Qt: QByteArray::toUpper() const &.
func (q *QByteArray) ToUpper() *QByteArray {
	r := wrapQByteArray(unsafe.Pointer(C.qbytearray_to_upper(q.handle())))
	runtime.KeepAlive(q)
	return r
}

// ToUpperBytes — Go-вариант ToUpper: копия в верхнем регистре в []byte.
// Qt: QByteArray::toUpper() const &, результат скопирован в []byte.
func (q *QByteArray) ToUpperBytes() []byte {
	r := q.ToUpper()
	defer r.Close()
	return r.Bytes()
}

// Trimmed возвращает копию без пробельных байт в начале и конце.
// Qt: QByteArray::trimmed() const &.
func (q *QByteArray) Trimmed() *QByteArray {
	r := wrapQByteArray(unsafe.Pointer(C.qbytearray_trimmed(q.handle())))
	runtime.KeepAlive(q)
	return r
}

// TrimmedBytes — Go-вариант Trimmed: копия без пробелов по краям в []byte.
// Qt: QByteArray::trimmed() const &, результат скопирован в []byte.
func (q *QByteArray) TrimmedBytes() []byte {
	r := q.Trimmed()
	defer r.Close()
	return r.Bytes()
}

// Simplified возвращает копию, в которой каждый блок пробельных байт заменён
// одним пробелом, а края обрезаны.
// Qt: QByteArray::simplified() const &.
func (q *QByteArray) Simplified() *QByteArray {
	r := wrapQByteArray(unsafe.Pointer(C.qbytearray_simplified(q.handle())))
	runtime.KeepAlive(q)
	return r
}

// SimplifiedBytes — Go-вариант Simplified: результат в []byte.
// Qt: QByteArray::simplified() const &, результат скопирован в []byte.
func (q *QByteArray) SimplifiedBytes() []byte {
	r := q.Simplified()
	defer r.Close()
	return r.Bytes()
}

// LeftJustified выравнивает по левому краю до ширины width, заполняя fill;
// truncate управляет усечением длинных строк.
// Qt: QByteArray::leftJustified(qsizetype, char, bool) const.
func (q *QByteArray) LeftJustified(width int, fill byte, truncate bool) *QByteArray {
	r := wrapQByteArray(unsafe.Pointer(
		C.qbytearray_left_justified(q.handle(), C.ptrdiff_t(width), C.char(fill), C.bool(truncate))))
	runtime.KeepAlive(q)
	return r
}

// LeftJustifiedBytes — Go-вариант LeftJustified: результат в []byte.
// Qt: QByteArray::leftJustified(qsizetype, char, bool) const, результат скопирован в []byte.
func (q *QByteArray) LeftJustifiedBytes(width int, fill byte, truncate bool) []byte {
	r := q.LeftJustified(width, fill, truncate)
	defer r.Close()
	return r.Bytes()
}

// RightJustified выравнивает по правому краю до ширины width, заполняя fill;
// truncate управляет усечением длинных строк.
// Qt: QByteArray::rightJustified(qsizetype, char, bool) const.
func (q *QByteArray) RightJustified(width int, fill byte, truncate bool) *QByteArray {
	r := wrapQByteArray(unsafe.Pointer(
		C.qbytearray_right_justified(q.handle(), C.ptrdiff_t(width), C.char(fill), C.bool(truncate))))
	runtime.KeepAlive(q)
	return r
}

// RightJustifiedBytes — Go-вариант RightJustified: результат в []byte.
// Qt: QByteArray::rightJustified(qsizetype, char, bool) const, результат скопирован в []byte.
func (q *QByteArray) RightJustifiedBytes(width int, fill byte, truncate bool) []byte {
	r := q.RightJustified(width, fill, truncate)
	defer r.Close()
	return r.Bytes()
}

// Repeated возвращает массив, повторённый times раз.
// Qt: QByteArray::repeated(qsizetype) const.
func (q *QByteArray) Repeated(times int) *QByteArray {
	r := wrapQByteArray(unsafe.Pointer(C.qbytearray_repeated(q.handle(), C.ptrdiff_t(times))))
	runtime.KeepAlive(q)
	return r
}

// RepeatedBytes — Go-вариант Repeated: повторённый массив в []byte.
// Qt: QByteArray::repeated(qsizetype) const, результат скопирован в []byte.
func (q *QByteArray) RepeatedBytes(times int) []byte {
	r := q.Repeated(times)
	defer r.Close()
	return r.Bytes()
}

// Concat возвращает новый массив q + other.
// Qt: operator+(const QByteArray &, const QByteArray &).
func (q *QByteArray) Concat(other *QByteArray) *QByteArray {
	r := wrapQByteArray(unsafe.Pointer(C.qbytearray_concat(q.handle(), other.handle())))
	runtime.KeepAlive(q)
	runtime.KeepAlive(other)
	return r
}

// ConcatBytes — Go-вариант Concat: возвращает q + b копией в []byte.
// Qt: operator+(const QByteArray &, QByteArrayView), результат скопирован в []byte.
func (q *QByteArray) ConcatBytes(b []byte) []byte {
	var p *C.char
	if len(b) > 0 {
		p = (*C.char)(unsafe.Pointer(&b[0]))
	}
	r := wrapQByteArray(unsafe.Pointer(
		C.qbytearray_concat_data(q.handle(), p, C.ptrdiff_t(len(b)))))
	runtime.KeepAlive(q)
	runtime.KeepAlive(b)
	defer r.Close()
	return r.Bytes()
}

// ---- append/prepend/insert ----------------------------------------------

// Append добавляет байт c в конец.
// Qt: QByteArray &QByteArray::append(char).
func (q *QByteArray) Append(c byte) {
	C.qbytearray_append_char(q.handle(), C.char(c))
	runtime.KeepAlive(q)
}

// AppendBytes добавляет b в конец.
// Qt: QByteArray &QByteArray::append(QByteArrayView).
func (q *QByteArray) AppendBytes(b []byte) {
	var p *C.char
	if len(b) > 0 {
		p = (*C.char)(unsafe.Pointer(&b[0]))
	}
	C.qbytearray_append_data(q.handle(), p, C.ptrdiff_t(len(b)))
	runtime.KeepAlive(q)
	runtime.KeepAlive(b)
}

// AppendString добавляет строку s (байты UTF-8) в конец.
// Qt: QByteArray &QByteArray::append(const char *).
func (q *QByteArray) AppendString(s string) { q.AppendBytes([]byte(s)) }

// AppendCount добавляет count байт c в конец.
// Qt: QByteArray &QByteArray::append(qsizetype, char).
func (q *QByteArray) AppendCount(count int, c byte) {
	C.qbytearray_append_count(q.handle(), C.ptrdiff_t(count), C.char(c))
	runtime.KeepAlive(q)
}

// AppendQByteArray добавляет other в конец.
// Qt: QByteArray &QByteArray::append(const QByteArray &).
func (q *QByteArray) AppendQByteArray(other *QByteArray) {
	C.qbytearray_append_ba(q.handle(), other.handle())
	runtime.KeepAlive(q)
	runtime.KeepAlive(other)
}

// Prepend добавляет байт c в начало.
// Qt: QByteArray &QByteArray::prepend(char).
func (q *QByteArray) Prepend(c byte) {
	C.qbytearray_prepend_char(q.handle(), C.char(c))
	runtime.KeepAlive(q)
}

// PrependBytes добавляет b в начало.
// Qt: QByteArray &QByteArray::prepend(QByteArrayView).
func (q *QByteArray) PrependBytes(b []byte) {
	var p *C.char
	if len(b) > 0 {
		p = (*C.char)(unsafe.Pointer(&b[0]))
	}
	C.qbytearray_prepend_data(q.handle(), p, C.ptrdiff_t(len(b)))
	runtime.KeepAlive(q)
	runtime.KeepAlive(b)
}

// PrependString добавляет строку s (байты UTF-8) в начало.
// Qt: QByteArray &QByteArray::prepend(const char *).
func (q *QByteArray) PrependString(s string) { q.PrependBytes([]byte(s)) }

// PrependCount добавляет count байт c в начало.
// Qt: QByteArray &QByteArray::prepend(qsizetype, char).
func (q *QByteArray) PrependCount(count int, c byte) {
	C.qbytearray_prepend_count(q.handle(), C.ptrdiff_t(count), C.char(c))
	runtime.KeepAlive(q)
}

// PrependQByteArray добавляет other в начало.
// Qt: QByteArray &QByteArray::prepend(const QByteArray &).
func (q *QByteArray) PrependQByteArray(other *QByteArray) {
	C.qbytearray_prepend_ba(q.handle(), other.handle())
	runtime.KeepAlive(q)
	runtime.KeepAlive(other)
}

// InsertBytes вставляет b в позицию i. Паникует, если i вне [0, Size()].
// Qt: QByteArray &QByteArray::insert(qsizetype, QByteArrayView).
func (q *QByteArray) InsertBytes(i int, b []byte) {
	if i < 0 || i > q.Size() {
		panic("qbytearray: InsertBytes out of range")
	}
	var p *C.char
	if len(b) > 0 {
		p = (*C.char)(unsafe.Pointer(&b[0]))
	}
	C.qbytearray_insert_data(q.handle(), C.ptrdiff_t(i), p, C.ptrdiff_t(len(b)))
	runtime.KeepAlive(q)
	runtime.KeepAlive(b)
}

// InsertString вставляет строку s в позицию i. Паникует, если i вне [0, Size()].
// Qt: QByteArray &QByteArray::insert(qsizetype, const char *).
func (q *QByteArray) InsertString(i int, s string) { q.InsertBytes(i, []byte(s)) }

// InsertByte вставляет байт c в позицию i. Паникует, если i вне [0, Size()].
// Qt: QByteArray &QByteArray::insert(qsizetype, char).
func (q *QByteArray) InsertByte(i int, c byte) {
	if i < 0 || i > q.Size() {
		panic("qbytearray: InsertByte out of range")
	}
	C.qbytearray_insert_char(q.handle(), C.ptrdiff_t(i), C.char(c))
	runtime.KeepAlive(q)
}

// InsertCount вставляет count байт c в позицию i. Паникует, если i вне [0, Size()].
// Qt: QByteArray &QByteArray::insert(qsizetype, qsizetype, char).
func (q *QByteArray) InsertCount(i, count int, c byte) {
	if i < 0 || i > q.Size() {
		panic("qbytearray: InsertCount out of range")
	}
	C.qbytearray_insert_count(q.handle(), C.ptrdiff_t(i), C.ptrdiff_t(count), C.char(c))
	runtime.KeepAlive(q)
}

// InsertQByteArray вставляет other в позицию i. Паникует, если i вне [0, Size()].
// Qt: QByteArray &QByteArray::insert(qsizetype, const QByteArray &).
func (q *QByteArray) InsertQByteArray(i int, other *QByteArray) {
	if i < 0 || i > q.Size() {
		panic("qbytearray: InsertQByteArray out of range")
	}
	C.qbytearray_insert_ba(q.handle(), C.ptrdiff_t(i), other.handle())
	runtime.KeepAlive(q)
	runtime.KeepAlive(other)
}

// ---- remove/replace ------------------------------------------------------

// Remove удаляет n байт начиная с позиции i. Паникует, если диапазон выходит за
// пределы массива.
// Qt: QByteArray &QByteArray::remove(qsizetype, qsizetype).
func (q *QByteArray) Remove(i, n int) {
	if i < 0 || n < 0 || i+n > q.Size() {
		panic("qbytearray: Remove out of range")
	}
	C.qbytearray_remove(q.handle(), C.ptrdiff_t(i), C.ptrdiff_t(n))
	runtime.KeepAlive(q)
}

// RemoveAt удаляет байт в позиции pos (если позиция вне массива — ничего не делает).
// Qt: QByteArray &QByteArray::removeAt(qsizetype).
func (q *QByteArray) RemoveAt(pos int) {
	C.qbytearray_remove_at(q.handle(), C.ptrdiff_t(pos))
	runtime.KeepAlive(q)
}

// RemoveFirst удаляет первый байт (если массив пуст — ничего не делает).
// Qt: QByteArray &QByteArray::removeFirst().
func (q *QByteArray) RemoveFirst() {
	C.qbytearray_remove_first(q.handle())
	runtime.KeepAlive(q)
}

// RemoveLast удаляет последний байт (если массив пуст — ничего не делает).
// Qt: QByteArray &QByteArray::removeLast().
func (q *QByteArray) RemoveLast() {
	C.qbytearray_remove_last(q.handle())
	runtime.KeepAlive(q)
}

// RemoveIf удаляет все байты, для которых pred возвращает true, и возвращает
// число удалённых байт. Это Go-реализация шаблонных removeIf/erase_if: готовый
// массив пересобирается и записывается через AssignBytes.
// Qt: QByteArray &QByteArray::removeIf(Predicate), qsizetype erase_if(QByteArray &, Predicate).
func (q *QByteArray) RemoveIf(pred func(byte) bool) int {
	if pred == nil {
		panic("qbytearray: RemoveIf with nil predicate")
	}
	src := q.Bytes()
	kept := src[:0]
	for _, b := range src {
		if !pred(b) {
			kept = append(kept, b)
		}
	}
	removed := len(src) - len(kept)
	if removed > 0 {
		q.AssignBytes(kept)
	}
	return removed
}

// ReplaceChar заменяет все байты before на after.
// Qt: QByteArray &QByteArray::replace(char, char).
func (q *QByteArray) ReplaceChar(before, after byte) {
	C.qbytearray_replace_char_char(q.handle(), C.char(before), C.char(after))
	runtime.KeepAlive(q)
}

// ReplaceCharBytes заменяет все вхождения байта before на after.
// Qt: QByteArray &QByteArray::replace(char, QByteArrayView).
func (q *QByteArray) ReplaceCharBytes(before byte, after []byte) {
	var p *C.char
	if len(after) > 0 {
		p = (*C.char)(unsafe.Pointer(&after[0]))
	}
	C.qbytearray_replace_char_data(q.handle(), C.char(before), p, C.ptrdiff_t(len(after)))
	runtime.KeepAlive(q)
	runtime.KeepAlive(after)
}

// ReplaceRange заменяет n байт с позиции i на with. Паникует, если диапазон
// выходит за пределы массива.
// Qt: QByteArray &QByteArray::replace(qsizetype, qsizetype, QByteArrayView).
func (q *QByteArray) ReplaceRange(i, n int, with []byte) {
	if i < 0 || n < 0 || i+n > q.Size() {
		panic("qbytearray: ReplaceRange out of range")
	}
	var p *C.char
	if len(with) > 0 {
		p = (*C.char)(unsafe.Pointer(&with[0]))
	}
	C.qbytearray_replace_range(q.handle(), C.ptrdiff_t(i), C.ptrdiff_t(n), p, C.ptrdiff_t(len(with)))
	runtime.KeepAlive(q)
	runtime.KeepAlive(with)
}

// ReplaceAll заменяет все вхождения before на after.
// Qt: QByteArray &QByteArray::replace(QByteArrayView, QByteArrayView).
func (q *QByteArray) ReplaceAll(before, after []byte) {
	var pb, pa *C.char
	if len(before) > 0 {
		pb = (*C.char)(unsafe.Pointer(&before[0]))
	}
	if len(after) > 0 {
		pa = (*C.char)(unsafe.Pointer(&after[0]))
	}
	C.qbytearray_replace_data_data(q.handle(), pb, C.ptrdiff_t(len(before)), pa, C.ptrdiff_t(len(after)))
	runtime.KeepAlive(q)
	runtime.KeepAlive(before)
	runtime.KeepAlive(after)
}

// ReplaceAllString заменяет все вхождения строки before на after.
// Qt: QByteArray &QByteArray::replace(QByteArrayView, QByteArrayView) с const char*.
func (q *QByteArray) ReplaceAllString(before, after string) {
	q.ReplaceAll([]byte(before), []byte(after))
}

// ---- numeric conversion --------------------------------------------------

// ToShort — Qt: QByteArray::toShort(bool *, int) const.
func (q *QByteArray) ToShort(base int) (int16, bool) {
	var ok C.bool
	v := C.qbytearray_to_short(q.handle(), &ok, C.int(base))
	runtime.KeepAlive(q)
	return int16(v), bool(ok)
}

// ToUShort — Qt: QByteArray::toUShort(bool *, int) const.
func (q *QByteArray) ToUShort(base int) (uint16, bool) {
	var ok C.bool
	v := C.qbytearray_to_ushort(q.handle(), &ok, C.int(base))
	runtime.KeepAlive(q)
	return uint16(v), bool(ok)
}

// ToInt — Qt: QByteArray::toInt(bool *, int) const.
func (q *QByteArray) ToInt(base int) (int32, bool) {
	var ok C.bool
	v := C.qbytearray_to_int(q.handle(), &ok, C.int(base))
	runtime.KeepAlive(q)
	return int32(v), bool(ok)
}

// ToUInt — Qt: QByteArray::toUInt(bool *, int) const.
func (q *QByteArray) ToUInt(base int) (uint32, bool) {
	var ok C.bool
	v := C.qbytearray_to_uint(q.handle(), &ok, C.int(base))
	runtime.KeepAlive(q)
	return uint32(v), bool(ok)
}

// ToLong — Qt: QByteArray::toLong(bool *, int) const.
func (q *QByteArray) ToLong(base int) (int64, bool) {
	var ok C.bool
	v := C.qbytearray_to_long(q.handle(), &ok, C.int(base))
	runtime.KeepAlive(q)
	return int64(v), bool(ok)
}

// ToULong — Qt: QByteArray::toULong(bool *, int) const.
func (q *QByteArray) ToULong(base int) (uint64, bool) {
	var ok C.bool
	v := C.qbytearray_to_ulong(q.handle(), &ok, C.int(base))
	runtime.KeepAlive(q)
	return uint64(v), bool(ok)
}

// ToLongLong — Qt: QByteArray::toLongLong(bool *, int) const.
func (q *QByteArray) ToLongLong(base int) (int64, bool) {
	var ok C.bool
	v := C.qbytearray_to_longlong(q.handle(), &ok, C.int(base))
	runtime.KeepAlive(q)
	return int64(v), bool(ok)
}

// ToULongLong — Qt: QByteArray::toULongLong(bool *, int) const.
func (q *QByteArray) ToULongLong(base int) (uint64, bool) {
	var ok C.bool
	v := C.qbytearray_to_ulonglong(q.handle(), &ok, C.int(base))
	runtime.KeepAlive(q)
	return uint64(v), bool(ok)
}

// ToFloat — Qt: QByteArray::toFloat(bool *) const.
func (q *QByteArray) ToFloat() (float32, bool) {
	var ok C.bool
	v := C.qbytearray_to_float(q.handle(), &ok)
	runtime.KeepAlive(q)
	return float32(v), bool(ok)
}

// ToDouble — Qt: QByteArray::toDouble(bool *) const.
func (q *QByteArray) ToDouble() (float64, bool) {
	var ok C.bool
	v := C.qbytearray_to_double(q.handle(), &ok)
	runtime.KeepAlive(q)
	return float64(v), bool(ok)
}

// SetNumLL — Qt: QByteArray::setNum(qlonglong, int).
func (q *QByteArray) SetNumLL(n int64, base int) {
	C.qbytearray_set_num_ll(q.handle(), C.int64_t(n), C.int(base))
	runtime.KeepAlive(q)
}

// SetNumULL — Qt: QByteArray::setNum(qulonglong, int).
func (q *QByteArray) SetNumULL(n uint64, base int) {
	C.qbytearray_set_num_ull(q.handle(), C.uint64_t(n), C.int(base))
	runtime.KeepAlive(q)
}

// SetNumInt — Qt: QByteArray::setNum(int, int).
func (q *QByteArray) SetNumInt(n int32, base int) {
	C.qbytearray_set_num_int(q.handle(), C.int32_t(n), C.int(base))
	runtime.KeepAlive(q)
}

// SetNumUInt — Qt: QByteArray::setNum(uint, int).
func (q *QByteArray) SetNumUInt(n uint32, base int) {
	C.qbytearray_set_num_uint(q.handle(), C.uint32_t(n), C.int(base))
	runtime.KeepAlive(q)
}

// SetNumFloat — Qt: QByteArray::setNum(float, char, int).
func (q *QByteArray) SetNumFloat(v float32, format byte, precision int) {
	C.qbytearray_set_num_float(q.handle(), C.float(v), C.char(format), C.int(precision))
	runtime.KeepAlive(q)
}

// SetNumDouble — Qt: QByteArray::setNum(double, char, int).
func (q *QByteArray) SetNumDouble(v float64, format byte, precision int) {
	C.qbytearray_set_num_double(q.handle(), C.double(v), C.char(format), C.int(precision))
	runtime.KeepAlive(q)
}

// ---- encoding ------------------------------------------------------------

// ToBase64 — Qt: QByteArray::toBase64(Base64Options) const.
func (q *QByteArray) ToBase64(options Base64Options) *QByteArray {
	r := wrapQByteArray(unsafe.Pointer(C.qbytearray_to_base64(q.handle(), C.int(options))))
	runtime.KeepAlive(q)
	return r
}

// ToBase64Bytes — Go-вариант ToBase64: Base64-текст копией в []byte.
// Qt: QByteArray::toBase64(Base64Options) const, результат скопирован в []byte.
func (q *QByteArray) ToBase64Bytes(options Base64Options) []byte {
	r := q.ToBase64(options)
	defer r.Close()
	return r.Bytes()
}

// ToHex — Qt: QByteArray::toHex(char) const.
func (q *QByteArray) ToHex(separator byte) *QByteArray {
	r := wrapQByteArray(unsafe.Pointer(C.qbytearray_to_hex(q.handle(), C.char(separator))))
	runtime.KeepAlive(q)
	return r
}

// ToHexBytes — Go-вариант ToHex: hex-текст копией в []byte.
// Qt: QByteArray::toHex(char) const, результат скопирован в []byte.
func (q *QByteArray) ToHexBytes(separator byte) []byte {
	r := q.ToHex(separator)
	defer r.Close()
	return r.Bytes()
}

// ToPercentEncoding возвращает percent-кодированную копию. exclude/include могут
// быть nil или пустыми (эквивалент пустого QByteArray).
// Qt: QByteArray::toPercentEncoding(const QByteArray &, const QByteArray &, char) const.
func (q *QByteArray) ToPercentEncoding(exclude, include []byte, percent byte) *QByteArray {
	exc := NewQByteArrayFromBytes(exclude)
	inc := NewQByteArrayFromBytes(include)
	defer exc.Close()
	defer inc.Close()

	var ep, ip *C.QByteArray
	if bp, ok := exc.tryHandle(); ok {
		ep = bp
	}
	if bp, ok := inc.tryHandle(); ok {
		ip = bp
	}

	r := wrapQByteArray(unsafe.Pointer(
		C.qbytearray_to_percent(q.handle(), ep, ip, C.char(percent))))
	runtime.KeepAlive(q)
	return r
}

// ToPercentEncodingBytes — Go-вариант ToPercentEncoding: результат в []byte.
// Qt: QByteArray::toPercentEncoding(const QByteArray &, const QByteArray &, char) const.
func (q *QByteArray) ToPercentEncodingBytes(exclude, include []byte, percent byte) []byte {
	r := q.ToPercentEncoding(exclude, include, percent)
	defer r.Close()
	return r.Bytes()
}

// PercentDecoded декодирует percent-кодирование в q (эквивалент
// FromPercentEncoding, применённого к q).
// Qt: QByteArray::percentDecoded(char) const &.
func (q *QByteArray) PercentDecoded(percent byte) *QByteArray {
	r := wrapQByteArray(unsafe.Pointer(C.qbytearray_percent_decoded(q.handle(), C.char(percent))))
	runtime.KeepAlive(q)
	return r
}

// PercentDecodedBytes — Go-вариант PercentDecoded: результат в []byte.
// Qt: QByteArray::percentDecoded(char) const &, результат скопирован в []byte.
func (q *QByteArray) PercentDecodedBytes(percent byte) []byte {
	r := q.PercentDecoded(percent)
	defer r.Close()
	return r.Bytes()
}

// ---- compression ---------------------------------------------------------

// Compress — Qt: qCompress(const QByteArray &, int).
func (q *QByteArray) Compress(level int) *QByteArray {
	r := wrapQByteArray(unsafe.Pointer(C.qbytearray_compress(q.handle(), C.int(level))))
	runtime.KeepAlive(q)
	return r
}

// CompressBytes — Go-вариант Compress: сжатые данные копией в []byte.
// Qt: qCompress(const QByteArray &, int), результат скопирован в []byte.
func (q *QByteArray) CompressBytes(level int) []byte {
	r := q.Compress(level)
	defer r.Close()
	return r.Bytes()
}

// Uncompress — Qt: qUncompress(const QByteArray &).
func (q *QByteArray) Uncompress() *QByteArray {
	r := wrapQByteArray(unsafe.Pointer(C.qbytearray_uncompress(q.handle())))
	runtime.KeepAlive(q)
	return r
}

// UncompressBytes — Go-вариант Uncompress: распакованные данные копией в []byte.
// Qt: qUncompress(const QByteArray &), результат скопирован в []byte.
func (q *QByteArray) UncompressBytes() []byte {
	r := q.Uncompress()
	defer r.Close()
	return r.Bytes()
}

// ---- split ---------------------------------------------------------------

// Split делит массив по байту sep и возвращает срез самостоятельно
// владеемых handle'ов. Пустые части сохраняются.
// Qt: QList<QByteArray> QByteArray::split(char) const.
func (q *QByteArray) Split(sep byte) []*QByteArray {
	var n C.ptrdiff_t
	arr := C.qbytearray_split(q.handle(), C.char(sep), &n)
	runtime.KeepAlive(q)
	if arr == nil || n == 0 {
		return nil
	}
	defer C.qbytearray_split_free_array(arr)

	cSlice := unsafe.Slice(arr, int(n))
	out := make([]*QByteArray, int(n))
	for i, p := range cSlice {
		out[i] = wrapQByteArray(unsafe.Pointer(p))
	}
	return out
}

// SplitBytes — Go-вариант Split: делит массив по sep и возвращает [][]byte.
// Qt: QList<QByteArray> QByteArray::split(char) const, элементы скопированы в [][]byte.
func (q *QByteArray) SplitBytes(sep byte) [][]byte {
	parts := q.Split(sep)
	if parts == nil {
		return nil
	}
	out := make([][]byte, len(parts))
	for i, p := range parts {
		out[i] = p.Bytes()
		p.Close()
	}
	return out
}
