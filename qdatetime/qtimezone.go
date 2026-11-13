package qdatetime

/*
 QTimeZone — биндинг Qt 6.11 QtCore/QTimeZone как HANDLE.

 QTimeZone использует COW (Q_DECLARE_SHARED) и НЕ является trivially copyable.
 Внутри хранит ShortData (для UTC / LocalTime / OffsetFromUTC) или указатель на
 QTimeZonePrivate (для полноценных зон). В Go хранится указатель на C++ QTimeZone.
 Все геттеры/сеттеры идут напрямую в C++. Для независимой копии — Clone().

 Семантика значений vs ссылок:
   - QTimeZone — ссылочный тип (handle). Копирование Go-указателя не создаёт новую зону.
   - Для независимой копии — Clone().
   - Go-нулевое значение QTimeZone{} невалидно. Используйте NewQTimeZone().

 Управление памятью:
   - Все конструкторы и Clone создают НОВЫЙ heap-объект.
   - Close() освобождает немедленно; иначе срабатывает finalizer.

 Соглашения:
   - Конструкторы: NewQTimeZone (null), NewQTimeZoneLocal, NewQTimeZoneUTC,
     NewQTimeZoneFromOffset, NewQTimeZoneFromIANA, NewQTimeZoneSystem, NewQTimeZoneStaticUTC.
   - Геттеры: Id, TimeSpec, FixedSecondsAheadOfUtc, IsUTCOrFixedOffset, IsValid,
     OffsetFromUtc, StandardTimeOffset, DaylightTimeOffset, IsDaylightTime,
     HasDaylightTime, HasTransitions, Territory.
   - Строки: DisplayName (две перегрузки), Abbreviation, SystemTimeZoneId.
   - Переходы: NextTransition, PreviousTransition → *QDateTime + смещения + аббревиатура.
   - Конвертации: AsBackendZone, ToQDateTime (через QDateTime).
   - Статические: SystemTimeZoneId, IsTimeZoneIdAvailable, SystemTimeZone, UTC.

 Что не реализовано (см. qtimezone_doc.go):
   - OffsetDataList transitions() — требует QList<OffsetData>.
   - availableTimeZoneIds() — требует QList<QByteArray>.
   - ianaIdToWindowsId / windowsIdToIanaIds / windowsIdToDefaultIanaId.
   - hasAlternativeName().
   - fromCFTimeZone / toCFTimeZone, fromNSTimeZone / toNSTimeZone (Darwin).
   - std::chrono interop.
   - QDataStream / QDebug операторы.
   - Deprecated country() (заменён на Territory()).
*/

/*
#cgo pkg-config: Qt6Core
#cgo CXXFLAGS: -std=c++17 -fPIC -Wno-sfinae-incomplete
#cgo LDFLAGS: -lstdc++
#include "qtimezone_c.h"
*/
import "C"

import (
	"runtime"
	"unsafe"

	"github.com/vadimlxx/go-qt6core/qtext"
)

// QTimeZone wraps a heap-allocated Qt6 QTimeZone (COW handle).
type QTimeZone struct {
	p C.go_qtimezone_t
}

func wrapQTimeZone(p C.go_qtimezone_t) *QTimeZone {
	if p == nil {
		return nil
	}
	z := &QTimeZone{p: p}
	runtime.SetFinalizer(z, finalizeQTimeZone)
	return z
}

func finalizeQTimeZone(z *QTimeZone) {
	if z.p != nil {
		C.qtimezone_delete(z.p)
		z.p = nil
	}
}

// Close releases the underlying C++ object immediately and disables the finalizer.
func (z *QTimeZone) Close() {
	if z.p != nil {
		C.qtimezone_delete(z.p)
		z.p = nil
		runtime.SetFinalizer(z, nil)
	}
}

// ---- Construction --------------------------------------------------------

// NewQTimeZone returns a null QTimeZone (mirrors Qt QTimeZone()).
func NewQTimeZone() *QTimeZone { return wrapQTimeZone(C.qtimezone_new_default()) }

// NewQTimeZoneLocal returns QTimeZone(QTimeZone::LocalTime).
func NewQTimeZoneLocal() *QTimeZone { return wrapQTimeZone(C.qtimezone_new_local()) }

// NewQTimeZoneUTC returns QTimeZone(QTimeZone::UTC).
func NewQTimeZoneUTC() *QTimeZone { return wrapQTimeZone(C.qtimezone_new_utc()) }

// NewQTimeZoneFromOffset returns QTimeZone(offsetSeconds).
func NewQTimeZoneFromOffset(offsetSeconds int) *QTimeZone {
	return wrapQTimeZone(C.qtimezone_new_from_offset(C.int(offsetSeconds)))
}

// NewQTimeZoneFromIANA returns QTimeZone(ianaId).
func NewQTimeZoneFromIANA(id string) *QTimeZone {
	var p *C.char
	if len(id) > 0 {
		p = (*C.char)(unsafe.Pointer(unsafe.StringData(id)))
	}
	z := wrapQTimeZone(C.qtimezone_new_from_iana(p, C.ptrdiff_t(len(id))))
	runtime.KeepAlive(id)
	return z
}

// Clone returns a heap copy of z (independent C++ object).
func (z *QTimeZone) Clone() *QTimeZone { return wrapQTimeZone(C.qtimezone_clone(z.p)) }

// ---- Validity / identification ------------------------------------------

func (z *QTimeZone) IsValid() bool { return bool(C.qtimezone_is_valid(z.p)) }

// Id returns a fresh *qtext.QByteArray with the IANA zone id.
func (z *QTimeZone) Id() *qtext.QByteArray {
	var out unsafe.Pointer
	C.qtimezone_id(z.p, &out)
	return qtext.QByteArrayFromUnsafePtr(out)
}

// ---- TimeSpec / offset ---------------------------------------------------

func (z *QTimeZone) TimeSpec() TimeSpec { return TimeSpec(C.qtimezone_time_spec(z.p)) }
func (z *QTimeZone) FixedSecondsAheadOfUtc() int {
	return int(C.qtimezone_fixed_seconds_ahead_of_utc(z.p))
}
func (z *QTimeZone) IsUTCOrFixedOffset() bool {
	return bool(C.qtimezone_is_utc_or_fixed_offset(z.p))
}

// ---- Offset queries ------------------------------------------------------

func (z *QTimeZone) OffsetFromUtc(at *QDateTime) int {
	return int(C.qtimezone_offset_from_utc(z.p, at.p))
}
func (z *QTimeZone) StandardTimeOffset(at *QDateTime) int {
	return int(C.qtimezone_standard_time_offset(z.p, at.p))
}
func (z *QTimeZone) DaylightTimeOffset(at *QDateTime) int {
	return int(C.qtimezone_daylight_time_offset(z.p, at.p))
}
func (z *QTimeZone) IsDaylightTime(at *QDateTime) bool {
	return bool(C.qtimezone_is_daylight_time(z.p, at.p))
}
func (z *QTimeZone) HasDaylightTime() bool { return bool(C.qtimezone_has_daylight_time(z.p)) }
func (z *QTimeZone) HasTransitions() bool  { return bool(C.qtimezone_has_transitions(z.p)) }

// ---- Display names -------------------------------------------------------

// DisplayNameAt returns a fresh *qtext.QString with the zone's display name at the given time.
// nameType selects the name variant (DefaultName, LongName, ShortName, OffsetName).
// locale may be nil to use the default QLocale().
func (z *QTimeZone) DisplayNameAt(at *QDateTime, nameType NameType, locale *QLocale) *qtext.QString {
	var loc unsafe.Pointer
	if locale != nil {
		loc = unsafe.Pointer(locale.p)
	}
	var out unsafe.Pointer
	C.qtimezone_display_name_at(z.p, at.p, C.int(nameType), loc, &out)
	return qtext.QStringFromUnsafePtr(out)
}

// DisplayNameType returns a fresh *qtext.QString with the zone's display name for the given
// TimeType (StandardTime, DaylightTime, GenericTime).
func (z *QTimeZone) DisplayNameType(timeType TimeType, nameType NameType,
	locale *QLocale) *qtext.QString {
	var loc unsafe.Pointer
	if locale != nil {
		loc = unsafe.Pointer(locale.p)
	}
	var out unsafe.Pointer
	C.qtimezone_display_name_type(z.p, C.int(timeType), C.int(nameType), loc, &out)
	return qtext.QStringFromUnsafePtr(out)
}

// Abbreviation returns a fresh *qtext.QString with the zone's abbreviation at the given time.
func (z *QTimeZone) Abbreviation(at *QDateTime) *qtext.QString {
	var out unsafe.Pointer
	C.qtimezone_abbreviation(z.p, at.p, &out)
	return qtext.QStringFromUnsafePtr(out)
}

// ---- Transitions ---------------------------------------------------------

// NextTransition returns the first transition after the given datetime.
// Returns (offset, stdOffset, dstOffset, atUtc, abbreviation).
func (z *QTimeZone) NextTransition(after *QDateTime) (int64, int64, int64, *QDateTime, *qtext.QString) {
	var offset, stdOffset, dstOffset C.int64_t
	var atUtc, abbr unsafe.Pointer
	C.qtimezone_next_transition(z.p, after.p, &offset, &stdOffset, &dstOffset, &atUtc, &abbr)
	return int64(offset), int64(stdOffset), int64(dstOffset),
		wrapQDateTime((*C.go_qdatetime_t)(atUtc)), qtext.QStringFromUnsafePtr(abbr)
}

// PreviousTransition returns the last transition before the given datetime.
// Returns (offset, stdOffset, dstOffset, atUtc, abbreviation).
func (z *QTimeZone) PreviousTransition(before *QDateTime) (int64, int64, int64, *QDateTime, *qtext.QString) {
	var offset, stdOffset, dstOffset C.int64_t
	var atUtc, abbr unsafe.Pointer
	C.qtimezone_previous_transition(z.p, before.p, &offset, &stdOffset, &dstOffset, &atUtc, &abbr)
	return int64(offset), int64(stdOffset), int64(dstOffset),
		wrapQDateTime((*C.go_qdatetime_t)(atUtc)), qtext.QStringFromUnsafePtr(abbr)
}

// ---- Static --------------------------------------------------------------

// SystemTimeZoneId returns the system's IANA zone id as a fresh *qtext.QByteArray.
func SystemTimeZoneId() *qtext.QByteArray {
	var out unsafe.Pointer
	C.qtimezone_system_id(&out)
	return qtext.QByteArrayFromUnsafePtr(out)
}

// SystemTimeZone returns a new QTimeZone for the system zone.
func SystemTimeZone() *QTimeZone { return wrapQTimeZone(C.qtimezone_system_zone()) }

// QTimeZoneUTC returns a new QTimeZone for UTC.
func QTimeZoneUTC() *QTimeZone { return wrapQTimeZone(C.qtimezone_static_utc()) }

// IsTimeZoneIdAvailable reports whether the given IANA id is known to Qt.
func IsTimeZoneIdAvailable(id string) bool {
	var p *C.char
	if len(id) > 0 {
		p = (*C.char)(unsafe.Pointer(unsafe.StringData(id)))
	}
	ok := bool(C.qtimezone_is_time_zone_id_available(p, C.ptrdiff_t(len(id))))
	runtime.KeepAlive(id)
	return ok
}

// ---- Conversions ---------------------------------------------------------

// AsBackendZone returns a new QTimeZone that uses the full backend (not the ShortData form).
func (z *QTimeZone) AsBackendZone() *QTimeZone {
	return wrapQTimeZone(C.qtimezone_as_backend_zone(z.p))
}

// Territory returns the territory associated with the zone (or AnyTerritory).
func (z *QTimeZone) Territory() Country {
	var out C.int
	C.qtimezone_territory(z.p, &out)
	return Country(out)
}