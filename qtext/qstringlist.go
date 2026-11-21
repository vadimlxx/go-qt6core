package qtext

/*
   QStringList — обёртка над C++ QList<QString> через непрозрачный указатель
   (C ABI). Go-структура содержит только этот указатель; владение обеспечивается
   финализатором и явным Close(), как в QString и QByteArray.

   В Go-API, как и в QString,
   case sensitivity передаётся как bool + хелпер csInt (существующий в пакете).

   Методы ToSlice / At / TakeAt / First / Last / TakeFirst / TakeLast
   возвращают новые *QString (копии). Владение — за вызывающим.

   Метод FromQStrings принимает []*QString; каждый элемент копируется
   в список. Метод FromGoStrings — удобный вариант для []string.

   JoinString принимает *QString-разделитель; JoinChar — QChar-разделитель;
   JoinGoString — Go string (UTF-8).

   Перегрузки filter / replaceInStrings / indexOf / lastIndexOf,
   принимающие QStringView / QLatin1StringView, заменены единым вариантом
   с *QString (или Go string). Аналогично для QRegularExpression.

 Что сознательно не вошло (пока):

   - QListIterator<QString> / QMutableListIterator<QString> —
     Java-style итераторы; используйте ToSlice().
   - STL-style итераторы (begin/end/rbegin/rend) — неприменимы в Go.
   - filter(const QStringMatcher&) / filter(const QLatin1StringMatcher&) —
     требует QRegularExpression / QStringMatcher, добавим по мере надобности.
   - Перегрузки с QStringView / QLatin1StringView — сейчас принимаются
     *QString или Go string; при необходимости добавляются с суффиксом.
   - operator+ / operator<< — заменяются Append / AppendList.
   - QDataStream operator<< / operator>> — не идиоматичны в Go.
   - QDebug operator<< — не идиоматично в Go.
   - QHashSpecialization / qHash для QStringList — не критично.
*/

/*
#cgo pkg-config: Qt6Core
#cgo CXXFLAGS: -std=c++17 -fPIC -Wno-sfinae-incomplete
#cgo LDFLAGS: -lstdc++
#include <stdint.h>
#include <stdbool.h>
#include <stddef.h>
#include "qstringlist_c.h"
*/
import "C"

import (
	"runtime"
	"unsafe"
)

// QStringList wraps a Qt6 QStringList (QList<QString>).
// Each *QStringList owns a heap-allocated C++ QStringList; the finalizer
// (and explicit Close) frees it.
type QStringList struct {
	p *C.QStringList
}

func wrapQStringList(p unsafe.Pointer) *QStringList {
	if p == nil {
		return nil
	}
	q := &QStringList{p: (*C.QStringList)(p)}
	runtime.SetFinalizer(q, finalizeQStringList)
	return q
}

func finalizeQStringList(q *QStringList) {
	if q.p != nil {
		C.qstringlist_delete(q.p)
		q.p = nil
	}
}

// Close releases the underlying C++ object immediately and disables the finalizer.
// Qt: ~QStringList() — explicit Close maps to the destructor call.
func (q *QStringList) Close() {
	if q.p != nil {
		C.qstringlist_delete(q.p)
		q.p = nil
		runtime.SetFinalizer(q, nil)
	}
}

// ---- construction --------------------------------------------------------

// NewQStringList returns an empty QStringList.
// Qt: QStringList() — default constructor.
// func NewQStringList() *QStringList {
// 	return wrapQStringList(unsafe.Pointer(C.qstringlist_new_empty()))
// }


// QStringListFromQString returns a QStringList with a single element.
// Qt: QStringList(const QString&) — explicit constructor.
func QStringListFromQString(s *QString) *QStringList {
	if s == nil || s.handle() == nil {
		return NewQStringList()
	}
	q := wrapQStringList(unsafe.Pointer(C.qstringlist_new_from_string(s.handle())))
	runtime.KeepAlive(s)
	return q
}

// QStringListFromGoString returns a QStringList with a single element from s (UTF-8).
// Qt: QStringList(const QString&) — with QString constructed via fromUtf8.
func QStringListFromGoString(s string) *QStringList {
	qs := NewQStringFromString(s)
	defer qs.Close()
	return QStringListFromQString(qs)
}

// QStringListFromQStrings returns a QStringList by copying each element of items.
// Qt: QStringList(const QList<QString>&) — copy constructor from the underlying list.
func QStringListFromQStrings(items []*QString) *QStringList {
	n := len(items)
	if n == 0 {
		return NewQStringList()
	}
	cArr := make([]*C.QString, n)
	for i, s := range items {
		if s != nil {
			p := s.handle()
			if p != nil { cArr[i] = p }
		}
	}
	q := wrapQStringList(unsafe.Pointer(C.qstringlist_from_array(
		(**C.QString)(unsafe.Pointer(&cArr[0])), C.ptrdiff_t(n))))
	runtime.KeepAlive(items)
	return q
}

// QStringListFromGoStrings returns a QStringList by decoding each Go string as UTF-8.
// Qt: QStringList(const QList<QString>&) — with each element from QString::fromUtf8.
func QStringListFromGoStrings(items []string) *QStringList {
	n := len(items)
	if n == 0 {
		return NewQStringList()
	}
	qs := make([]*QString, n)
	for i, s := range items {
		qs[i] = NewQStringFromString(s)
	}
	out := QStringListFromQStrings(qs)
	for _, s := range qs {
		s.Close()
	}
	return out
}

func NewQStringList(arg ...any) *QStringList {
	if len(arg) == 0 || arg[0] == nil {
		return wrapQStringList(unsafe.Pointer(C.qstringlist_new_empty()))
	}

	switch v := arg[0].(type) {
	case *QString:
		return QStringListFromQString(v)
	case string:
		return QStringListFromGoString(v)
	case []*QString:
		return QStringListFromQStrings(v)
	case []string:
		return QStringListFromGoStrings(v)
	case *QStringList:
		return v.Clone()
	}

	panic("NewQStringList invalid argument")
}

// Clone returns a deep copy of the list (each QString is copied).
// Qt: QStringList(const QStringList&) — copy constructor.
func (q *QStringList) Clone() *QStringList {
	return wrapQStringList(unsafe.Pointer(C.qstringlist_new_copy(q.p)))
}

// Swap exchanges the contents of q and other.
// Qt: QStringList::swap(QStringList&).
func (q *QStringList) Swap(other *QStringList) {
	C.qstringlist_swap(q.p, other.p)
	runtime.KeepAlive(other)
}

// ---- accessors -----------------------------------------------------------

// Size returns the number of elements.
// Qt: QList<QString>::size() const.
func (q *QStringList) Size() int { return int(C.qstringlist_size(q.p)) }

// Len is an alias for Size.
// Qt: QList<QString>::length() const (which forwards to size()).
func (q *QStringList) Len() int { return q.Size() }

// IsEmpty reports whether the list has zero elements.
// Qt: QList<QString>::isEmpty() const.
func (q *QStringList) IsEmpty() bool { return bool(C.qstringlist_is_empty(q.p)) }

// ---- modifiers -----------------------------------------------------------

// Clear removes all elements.
// Qt: QList<QString>::clear().
func (q *QStringList) Clear() { C.qstringlist_clear(q.p) }

// Reserve pre-allocates storage for n elements.
// Qt: QList<QString>::reserve(qsizetype).
func (q *QStringList) Reserve(n int) { C.qstringlist_reserve(q.p, C.ptrdiff_t(n)) }

// Squeeze releases unused capacity.
// Qt: QList<QString>::squeeze().
func (q *QStringList) Squeeze() { C.qstringlist_squeeze(q.p) }

// ---- element access -----------------------------------------------------

// At returns a copy of the i-th element. Panics if i is out of range.
// Qt: QList<QString>::at(qsizetype) const.
func (q *QStringList) At(i int) *QString {
	if i < 0 || i >= q.Size() {
		panic("qstringlist: At out of range")
	}
	r := wrapQString(unsafe.Pointer(C.qstringlist_at(q.p, C.ptrdiff_t(i))))
	runtime.KeepAlive(q)
	return r
}

// SetAt replaces the i-th element with s. Panics if i is out of range.
// Qt: QList<QString>::operator[](qsizetype) (non-const) — assignable reference.
func (q *QStringList) SetAt(i int, s *QString) {
	if i < 0 || i >= q.Size() {
		panic("qstringlist: SetAt out of range")
	}
	C.qstringlist_set_at(q.p, C.ptrdiff_t(i), s.handle())
	runtime.KeepAlive(s)
}

// Front returns a copy of the first element. Panics if empty.
// Qt: QList<QString>::first() const.
func (q *QStringList) Front() *QString {
	if q.IsEmpty() {
		panic("qstringlist: Front on empty QStringList")
	}
	r := wrapQString(unsafe.Pointer(C.qstringlist_first(q.p)))
	runtime.KeepAlive(q)
	return r
}

// Back returns a copy of the last element. Panics if empty.
// Qt: QList<QString>::last() const.
func (q *QStringList) Back() *QString {
	if q.IsEmpty() {
		panic("qstringlist: Back on empty QStringList")
	}
	r := wrapQString(unsafe.Pointer(C.qstringlist_last(q.p)))
	runtime.KeepAlive(q)
	return r
}

// ToSlice returns a new []*QString with copies of all elements.
// Qt: QList<QString>::toVector() (or STL begin/end iteration).
func (q *QStringList) ToSlice() []*QString {
	var count C.ptrdiff_t
	arr := C.qstringlist_to_array(q.p, &count)
	if arr == nil {
		return nil
	}
	defer C.qstringlist_free_string_array(arr)
	n := int(count)
	if n == 0 {
		return nil
	}
	ptrs := unsafe.Slice(arr, n)
	result := make([]*QString, n)
	for i := 0; i < n; i++ {
		result[i] = wrapQString(unsafe.Pointer(ptrs[i]))
	}
	return result
}

// ---- appending / inserting / removing -----------------------------------

// Append adds s to the end of the list.
// Qt: QStringList::operator<<(const QString&) / QList<QString>::append(const QString&).
func (q *QStringList) Append(s *QString) {
	C.qstringlist_append(q.p, s.handle())
	runtime.KeepAlive(s)
}

// AppendGoString adds a Go string (UTF-8) to the end of the list.
// Qt: QStringList::operator<<(const QString&) — with QString::fromUtf8.
func (q *QStringList) AppendGoString(s string) {
	qs := NewQStringFromString(s)
	defer qs.Close()
	q.Append(qs)
}

// AppendList adds all elements of other to the end of q.
// Qt: QStringList::operator<<(const QStringList&).
func (q *QStringList) AppendList(other *QStringList) {
	for _, item := range other.ToSlice() {
		q.Append(item)
		item.Close()
	}
}

// Prepend inserts s at the beginning of the list.
// Qt: QList<QString>::prepend(const QString&).
func (q *QStringList) Prepend(s *QString) {
	C.qstringlist_prepend(q.p, s.handle())
	runtime.KeepAlive(s)
}

// Insert inserts s at position i (0 <= i <= Size()). Panics if i is out of range.
// Qt: QList<QString>::insert(qsizetype, const QString&).
func (q *QStringList) Insert(i int, s *QString) {
	if i < 0 || i > q.Size() {
		panic("qstringlist: Insert out of range")
	}
	C.qstringlist_insert(q.p, C.ptrdiff_t(i), s.handle())
	runtime.KeepAlive(s)
}

// RemoveAt removes the element at position i. Panics if i is out of range.
// Qt: QList<QString>::removeAt(qsizetype).
func (q *QStringList) RemoveAt(i int) {
	if i < 0 || i >= q.Size() {
		panic("qstringlist: RemoveAt out of range")
	}
	C.qstringlist_remove_at(q.p, C.ptrdiff_t(i))
}

// RemoveFirst removes the first element (no-op if empty).
// Qt: QList<QString>::removeFirst().
func (q *QStringList) RemoveFirst() { C.qstringlist_remove_first(q.p) }

// RemoveLast removes the last element (no-op if empty).
// Qt: QList<QString>::removeLast().
func (q *QStringList) RemoveLast() { C.qstringlist_remove_last(q.p) }

// TakeAt removes the element at position i and returns it. Panics if i is out of range.
// Qt: QList<QString>::takeAt(qsizetype).
func (q *QStringList) TakeAt(i int) *QString {
	if i < 0 || i >= q.Size() {
		panic("qstringlist: TakeAt out of range")
	}
	r := wrapQString(unsafe.Pointer(C.qstringlist_take_at(q.p, C.ptrdiff_t(i))))
	runtime.KeepAlive(q)
	return r
}

// TakeFirst removes the first element and returns it. Panics if empty.
// Qt: QList<QString>::takeFirst().
func (q *QStringList) TakeFirst() *QString {
	if q.IsEmpty() {
		panic("qstringlist: TakeFirst on empty QStringList")
	}
	r := wrapQString(unsafe.Pointer(C.qstringlist_take_first(q.p)))
	runtime.KeepAlive(q)
	return r
}

// TakeLast removes the last element and returns it. Panics if empty.
// Qt: QList<QString>::takeLast().
func (q *QStringList) TakeLast() *QString {
	if q.IsEmpty() {
		panic("qstringlist: TakeLast on empty QStringList")
	}
	r := wrapQString(unsafe.Pointer(C.qstringlist_take_last(q.p)))
	runtime.KeepAlive(q)
	return r
}

// ---- QStringList-specific ------------------------------------------------

// Sort sorts the list in ascending order.
// Qt: QStringList::sort(Qt::CaseSensitivity) — with Qt::CaseSensitive / Qt::CaseInsensitive.
func (q *QStringList) Sort(caseSensitive bool) {
	C.qstringlist_sort(q.p, csInt(caseSensitive))
}

// RemoveDuplicates removes duplicate entries and returns the number removed.
// The list must be sorted first for all duplicates to be adjacent.
// Qt: QStringList::removeDuplicates().
func (q *QStringList) RemoveDuplicates() int {
	return int(C.qstringlist_remove_duplicates(q.p))
}

// JoinString joins all elements with sep.
// Qt: QStringList::join(const QString&).
func (q *QStringList) JoinString(sep *QString) *QString {
	r := wrapQString(unsafe.Pointer(C.qstringlist_join_qstring(q.p, sep.handle())))
	runtime.KeepAlive(sep)
	return r
}

// JoinGoString joins all elements with sep (UTF-8).
// Qt: QStringList::join(const QString&) — with QString::fromUtf8.
func (q *QStringList) JoinGoString(sep string) *QString {
	qs := NewQStringFromString(sep)
	defer qs.Close()
	return q.JoinString(qs)
}

// JoinChar joins all elements with a single QChar.
// Qt: QStringList::join(QChar).
func (q *QStringList) JoinChar(sep QChar) *QString {
	return wrapQString(unsafe.Pointer(C.qstringlist_join_qchar(q.p, C.uint16_t(sep.Unicode()))))
}

// FilterString returns a new list with elements containing needle.
// Qt: QStringList::filter(const QString&, Qt::CaseSensitivity).
func (q *QStringList) FilterString(needle *QString, caseSensitive bool) *QStringList {
	r := wrapQStringList(unsafe.Pointer(C.qstringlist_filter_qstring(q.p, needle.handle(), csInt(caseSensitive))))
	runtime.KeepAlive(needle)
	return r
}

// FilterGoString returns a new list with elements containing needle (UTF-8).
// Qt: QStringList::filter(const QString&, Qt::CaseSensitivity) — with QString::fromUtf8.
func (q *QStringList) FilterGoString(needle string, caseSensitive bool) *QStringList {
	qs := NewQStringFromString(needle)
	defer qs.Close()
	return q.FilterString(qs, caseSensitive)
}

// FilterQRegularExpression returns a new list with elements matching re.
// Qt: QStringList::filter(const QRegularExpression&).
func (q *QStringList) FilterQRegularExpression(re *QRegularExpression) *QStringList {
	p := re.handle()
	r := wrapQStringList(unsafe.Pointer(C.qstringlist_filter_qregexp(q.p, p)))
	runtime.KeepAlive(re)
	return r
}

// ReplaceInStringsString replaces every occurrence of before with after in each element.
// Qt: QStringList::replaceInStrings(const QString&, const QString&, Qt::CaseSensitivity).
func (q *QStringList) ReplaceInStringsString(before, after *QString, caseSensitive bool) {
	C.qstringlist_replace_in_strings_qstring(q.p, before.handle(), after.handle(), csInt(caseSensitive))
	runtime.KeepAlive(before)
	runtime.KeepAlive(after)
}

// ReplaceInStringsGoString replaces every occurrence of before with after (UTF-8).
// Qt: QStringList::replaceInStrings(const QString&, const QString&, Qt::CaseSensitivity).
func (q *QStringList) ReplaceInStringsGoString(before, after string, caseSensitive bool) {
	qsBefore := NewQStringFromString(before)
	qsAfter := NewQStringFromString(after)
	defer qsBefore.Close()
	defer qsAfter.Close()
	q.ReplaceInStringsString(qsBefore, qsAfter, caseSensitive)
}

// ReplaceInStringsQRegularExpression replaces every match of re with after in each element.
// Qt: QStringList::replaceInStrings(const QRegularExpression&, const QString&).
func (q *QStringList) ReplaceInStringsQRegularExpression(re *QRegularExpression, after *QString) {
	p := re.handle()
	C.qstringlist_replace_in_strings_qregexp(q.p, p, after.handle())
	runtime.KeepAlive(re)
	runtime.KeepAlive(after)
}

// ContainsQString reports whether the list contains s (comparing each element).
// Qt: QStringList::contains(const QString&, Qt::CaseSensitivity).
func (q *QStringList) ContainsQString(s *QString, caseSensitive bool) bool {
	r := bool(C.qstringlist_contains_qstring(q.p, s.handle(), csInt(caseSensitive)))
	runtime.KeepAlive(s)
	return r
}

// ContainsGoString reports whether the list contains s (UTF-8).
// Qt: QStringList::contains(const QString&, Qt::CaseSensitivity).
func (q *QStringList) ContainsGoString(s string, caseSensitive bool) bool {
	qs := NewQStringFromString(s)
	defer qs.Close()
	return q.ContainsQString(qs, caseSensitive)
}

// IndexOfQString returns the index of the first element equal to s at or after from (or -1).
// Qt: QStringList::indexOf(const QString&, qsizetype, Qt::CaseSensitivity).
func (q *QStringList) IndexOfQString(s *QString, from int, caseSensitive bool) int {
	r := C.qstringlist_index_of_qstring(q.p, s.handle(), C.ptrdiff_t(from), csInt(caseSensitive))
	runtime.KeepAlive(s)
	return int(r)
}

// IndexOfGoString returns the index of the first element equal to s at or after from (or -1).
// Qt: QStringList::indexOf(const QString&, qsizetype, Qt::CaseSensitivity).
func (q *QStringList) IndexOfGoString(s string, from int, caseSensitive bool) int {
	qs := NewQStringFromString(s)
	defer qs.Close()
	return q.IndexOfQString(qs, from, caseSensitive)
}

// IndexOfQRegularExpression returns the index of the first element matching re at or after from (or -1).
// Qt: QStringList::indexOf(const QRegularExpression&, qsizetype).
func (q *QStringList) IndexOfQRegularExpression(re *QRegularExpression, from int) int {
	r := C.qstringlist_index_of_qregexp(q.p, re.handle(), C.ptrdiff_t(from))
	runtime.KeepAlive(re)
	return int(r)
}

// LastIndexOfQString returns the index of the last element equal to s at or before from (or -1).
// Qt: QStringList::lastIndexOf(const QString&, qsizetype, Qt::CaseSensitivity).
func (q *QStringList) LastIndexOfQString(s *QString, from int, caseSensitive bool) int {
	r := C.qstringlist_last_index_of_qstring(q.p, s.handle(), C.ptrdiff_t(from), csInt(caseSensitive))
	runtime.KeepAlive(s)
	return int(r)
}

// LastIndexOfGoString returns the index of the last element equal to s at or before from (or -1).
// Qt: QStringList::lastIndexOf(const QString&, qsizetype, Qt::CaseSensitivity).
func (q *QStringList) LastIndexOfGoString(s string, from int, caseSensitive bool) int {
	qs := NewQStringFromString(s)
	defer qs.Close()
	return q.LastIndexOfQString(qs, from, caseSensitive)
}

// LastIndexOfQRegularExpression returns the index of the last element matching re at or before from (or -1).
// Qt: QStringList::lastIndexOf(const QRegularExpression&, qsizetype).
func (q *QStringList) LastIndexOfQRegularExpression(re *QRegularExpression, from int) int {
	r := C.qstringlist_last_index_of_qregexp(q.p, re.handle(), C.ptrdiff_t(from))
	runtime.KeepAlive(re)
	return int(r)
}
