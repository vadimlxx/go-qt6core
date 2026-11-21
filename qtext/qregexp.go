package qtext

/*
	QRegularExpression — handle над C++ QRegularExpression; владение через
	runtime.AddCleanup + Close().

	# Семантика

	Каждый *QRegularExpression — handle на C++ объект в куче. Копирование
	(Clone) дёшево. Строки шаблона передаются как *QString или Go string
	(неявная UTF-8 → QString конверсия).

	# Владение

	C++ объект освобождается одним из двух способов:
	- Close() — немедленно, детерминированно;
	- runtime.AddCleanup — когда *QRegularExpression становится недостижим.

	Close() опционален, идемпотентен и безопасен на nil. Если Close() и cleanup
	гоняются, владение захватывает ровно один участник (атомарный Swap(nil)),
	двойного free нет.

	# Use after Close

	После Close() любой метод, кроме Close(), паникует ("qregexp: use after
	Close (or nil receiver)"). Nil-safe методы используют tryHandle().

	# Соглашения по именам

	Все конструкторы, возвращающие *QRegularExpression или *QString, начинаются
	с NewQRegularExpression* / NewQStringFromRegexp*.

	# Enum'ы

	Enum-ы Qt (PatternOption, MatchType, MatchOption,
	WildcardConversionOption) проброшены как int-совместимые Go-типы с
	константами, повторяющими значения из qregularexpression.h. Комбинация
	флагов — через побитовое ИЛИ.

	# Что сознательно не вошло

	- namedCaptureGroups() — возвращает QStringList.
	- Операторы сравнения == и != — заменяются на будущий Compare().
	- qHash(QRegularExpression) — не идиоматично для Go.
	- optimize() — no-op (в Qt это подсказка движку).
	- swap() — реализован как метод Swap().
	- Range-based for по итератору — используйте for it.HasNext() { it.Next() }.
	- isNull/isEmpty-семантика QRegularExpression — в Qt выражается через
		pattern().isEmpty(), отдельных методов не заводим.
*/

/*
#cgo pkg-config: Qt6Core
#cgo CXXFLAGS: -std=c++17 -fPIC
#cgo LDFLAGS: -lstdc++
#include <stdint.h>
#include <stdbool.h>
#include <stddef.h>
#include "qregexp_c.h"
*/
import "C"

import (
	"runtime"
	"sync/atomic"
	"unsafe"
)

// PatternOption mirrors QRegularExpression::PatternOption.
type PatternOption int

const (
	NoPatternOption             PatternOption = 0x0000
	CaseInsensitiveOption       PatternOption = 0x0001
	DotMatchesEverythingOption  PatternOption = 0x0002
	MultilineOption             PatternOption = 0x0004
	ExtendedPatternSyntaxOption PatternOption = 0x0008
	InvertedGreedinessOption    PatternOption = 0x0010
	DontCaptureOption           PatternOption = 0x0020
	UseUnicodePropertiesOption  PatternOption = 0x0040
)

// MatchType mirrors QRegularExpression::MatchType.
type MatchType int

const (
	NormalMatch                MatchType = 0
	PartialPreferCompleteMatch MatchType = 1
	PartialPreferFirstMatch    MatchType = 2
	NoMatch                    MatchType = 3
)

// MatchOption mirrors QRegularExpression::MatchOption.
type MatchOption int

const (
	NoMatchOption             MatchOption = 0x0000
	AnchorAtOffsetMatchOption MatchOption = 0x0001
	// Deprecated: use AnchorAtOffsetMatchOption.
	AnchoredMatchOption               MatchOption = AnchorAtOffsetMatchOption
	DontCheckSubjectStringMatchOption MatchOption = 0x0002
)

// WildcardConversionOption mirrors QRegularExpression::WildcardConversionOption.
type WildcardConversionOption int

const (
	DefaultWildcardConversion    WildcardConversionOption = 0x0
	UnanchoredWildcardConversion WildcardConversionOption = 0x1
	NonPathWildcardConversion    WildcardConversionOption = 0x2
)

// qregexpState holds the C++ QRegularExpression pointer in an atomic slot so
// that Close() and the runtime cleanup can race safely: whoever swaps the
// slot to nil owns the C++ object and is responsible for deleting it.
//
// Why *byte and not *C.QRegularExpression: cgo exposes C.QRegularExpression
// as an incomplete Go type (forward-declared only under C mode), and
// atomic.Pointer[T] requires T to be a complete type. The pointer value is
// unchanged; only its Go-side static type differs. Cast back to
// *C.QRegularExpression at C-call boundaries.
type qregexpState struct {
	p atomic.Pointer[byte]
}

// cpp returns the underlying *C.QRegularExpression, or nil if the slot was
// claimed (closed) already. Read-only, does not take ownership.
func (s *qregexpState) cpp() *C.QRegularExpression {
	if bp := s.p.Load(); bp != nil {
		return (*C.QRegularExpression)(unsafe.Pointer(bp))
	}
	return nil
}

// claim atomically sets the slot to nil and returns the *C.QRegularExpression
// that was there (if any). Exactly one caller gets a non-nil result for a
// given stored pointer; subsequent claimers see nil. This is what makes
// Close() vs cleanup safe against double free.
func (s *qregexpState) claim() *C.QRegularExpression {
	if bp := s.p.Swap(nil); bp != nil {
		return (*C.QRegularExpression)(unsafe.Pointer(bp))
	}
	return nil
}

// QRegularExpression wraps a Qt6 QRegularExpression object.
// Call Close() to release the underlying C++ object immediately;
// otherwise a runtime cleanup handles it.
type QRegularExpression struct {
	state   *qregexpState
	cleanup runtime.Cleanup
}

// wrapQRegularExpression wraps an opaque C++ QRegularExpression* into
// *QRegularExpression. Takes ownership: cleanup (or Close) frees the C++
// object.
func wrapQRegularExpression(p unsafe.Pointer) *QRegularExpression {
	if p == nil {
		return nil
	}
	st := &qregexpState{}
	st.p.Store((*byte)(p))
	q := &QRegularExpression{state: st}
	q.cleanup = runtime.AddCleanup(q, func(st *qregexpState) {
		if cpp := st.claim(); cpp != nil {
			C.qregexp_delete(cpp)
		}
	}, st)
	return q
}

// handle returns the underlying C++ pointer, panicking if q is nil or closed.
func (q *QRegularExpression) handle() *C.QRegularExpression {
	if cpp, ok := q.tryHandle(); ok {
		return cpp
	}
	panic("qregexp: use after Close (or nil receiver)")
}

// tryHandle is the non-panicking variant of handle. It returns nil, false
// if q is nil, has no state, or has already been closed.
func (q *QRegularExpression) tryHandle() (*C.QRegularExpression, bool) {
	if q == nil || q.state == nil {
		return nil, false
	}
	if cpp := q.state.cpp(); cpp != nil {
		return cpp, true
	}
	return nil, false
}

// Close releases the underlying C++ object immediately and cancels the
// cleanup. Close is optional, idempotent, and nil-safe.
// Qt: QRegularExpression::~QRegularExpression().
func (q *QRegularExpression) Close() {
	if q == nil || q.state == nil {
		return
	}
	q.cleanup.Stop()
	if cpp := q.state.claim(); cpp != nil {
		C.qregexp_delete(cpp)
	}
}

// NewQRegularExpression dispatches to the matching constructor based on the
// concrete Go argument types. Qt: no single counterpart.
//
//	NewQRegularExpression()                       → empty pattern, no options
//	NewQRegularExpression(*QRegularExpression)    → Clone
//	NewQRegularExpression(*QString)               → from pattern, no options
//	NewQRegularExpression(string)                 → from pattern (UTF-8), no options
//	NewQRegularExpression(*QString, PatternOption) → from pattern + options
//	NewQRegularExpression(string, PatternOption)   → from pattern (UTF-8) + options
func NewQRegularExpression(arg ...any) *QRegularExpression {
	if len(arg) == 0 || arg[0] == nil {
		return wrapQRegularExpression(unsafe.Pointer(C.qregexp_new()))
	}

	if len(arg) == 1 {
		switch v := arg[0].(type) {
		case *QRegularExpression:
			return v.Clone()
		case *QString:
			return NewQRegularExpressionFromPattern(v, NoPatternOption)
		case string:
			return NewQRegularExpressionFromString(v, NoPatternOption)
		}
	}

	if len(arg) == 2 {
		opt, ok := arg[1].(PatternOption)
		if ok {
			switch s := arg[0].(type) {
			case *QString:
				return NewQRegularExpressionFromPattern(s, opt)
			case string:
				return NewQRegularExpressionFromString(s, opt)
			}
		}
	}

	panic("NewQRegularExpression: unsupported argument combination")
}

// NewQRegularExpressionFromPattern — Qt: QRegularExpression(const QString &pattern, PatternOptions).
func NewQRegularExpressionFromPattern(pattern *QString, options PatternOption) *QRegularExpression {
	r := wrapQRegularExpression(unsafe.Pointer(
		C.qregexp_new_from_pattern(pattern.handle(), C.int(options))))
	runtime.KeepAlive(pattern)
	return r
}

// NewQRegularExpressionFromString — Qt: QRegularExpression(const QString &pattern, PatternOptions)
// with the string decoded from UTF-8.
func NewQRegularExpressionFromString(pattern string, options PatternOption) *QRegularExpression {
	p := NewQStringFromString(pattern)
	defer p.Close()
	return NewQRegularExpressionFromPattern(p, options)
}

// Clone returns a COW copy.
// Qt: QRegularExpression(const QRegularExpression&).
func (q *QRegularExpression) Clone() *QRegularExpression {
	r := wrapQRegularExpression(unsafe.Pointer(C.qregexp_new_copy(q.handle())))
	runtime.KeepAlive(q)
	return r
}

// Pattern — Qt: QRegularExpression::pattern() const.
func (q *QRegularExpression) Pattern() *QString {
	r := wrapQString(unsafe.Pointer(C.qregexp_pattern(q.handle())))
	runtime.KeepAlive(q)
	return r
}

// SetPattern — Qt: QRegularExpression::setPattern(const QString&).
func (q *QRegularExpression) SetPattern(pattern *QString) {
	C.qregexp_set_pattern(q.handle(), pattern.handle())
	runtime.KeepAlive(q)
	runtime.KeepAlive(pattern)
}

// IsValid — Qt: QRegularExpression::isValid() const.
func (q *QRegularExpression) IsValid() bool {
	r := bool(C.qregexp_is_valid(q.handle()))
	runtime.KeepAlive(q)
	return r
}

// PatternErrorOffset — Qt: QRegularExpression::patternErrorOffset() const.
func (q *QRegularExpression) PatternErrorOffset() int {
	r := int(C.qregexp_pattern_error_offset(q.handle()))
	runtime.KeepAlive(q)
	return r
}

// ErrorString — Qt: QRegularExpression::errorString() const.
func (q *QRegularExpression) ErrorString() *QString {
	r := wrapQString(unsafe.Pointer(C.qregexp_error_string(q.handle())))
	runtime.KeepAlive(q)
	return r
}

// CaptureCount — Qt: QRegularExpression::captureCount() const.
func (q *QRegularExpression) CaptureCount() int {
	r := int(C.qregexp_capture_count(q.handle()))
	runtime.KeepAlive(q)
	return r
}

// PatternOptions — Qt: QRegularExpression::patternOptions() const.
func (q *QRegularExpression) PatternOptions() PatternOption {
	r := PatternOption(C.qregexp_pattern_options(q.handle()))
	runtime.KeepAlive(q)
	return r
}

// SetPatternOptions — Qt: QRegularExpression::setPatternOptions(PatternOptions).
func (q *QRegularExpression) SetPatternOptions(options PatternOption) {
	C.qregexp_set_pattern_options(q.handle(), C.int(options))
	runtime.KeepAlive(q)
}

// MatchString — Qt: QRegularExpression::match(const QString&, qsizetype, MatchType, MatchOptions) const.
func (q *QRegularExpression) MatchString(subject *QString, offset int,
	matchType MatchType, matchOptions MatchOption) *QRegularExpressionMatch {
	m := wrapQRegularExpressionMatch(unsafe.Pointer(C.qregexp_match_string(
		q.handle(), subject.handle(),
		C.ptrdiff_t(offset),
		C.int(matchType),
		C.int(matchOptions),
	)))
	if m != nil {
		m.keepAlive = subject
	}
	runtime.KeepAlive(q)
	runtime.KeepAlive(subject)
	return m
}

// MatchView — Qt: QRegularExpression::match(QStringView, qsizetype, MatchType, MatchOptions) const.
func (q *QRegularExpression) MatchView(subject QStringView, offset int,
	matchType MatchType, matchOptions MatchOption) *QRegularExpressionMatch {
	m := wrapQRegularExpressionMatch(unsafe.Pointer(C.qregexp_match_view(q.handle(),
		(*C.uint16_t)(subject.Utf16Unsafe()), C.ptrdiff_t(subject.Size()),
		C.ptrdiff_t(offset), C.int(matchType), C.int(matchOptions))))
	if m != nil {
		m.keepAlive = subject
	}
	runtime.KeepAlive(q)
	runtime.KeepAlive(subject)
	return m
}

// GlobalMatchString — Qt: QRegularExpression::globalMatch(const QString&, qsizetype, MatchType, MatchOptions) const.
func (q *QRegularExpression) GlobalMatchString(subject *QString, offset int,
	matchType MatchType, matchOptions MatchOption) *QRegularExpressionMatchIterator {
	r := wrapQRegularExpressionMatchIterator(unsafe.Pointer(C.qregexp_global_match_string(
		q.handle(), subject.handle(),
		C.ptrdiff_t(offset), C.int(matchType), C.int(matchOptions))))
	runtime.KeepAlive(q)
	runtime.KeepAlive(subject)
	return r
}

// GlobalMatchView — Qt: QRegularExpression::globalMatch(QStringView, qsizetype, MatchType, MatchOptions) const.
func (q *QRegularExpression) GlobalMatchView(subject QStringView, offset int,
	matchType MatchType, matchOptions MatchOption) *QRegularExpressionMatchIterator {
	r := wrapQRegularExpressionMatchIterator(unsafe.Pointer(C.qregexp_global_match_view(q.handle(),
		(*C.uint16_t)(subject.Utf16Unsafe()), C.ptrdiff_t(subject.Size()),
		C.ptrdiff_t(offset), C.int(matchType), C.int(matchOptions))))
	runtime.KeepAlive(q)
	runtime.KeepAlive(subject)
	return r
}

// Optimize is a no-op in this binding; the C ABI does not expose
// QRegularExpression::optimize().
func (q *QRegularExpression) Optimize() {}

// Swap swaps q and other (no data copy).
// Qt: QRegularExpression::swap(QRegularExpression&).
func (q *QRegularExpression) Swap(other *QRegularExpression) {
	C.qregexp_swap(q.handle(), other.handle())
	runtime.KeepAlive(q)
	runtime.KeepAlive(other)
}

// NewQStringFromRegexpEscape — Qt: QRegularExpression::escape(QStringView) → QString.
func NewQStringFromRegexpEscape(str *QString) *QString {
	r := wrapQString(unsafe.Pointer(C.qregexp_escape(str.handle())))
	runtime.KeepAlive(str)
	return r
}

// NewQStringFromRegexpWildcard — Qt: QRegularExpression::wildcardToRegularExpression(QStringView, WildcardConversionOptions) → QString.
func NewQStringFromRegexpWildcard(str *QString, options WildcardConversionOption) *QString {
	r := wrapQString(unsafe.Pointer(C.qregexp_wildcard_to_regular_expression(
		str.handle(), C.int(options))))
	runtime.KeepAlive(str)
	return r
}

// NewQStringFromRegexpAnchoredPattern — Qt: QRegularExpression::anchoredPattern(QStringView) → QString.
func NewQStringFromRegexpAnchoredPattern(expression *QString) *QString {
	r := wrapQString(unsafe.Pointer(C.qregexp_anchored_pattern(expression.handle())))
	runtime.KeepAlive(expression)
	return r
}

// NewQRegularExpressionFromWildcard — Qt: QRegularExpression::fromWildcard(QStringView, Qt::CaseSensitivity, WildcardConversionOptions).
func NewQRegularExpressionFromWildcard(pattern *QString, cs bool,
	options WildcardConversionOption) *QRegularExpression {
	csInt := 0
	if !cs {
		csInt = 1
	}
	r := wrapQRegularExpression(unsafe.Pointer(C.qregexp_from_wildcard(
		pattern.handle(), C.int(csInt), C.int(options))))
	runtime.KeepAlive(pattern)
	return r
}
