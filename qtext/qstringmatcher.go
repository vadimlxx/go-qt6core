package qtext

/*
QStringMatcher — handle над C++ QStringMatcher; владение через
runtime.AddCleanup + Close().

# Семантика

QStringMatcher — быстрый поиск подстроки по алгоритму
Boyer-Moore-Horspool со skip-таблицей. Эффективен, когда один и тот же
шаблон ищется многократно в разных строках.

# Владение

C++ объект освобождается одним из двух способов:
  - Close() — немедленно, детерминированно;
  - runtime.AddCleanup — когда *QStringMatcher становится недостижим.

Close() опционален, идемпотентен и безопасен на nil. Если Close() и
cleanup гоняются, владение захватывает ровно один участник (атомарный
Swap(nil)), двойного free нет.

# Use after Close

После Close() любой метод, кроме Close(), паникует
("qstringmatcher: use after Close (or nil receiver)").

# Case sensitivity

Параметр cs имеет тип qdata.CaseSensitivity:
  - qdata.CaseInsensitive — поиск без учёта регистра;
  - qdata.CaseSensitive   — поиск с учётом регистра (по умолчанию в Qt).

Этот класс использует qdata-константы напрямую (в отличие от QString,
где в API используется bool caseSensitive — так сложилось исторически).

# PatternView lifetime

PatternView() возвращает невладеющий QStringView на внутреннее хранилище
шаблона. View инвалидируется при любом последующем вызове
SetPattern / SetPatternString / SetPatternView / Assign, а также после
Close. Если нужна стабильная копия шаблона — используйте Pattern() или
PatternString().

# Что сознательно не вошло

  - updateSkipTable() — protected в Qt, вызывается из setPattern
    автоматически; отдельно не выносим.
  - operator= как Go-синтаксис — используйте Assign.
  - Конструктор QStringMatcher(const QChar *, qsizetype, cs) — покрывается
    NewQStringMatcherFromQChars (эквивалентен NewQStringMatcherFromQStringView).
  - indexIn(const QChar *, qsizetype, from) — покрывается IndexInQChars
    (эквивалентен IndexInView).

См. qstringmatcher_doc.go для подробностей на английском.
*/

/*
#cgo pkg-config: Qt6Core
#cgo CXXFLAGS: -std=c++17 -fPIC
#cgo LDFLAGS: -lstdc++
#include <stdint.h>
#include <stdbool.h>
#include <stddef.h>
#include "qstringmatcher_c.h"
*/
import "C"

import (
	"runtime"
	"sync/atomic"
	"unsafe"

	"github.com/vadimlxx/go-qt6core/qdata"
)

// qstringmatcherState holds the C++ QStringMatcher pointer in an atomic slot
// so that Close() and the runtime cleanup can race safely: whoever swaps the
// slot to nil owns the C++ object and is responsible for deleting it.
//
// Why *byte and not *C.QStringMatcher: cgo exposes C.QStringMatcher as an
// incomplete Go type (forward-declared only under C mode), and
// atomic.Pointer[T] requires T to be a complete type. The pointer value is
// unchanged; only its Go-side static type differs. Cast back to
// *C.QStringMatcher at C-call boundaries.
type qstringmatcherState struct {
	p atomic.Pointer[byte]
}

// cpp returns the underlying *C.QStringMatcher, or nil if the slot was
// claimed (closed) already. Read-only, does not take ownership.
func (s *qstringmatcherState) cpp() *C.QStringMatcher {
	if bp := s.p.Load(); bp != nil {
		return (*C.QStringMatcher)(unsafe.Pointer(bp))
	}
	return nil
}

// claim atomically sets the slot to nil and returns the *C.QStringMatcher
// that was there (if any). Exactly one caller gets a non-nil result for a
// given stored pointer; subsequent claimers see nil. This is what makes
// Close() vs cleanup safe against double free.
func (s *qstringmatcherState) claim() *C.QStringMatcher {
	if bp := s.p.Swap(nil); bp != nil {
		return (*C.QStringMatcher)(unsafe.Pointer(bp))
	}
	return nil
}

// QStringMatcher wraps a Qt6 QStringMatcher object.
type QStringMatcher struct {
	state   *qstringmatcherState
	cleanup runtime.Cleanup
}

// wrapQStringMatcher wraps an opaque C++ QStringMatcher* into *QStringMatcher.
// Takes ownership: cleanup (or Close) frees the C++ object.
func wrapQStringMatcher(p unsafe.Pointer) *QStringMatcher {
	if p == nil {
		return nil
	}
	st := &qstringmatcherState{}
	st.p.Store((*byte)(p))
	m := &QStringMatcher{state: st}
	m.cleanup = runtime.AddCleanup(m, func(st *qstringmatcherState) {
		if cpp := st.claim(); cpp != nil {
			C.qstringmatcher_delete(cpp)
		}
	}, st)
	return m
}

// handle returns the underlying C++ pointer, panicking if m is nil or closed.
func (m *QStringMatcher) handle() *C.QStringMatcher {
	if cpp, ok := m.tryHandle(); ok {
		return cpp
	}
	panic("qstringmatcher: use after Close (or nil receiver)")
}

// tryHandle is the non-panicking variant of handle. It returns nil, false
// if m is nil, has no state, or has already been closed.
func (m *QStringMatcher) tryHandle() (*C.QStringMatcher, bool) {
	if m == nil || m.state == nil {
		return nil, false
	}
	if cpp := m.state.cpp(); cpp != nil {
		return cpp, true
	}
	return nil, false
}

// Close releases the underlying C++ object immediately and cancels the
// cleanup. Close is optional, idempotent, and nil-safe.
// Qt: QStringMatcher::~QStringMatcher().
func (m *QStringMatcher) Close() {
	if m == nil || m.state == nil {
		return
	}
	m.cleanup.Stop()
	if cpp := m.state.claim(); cpp != nil {
		C.qstringmatcher_delete(cpp)
	}
}

// ---- constructors --------------------------------------------------------

// NewQStringMatcher returns a default-constructed matcher (empty pattern,
// CaseSensitive).

// NewQStringMatcherFromQString constructs a matcher with the given QString
// pattern and case sensitivity.
// Qt: QStringMatcher(const QString &pattern, Qt::CaseSensitivity cs).
func NewQStringMatcherFromQString(pattern *QString, cs qdata.CaseSensitivity) *QStringMatcher {
	var pp *C.QString
	if pattern != nil {
		pp, _ = pattern.tryHandle()
	}
	m := wrapQStringMatcher(unsafe.Pointer(
		C.qstringmatcher_new_from_qstring(pp, C.int(cs))))
	runtime.KeepAlive(pattern)
	return m
}

// NewQStringMatcherFromString constructs a matcher from a Go string (UTF-8
// decoded into QString).
// Qt: QStringMatcher(const QString &pattern, Qt::CaseSensitivity cs).
func NewQStringMatcherFromString(pattern string, cs qdata.CaseSensitivity) *QStringMatcher {
	qs := NewQStringFromString(pattern)
	defer qs.Close()
	return NewQStringMatcherFromQString(qs, cs)
}

// NewQStringMatcherFromQStringView constructs a matcher from a QStringView.
// Qt: QStringMatcher(QStringView pattern, Qt::CaseSensitivity cs).
func NewQStringMatcherFromQStringView(pattern QStringView, cs qdata.CaseSensitivity) *QStringMatcher {
	var p *C.uint16_t
	if pattern.Size() > 0 {
		p = (*C.uint16_t)(pattern.Utf16Unsafe())
	}
	m := wrapQStringMatcher(unsafe.Pointer(C.qstringmatcher_new_from_view(
		p, C.ptrdiff_t(pattern.Size()), C.int(cs))))
	runtime.KeepAlive(pattern)
	return m
}

// NewQStringMatcherFromQChars constructs a matcher from raw UTF-16 code units.
// Qt: QStringMatcher(const QChar *uc, qsizetype len, Qt::CaseSensitivity cs).
func NewQStringMatcherFromQChars(units []uint16, cs qdata.CaseSensitivity) *QStringMatcher {
	var p *C.uint16_t
	if len(units) > 0 {
		p = (*C.uint16_t)(unsafe.Pointer(&units[0]))
	}
	m := wrapQStringMatcher(unsafe.Pointer(C.qstringmatcher_new_from_view(
		p, C.ptrdiff_t(len(units)), C.int(cs))))
	runtime.KeepAlive(units)
	return m
}

// NewQStringMatcher dispatches to the matching constructor based on the
// concrete Go argument types. Qt: no single counterpart.
//
//	NewQStringMatcher()                                    → default
//	NewQStringMatcher(*QStringMatcher)                     → Clone
//	NewQStringMatcher(*QString)                            → from QString, CaseSensitive
//	NewQStringMatcher(string)                              → from Go string, CaseSensitive
//	NewQStringMatcher(QStringView)                         → from view, CaseSensitive
//	NewQStringMatcher([]uint16)                            → from code units, CaseSensitive
//	NewQStringMatcher(*QString, qdata.CaseSensitivity)     → from QString + cs
//	NewQStringMatcher(string, qdata.CaseSensitivity)       → from Go string + cs
//	NewQStringMatcher(QStringView, qdata.CaseSensitivity)  → from view + cs
//	NewQStringMatcher([]uint16, qdata.CaseSensitivity)     → from code units + cs
func NewQStringMatcher(arg ...any) *QStringMatcher {
	if len(arg) == 0 || arg[0] == nil {
		return wrapQStringMatcher(unsafe.Pointer(C.qstringmatcher_new()))
	}

	if len(arg) == 1 {
		switch v := arg[0].(type) {
		case *QStringMatcher:
			return v.Clone()
		case *QString:
			return NewQStringMatcherFromQString(v, qdata.CaseSensitive)
		case string:
			return NewQStringMatcherFromString(v, qdata.CaseSensitive)
		case QStringView:
			return NewQStringMatcherFromQStringView(v, qdata.CaseSensitive)
		case []uint16:
			return NewQStringMatcherFromQChars(v, qdata.CaseSensitive)
		}
	}

	if len(arg) == 2 {
		cs, ok := arg[1].(qdata.CaseSensitivity)
		if ok {
			switch v := arg[0].(type) {
			case *QString:
				return NewQStringMatcherFromQString(v, cs)
			case string:
				return NewQStringMatcherFromString(v, cs)
			case QStringView:
				return NewQStringMatcherFromQStringView(v, cs)
			case []uint16:
				return NewQStringMatcherFromQChars(v, cs)
			}
		}
	}

	panic("NewQStringMatcher: unsupported argument combination")
}

// ---- copy / assignment --------------------------------------------------

// Clone returns a copy of the matcher.
// Qt: QStringMatcher(const QStringMatcher &other).
func (m *QStringMatcher) Clone() *QStringMatcher {
	r := wrapQStringMatcher(unsafe.Pointer(C.qstringmatcher_new_copy(m.handle())))
	runtime.KeepAlive(m)
	return r
}

// Assign copies the contents of other into m.
// Qt: QStringMatcher &operator=(const QStringMatcher &other).
func (m *QStringMatcher) Assign(other *QStringMatcher) {
	C.qstringmatcher_assign(m.handle(), other.handle())
	runtime.KeepAlive(m)
	runtime.KeepAlive(other)
}

// ---- setters -------------------------------------------------------------

// SetPattern replaces the search pattern.
// Qt: QStringMatcher::setPattern(const QString &pattern).
func (m *QStringMatcher) SetPattern(pattern *QString) {
	var pp *C.QString
	if pattern != nil {
		pp, _ = pattern.tryHandle()
	}
	if pp == nil {
		return
	}
	C.qstringmatcher_set_pattern(m.handle(), pp)
	runtime.KeepAlive(m)
	runtime.KeepAlive(pattern)
}

// SetPatternString replaces the search pattern with a Go string (UTF-8 decoded).
// Qt: QStringMatcher::setPattern(const QString &pattern).
func (m *QStringMatcher) SetPatternString(s string) {
	qs := NewQStringFromString(s)
	defer qs.Close()
	m.SetPattern(qs)
}

// SetPatternView replaces the search pattern with the contents of a QStringView.
// Convenience — Qt has no direct setPattern(QStringView) overload.
func (m *QStringMatcher) SetPatternView(v QStringView) {
	qs := NewQStringFromQStringView(v)
	defer qs.Close()
	m.SetPattern(qs)
}

// SetCaseSensitivity changes the case sensitivity.
// Qt: QStringMatcher::setCaseSensitivity(Qt::CaseSensitivity cs).
func (m *QStringMatcher) SetCaseSensitivity(cs qdata.CaseSensitivity) {
	C.qstringmatcher_set_case_sensitivity(m.handle(), C.int(cs))
	runtime.KeepAlive(m)
}

// ---- indexIn -------------------------------------------------------------

// IndexIn returns the index of the first occurrence of the pattern in str,
// starting at offset from, or -1 if not found.
// Qt: QStringMatcher::indexIn(const QString &str, qsizetype from) const.
//
// Returns -1 if str is nil.
func (m *QStringMatcher) IndexIn(str *QString, from int) int {
	var sp *C.QString
	if str != nil {
		sp, _ = str.tryHandle()
	}
	if sp == nil {
		return -1
	}
	r := int(C.qstringmatcher_index_in_qstring(m.handle(), sp, C.ptrdiff_t(from)))
	runtime.KeepAlive(m)
	runtime.KeepAlive(str)
	return r
}

// IndexInString is a Go-string variant of IndexIn (UTF-8 decoded into QString).
// Qt: QStringMatcher::indexIn(const QString &str, qsizetype from) const.
func (m *QStringMatcher) IndexInString(s string, from int) int {
	qs := NewQStringFromString(s)
	defer qs.Close()
	return m.IndexIn(qs, from)
}

// IndexInView returns the index of the first occurrence of the pattern in v,
// starting at offset from, or -1 if not found.
// Qt: QStringMatcher::indexIn(QStringView str, qsizetype from) const.
func (m *QStringMatcher) IndexInView(v QStringView, from int) int {
	var p *C.uint16_t
	if v.Size() > 0 {
		p = (*C.uint16_t)(v.Utf16Unsafe())
	}
	r := int(C.qstringmatcher_index_in_view(m.handle(), p,
		C.ptrdiff_t(v.Size()), C.ptrdiff_t(from)))
	runtime.KeepAlive(m)
	runtime.KeepAlive(v)
	return r
}

// IndexInQChars returns the index of the first occurrence of the pattern in
// the UTF-16 code units units, starting at offset from, or -1 if not found.
// Qt: QStringMatcher::indexIn(const QChar *str, qsizetype length, qsizetype from) const.
func (m *QStringMatcher) IndexInQChars(units []uint16, from int) int {
	var p *C.uint16_t
	if len(units) > 0 {
		p = (*C.uint16_t)(unsafe.Pointer(&units[0]))
	}
	r := int(C.qstringmatcher_index_in_view(m.handle(), p,
		C.ptrdiff_t(len(units)), C.ptrdiff_t(from)))
	runtime.KeepAlive(m)
	runtime.KeepAlive(units)
	return r
}

// ---- getters -------------------------------------------------------------

// Pattern returns the current search pattern as a fresh *QString.
// Qt: QStringMatcher::pattern() const.
func (m *QStringMatcher) Pattern() *QString {
	r := wrapQString(unsafe.Pointer(C.qstringmatcher_pattern(m.handle())))
	runtime.KeepAlive(m)
	return r
}

// PatternString returns the current search pattern as a Go string (UTF-8).
// Qt: QStringMatcher::pattern() const (encoded as UTF-8).
func (m *QStringMatcher) PatternString() string {
	p := m.Pattern()
	if p == nil {
		return ""
	}
	defer p.Close()
	return p.GoString()
}

// PatternView returns a non-owning QStringView over the internal pattern.
// Qt: QStringMatcher::patternView() const noexcept.
//
// The view is invalidated by any subsequent SetPattern*, Assign, or Close.
// Use Pattern() or PatternString() if you need a stable copy.
func (m *QStringMatcher) PatternView() QStringView {
	p := m.handle()
	var n C.ptrdiff_t
	data := C.qstringmatcher_pattern_view(p, &n)
	var viewPtr unsafe.Pointer
	if data != nil {
		viewPtr = unsafe.Pointer(data)
	}
	v := newQStringViewWithOwner(viewPtr, int(n), m)
	runtime.KeepAlive(m)
	return v
}

// CaseSensitivity returns the current case sensitivity.
// Qt: QStringMatcher::caseSensitivity() const.
func (m *QStringMatcher) CaseSensitivity() qdata.CaseSensitivity {
	r := qdata.CaseSensitivity(C.qstringmatcher_case_sensitivity(m.handle()))
	runtime.KeepAlive(m)
	return r
}
