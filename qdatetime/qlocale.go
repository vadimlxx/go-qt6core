package qdatetime

/*
 QLocale — биндинг Qt 6.11 QtCore/QLocale как HANDLE.

 QLocale — value-тип в Qt, использующий COW (Copy-on-Write) через QSharedDataPointer.
 В Go хранится указатель на C++ QLocale. Все геттеры идут напрямую в C++.
 Для получения независимой копии используйте Clone().

 Пакет qdatetime:

   В этом же пакете будут реализованы QDate, QTime, QDateTime. Это снимает проблему
   импортных циклов: методы QLocale::toString(QDate) / toDate() / toTime() /
   toDateTime() / QCalendar-aware варианты можно будет добавить напрямую, без
   промежуточных пакетов и без unsafe.Pointer-мостов.

 Семантика значений vs ссылок:
   - QLocale — ссылочный тип (handle). Копирование Go-указателя не создаёт новую локаль.
   - Для независимой копии — Clone().
   - Go-нулевое значение QLocale{} невалидно. Используйте NewQLocale / NewQLocaleFromName /
     NewQLocaleFromLanguageCountry.

 Управление памятью:
   - Конструкторы и Clone / CLocale / SystemLocale создают НОВЫЙ heap-объект.
   - Close() освобождает немедленно; иначе срабатывает finalizer.

 Зависимости:
   - qtext.QString — все строковые методы возвращают *qtext.QString.
   - QDate / QTime / QDateTime — будут в этом же пакете (пока не реализованы).
   - QStringList — для uiLanguages / createSeparatedList (пока не реализован в qtext).

 Соглашения:
   - Геттеры: Language(), Script(), Country(), MeasurementSystem(), TextDirection().
   - Строковые методы: Name(), Bcp47Name(), NativeLanguageName(), NativeTerritoryName(),
     DecimalPoint(), GroupSeparator(), Percent(), ZeroDigit(), NegativeSign(), PositiveSign(),
     Exponential(), AmText(), PmText() — все возвращают *qtext.QString.
   - Форматирование чисел: ToStringInt64(v), ToStringUInt64(v),
     ToStringDouble(v, format, precision) → *qtext.QString.
   - Парсинг чисел: ToInt64(s) (int64, bool), ToDouble(s) (float64, bool).
   - Регистр: ToUpper(s), ToLower(s) → *qtext.QString.
   - Имена месяцев/дней: MonthName(m, fmt), StandaloneMonthName(m, fmt),
     DayName(d, fmt), StandaloneDayName(d, fmt) → *qtext.QString.
   - Статические: SetDefaultLocale(locale), CLocale(), SystemLocale().

 Что не реализовано (см. qdatetime_doc.go):

   - QDate / QTime / QDateTime методы (toString, toDate, toTime, toDateTime)
     — будут добавлены вместе с соответствующими классами в этом пакете.
   - QCalendar-зависимые методы — вместе с классом QCalendar в этом пакете.
   - QList<Qt::DayOfWeek> weekdays() — нет контейнера.
   - QList<QLocale> matchingLocales() — нет контейнера.
   - QDataStream / QDebug операторы.
   - Статические languageToCode / codeToLanguage / territoryToCode / codeToTerritory /
     scriptToCode / codeToScript / languageToString / territoryToString / scriptToString.
   - formattedDataSize, toCurrencyString, currencySymbol.
   - quoteString, createSeparatedList.
   - NumberOptions и DataSizeFormats — типы определены, но не интегрированы в сигнатуры.
   - Protected-секций у QLocale нет — открывать нечего.
*/

/*
#cgo pkg-config: Qt6Core
#cgo CXXFLAGS: -std=c++17 -fPIC -Wno-sfinae-incomplete
#cgo LDFLAGS: -lstdc++
#include "qlocale_c.h"
*/
import "C"

import (
	"runtime"
	"unsafe"

	"github.com/vadimlxx/go-qt6core/qtext"
)

// QLocale wraps a heap-allocated Qt6 QLocale.
type QLocale struct {
	p C.go_qlocale_t
}

func wrapQLocale(p C.go_qlocale_t) *QLocale {
	if p == nil {
		return nil
	}
	q := &QLocale{p: p}
	runtime.SetFinalizer(q, finalizeQLocale)
	return q
}

func finalizeQLocale(q *QLocale) {
	if q.p != nil {
		C.qlocale_delete(q.p)
		q.p = nil
	}
}

// Close releases the underlying C++ object immediately and disables the finalizer.
func (q *QLocale) Close() {
	if q.p != nil {
		C.qlocale_delete(q.p)
		q.p = nil
		runtime.SetFinalizer(q, nil)
	}
}

// ---- Construction --------------------------------------------------------

// NewQLocale returns a handle to the default-constructed QLocale().
func NewQLocale() *QLocale { return wrapQLocale(C.qlocale_new_default()) }

// NewQLocaleFromName returns a handle to QLocale(name).
func NewQLocaleFromName(name string) *QLocale {
	var p *C.char
	if len(name) > 0 {
		p = (*C.char)(unsafe.Pointer(unsafe.StringData(name)))
	}
	q := wrapQLocale(C.qlocale_new_from_name(p, C.ptrdiff_t(len(name))))
	runtime.KeepAlive(name)
	return q
}

// NewQLocaleFromLanguageCountry returns QLocale(language, country).
func NewQLocaleFromLanguageCountry(lang Language, country Country) *QLocale {
	return wrapQLocale(C.qlocale_new_from_language_country(C.int(lang), C.int(country)))
}

// Clone returns a heap copy of q (independent C++ object).
func (q *QLocale) Clone() *QLocale { return wrapQLocale(C.qlocale_clone(q.p)) }

// ---- Accessors -----------------------------------------------------------

func (q *QLocale) Language() Language {
	return Language(C.qlocale_language(q.p))
}
func (q *QLocale) Script() Script {
	return Script(C.qlocale_script(q.p))
}
func (q *QLocale) Country() Country {
	return Country(C.qlocale_country(q.p))
}
func (q *QLocale) MeasurementSystem() MeasurementSystem {
	return MeasurementSystem(C.qlocale_measurement_system(q.p))
}
func (q *QLocale) TextDirection() int {
	return int(C.qlocale_text_direction(q.p))
}

// ---- String representation ----------------------------------------------

func (q *QLocale) Name() *qtext.QString {
	var out unsafe.Pointer
	C.qlocale_name(q.p, &out)
	return qtext.QStringFromUnsafePtr(out)
}
func (q *QLocale) Bcp47Name() *qtext.QString {
	var out unsafe.Pointer
	C.qlocale_bcp47_name(q.p, &out)
	return qtext.QStringFromUnsafePtr(out)
}
func (q *QLocale) NativeLanguageName() *qtext.QString {
	var out unsafe.Pointer
	C.qlocale_native_language_name(q.p, &out)
	return qtext.QStringFromUnsafePtr(out)
}
func (q *QLocale) NativeTerritoryName() *qtext.QString {
	var out unsafe.Pointer
	C.qlocale_native_territory_name(q.p, &out)
	return qtext.QStringFromUnsafePtr(out)
}

// ---- Number parsing ------------------------------------------------------

func (q *QLocale) ToInt64(s string) (int64, bool) {
	var p *C.char
	if len(s) > 0 {
		p = (*C.char)(unsafe.Pointer(unsafe.StringData(s)))
	}
	var out C.int64_t
	ok := bool(C.qlocale_to_int64(q.p, p, C.ptrdiff_t(len(s)), &out))
	runtime.KeepAlive(s)
	return int64(out), ok
}

func (q *QLocale) ToDouble(s string) (float64, bool) {
	var p *C.char
	if len(s) > 0 {
		p = (*C.char)(unsafe.Pointer(unsafe.StringData(s)))
	}
	var out C.double
	ok := bool(C.qlocale_to_double(q.p, p, C.ptrdiff_t(len(s)), &out))
	runtime.KeepAlive(s)
	return float64(out), ok
}

// ---- Number formatting ---------------------------------------------------

func (q *QLocale) ToStringInt64(value int64) *qtext.QString {
	var out unsafe.Pointer
	C.qlocale_to_string_int64(q.p, C.int64_t(value), &out)
	return qtext.QStringFromUnsafePtr(out)
}

func (q *QLocale) ToStringUInt64(value uint64) *qtext.QString {
	var out unsafe.Pointer
	C.qlocale_to_string_uint64(q.p, C.uint64_t(value), &out)
	return qtext.QStringFromUnsafePtr(out)
}

func (q *QLocale) ToStringDouble(value float64, format byte, precision int) *qtext.QString {
	var out unsafe.Pointer
	C.qlocale_to_string_double(q.p, C.double(value), C.char(format), C.int(precision), &out)
	return qtext.QStringFromUnsafePtr(out)
}

// ---- Locale-specific strings --------------------------------------------

func (q *QLocale) DecimalPoint() *qtext.QString {
	var out unsafe.Pointer
	C.qlocale_decimal_point(q.p, &out)
	return qtext.QStringFromUnsafePtr(out)
}
func (q *QLocale) GroupSeparator() *qtext.QString {
	var out unsafe.Pointer
	C.qlocale_group_separator(q.p, &out)
	return qtext.QStringFromUnsafePtr(out)
}
func (q *QLocale) Percent() *qtext.QString {
	var out unsafe.Pointer
	C.qlocale_percent(q.p, &out)
	return qtext.QStringFromUnsafePtr(out)
}
func (q *QLocale) ZeroDigit() *qtext.QString {
	var out unsafe.Pointer
	C.qlocale_zero_digit(q.p, &out)
	return qtext.QStringFromUnsafePtr(out)
}
func (q *QLocale) NegativeSign() *qtext.QString {
	var out unsafe.Pointer
	C.qlocale_negative_sign(q.p, &out)
	return qtext.QStringFromUnsafePtr(out)
}
func (q *QLocale) PositiveSign() *qtext.QString {
	var out unsafe.Pointer
	C.qlocale_positive_sign(q.p, &out)
	return qtext.QStringFromUnsafePtr(out)
}
func (q *QLocale) Exponential() *qtext.QString {
	var out unsafe.Pointer
	C.qlocale_exponential(q.p, &out)
	return qtext.QStringFromUnsafePtr(out)
}

// ---- Month / day names ---------------------------------------------------

func (q *QLocale) MonthName(month int, format FormatType) *qtext.QString {
	var out unsafe.Pointer
	C.qlocale_month_name(q.p, C.int(month), C.int(format), &out)
	return qtext.QStringFromUnsafePtr(out)
}

func (q *QLocale) StandaloneMonthName(month int, format FormatType) *qtext.QString {
	var out unsafe.Pointer
	C.qlocale_standalone_month_name(q.p, C.int(month), C.int(format), &out)
	return qtext.QStringFromUnsafePtr(out)
}

func (q *QLocale) DayName(day int, format FormatType) *qtext.QString {
	var out unsafe.Pointer
	C.qlocale_day_name(q.p, C.int(day), C.int(format), &out)
	return qtext.QStringFromUnsafePtr(out)
}

func (q *QLocale) StandaloneDayName(day int, format FormatType) *qtext.QString {
	var out unsafe.Pointer
	C.qlocale_standalone_day_name(q.p, C.int(day), C.int(format), &out)
	return qtext.QStringFromUnsafePtr(out)
}

// ---- AM/PM ---------------------------------------------------------------

func (q *QLocale) AmText() *qtext.QString {
	var out unsafe.Pointer
	C.qlocale_am_text(q.p, &out)
	return qtext.QStringFromUnsafePtr(out)
}
func (q *QLocale) PmText() *qtext.QString {
	var out unsafe.Pointer
	C.qlocale_pm_text(q.p, &out)
	return qtext.QStringFromUnsafePtr(out)
}

// ---- Case conversion -----------------------------------------------------

func (q *QLocale) ToUpper(s string) *qtext.QString {
	var p *C.char
	if len(s) > 0 {
		p = (*C.char)(unsafe.Pointer(unsafe.StringData(s)))
	}
	var out unsafe.Pointer
	C.qlocale_to_upper(q.p, p, C.ptrdiff_t(len(s)), &out)
	runtime.KeepAlive(s)
	return qtext.QStringFromUnsafePtr(out)
}

func (q *QLocale) ToLower(s string) *qtext.QString {
	var p *C.char
	if len(s) > 0 {
		p = (*C.char)(unsafe.Pointer(unsafe.StringData(s)))
	}
	var out unsafe.Pointer
	C.qlocale_to_lower(q.p, p, C.ptrdiff_t(len(s)), &out)
	runtime.KeepAlive(s)
	return qtext.QStringFromUnsafePtr(out)
}

// ---- Static methods ------------------------------------------------------

// SetDefaultLocale sets the global default locale. Mirrors QLocale::setDefault(locale).
func SetDefaultLocale(locale *QLocale) { C.qlocale_set_default(locale.p) }

// CLocale returns the C locale (QLocale::c()).
func CLocale() *QLocale { return wrapQLocale(C.qlocale_c()) }

// SystemLocale returns the system locale (QLocale::system()).
func SystemLocale() *QLocale { return wrapQLocale(C.qlocale_system()) }
