package qdatetime

/*
QLocale — handle над C++ QLocale; владение через runtime.AddCleanup + Close().

# Семантика

QLocale — value-тип в Qt (COW через QSharedDataPointer). В Go — handle:
*QLocale указывает на C++ объект; копирование Go-указателя не создаёт
новую локаль. Для независимой копии — Clone().

# Владение

C++ объект освобождается одним из двух способов:
  - Close() — немедленно, детерминированно;
  - runtime.AddCleanup — когда *QLocale становится недостижим.

Close() опционален, идемпотентен и безопасен на nil. Если Close() и
cleanup гоняются, владение захватывает ровно один участник (атомарный
Swap(nil)), двойного free нет.

# Use after Close

После Close() любой метод, кроме Close(), паникует
("qlocale: use after Close (or nil receiver)").

# Соглашения по именам

  - Все строковые геттеры возвращают *qtext.QString (владение у
    вызывающего; Close()). Параллельно XxxString() возвращает Go string
    и закрывает промежуточный QString автоматически.
  - Функции, принимающие на входе строку: Xxx(s string) — Go-строка
    (UTF-8); если есть QString-вариант — XxxQString(s *qtext.QString).
  - QDate / QTime — value-типы, передаются по значению. QDateTime /
    QCalendar — handle-типы (через внутренний handle()).

# Дата/время-методы

  - ToStringDateFormat(d, ft), ToStringDateString(d, fmt)
  - ToStringTimeFormat(t, ft), ToStringTimeString(t, fmt)
  - ToStringDateTimeFormat(dt, ft), ToStringDateTimeString(dt, fmt)
  - Календарные варианты: суффикс ...Cal(..., cal *QCalendar)
  - DateFormat(ft), TimeFormat(ft), DateTimeFormat(ft)
  - ToDateFormat(s, ft, baseYear) → QDate (нулевой QDate при parse error)
  - ToDateString(s, fmt, baseYear) → QDate
  - ToDateFormatCal(s, ft, cal, baseYear) → QDate
  - ToDateStringCal(s, fmt, cal, baseYear) → QDate
  - ToTimeFormat(s, ft), ToTimeString(s, fmt)
  - ToDateTimeFormat(s, ft, baseYear) → *QDateTime (проверяйте IsValid())
  - ToDateTimeString(s, fmt, baseYear) → *QDateTime
  - ToDateTimeFormatCal(s, ft, cal, baseYear) → *QDateTime
  - ToDateTimeStringCal(s, fmt, cal, baseYear) → *QDateTime

# baseYear

DefaultTwoDigitBaseYear = 1900 — соответствует
QLocale::DefaultTwoDigitBaseYear в Qt. Влияет только на разбор
двухзначных годов в строке.

# Что не реализовано (TODO)

  - QDataStream / QDebug операторы.
  - qHash(QLocale).
  - Deprecated-члены Qt (country(), countriesForLanguage(), name()/bcp47Name()
    без TagSeparator) намеренно не экспонируются: используйте Territory(),
    MatchingLocales() и Name/Bcp47Name с явным TagSeparator.
*/

/*
#cgo pkg-config: Qt6Core
#cgo CXXFLAGS: -std=c++17 -fPIC
#cgo LDFLAGS: -lstdc++
#include <stdint.h>
#include <stdbool.h>
#include <stddef.h>
#include "qlocale_c.h"
*/
import "C"

import (
	"runtime"
	"sync/atomic"
	"unsafe"

	"github.com/vadimlxx/go-qt6core/qtext"
)

// DefaultTwoDigitBaseYear — mirrors QLocale::DefaultTwoDigitBaseYear.
// Used by toDate/toDateTime when interpreting two-digit years.
const DefaultTwoDigitBaseYear = 1900

// ---- helpers -------------------------------------------------------------

// wrapQStr wraps a raw *C.QString into *qtext.QString (ownership: caller).
func wrapQStr(p unsafe.Pointer) *qtext.QString {
	if p == nil {
		return nil
	}
	return qtext.QStringFromUnsafePtr(p)
}

// qstrToGo reads UTF-8 string from q and closes q.
func qstrToGo(q *qtext.QString) string {
	if q == nil {
		return ""
	}
	defer q.Close()
	return q.GoString()
}

// ---- state ---------------------------------------------------------------

type qlocaleState struct {
	p atomic.Pointer[byte]
}

func (s *qlocaleState) cpp() C.go_qlocale_t {
	if bp := s.p.Load(); bp != nil {
		return C.go_qlocale_t(unsafe.Pointer(bp))
	}
	return nil
}

func (s *qlocaleState) claim() C.go_qlocale_t {
	if bp := s.p.Swap(nil); bp != nil {
		return C.go_qlocale_t(unsafe.Pointer(bp))
	}
	return nil
}

// QLocale wraps a heap-allocated Qt6 QLocale.
type QLocale struct {
	state   *qlocaleState
	cleanup runtime.Cleanup
}

func wrapQLocale(p C.go_qlocale_t) *QLocale {
	if p == nil {
		return nil
	}
	st := &qlocaleState{}
	st.p.Store((*byte)(unsafe.Pointer(p)))
	q := &QLocale{state: st}
	q.cleanup = runtime.AddCleanup(q, func(st *qlocaleState) {
		if cpp := st.claim(); cpp != nil {
			C.qlocale_delete(cpp)
		}
	}, st)
	return q
}

func (q *QLocale) handle() C.go_qlocale_t {
	if cpp, ok := q.tryHandle(); ok {
		return cpp
	}
	panic("qlocale: use after Close (or nil receiver)")
}

func (q *QLocale) tryHandle() (C.go_qlocale_t, bool) {
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
// Qt: QLocale::~QLocale().
func (q *QLocale) Close() {
	if q == nil || q.state == nil {
		return
	}
	q.cleanup.Stop()
	if cpp := q.state.claim(); cpp != nil {
		C.qlocale_delete(cpp)
	}
}

// ---- construction --------------------------------------------------------

// NewQLocale — Qt: QLocale().
func NewQLocale() *QLocale { return wrapQLocale(C.qlocale_new_default()) }

// NewQLocaleFromName — Qt: QLocale(const QString&) / QLocale(QStringView).
func NewQLocaleFromName(name string) *QLocale {
	q := wrapQLocale(C.qlocale_new_from_name(cstr(name), C.ptrdiff_t(len(name))))
	runtime.KeepAlive(name)
	return q
}

// NewQLocaleFromLanguageTerritory — Qt: QLocale(Language, Territory).
func NewQLocaleFromLanguageTerritory(lang Language, territory Territory) *QLocale {
	return wrapQLocale(C.qlocale_new_from_language_territory(C.int(lang), C.int(territory)))
}

// NewQLocaleFromLanguageScriptTerritory — Qt: QLocale(Language, Script, Territory).
func NewQLocaleFromLanguageScriptTerritory(lang Language, script Script, territory Territory) *QLocale {
	return wrapQLocale(C.qlocale_new_from_language_script_territory(
		C.int(lang), C.int(script), C.int(territory)))
}

// Clone — independent copy of q. Qt: QLocale(const QLocale&).
func (q *QLocale) Clone() *QLocale {
	r := wrapQLocale(C.qlocale_new_copy(q.handle()))
	runtime.KeepAlive(q)
	return r
}

// Swap — Qt: QLocale::swap(QLocale&).
func (q *QLocale) Swap(other *QLocale) {
	if op, ok := other.tryHandle(); ok {
		C.qlocale_swap(q.handle(), op)
		runtime.KeepAlive(q)
		runtime.KeepAlive(other)
	}
}

// ---- value accessors -----------------------------------------------------

// Language — Qt: QLocale::language() const.
func (q *QLocale) Language() Language {
	r := Language(C.qlocale_language(q.handle()))
	runtime.KeepAlive(q)
	return r
}

// Script — Qt: QLocale::script() const.
func (q *QLocale) Script() Script {
	r := Script(C.qlocale_script(q.handle()))
	runtime.KeepAlive(q)
	return r
}

// Territory — Qt: QLocale::territory() const.
func (q *QLocale) Territory() Territory {
	r := Territory(C.qlocale_territory(q.handle()))
	runtime.KeepAlive(q)
	return r
}

// MeasurementSystem — Qt: QLocale::measurementSystem() const.
func (q *QLocale) MeasurementSystem() MeasurementSystem {
	r := MeasurementSystem(C.qlocale_measurement_system(q.handle()))
	runtime.KeepAlive(q)
	return r
}

// TextDirection — Qt: QLocale::textDirection() const.
func (q *QLocale) TextDirection() int {
	r := int(C.qlocale_text_direction(q.handle()))
	runtime.KeepAlive(q)
	return r
}

// FirstDayOfWeek — Qt: QLocale::firstDayOfWeek() const.
func (q *QLocale) FirstDayOfWeek() int {
	r := int(C.qlocale_first_day_of_week(q.handle()))
	runtime.KeepAlive(q)
	return r
}

// NumberOptions — Qt: QLocale::numberOptions() const.
func (q *QLocale) NumberOptions() NumberOptions {
	r := NumberOptions(C.qlocale_number_options(q.handle()))
	runtime.KeepAlive(q)
	return r
}

// SetNumberOptions — Qt: QLocale::setNumberOptions(NumberOptions).
func (q *QLocale) SetNumberOptions(opts NumberOptions) {
	C.qlocale_set_number_options(q.handle(), C.int(opts))
	runtime.KeepAlive(q)
}

// ---- string representation ----------------------------------------------

// Name — Qt: QLocale::name(TagSeparator = Underscore).
func (q *QLocale) Name(sep TagSeparator) *qtext.QString {
	r := wrapQStr(unsafe.Pointer(C.qlocale_name(q.handle(), C.int(sep))))
	runtime.KeepAlive(q)
	return r
}
func (q *QLocale) NameString(sep TagSeparator) string { return qstrToGo(q.Name(sep)) }

// Bcp47Name — Qt: QLocale::bcp47Name(TagSeparator = Dash).
func (q *QLocale) Bcp47Name(sep TagSeparator) *qtext.QString {
	r := wrapQStr(unsafe.Pointer(C.qlocale_bcp47_name(q.handle(), C.int(sep))))
	runtime.KeepAlive(q)
	return r
}
func (q *QLocale) Bcp47NameString(sep TagSeparator) string { return qstrToGo(q.Bcp47Name(sep)) }

// NativeLanguageName — Qt: QLocale::nativeLanguageName() const.
func (q *QLocale) NativeLanguageName() *qtext.QString {
	r := wrapQStr(unsafe.Pointer(C.qlocale_native_language_name(q.handle())))
	runtime.KeepAlive(q)
	return r
}
func (q *QLocale) NativeLanguageNameString() string { return qstrToGo(q.NativeLanguageName()) }

// NativeTerritoryName — Qt: QLocale::nativeTerritoryName() const.
func (q *QLocale) NativeTerritoryName() *qtext.QString {
	r := wrapQStr(unsafe.Pointer(C.qlocale_native_territory_name(q.handle())))
	runtime.KeepAlive(q)
	return r
}
func (q *QLocale) NativeTerritoryNameString() string { return qstrToGo(q.NativeTerritoryName()) }

// ---- number parsing ------------------------------------------------------

// ToShort — Qt: QLocale::toShort(QStringView, bool*) const.
func (q *QLocale) ToShort(s string) (int16, bool) {
	p := q.handle()
	var v C.int16_t
	ok := bool(C.qlocale_to_short(p, cstr(s), C.ptrdiff_t(len(s)), &v))
	runtime.KeepAlive(q)
	runtime.KeepAlive(s)
	return int16(v), ok
}

// ToUShort — Qt: QLocale::toUShort(QStringView, bool*).
func (q *QLocale) ToUShort(s string) (uint16, bool) {
	p := q.handle()
	var v C.uint16_t
	ok := bool(C.qlocale_to_ushort(p, cstr(s), C.ptrdiff_t(len(s)), &v))
	runtime.KeepAlive(q)
	runtime.KeepAlive(s)
	return uint16(v), ok
}

// ToInt — Qt: QLocale::toInt(QStringView, bool*).
func (q *QLocale) ToInt(s string) (int32, bool) {
	p := q.handle()
	var v C.int32_t
	ok := bool(C.qlocale_to_int(p, cstr(s), C.ptrdiff_t(len(s)), &v))
	runtime.KeepAlive(q)
	runtime.KeepAlive(s)
	return int32(v), ok
}

// ToUInt — Qt: QLocale::toUInt(QStringView, bool*).
func (q *QLocale) ToUInt(s string) (uint32, bool) {
	p := q.handle()
	var v C.uint32_t
	ok := bool(C.qlocale_to_uint(p, cstr(s), C.ptrdiff_t(len(s)), &v))
	runtime.KeepAlive(q)
	runtime.KeepAlive(s)
	return uint32(v), ok
}

// ToLong — Qt: QLocale::toLong(QStringView, bool*).
func (q *QLocale) ToLong(s string) (int64, bool) {
	p := q.handle()
	var v C.int64_t
	ok := bool(C.qlocale_to_long(p, cstr(s), C.ptrdiff_t(len(s)), &v))
	runtime.KeepAlive(q)
	runtime.KeepAlive(s)
	return int64(v), ok
}

// ToULong — Qt: QLocale::toULong(QStringView, bool*).
func (q *QLocale) ToULong(s string) (uint64, bool) {
	p := q.handle()
	var v C.uint64_t
	ok := bool(C.qlocale_to_ulong(p, cstr(s), C.ptrdiff_t(len(s)), &v))
	runtime.KeepAlive(q)
	runtime.KeepAlive(s)
	return uint64(v), ok
}

// ToLongLong — Qt: QLocale::toLongLong(QStringView, bool*).
func (q *QLocale) ToLongLong(s string) (int64, bool) {
	p := q.handle()
	var v C.int64_t
	ok := bool(C.qlocale_to_longlong(p, cstr(s), C.ptrdiff_t(len(s)), &v))
	runtime.KeepAlive(q)
	runtime.KeepAlive(s)
	return int64(v), ok
}

// ToULongLong — Qt: QLocale::toULongLong(QStringView, bool*).
func (q *QLocale) ToULongLong(s string) (uint64, bool) {
	p := q.handle()
	var v C.uint64_t
	ok := bool(C.qlocale_to_ulonglong(p, cstr(s), C.ptrdiff_t(len(s)), &v))
	runtime.KeepAlive(q)
	runtime.KeepAlive(s)
	return uint64(v), ok
}

// ToFloat — Qt: QLocale::toFloat(QStringView, bool*).
func (q *QLocale) ToFloat(s string) (float32, bool) {
	p := q.handle()
	var v C.float
	ok := bool(C.qlocale_to_float(p, cstr(s), C.ptrdiff_t(len(s)), &v))
	runtime.KeepAlive(q)
	runtime.KeepAlive(s)
	return float32(v), ok
}

// ToDouble — Qt: QLocale::toDouble(QStringView, bool*).
func (q *QLocale) ToDouble(s string) (float64, bool) {
	p := q.handle()
	var v C.double
	ok := bool(C.qlocale_to_double(p, cstr(s), C.ptrdiff_t(len(s)), &v))
	runtime.KeepAlive(q)
	runtime.KeepAlive(s)
	return float64(v), ok
}

// ---- number formatting ---------------------------------------------------

// ToStringLL — Qt: QLocale::toString(qlonglong) const.
func (q *QLocale) ToStringLL(v int64) *qtext.QString {
	r := wrapQStr(unsafe.Pointer(C.qlocale_to_string_ll(q.handle(), C.int64_t(v))))
	runtime.KeepAlive(q)
	return r
}
func (q *QLocale) ToStringLLString(v int64) string { return qstrToGo(q.ToStringLL(v)) }

// ToStringULL — Qt: QLocale::toString(qulonglong) const.
func (q *QLocale) ToStringULL(v uint64) *qtext.QString {
	r := wrapQStr(unsafe.Pointer(C.qlocale_to_string_ull(q.handle(), C.uint64_t(v))))
	runtime.KeepAlive(q)
	return r
}
func (q *QLocale) ToStringULLString(v uint64) string { return qstrToGo(q.ToStringULL(v)) }

// ToStringLLW — Qt: QLocale::toString(qlonglong, int, char32_t) const.
func (q *QLocale) ToStringLLW(v int64, fieldWidth int, fillChar rune) *qtext.QString {
	r := wrapQStr(unsafe.Pointer(C.qlocale_to_string_ll_w(
		q.handle(), C.int64_t(v), C.int(fieldWidth), C.uint32_t(fillChar))))
	runtime.KeepAlive(q)
	return r
}

// ToStringULLW — Qt: QLocale::toString(qulonglong, int, char32_t) const.
func (q *QLocale) ToStringULLW(v uint64, fieldWidth int, fillChar rune) *qtext.QString {
	r := wrapQStr(unsafe.Pointer(C.qlocale_to_string_ull_w(
		q.handle(), C.uint64_t(v), C.int(fieldWidth), C.uint32_t(fillChar))))
	runtime.KeepAlive(q)
	return r
}

// ToStringDouble — Qt: QLocale::toString(double, char, int) const.
func (q *QLocale) ToStringDouble(v float64, format byte, precision int) *qtext.QString {
	r := wrapQStr(unsafe.Pointer(C.qlocale_to_string_double(
		q.handle(), C.double(v), C.char(format), C.int(precision))))
	runtime.KeepAlive(q)
	return r
}
func (q *QLocale) ToStringDoubleString(v float64, format byte, precision int) string {
	return qstrToGo(q.ToStringDouble(v, format, precision))
}

// ToStringFloat — Qt: QLocale::toString(float, char, int) const.
func (q *QLocale) ToStringFloat(v float32, format byte, precision int) *qtext.QString {
	r := wrapQStr(unsafe.Pointer(C.qlocale_to_string_float(
		q.handle(), C.float(v), C.char(format), C.int(precision))))
	runtime.KeepAlive(q)
	return r
}
func (q *QLocale) ToStringFloatString(v float32, format byte, precision int) string {
	return qstrToGo(q.ToStringFloat(v, format, precision))
}

// ---- locale-specific strings --------------------------------------------

// DecimalPoint — Qt: QLocale::decimalPoint() const.
func (q *QLocale) DecimalPoint() *qtext.QString {
	r := wrapQStr(unsafe.Pointer(C.qlocale_decimal_point(q.handle())))
	runtime.KeepAlive(q)
	return r
}
func (q *QLocale) DecimalPointString() string { return qstrToGo(q.DecimalPoint()) }

// GroupSeparator — Qt: QLocale::groupSeparator() const.
func (q *QLocale) GroupSeparator() *qtext.QString {
	r := wrapQStr(unsafe.Pointer(C.qlocale_group_separator(q.handle())))
	runtime.KeepAlive(q)
	return r
}
func (q *QLocale) GroupSeparatorString() string { return qstrToGo(q.GroupSeparator()) }

// Percent — Qt: QLocale::percent() const.
func (q *QLocale) Percent() *qtext.QString {
	r := wrapQStr(unsafe.Pointer(C.qlocale_percent(q.handle())))
	runtime.KeepAlive(q)
	return r
}
func (q *QLocale) PercentString() string { return qstrToGo(q.Percent()) }

// ZeroDigit — Qt: QLocale::zeroDigit() const.
func (q *QLocale) ZeroDigit() *qtext.QString {
	r := wrapQStr(unsafe.Pointer(C.qlocale_zero_digit(q.handle())))
	runtime.KeepAlive(q)
	return r
}
func (q *QLocale) ZeroDigitString() string { return qstrToGo(q.ZeroDigit()) }

// NegativeSign — Qt: QLocale::negativeSign() const.
func (q *QLocale) NegativeSign() *qtext.QString {
	r := wrapQStr(unsafe.Pointer(C.qlocale_negative_sign(q.handle())))
	runtime.KeepAlive(q)
	return r
}
func (q *QLocale) NegativeSignString() string { return qstrToGo(q.NegativeSign()) }

// PositiveSign — Qt: QLocale::positiveSign() const.
func (q *QLocale) PositiveSign() *qtext.QString {
	r := wrapQStr(unsafe.Pointer(C.qlocale_positive_sign(q.handle())))
	runtime.KeepAlive(q)
	return r
}
func (q *QLocale) PositiveSignString() string { return qstrToGo(q.PositiveSign()) }

// Exponential — Qt: QLocale::exponential() const.
func (q *QLocale) Exponential() *qtext.QString {
	r := wrapQStr(unsafe.Pointer(C.qlocale_exponential(q.handle())))
	runtime.KeepAlive(q)
	return r
}
func (q *QLocale) ExponentialString() string { return qstrToGo(q.Exponential()) }

// ---- month / day names ---------------------------------------------------

// MonthName — Qt: QLocale::monthName(int, FormatType) const.
func (q *QLocale) MonthName(month int, format FormatType) *qtext.QString {
	r := wrapQStr(unsafe.Pointer(C.qlocale_month_name(q.handle(), C.int(month), C.int(format))))
	runtime.KeepAlive(q)
	return r
}
func (q *QLocale) MonthNameString(month int, format FormatType) string {
	return qstrToGo(q.MonthName(month, format))
}

// StandaloneMonthName — Qt: QLocale::standaloneMonthName(int, FormatType) const.
func (q *QLocale) StandaloneMonthName(month int, format FormatType) *qtext.QString {
	r := wrapQStr(unsafe.Pointer(C.qlocale_standalone_month_name(
		q.handle(), C.int(month), C.int(format))))
	runtime.KeepAlive(q)
	return r
}
func (q *QLocale) StandaloneMonthNameString(month int, format FormatType) string {
	return qstrToGo(q.StandaloneMonthName(month, format))
}

// DayName — Qt: QLocale::dayName(int, FormatType) const.
func (q *QLocale) DayName(day int, format FormatType) *qtext.QString {
	r := wrapQStr(unsafe.Pointer(C.qlocale_day_name(q.handle(), C.int(day), C.int(format))))
	runtime.KeepAlive(q)
	return r
}
func (q *QLocale) DayNameString(day int, format FormatType) string {
	return qstrToGo(q.DayName(day, format))
}

// StandaloneDayName — Qt: QLocale::standaloneDayName(int, FormatType) const.
func (q *QLocale) StandaloneDayName(day int, format FormatType) *qtext.QString {
	r := wrapQStr(unsafe.Pointer(C.qlocale_standalone_day_name(
		q.handle(), C.int(day), C.int(format))))
	runtime.KeepAlive(q)
	return r
}
func (q *QLocale) StandaloneDayNameString(day int, format FormatType) string {
	return qstrToGo(q.StandaloneDayName(day, format))
}

// Weekdays — Qt: QLocale::weekdays() const (as []int, DayOfWeek values).
func (q *QLocale) Weekdays() []int {
	var count C.ptrdiff_t
	arr := C.qlocale_weekdays(q.handle(), &count)
	runtime.KeepAlive(q)
	if arr == nil {
		return nil
	}
	defer C.qlocale_free_int_array(arr)
	n := int(count)
	if n == 0 {
		return nil
	}
	src := unsafe.Slice((*C.int)(unsafe.Pointer(arr)), n)
	out := make([]int, n)
	for i := 0; i < n; i++ {
		out[i] = int(src[i])
	}
	return out
}

// ---- AM/PM ---------------------------------------------------------------

// AmText — Qt: QLocale::amText() const.
func (q *QLocale) AmText() *qtext.QString {
	r := wrapQStr(unsafe.Pointer(C.qlocale_am_text(q.handle())))
	runtime.KeepAlive(q)
	return r
}
func (q *QLocale) AmTextString() string { return qstrToGo(q.AmText()) }

// PmText — Qt: QLocale::pmText() const.
func (q *QLocale) PmText() *qtext.QString {
	r := wrapQStr(unsafe.Pointer(C.qlocale_pm_text(q.handle())))
	runtime.KeepAlive(q)
	return r
}
func (q *QLocale) PmTextString() string { return qstrToGo(q.PmText()) }

// ---- case conversion -----------------------------------------------------

// ToUpper — Qt: QLocale::toUpper(const QString&) const.
func (q *QLocale) ToUpper(s string) *qtext.QString {
	r := wrapQStr(unsafe.Pointer(C.qlocale_to_upper(q.handle(), cstr(s), C.ptrdiff_t(len(s)))))
	runtime.KeepAlive(q)
	runtime.KeepAlive(s)
	return r
}
func (q *QLocale) ToUpperString(s string) string { return qstrToGo(q.ToUpper(s)) }

// ToLower — Qt: QLocale::toLower(const QString&) const.
func (q *QLocale) ToLower(s string) *qtext.QString {
	r := wrapQStr(unsafe.Pointer(C.qlocale_to_lower(q.handle(), cstr(s), C.ptrdiff_t(len(s)))))
	runtime.KeepAlive(q)
	runtime.KeepAlive(s)
	return r
}
func (q *QLocale) ToLowerString(s string) string { return qstrToGo(q.ToLower(s)) }

// ---- currency ------------------------------------------------------------

// CurrencySymbol — Qt: QLocale::currencySymbol(CurrencySymbolFormat) const.
func (q *QLocale) CurrencySymbol(format CurrencySymbolFormat) *qtext.QString {
	r := wrapQStr(unsafe.Pointer(C.qlocale_currency_symbol(q.handle(), C.int(format))))
	runtime.KeepAlive(q)
	return r
}
func (q *QLocale) CurrencySymbolString(format CurrencySymbolFormat) string {
	return qstrToGo(q.CurrencySymbol(format))
}

// ToCurrencyStringLL — Qt: QLocale::toCurrencyString(qlonglong, const QString&) const.
func (q *QLocale) ToCurrencyStringLL(v int64, symbol *qtext.QString) *qtext.QString {
	var sp *C.QString
	if symbol != nil {
		if p := symbol.UnsafePtr(); p != nil {
			sp = (*C.QString)(p)
		}
	}
	r := wrapQStr(unsafe.Pointer(C.qlocale_to_currency_string_ll(q.handle(), C.int64_t(v), sp)))
	runtime.KeepAlive(q)
	runtime.KeepAlive(symbol)
	return r
}

// ToCurrencyStringULL — Qt: QLocale::toCurrencyString(qulonglong, const QString&) const.
func (q *QLocale) ToCurrencyStringULL(v uint64, symbol *qtext.QString) *qtext.QString {
	var sp *C.QString
	if symbol != nil {
		if p := symbol.UnsafePtr(); p != nil {
			sp = (*C.QString)(p)
		}
	}
	r := wrapQStr(unsafe.Pointer(C.qlocale_to_currency_string_ull(q.handle(), C.uint64_t(v), sp)))
	runtime.KeepAlive(q)
	runtime.KeepAlive(symbol)
	return r
}

// ToCurrencyStringDouble — Qt: QLocale::toCurrencyString(double, const QString&, int) const.
func (q *QLocale) ToCurrencyStringDouble(v float64, symbol *qtext.QString, precision int) *qtext.QString {
	var sp *C.QString
	if symbol != nil {
		if p := symbol.UnsafePtr(); p != nil {
			sp = (*C.QString)(p)
		}
	}
	r := wrapQStr(unsafe.Pointer(C.qlocale_to_currency_string_double(
		q.handle(), C.double(v), sp, C.int(precision))))
	runtime.KeepAlive(q)
	runtime.KeepAlive(symbol)
	return r
}

// ---- data size -----------------------------------------------------------

// FormattedDataSize — Qt: QLocale::formattedDataSize(qint64, int, DataSizeFormats) const.
func (q *QLocale) FormattedDataSize(bytes int64, precision int, format DataSizeFormats) *qtext.QString {
	r := wrapQStr(unsafe.Pointer(C.qlocale_formatted_data_size(
		q.handle(), C.int64_t(bytes), C.int(precision), C.int(format))))
	runtime.KeepAlive(q)
	return r
}
func (q *QLocale) FormattedDataSizeString(bytes int64, precision int, format DataSizeFormats) string {
	return qstrToGo(q.FormattedDataSize(bytes, precision, format))
}

// ---- quotation -----------------------------------------------------------

// QuoteString — Qt: QLocale::quoteString(QStringView, QuotationStyle) const.
func (q *QLocale) QuoteString(s string, style QuotationStyle) *qtext.QString {
	r := wrapQStr(unsafe.Pointer(C.qlocale_quote_string(
		q.handle(), cstr(s), C.ptrdiff_t(len(s)), C.int(style))))
	runtime.KeepAlive(q)
	runtime.KeepAlive(s)
	return r
}
func (q *QLocale) QuoteStringString(s string, style QuotationStyle) string {
	return qstrToGo(q.QuoteString(s, style))
}

// ---- separated list ------------------------------------------------------

// CreateSeparatedList — Qt: QLocale::createSeparatedList(const QStringList&) const.
func (q *QLocale) CreateSeparatedList(strings []string) *qtext.QString {
	n := len(strings)
	qs := make([]*qtext.QString, n)
	defer func() {
		for _, s := range qs {
			if s != nil {
				s.Close()
			}
		}
	}()
	for i, s := range strings {
		qs[i] = qtext.NewQStringFromString(s)
	}
	return q.CreateSeparatedListQStrings(qs)
}
func (q *QLocale) CreateSeparatedListString(strings []string) string {
	return qstrToGo(q.CreateSeparatedList(strings))
}

// CreateSeparatedListQStrings — QString variant of CreateSeparatedList.
func (q *QLocale) CreateSeparatedListQStrings(qstrs []*qtext.QString) *qtext.QString {
	n := len(qstrs)
	carr := make([]*C.QString, n)
	for i, s := range qstrs {
		if s != nil {
			if p := s.UnsafePtr(); p != nil {
				carr[i] = (*C.QString)(p)
			}
		}
	}
	var arrPtr **C.QString
	if n > 0 {
		arrPtr = (**C.QString)(unsafe.Pointer(&carr[0]))
	}
	r := wrapQStr(unsafe.Pointer(C.qlocale_create_separated_list(
		q.handle(), arrPtr, C.ptrdiff_t(n))))
	runtime.KeepAlive(q)
	runtime.KeepAlive(qstrs)
	runtime.KeepAlive(carr)
	return r
}

// ---- collation -----------------------------------------------------------

// Collation — Qt: QLocale::collation() const.
func (q *QLocale) Collation() *QLocale {
	r := wrapQLocale(C.qlocale_collation(q.handle()))
	runtime.KeepAlive(q)
	return r
}

// ---- uiLanguages ---------------------------------------------------------

// UILanguages — Qt: QLocale::uiLanguages(TagSeparator) const.
func (q *QLocale) UILanguages(sep TagSeparator) []*qtext.QString {
	var count C.ptrdiff_t
	arr := C.qlocale_ui_languages(q.handle(), C.int(sep), &count)
	runtime.KeepAlive(q)
	if arr == nil {
		return nil
	}
	defer C.qlocale_free_string_array(arr)
	n := int(count)
	if n == 0 {
		return nil
	}
	src := unsafe.Slice((**C.QString)(unsafe.Pointer(arr)), n)
	out := make([]*qtext.QString, n)
	for i := 0; i < n; i++ {
		out[i] = wrapQStr(unsafe.Pointer(src[i]))
	}
	return out
}

// UILanguagesStrings — Go-string variant of UILanguages.
func (q *QLocale) UILanguagesStrings(sep TagSeparator) []string {
	qs := q.UILanguages(sep)
	out := make([]string, len(qs))
	for i, s := range qs {
		out[i] = qstrToGo(s)
	}
	return out
}

// =================================================================
// Date/Time-aware API
// =================================================================

// ToStringDateFormat — Qt: QLocale::toString(QDate, FormatType) const.
func (q *QLocale) ToStringDateFormat(d QDate, format FormatType) *qtext.QString {
	r := wrapQStr(unsafe.Pointer(C.qlocale_to_string_date_fmt(
		q.handle(), C.int64_t(qdateToJD(d)), C.int(format))))
	runtime.KeepAlive(q)
	return r
}
func (q *QLocale) ToStringDateFormatString(d QDate, format FormatType) string {
	return qstrToGo(q.ToStringDateFormat(d, format))
}

// ToStringDateString — Qt: QLocale::toString(QDate, QStringView) const.
func (q *QLocale) ToStringDateString(d QDate, format string) *qtext.QString {
	r := wrapQStr(unsafe.Pointer(C.qlocale_to_string_date_str(
		q.handle(), C.int64_t(qdateToJD(d)), cstr(format), C.ptrdiff_t(len(format)))))
	runtime.KeepAlive(q)
	runtime.KeepAlive(format)
	return r
}
func (q *QLocale) ToStringDateStringString(d QDate, format string) string {
	return qstrToGo(q.ToStringDateString(d, format))
}

// ToStringDateFormatCal — Qt: QLocale::toString(QDate, FormatType, QCalendar) const.
func (q *QLocale) ToStringDateFormatCal(d QDate, format FormatType, cal *QCalendar) *qtext.QString {
	var cp C.go_qcalendar_t
	if cal != nil {
		cp = cal.handle()
	}
	r := wrapQStr(unsafe.Pointer(C.qlocale_to_string_date_fmt_cal(
		q.handle(), C.int64_t(qdateToJD(d)), C.int(format), cp)))
	runtime.KeepAlive(q)
	runtime.KeepAlive(cal)
	return r
}
func (q *QLocale) ToStringDateFormatCalString(d QDate, format FormatType, cal *QCalendar) string {
	return qstrToGo(q.ToStringDateFormatCal(d, format, cal))
}

// ToStringDateStringCal — Qt: QLocale::toString(QDate, QStringView, QCalendar) const.
func (q *QLocale) ToStringDateStringCal(d QDate, format string, cal *QCalendar) *qtext.QString {
	var cp C.go_qcalendar_t
	if cal != nil {
		cp = cal.handle()
	}
	r := wrapQStr(unsafe.Pointer(C.qlocale_to_string_date_str_cal(
		q.handle(), C.int64_t(qdateToJD(d)), cstr(format), C.ptrdiff_t(len(format)), cp)))
	runtime.KeepAlive(q)
	runtime.KeepAlive(cal)
	runtime.KeepAlive(format)
	return r
}
func (q *QLocale) ToStringDateStringCalString(d QDate, format string, cal *QCalendar) string {
	return qstrToGo(q.ToStringDateStringCal(d, format, cal))
}

// ToStringTimeFormat — Qt: QLocale::toString(QTime, FormatType) const.
func (q *QLocale) ToStringTimeFormat(t QTime, format FormatType) *qtext.QString {
	r := wrapQStr(unsafe.Pointer(C.qlocale_to_string_time_fmt(
		q.handle(), C.int32_t(qtimeToMSecs(t)), C.int(format))))
	runtime.KeepAlive(q)
	return r
}
func (q *QLocale) ToStringTimeFormatString(t QTime, format FormatType) string {
	return qstrToGo(q.ToStringTimeFormat(t, format))
}

// ToStringTimeString — Qt: QLocale::toString(QTime, QStringView) const.
func (q *QLocale) ToStringTimeString(t QTime, format string) *qtext.QString {
	r := wrapQStr(unsafe.Pointer(C.qlocale_to_string_time_str(
		q.handle(), C.int32_t(qtimeToMSecs(t)), cstr(format), C.ptrdiff_t(len(format)))))
	runtime.KeepAlive(q)
	runtime.KeepAlive(format)
	return r
}
func (q *QLocale) ToStringTimeStringString(t QTime, format string) string {
	return qstrToGo(q.ToStringTimeString(t, format))
}

// ToStringDateTimeFormat — Qt: QLocale::toString(const QDateTime&, FormatType) const.
func (q *QLocale) ToStringDateTimeFormat(dt *QDateTime, format FormatType) *qtext.QString {
	if dt == nil {
		return nil
	}
	dp := dt.handle()
	r := wrapQStr(unsafe.Pointer(C.qlocale_to_string_datetime_fmt(
		q.handle(), dp, C.int(format))))
	runtime.KeepAlive(q)
	runtime.KeepAlive(dt)
	return r
}
func (q *QLocale) ToStringDateTimeFormatString(dt *QDateTime, format FormatType) string {
	return qstrToGo(q.ToStringDateTimeFormat(dt, format))
}

// ToStringDateTimeString — Qt: QLocale::toString(const QDateTime&, QStringView) const.
func (q *QLocale) ToStringDateTimeString(dt *QDateTime, format string) *qtext.QString {
	if dt == nil {
		return nil
	}
	dp := dt.handle()
	r := wrapQStr(unsafe.Pointer(C.qlocale_to_string_datetime_str(
		q.handle(), dp, cstr(format), C.ptrdiff_t(len(format)))))
	runtime.KeepAlive(q)
	runtime.KeepAlive(dt)
	runtime.KeepAlive(format)
	return r
}
func (q *QLocale) ToStringDateTimeStringString(dt *QDateTime, format string) string {
	return qstrToGo(q.ToStringDateTimeString(dt, format))
}

// ToStringDateTimeFormatCal — Qt: QLocale::toString(const QDateTime&, FormatType, QCalendar) const.
func (q *QLocale) ToStringDateTimeFormatCal(dt *QDateTime, format FormatType, cal *QCalendar) *qtext.QString {
	if dt == nil {
		return nil
	}
	var cp C.go_qcalendar_t
	if cal != nil {
		cp = cal.handle()
	}
	dp := dt.handle()
	r := wrapQStr(unsafe.Pointer(C.qlocale_to_string_datetime_fmt_cal(
		q.handle(), dp, C.int(format), cp)))
	runtime.KeepAlive(q)
	runtime.KeepAlive(dt)
	runtime.KeepAlive(cal)
	return r
}
func (q *QLocale) ToStringDateTimeFormatCalString(dt *QDateTime, format FormatType, cal *QCalendar) string {
	return qstrToGo(q.ToStringDateTimeFormatCal(dt, format, cal))
}

// ToStringDateTimeStringCal — Qt: QLocale::toString(const QDateTime&, QStringView, QCalendar) const.
func (q *QLocale) ToStringDateTimeStringCal(dt *QDateTime, format string, cal *QCalendar) *qtext.QString {
	if dt == nil {
		return nil
	}
	var cp C.go_qcalendar_t
	if cal != nil {
		cp = cal.handle()
	}
	dp := dt.handle()
	r := wrapQStr(unsafe.Pointer(C.qlocale_to_string_datetime_str_cal(
		q.handle(), dp, cstr(format), C.ptrdiff_t(len(format)), cp)))
	runtime.KeepAlive(q)
	runtime.KeepAlive(dt)
	runtime.KeepAlive(cal)
	runtime.KeepAlive(format)
	return r
}
func (q *QLocale) ToStringDateTimeStringCalString(dt *QDateTime, format string, cal *QCalendar) string {
	return qstrToGo(q.ToStringDateTimeStringCal(dt, format, cal))
}

// DateFormat — Qt: QLocale::dateFormat(FormatType) const.
func (q *QLocale) DateFormat(format FormatType) *qtext.QString {
	r := wrapQStr(unsafe.Pointer(C.qlocale_date_format(q.handle(), C.int(format))))
	runtime.KeepAlive(q)
	return r
}
func (q *QLocale) DateFormatString(format FormatType) string { return qstrToGo(q.DateFormat(format)) }

// TimeFormat — Qt: QLocale::timeFormat(FormatType) const.
func (q *QLocale) TimeFormat(format FormatType) *qtext.QString {
	r := wrapQStr(unsafe.Pointer(C.qlocale_time_format(q.handle(), C.int(format))))
	runtime.KeepAlive(q)
	return r
}
func (q *QLocale) TimeFormatString(format FormatType) string { return qstrToGo(q.TimeFormat(format)) }

// DateTimeFormat — Qt: QLocale::dateTimeFormat(FormatType) const.
func (q *QLocale) DateTimeFormat(format FormatType) *qtext.QString {
	r := wrapQStr(unsafe.Pointer(C.qlocale_datetime_format(q.handle(), C.int(format))))
	runtime.KeepAlive(q)
	return r
}
func (q *QLocale) DateTimeFormatString(format FormatType) string {
	return qstrToGo(q.DateTimeFormat(format))
}

// ToDateFormat — Qt: QLocale::toDate(const QString&, FormatType, int baseYear) const.
// Returns null QDate on parse failure; check IsValid().
func (q *QLocale) ToDateFormat(s string, format FormatType, baseYear int) QDate {
	jd := int64(C.qlocale_to_date_fmt(q.handle(), cstr(s), C.ptrdiff_t(len(s)),
		C.int(format), C.int(baseYear)))
	runtime.KeepAlive(q)
	runtime.KeepAlive(s)
	return qdateFromJD(jd)
}

// ToDateString — Qt: QLocale::toDate(const QString&, const QString&, int baseYear) const.
func (q *QLocale) ToDateString(s string, format string, baseYear int) QDate {
	jd := int64(C.qlocale_to_date_str(q.handle(), cstr(s), C.ptrdiff_t(len(s)),
		cstr(format), C.ptrdiff_t(len(format)), C.int(baseYear)))
	runtime.KeepAlive(q)
	runtime.KeepAlive(s)
	runtime.KeepAlive(format)
	return qdateFromJD(jd)
}

// ToDateFormatCal — Qt: QLocale::toDate(const QString&, FormatType, QCalendar, int baseYear) const.
func (q *QLocale) ToDateFormatCal(s string, format FormatType, cal *QCalendar, baseYear int) QDate {
	var cp C.go_qcalendar_t
	if cal != nil {
		cp = cal.handle()
	}
	jd := int64(C.qlocale_to_date_fmt_cal(q.handle(), cstr(s), C.ptrdiff_t(len(s)),
		C.int(format), cp, C.int(baseYear)))
	runtime.KeepAlive(q)
	runtime.KeepAlive(s)
	runtime.KeepAlive(cal)
	return qdateFromJD(jd)
}

// ToDateStringCal — Qt: QLocale::toDate(const QString&, const QString&, QCalendar, int baseYear) const.
func (q *QLocale) ToDateStringCal(s string, format string, cal *QCalendar, baseYear int) QDate {
	var cp C.go_qcalendar_t
	if cal != nil {
		cp = cal.handle()
	}
	jd := int64(C.qlocale_to_date_str_cal(q.handle(), cstr(s), C.ptrdiff_t(len(s)),
		cstr(format), C.ptrdiff_t(len(format)), cp, C.int(baseYear)))
	runtime.KeepAlive(q)
	runtime.KeepAlive(s)
	runtime.KeepAlive(format)
	runtime.KeepAlive(cal)
	return qdateFromJD(jd)
}

// ToTimeFormat — Qt: QLocale::toTime(const QString&, FormatType) const.
// Returns null QTime on parse failure; check IsValid().
func (q *QLocale) ToTimeFormat(s string, format FormatType) QTime {
	ms := int32(C.qlocale_to_time_fmt(q.handle(), cstr(s), C.ptrdiff_t(len(s)), C.int(format)))
	runtime.KeepAlive(q)
	runtime.KeepAlive(s)
	return qtimeFromMSecs(ms)
}

// ToTimeString — Qt: QLocale::toTime(const QString&, const QString&) const.
func (q *QLocale) ToTimeString(s string, format string) QTime {
	ms := int32(C.qlocale_to_time_str(q.handle(), cstr(s), C.ptrdiff_t(len(s)),
		cstr(format), C.ptrdiff_t(len(format))))
	runtime.KeepAlive(q)
	runtime.KeepAlive(s)
	runtime.KeepAlive(format)
	return qtimeFromMSecs(ms)
}

// ToDateTimeFormat — Qt: QLocale::toDateTime(const QString&, FormatType, int baseYear) const.
// Returns a *QDateTime (possibly null/invalid); check IsValid().
func (q *QLocale) ToDateTimeFormat(s string, format FormatType, baseYear int) *QDateTime {
	p := C.qlocale_to_datetime_fmt(q.handle(), cstr(s), C.ptrdiff_t(len(s)),
		C.int(format), C.int(baseYear))
	runtime.KeepAlive(q)
	runtime.KeepAlive(s)
	return wrapQDateTime(p)
}

// ToDateTimeString — Qt: QLocale::toDateTime(const QString&, const QString&, int baseYear) const.
func (q *QLocale) ToDateTimeString(s string, format string, baseYear int) *QDateTime {
	p := C.qlocale_to_datetime_str(q.handle(), cstr(s), C.ptrdiff_t(len(s)),
		cstr(format), C.ptrdiff_t(len(format)), C.int(baseYear))
	runtime.KeepAlive(q)
	runtime.KeepAlive(s)
	runtime.KeepAlive(format)
	return wrapQDateTime(p)
}

// ToDateTimeFormatCal — Qt: QLocale::toDateTime(const QString&, FormatType, QCalendar, int baseYear) const.
func (q *QLocale) ToDateTimeFormatCal(s string, format FormatType, cal *QCalendar, baseYear int) *QDateTime {
	var cp C.go_qcalendar_t
	if cal != nil {
		cp = cal.handle()
	}
	p := C.qlocale_to_datetime_fmt_cal(q.handle(), cstr(s), C.ptrdiff_t(len(s)),
		C.int(format), cp, C.int(baseYear))
	runtime.KeepAlive(q)
	runtime.KeepAlive(s)
	runtime.KeepAlive(cal)
	return wrapQDateTime(p)
}

// ToDateTimeStringCal — Qt: QLocale::toDateTime(const QString&, const QString&, QCalendar, int baseYear) const.
func (q *QLocale) ToDateTimeStringCal(s string, format string, cal *QCalendar, baseYear int) *QDateTime {
	var cp C.go_qcalendar_t
	if cal != nil {
		cp = cal.handle()
	}
	p := C.qlocale_to_datetime_str_cal(q.handle(), cstr(s), C.ptrdiff_t(len(s)),
		cstr(format), C.ptrdiff_t(len(format)), cp, C.int(baseYear))
	runtime.KeepAlive(q)
	runtime.KeepAlive(s)
	runtime.KeepAlive(format)
	runtime.KeepAlive(cal)
	return wrapQDateTime(p)
}

// ---- static conversion methods ------------------------------------------

// LanguageToCode — Qt: static QLocale::languageToCode(Language, LanguageCodeTypes).
func LanguageToCode(lang Language, codeTypes LanguageCodeTypes) *qtext.QString {
	return wrapQStr(unsafe.Pointer(C.qlocale_language_to_code(C.int(lang), C.int(codeTypes))))
}
func LanguageToCodeString(lang Language, codeTypes LanguageCodeTypes) string {
	return qstrToGo(LanguageToCode(lang, codeTypes))
}

// CodeToLanguage — Qt: static QLocale::codeToLanguage(QStringView, LanguageCodeTypes).
func CodeToLanguage(code string, codeTypes LanguageCodeTypes) Language {
	r := Language(C.qlocale_code_to_language(cstr(code), C.ptrdiff_t(len(code)), C.int(codeTypes)))
	runtime.KeepAlive(code)
	return r
}

// TerritoryToCode — Qt: static QLocale::territoryToCode(Territory).
func TerritoryToCode(territory Territory) *qtext.QString {
	return wrapQStr(unsafe.Pointer(C.qlocale_territory_to_code(C.int(territory))))
}
func TerritoryToCodeString(territory Territory) string { return qstrToGo(TerritoryToCode(territory)) }

// CodeToTerritory — Qt: static QLocale::codeToTerritory(QStringView).
func CodeToTerritory(code string) Territory {
	r := Territory(C.qlocale_code_to_territory(cstr(code), C.ptrdiff_t(len(code))))
	runtime.KeepAlive(code)
	return r
}

// ScriptToCode — Qt: static QLocale::scriptToCode(Script).
func ScriptToCode(script Script) *qtext.QString {
	return wrapQStr(unsafe.Pointer(C.qlocale_script_to_code(C.int(script))))
}
func ScriptToCodeString(script Script) string { return qstrToGo(ScriptToCode(script)) }

// CodeToScript — Qt: static QLocale::codeToScript(QStringView).
func CodeToScript(code string) Script {
	r := Script(C.qlocale_code_to_script(cstr(code), C.ptrdiff_t(len(code))))
	runtime.KeepAlive(code)
	return r
}

// LanguageToString — Qt: static QLocale::languageToString(Language).
func LanguageToString(lang Language) *qtext.QString {
	return wrapQStr(unsafe.Pointer(C.qlocale_language_to_string(C.int(lang))))
}
func LanguageToStringString(lang Language) string { return qstrToGo(LanguageToString(lang)) }

// TerritoryToString — Qt: static QLocale::territoryToString(Territory).
func TerritoryToString(territory Territory) *qtext.QString {
	return wrapQStr(unsafe.Pointer(C.qlocale_territory_to_string(C.int(territory))))
}
func TerritoryToStringString(territory Territory) string {
	return qstrToGo(TerritoryToString(territory))
}

// ScriptToString — Qt: static QLocale::scriptToString(Script).
func ScriptToString(script Script) *qtext.QString {
	return wrapQStr(unsafe.Pointer(C.qlocale_script_to_string(C.int(script))))
}
func ScriptToStringString(script Script) string { return qstrToGo(ScriptToString(script)) }

// ---- static default / system --------------------------------------------

// SetDefaultLocale — Qt: static QLocale::setDefault(const QLocale&).
func SetDefaultLocale(locale *QLocale) {
	if p, ok := locale.tryHandle(); ok {
		C.qlocale_set_default(p)
	}
}

// CLocale — Qt: static QLocale::c().
func CLocale() *QLocale { return wrapQLocale(C.qlocale_c()) }

// SystemLocale — Qt: static QLocale::system().
func SystemLocale() *QLocale { return wrapQLocale(C.qlocale_system()) }

// ---- matching locales ---------------------------------------------------

// MatchingLocales — Qt: static QLocale::matchingLocales(Language, Script, Territory).
func MatchingLocales(lang Language, script Script, territory Territory) []*QLocale {
	var count C.ptrdiff_t
	arr := C.qlocale_matching_locales(C.int(lang), C.int(script), C.int(territory), &count)
	if arr == nil {
		return nil
	}
	defer C.qlocale_free_locale_array(arr)
	n := int(count)
	if n == 0 {
		return nil
	}
	src := unsafe.Slice((*C.go_qlocale_t)(unsafe.Pointer(arr)), n)
	out := make([]*QLocale, n)
	for i := 0; i < n; i++ {
		out[i] = wrapQLocale(src[i])
	}
	return out
}
