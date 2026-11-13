package qtext

/*
   QRegularExpression оборачивает C++ QRegularExpression через непрозрачный
   указатель (C ABI). Go-структура содержит только этот указатель; владение
   обеспечивается финализатором и явным Close(), как в QString и QByteArray.

   Enum-ы Qt (PatternOption, MatchType, MatchOption, WildcardConversionOption)
   проброшены как int-совместимые Go-типы с константами, повторяющими
   значения из qregularexpression.h. Комбинация флагов — через побитовое ИЛИ
   (в C++ это Q_DECLARE_FLAGS; в Go-биндинге отдельный тип для набора флагов
   не заводим, поскольку методы принимают int-совместимые значения).

   C ABI-заголовок объединённый (qregexp_c.h / qregexp_c.cpp) — как в Qt6
   qregularexpression.h, где три класса QRegularExpression, Match и MatchIterator
   лежат рядом. Go-файлы при этом раздельные (qregexp.go, qregexpmatch.go,
   qregexpmatchiter.go) — по одному на кластер типа, что соответствует
   идиоматике Go и gopls/go doc.

   Все match-методы (MatchString / MatchView / GlobalMatchString /
   GlobalMatchView) возвращают уже сконструированные C++-объекты
   QRegularExpressionMatch и QRegularExpressionMatchIterator — их
   владение переходит вызывающему.

   Конструкторы принимают либо *QString (сохраняет семантическое различие
   null/empty), либо Go string — удобный вариант с неявной конвертацией
   UTF-8 → QString.

 Что сознательно не вошло (пока):

   - namedCaptureGroups() — возвращает QStringList (нужна специализация
     контейнера QStringList).
   - Операторы сравнения == и != — заменяются на будущий Compare()
     (потребует is_equal-хелпер, повторяющий Q_DECLARE_EQUALITY_COMPARABLE).
   - qHash(QRegularExpression) — не идиоматично для Go.
   - optimize() — оставлен как no-op, поскольку вызов из C ABI не даёт
     наблюдаемого эффекта; в Qt это подсказка внутреннему движку.
   - swap() — реализован как метод Swap() (не как free-функция,
     чтобы не загромождать API).
   - Range-based for по итератору — неприменимо в Go; используйте цикл
     for it.HasNext() { it.Next() }.
   - isNull/isEmpty-семантика QRegularExpression — в Qt выражается
     через pattern().isEmpty(), отдельных методов не заводим.
*/

/*
#cgo pkg-config: Qt6Core
#cgo CXXFLAGS: -std=c++17 -fPIC -Wno-sfinae-incomplete
#cgo LDFLAGS: -lstdc++
#include <stdint.h>
#include <stdbool.h>
#include <stddef.h>
#include "qregexp_c.h"
*/
import "C"

import (
	"runtime"
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
	NoMatchOption                     MatchOption = 0x0000
	AnchorAtOffsetMatchOption         MatchOption = 0x0001
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

// QRegularExpression wraps a Qt6 QRegularExpression object.
// Call Close() to release the underlying C++ object immediately;
// otherwise a finalizer handles cleanup.
type QRegularExpression struct {
	p *C.QRegularExpression
}

func wrapQRegularExpression(p unsafe.Pointer) *QRegularExpression {
	if p == nil {
		return nil
	}
	q := &QRegularExpression{p: (*C.QRegularExpression)(p)}
	runtime.SetFinalizer(q, finalizeQRegularExpression)
	return q
}

func finalizeQRegularExpression(q *QRegularExpression) {
	if q.p != nil {
		C.qregexp_delete(q.p)
		q.p = nil
	}
}

func (q *QRegularExpression) Close() {
	if q.p != nil {
		C.qregexp_delete(q.p)
		q.p = nil
		runtime.SetFinalizer(q, nil)
	}
}

func NewQRegularExpression() *QRegularExpression {
	return wrapQRegularExpression(unsafe.Pointer(C.qregexp_new()))
}

func NewQRegularExpressionFromPattern(pattern *QString, options PatternOption) *QRegularExpression {
	return wrapQRegularExpression(unsafe.Pointer(C.qregexp_new_from_pattern(pattern.p, C.int(options))))
}

func NewQRegularExpressionFromString(pattern string, options PatternOption) *QRegularExpression {
	p := QStringFromString(pattern)
	defer p.Close()
	return NewQRegularExpressionFromPattern(p, options)
}

func (q *QRegularExpression) Clone() *QRegularExpression {
	return wrapQRegularExpression(unsafe.Pointer(C.qregexp_new_copy(q.p)))
}

func (q *QRegularExpression) Pattern() *QString {
	return wrapQString(unsafe.Pointer(C.qregexp_pattern(q.p)))
}

func (q *QRegularExpression) SetPattern(pattern *QString) {
	C.qregexp_set_pattern(q.p, pattern.p)
}

func (q *QRegularExpression) IsValid() bool {
	return bool(C.qregexp_is_valid(q.p))
}

func (q *QRegularExpression) PatternErrorOffset() int {
	return int(C.qregexp_pattern_error_offset(q.p))
}

func (q *QRegularExpression) ErrorString() *QString {
	return wrapQString(unsafe.Pointer(C.qregexp_error_string(q.p)))
}

func (q *QRegularExpression) CaptureCount() int {
	return int(C.qregexp_capture_count(q.p))
}

func (q *QRegularExpression) PatternOptions() PatternOption {
	return PatternOption(C.qregexp_pattern_options(q.p))
}

func (q *QRegularExpression) SetPatternOptions(options PatternOption) {
	C.qregexp_set_pattern_options(q.p, C.int(options))
}

func (q *QRegularExpression) MatchString(subject *QString, offset int,
	matchType MatchType, matchOptions MatchOption) *QRegularExpressionMatch {
	m := wrapQRegularExpressionMatch(unsafe.Pointer(C.qregexp_match_string(q.p, subject.p,
		C.ptrdiff_t(offset), C.int(matchType), C.int(matchOptions))))
	if m != nil {
		m.keepAlive = subject
	}
	return m
}

func (q *QRegularExpression) MatchView(subject QStringView, offset int,
	matchType MatchType, matchOptions MatchOption) *QRegularExpressionMatch {
	m := wrapQRegularExpressionMatch(unsafe.Pointer(C.qregexp_match_view(q.p,
		(*C.uint16_t)(subject.Utf16Unsafe()), C.ptrdiff_t(subject.Size()),
		C.ptrdiff_t(offset), C.int(matchType), C.int(matchOptions))))
	if m != nil {
		m.keepAlive = subject
	}
	return m
}

func (q *QRegularExpression) GlobalMatchString(subject *QString, offset int,
	matchType MatchType, matchOptions MatchOption) *QRegularExpressionMatchIterator {
	return wrapQRegularExpressionMatchIterator(unsafe.Pointer(C.qregexp_global_match_string(q.p, subject.p,
		C.ptrdiff_t(offset), C.int(matchType), C.int(matchOptions))))
}

func (q *QRegularExpression) GlobalMatchView(subject QStringView, offset int,
	matchType MatchType, matchOptions MatchOption) *QRegularExpressionMatchIterator {
	return wrapQRegularExpressionMatchIterator(unsafe.Pointer(C.qregexp_global_match_view(q.p,
		(*C.uint16_t)(subject.Utf16Unsafe()), C.ptrdiff_t(subject.Size()),
		C.ptrdiff_t(offset), C.int(matchType), C.int(matchOptions))))
}

// Optimize is a no-op in this binding; the C ABI does not expose
// QRegularExpression::optimize().
func (q *QRegularExpression) Optimize() {}

func (q *QRegularExpression) Swap(other *QRegularExpression) {
	C.qregexp_swap(q.p, other.p)
}

func QRegularExpressionEscape(str *QString) *QString {
	return wrapQString(unsafe.Pointer(C.qregexp_escape(str.p)))
}

func QRegularExpressionWildcardToRegularExpression(str *QString, options WildcardConversionOption) *QString {
	return wrapQString(unsafe.Pointer(C.qregexp_wildcard_to_regular_expression(str.p, C.int(options))))
}

func QRegularExpressionAnchoredPattern(expression *QString) *QString {
	return wrapQString(unsafe.Pointer(C.qregexp_anchored_pattern(expression.p)))
}

func QRegularExpressionFromWildcard(pattern *QString, cs bool,
	options WildcardConversionOption) *QRegularExpression {
	csInt := 0
	if !cs {
		csInt = 1
	}
	return wrapQRegularExpression(unsafe.Pointer(C.qregexp_from_wildcard(pattern.p, C.int(csInt), C.int(options))))
}
