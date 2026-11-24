package qdatetime

/*
	QCollator — биндинг Qt 6.11 QtCore/QCollator (locale-aware сравнение строк).

	QCollator — heap-объект с pimpl; в Go он представлен handle-обёрткой с
	владением через runtime.AddCleanup. Close() опционален: если его не вызывать,
	C++ объект удалится при сборке мусора. Close() идемпотентен и безопасен для nil.

	Реализовано:

	- конструкторы QCollator(), QCollator(const QLocale &), QCollator(const QCollator &);
	- setLocale/locale, caseSensitivity/setCaseSensitivity (константы берутся из
		qdata: qdata.CaseSensitivity, qdata.CaseInsensitive, qdata.CaseSensitive);
	- numericMode/setNumericMode, ignorePunctuation/setIgnorePunctuation;
	- compare(const QString &, const QString &) и compare(QStringView, QStringView)
		плюс Go-варианты со string;
	- operator() в виде Less/LessString/LessView/LessViewString;
	- sortKey(const QString &) плюс string-вариант;
	- статические defaultCompare/defaultSortKey (пакетные функции).

	Не реализовано (оставлено на будущее):

	- compare(const QChar *, len, const QChar *, len) — Go-пользователям доступны
		Compare/CompareView со string/[]uint16 через QStringView;
	- operator= (copy assignment) — используйте Clone()/Swap();
	- move-конструктор/присваивание — не нужны для GC-управляемых handle;
	- detach() — приватный метод, не экспонируется;
	- QDataStream/qHash/QDebug-интеграция.

	Перечисления: собственных enum у QCollator нет; Qt::CaseSensitivity уже
	представлен в qdata и переиспользуется.
*/

/*
#cgo pkg-config: Qt6Core
#cgo CXXFLAGS: -std=c++17 -fPIC -Wno-sfinae-incomplete
#cgo LDFLAGS: -lstdc++
#include <stdint.h>
#include <stdbool.h>
#include <stddef.h>
#include "qcollator_c.h"
*/
import "C"

import (
	"runtime"
	"sync/atomic"
	"unsafe"

	"github.com/vadimlxx/go-qt6core/qdata"
	"github.com/vadimlxx/go-qt6core/qtext"
)

// qcollatorState holds the C++ QCollator pointer in an atomic slot so Close and
// the runtime cleanup can race safely.
//
// Qt: no direct counterpart — internal ownership helper.
type qcollatorState struct {
	p atomic.Pointer[byte]
}

// cpp returns the stored QCollator pointer, or nil if the handle was claimed.
// Qt: no direct counterpart — internal helper.
func (s *qcollatorState) cpp() *C.go_qcollator_t {
	if bp := s.p.Load(); bp != nil {
		return (*C.go_qcollator_t)(unsafe.Pointer(bp))
	}
	return nil
}

// claim atomically clears the slot and returns the pointer stored there.
// Qt: no direct counterpart — internal helper.
func (s *qcollatorState) claim() *C.go_qcollator_t {
	if bp := s.p.Swap(nil); bp != nil {
		return (*C.go_qcollator_t)(unsafe.Pointer(bp))
	}
	return nil
}

// QCollator wraps a heap-allocated Qt6 QCollator.
//
// Qt: QCollator.
type QCollator struct {
	state   *qcollatorState
	cleanup runtime.Cleanup
}

// wrapQCollator takes ownership of a raw QCollator pointer.
// Qt: no direct counterpart — internal constructor helper.
func wrapQCollator(p *C.go_qcollator_t) *QCollator {
	if p == nil {
		return nil
	}

	st := &qcollatorState{}
	st.p.Store((*byte)(unsafe.Pointer(p)))
	c := &QCollator{state: st}
	c.cleanup = runtime.AddCleanup(c, func(st *qcollatorState) {
		if cpp := st.claim(); cpp != nil {
			C.qcollator_delete(cpp)
		}
	}, st)

	return c
}

// handle returns the underlying C++ pointer, panicking if the handle is nil or
// closed.
// Qt: no direct counterpart — internal handle helper.
func (c *QCollator) handle() *C.go_qcollator_t {
	if p, ok := c.tryHandle(); ok {
		return p
	}
	panic("qcollator: use after Close (or nil receiver)")
}

// tryHandle is the non-panicking variant of handle.
// Qt: no direct counterpart — internal handle helper.
func (c *QCollator) tryHandle() (*C.go_qcollator_t, bool) {
	if c == nil || c.state == nil {
		return nil, false
	}
	if p := c.state.cpp(); p != nil {
		return p, true
	}
	return nil, false
}

// Close releases the C++ object immediately and cancels the cleanup. Close is
// optional, idempotent and nil-safe.
// Qt: QCollator::~QCollator().
func (c *QCollator) Close() {
	if c == nil || c.state == nil {
		return
	}
	c.cleanup.Stop()
	if p := c.state.claim(); p != nil {
		C.qcollator_delete(p)
	}
}

// ---- construction -------------------------------------------------------

// NewQCollator is the generic constructor. It dispatches to the matching
// constructor based on the Go argument types:
//
//	NewQCollator()            -> QCollator()
//	NewQCollator(*QLocale)    -> QCollator(const QLocale &)
//	NewQCollator(*QCollator)  -> QCollator(const QCollator &)
//
// Qt: QCollator(), QCollator(const QLocale &), QCollator(const QCollator &).
func NewQCollator(arg ...any) *QCollator {
	if len(arg) == 0 || arg[0] == nil {
		return wrapQCollator(C.qcollator_new_default())
	}

	switch v := arg[0].(type) {
	case *QLocale:
		return NewQCollatorFromLocale(v)
	case *QCollator:
		return NewQCollatorCopy(v)
	}

	panic("NewQCollator: unsupported argument combination")
}

// NewQCollatorFromLocale creates a collator for the given locale.
// Qt: QCollator::QCollator(const QLocale &locale).
func NewQCollatorFromLocale(locale *QLocale) *QCollator {
	if locale == nil {
		panic("qcollator: NewQCollatorFromLocale with nil QLocale")
	}
	c := wrapQCollator(C.qcollator_new_from_locale(locale.handle()))
	runtime.KeepAlive(locale)
	return c
}

// NewQCollatorCopy creates an independent copy of other.
// Qt: QCollator::QCollator(const QCollator &other).
func NewQCollatorCopy(other *QCollator) *QCollator {
	c := wrapQCollator(C.qcollator_new_copy(other.handle()))
	runtime.KeepAlive(other)
	return c
}

// Clone returns an independent copy of the collator.
// Qt: QCollator::QCollator(const QCollator &other).
func (c *QCollator) Clone() *QCollator {
	return NewQCollatorCopy(c)
}

// Swap exchanges the contents of c and other.
// Qt: QCollator::swap(QCollator &other).
func (c *QCollator) Swap(other *QCollator) {
	C.qcollator_swap(c.handle(), other.handle())
	runtime.KeepAlive(c)
	runtime.KeepAlive(other)
}

// ---- locale and options -------------------------------------------------

// SetLocale sets the collation locale.
// Qt: QCollator::setLocale(const QLocale &locale).
func (c *QCollator) SetLocale(locale *QLocale) {
	if locale == nil {
		panic("qcollator: SetLocale with nil QLocale")
	}
	C.qcollator_set_locale(c.handle(), locale.handle())
	runtime.KeepAlive(c)
	runtime.KeepAlive(locale)
}

// Locale returns the collation locale as a new QLocale handle.
// Qt: QCollator::locale() const.
func (c *QCollator) Locale() *QLocale {
	p := C.qcollator_locale(c.handle())
	runtime.KeepAlive(c)
	return wrapQLocale(p)
}

// CaseSensitivity returns the current case sensitivity.
// Qt: QCollator::caseSensitivity() const.
func (c *QCollator) CaseSensitivity() qdata.CaseSensitivity {
	return qdata.CaseSensitivity(C.qcollator_case_sensitivity(c.handle()))
}

// SetCaseSensitivity sets the case sensitivity.
// Qt: QCollator::setCaseSensitivity(Qt::CaseSensitivity cs).
func (c *QCollator) SetCaseSensitivity(cs qdata.CaseSensitivity) {
	C.qcollator_set_case_sensitivity(c.handle(), C.int(cs))
	runtime.KeepAlive(c)
}

// NumericMode reports whether numeric collation is enabled.
// Qt: QCollator::numericMode() const.
func (c *QCollator) NumericMode() bool {
	return bool(C.qcollator_numeric_mode(c.handle()))
}

// SetNumericMode enables or disables numeric collation.
// Qt: QCollator::setNumericMode(bool on).
func (c *QCollator) SetNumericMode(on bool) {
	C.qcollator_set_numeric_mode(c.handle(), C.bool(on))
	runtime.KeepAlive(c)
}

// IgnorePunctuation reports whether punctuation is ignored.
// Qt: QCollator::ignorePunctuation() const.
func (c *QCollator) IgnorePunctuation() bool {
	return bool(C.qcollator_ignore_punctuation(c.handle()))
}

// SetIgnorePunctuation enables or disables punctuation ignoring.
// Qt: QCollator::setIgnorePunctuation(bool on).
func (c *QCollator) SetIgnorePunctuation(on bool) {
	C.qcollator_set_ignore_punctuation(c.handle(), C.bool(on))
	runtime.KeepAlive(c)
}

// ---- comparison ---------------------------------------------------------

// qstringArgs returns the UTF-16 (data, length) pair of s; nil is treated as an
// empty string.
// Qt: no direct counterpart — internal helper.
func qstringArgs(s *qtext.QString) (*C.uint16_t, C.ptrdiff_t) {
	if s == nil || s.Size() == 0 {
		return nil, 0
	}
	return (*C.uint16_t)(s.Utf16Unsafe()), C.ptrdiff_t(s.Size())
}

// qstringViewArgs returns the UTF-16 (data, length) pair of v.
// Qt: no direct counterpart — internal helper.
func qstringViewArgs(v qtext.QStringView) (*C.uint16_t, C.ptrdiff_t) {
	if v.Size() == 0 {
		return nil, 0
	}
	return (*C.uint16_t)(v.Utf16Unsafe()), C.ptrdiff_t(v.Size())
}

// Compare compares s1 and s2 using the collator's locale and options.
// Returns <0, 0 or >0. Qt: QCollator::compare(const QString &s1, const QString &s2) const.
func (c *QCollator) Compare(s1, s2 *qtext.QString) int {
	p1, n1 := qstringArgs(s1)
	p2, n2 := qstringArgs(s2)
	r := int(C.qcollator_compare(c.handle(), p1, n1, p2, n2))
	runtime.KeepAlive(c)
	runtime.KeepAlive(s1)
	runtime.KeepAlive(s2)
	return r
}

// CompareString compares s1 and s2 using the collator's locale and options.
// Qt: QCollator::compare(const QString &s1, const QString &s2) const — Go string form.
func (c *QCollator) CompareString(s1, s2 string) int {
	q1 := qtext.NewQStringFromString(s1)
	defer q1.Close()
	q2 := qtext.NewQStringFromString(s2)
	defer q2.Close()
	return c.Compare(q1, q2)
}

// CompareView compares the views v1 and v2 using the collator's locale and
// options. Qt: QCollator::compare(QStringView s1, QStringView s2) const.
func (c *QCollator) CompareView(v1, v2 qtext.QStringView) int {
	p1, n1 := qstringViewArgs(v1)
	p2, n2 := qstringViewArgs(v2)
	r := int(C.qcollator_compare(c.handle(), p1, n1, p2, n2))
	runtime.KeepAlive(c)
	runtime.KeepAlive(v1)
	runtime.KeepAlive(v2)
	return r
}

// CompareViewString compares s1 and s2 as QStringViews.
// Qt: QCollator::compare(QStringView s1, QStringView s2) const — Go string form.
func (c *QCollator) CompareViewString(s1, s2 string) int {
	q1 := qtext.NewQStringFromString(s1)
	defer q1.Close()
	q2 := qtext.NewQStringFromString(s2)
	defer q2.Close()
	return c.CompareView(q1.ToQStringView(), q2.ToQStringView())
}

// Less reports whether s1 sorts before s2.
// Qt: QCollator::operator()(const QString &s1, const QString &s2) const.
func (c *QCollator) Less(s1, s2 *qtext.QString) bool { return c.Compare(s1, s2) < 0 }

// LessString reports whether s1 sorts before s2.
// Qt: QCollator::operator()(const QString &s1, const QString &s2) const — Go string form.
func (c *QCollator) LessString(s1, s2 string) bool { return c.CompareString(s1, s2) < 0 }

// LessView reports whether v1 sorts before v2.
// Qt: QCollator::operator()(QStringView s1, QStringView s2) const.
func (c *QCollator) LessView(v1, v2 qtext.QStringView) bool { return c.CompareView(v1, v2) < 0 }

// LessViewString reports whether s1 sorts before s2 when compared as views.
// Qt: QCollator::operator()(QStringView s1, QStringView s2) const — Go string form.
func (c *QCollator) LessViewString(s1, s2 string) bool { return c.CompareViewString(s1, s2) < 0 }

// ---- sort keys ----------------------------------------------------------

// SortKey returns the collation sort key for s.
// Qt: QCollator::sortKey(const QString &string) const.
func (c *QCollator) SortKey(s *qtext.QString) *QCollatorSortKey {
	p, n := qstringArgs(s)
	k := wrapQCollatorSortKey(C.qcollator_sort_key(c.handle(), p, n))
	runtime.KeepAlive(c)
	runtime.KeepAlive(s)
	return k
}

// SortKeyString returns the collation sort key for s.
// Qt: QCollator::sortKey(const QString &string) const — Go string form.
func (c *QCollator) SortKeyString(s string) *QCollatorSortKey {
	q := qtext.NewQStringFromString(s)
	defer q.Close()
	return c.SortKey(q)
}

// ---- static helpers -----------------------------------------------------

// DefaultCompare compares v1 and v2 using the default (C locale-independent)
// ordering. Qt: static QCollator::defaultCompare(QStringView s1, QStringView s2).
func DefaultCompare(v1, v2 qtext.QStringView) int {
	p1, n1 := qstringViewArgs(v1)
	p2, n2 := qstringViewArgs(v2)
	r := int(C.qcollator_default_compare(p1, n1, p2, n2))
	runtime.KeepAlive(v1)
	runtime.KeepAlive(v2)
	return r
}

// DefaultCompareString compares s1 and s2 using the default ordering.
// Qt: static QCollator::defaultCompare(QStringView s1, QStringView s2) — Go string form.
func DefaultCompareString(s1, s2 string) int {
	q1 := qtext.NewQStringFromString(s1)
	defer q1.Close()
	q2 := qtext.NewQStringFromString(s2)
	defer q2.Close()
	return DefaultCompare(q1.ToQStringView(), q2.ToQStringView())
}

// DefaultSortKey returns the default sort key for v.
// Qt: static QCollator::defaultSortKey(QStringView key).
func DefaultSortKey(v qtext.QStringView) *QCollatorSortKey {
	p, n := qstringViewArgs(v)
	k := wrapQCollatorSortKey(C.qcollator_default_sort_key(p, n))
	runtime.KeepAlive(v)
	return k
}

// DefaultSortKeyString returns the default sort key for s.
// Qt: static QCollator::defaultSortKey(QStringView key) — Go string form.
func DefaultSortKeyString(s string) *QCollatorSortKey {
	q := qtext.NewQStringFromString(s)
	defer q.Close()
	return DefaultSortKey(q.ToQStringView())
}
