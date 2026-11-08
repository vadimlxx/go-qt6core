package qstringmatcher

// Package qstringmatcher: QStringMatcher — эффективный поиск подстроки с предвычислением.

/*
 ВАЖНО (Qt 6.11.2):
   - QStringMatcher — НЕ QObject, а обычный C++ value-класс с приватной реализацией
     (d_ptr), публичным копирующим конструктором и публичным деструктором. Копирование
     дешёвое (QString — COW, QStringView указывает в собственный буфер q_pattern).
   - Qt::CaseSensitivity передаётся как `caseSensitive bool` (true = CaseSensitive),
     в стиле остальных пакетов (qstring, qstringview, qbytearray).
   - PatternView() возвращает zero-copy QStringView в внутренний буфер matcher'а;
     view держит *QStringMatcher живым. Не вызывайте SetPattern или Close на matcher'е,
     пока view используется — буфер может быть пересоздан или освобождён. Для
     независимой копии используйте Pattern().

 не реализовано:
   - QStringMatcher(const QChar*, len, cs) — в Qt это inline-форвардер к
     QStringMatcher(QStringView(uc, len), cs); покрыт FromUtf16.
   - updateSkipTable() — private helper Qt, не часть публичного API, не экспонируется.
   - Qt::CaseSensitivity как отдельный Go-тип — используется `bool`, как в соседних
     пакетах. Для явной семантики добавьте константы в вызывающем коде, если нужно.
*/

/*
#cgo pkg-config: Qt6Core
#cgo CXXFLAGS: -std=c++17 -fPIC -Wno-sfinae-incomplete
#cgo LDFLAGS: -lstdc++
#include <stdint.h>
#include <stddef.h>
#include "qstringmatcher_c.h"
*/
import "C"

import (
	"runtime"
	"unsafe"

	"github.com/vadimlxx/go-qt6core/qstring"
	"github.com/vadimlxx/go-qt6core/qstringview"
)

// QStringMatcher mirrors Qt 6.11.2 QStringMatcher.
//
// It is a value-semantics class: copy is cheap (Clone / Assign). Use Close() to free the
// underlying C++ object immediately; otherwise the finalizer releases it.
type QStringMatcher struct {
	p *C.go_qsm_t
}

func wrap(p *C.go_qsm_t) *QStringMatcher {
	if p == nil {
		return nil
	}
	m := &QStringMatcher{p: p}
	runtime.SetFinalizer(m, finalize)
	return m
}

func finalize(m *QStringMatcher) {
	if m.p != nil {
		C.qstringmatcher_delete(m.p)
		m.p = nil
	}
}

func csInt(caseSensitive bool) C.int {
	if caseSensitive {
		return 1 // Qt::CaseSensitive
	}
	return 0 // Qt::CaseInsensitive
}

// ---- constructors --------------------------------------------------------

// New creates an empty QStringMatcher (empty pattern, Qt::CaseSensitive).
// Mirrors QStringMatcher().
func New() *QStringMatcher {
	return wrap(C.qstringmatcher_new_default())
}

// FromQString creates a QStringMatcher from a QString pattern.
// Mirrors QStringMatcher(const QString &pattern, Qt::CaseSensitivity cs = CaseSensitive).
// Passing a nil *qstring.QString is equivalent to an empty QString pattern.
func FromQString(pattern *qstring.QString, caseSensitive bool) *QStringMatcher {
	var p unsafe.Pointer
	if pattern != nil {
		p = pattern.UnsafePtr()
	}
	m := wrap(C.qstringmatcher_new_from_qstring(p, csInt(caseSensitive)))
	runtime.KeepAlive(pattern)
	return m
}

// FromQStringView creates a QStringMatcher from a QStringView pattern.
// Mirrors QStringMatcher(QStringView pattern, Qt::CaseSensitivity cs = CaseSensitive).
// The data is copied (QStringMatcher stores its own QString q_pattern), so the view does
// not need to outlive the matcher.
func FromQStringView(pattern qstringview.QStringView, caseSensitive bool) *QStringMatcher {
	n := pattern.Size()
	var p *C.uint16_t
	if n > 0 {
		p = (*C.uint16_t)(pattern.Utf16Unsafe())
	}
	m := wrap(C.qstringmatcher_new_from_utf16(p, C.ptrdiff_t(n), csInt(caseSensitive)))
	runtime.KeepAlive(pattern)
	return m
}

// FromUtf16 creates a QStringMatcher from raw UTF-16 code units (QChar*).
// Mirrors QStringMatcher(const QChar *uc, qsizetype len, Qt::CaseSensitivity cs).
// The data is copied into an internal QString.
func FromUtf16(units []uint16, caseSensitive bool) *QStringMatcher {
	var p *C.uint16_t
	if len(units) > 0 {
		p = (*C.uint16_t)(unsafe.Pointer(&units[0]))
	}
	m := wrap(C.qstringmatcher_new_from_utf16(p, C.ptrdiff_t(len(units)), csInt(caseSensitive)))
	runtime.KeepAlive(units)
	return m
}

// ---- lifecycle -----------------------------------------------------------

// Close frees the underlying C++ object immediately and disables the finalizer.
// After Close, all methods on m are UB.
func (m *QStringMatcher) Close() {
	if m.p != nil {
		C.qstringmatcher_delete(m.p)
		m.p = nil
		runtime.SetFinalizer(m, nil)
	}
}

// Clone returns a copy of m (via the C++ copy constructor).
func (m *QStringMatcher) Clone() *QStringMatcher {
	r := wrap(C.qstringmatcher_new_copy(m.p))
	runtime.KeepAlive(m)
	return r
}

// Assign copies the contents of other into m (via C++ operator=).
// After Assign, m's pattern and case sensitivity match other's.
func (m *QStringMatcher) Assign(other *QStringMatcher) {
	if other == nil || other.p == nil {
		return
	}
	C.qstringmatcher_assign(m.p, other.p)
	runtime.KeepAlive(m)
	runtime.KeepAlive(other)
}

// ---- pattern / sensitivity ----------------------------------------------

// SetPattern replaces m's pattern with q (copies q via QString COW).
// Passing nil is equivalent to an empty QString pattern.
// Invalidates any QStringView previously returned by PatternView.
func (m *QStringMatcher) SetPattern(q *qstring.QString) {
	var p unsafe.Pointer
	if q != nil {
		p = q.UnsafePtr()
	}
	C.qstringmatcher_set_pattern(m.p, p)
	runtime.KeepAlive(m)
	runtime.KeepAlive(q)
}

// SetCaseSensitivity sets case sensitivity (true = Qt::CaseSensitive, false =
// Qt::CaseInsensitive). Mirrors setCaseSensitivity(Qt::CaseSensitivity).
func (m *QStringMatcher) SetCaseSensitivity(caseSensitive bool) {
	C.qstringmatcher_set_case_sensitivity(m.p, csInt(caseSensitive))
	runtime.KeepAlive(m)
}

// CaseSensitivity returns the current case sensitivity.
// Mirrors caseSensitivity().
func (m *QStringMatcher) CaseSensitivity() bool {
	r := C.qstringmatcher_case_sensitivity(m.p)
	runtime.KeepAlive(m)
	return r != 0
}

// Pattern returns a fresh owned copy of the matcher's pattern as *qstring.QString.
// Mirrors QString::pattern().
func (m *QStringMatcher) Pattern() *qstring.QString {
	var out unsafe.Pointer
	C.qstringmatcher_pattern(m.p, &out)
	runtime.KeepAlive(m)
	return qstring.AdoptUnsafePtr(out)
}

// PatternView returns a zero-copy view into the matcher's internal pattern buffer.
// The view keeps m alive via runtime.SetFinalizer-friendly keepAlive.
//
// IMPORTANT: the view is invalidated by SetPattern, Assign, Close, or any other operation
// that modifies the matcher's q_pattern. If you need a stable snapshot, use Pattern().
// Mirrors patternView().
func (m *QStringMatcher) PatternView() qstringview.QStringView {
	var data *C.uint16_t
	var size C.ptrdiff_t
	C.qstringmatcher_pattern_view(m.p, &data, &size)
	runtime.KeepAlive(m)
	if data == nil || size == 0 {
		return qstringview.Empty()
	}
	return qstringview.NewUnsafeWithOwner(unsafe.Pointer(data), int(size), m)
}

// ---- search --------------------------------------------------------------

// IndexIn searches for the pattern in str starting at position from.
// Returns the position of the first match, or -1 if not found.
// Mirrors indexIn(const QString &str, qsizetype from = 0).
func (m *QStringMatcher) IndexIn(str *qstring.QString, from int) int {
	if str == nil {
		return -1
	}
	r := C.qstringmatcher_index_in_utf16(m.p,
		(*C.uint16_t)(str.Utf16Unsafe()),
		C.ptrdiff_t(str.Size()),
		C.ptrdiff_t(from))
	runtime.KeepAlive(m)
	runtime.KeepAlive(str)
	return int(r)
}

// IndexInView searches for the pattern in str starting at position from.
// Returns the position of the first match, or -1 if not found.
// Mirrors indexIn(QStringView str, qsizetype from = 0).
func (m *QStringMatcher) IndexInView(str qstringview.QStringView, from int) int {
	n := str.Size()
	var p *C.uint16_t
	if n > 0 {
		p = (*C.uint16_t)(str.Utf16Unsafe())
	}
	r := C.qstringmatcher_index_in_utf16(m.p, p, C.ptrdiff_t(n), C.ptrdiff_t(from))
	runtime.KeepAlive(m)
	runtime.KeepAlive(str)
	return int(r)
}

// IndexInUtf16 searches for the pattern in raw UTF-16 code units starting at position from.
// Returns the position of the first match, or -1 if not found.
// Mirrors indexIn(const QChar *str, qsizetype length, qsizetype from = 0).
func (m *QStringMatcher) IndexInUtf16(units []uint16, from int) int {
	var p *C.uint16_t
	if len(units) > 0 {
		p = (*C.uint16_t)(unsafe.Pointer(&units[0]))
	}
	r := C.qstringmatcher_index_in_utf16(m.p, p, C.ptrdiff_t(len(units)), C.ptrdiff_t(from))
	runtime.KeepAlive(m)
	runtime.KeepAlive(units)
	return int(r)
}
