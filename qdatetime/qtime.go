package qdatetime

/*
	QTime — биндинг Qt 6.11 QtCore/QTime как value-типа.

	QTime хранит миллисекунды от начала суток (int32): -1 — null, 0..86399999 —
	валидное время. Указателей, финализаторов и Close() нет; присваивание создаёт
	независимую копию.

	Важно:
	- NewQTime() возвращает null QTime (Qt QTime(), mds = -1);
	- нулевое значение QTime{} == 0 — валидное время 00:00:00.000.

	Реализовано:
	- конструкторы: NewQTime, NewQTimeFromHMS, NewQTimeCurrentTime,
		NewQTimeFromString, NewQTimeFromStringFormat,
		NewQTimeFromMSecsSinceStartOfDay, NewQTimeFromTime;
	- геттеры, арифметика, форматирование, SetHMS, QTimeIsValid;
	- ToTime() time.Time (дата-заглушка 1970-01-01 UTC).

	Не реализовано:
	- fromStdTimePoint/toStdTimePoint и прочий std::chrono interop;
	- QDataStream/qHash/QDebug-интеграция.

	Константы Qt::DateFormat берутся из qdata; собственных enum у QTime нет.
*/

/*
#cgo pkg-config: Qt6Core
#cgo CXXFLAGS: -std=c++17 -fPIC -Wno-sfinae-incomplete
#cgo LDFLAGS: -lstdc++
#include <stdint.h>
#include <stdbool.h>
#include <stddef.h>
#include "qdatetime_c.h"
*/
import "C"

import (
	"runtime"
	"time"
	"unsafe"

	"github.com/vadimlxx/go-qt6core/qdata"
	"github.com/vadimlxx/go-qt6core/qtext"
)

// QTime is a value type mirroring Qt6 QTime (int32 msecs since start of day,
// -1 means null).
//
// Qt: QTime.
type QTime int32

// Qt: QTime::QTime()
func NewQTime(arg ...any) QTime {
	if len(arg) == 0 || arg[0] == nil {
		// NewQTime() returns a null QTime.
		return QTime(C.qtime_default())
	}

	switch len(arg) {
	case 1:
		switch t := arg[0].(type) {
		case int:
			return NewQTimeFromMSecsSinceStartOfDay(t)
		case time.Time:
			return NewQTimeFromTime(t)
		}
	case 2:
		str, ok := arg[0].(string)
		if ok {
			switch format := arg[1].(type) {
			case qdata.DateFormat:
				return NewQTimeFromString(str, format)
			case string:
				return NewQTimeFromStringFormat(str, format)
			}
		}
	case 4:
		h,  ok1 := arg[0].(int)
		m,  ok2 := arg[1].(int)
		s,  ok3 := arg[2].(int)
		ms, ok4 := arg[3].(int)

		if ok1 && ok2 && ok3 && ok4 {
			return NewQTimeFromHMS(h, m, s, ms)
		}
	}

	panic("NewQTime invalid argument")
}

// NewQTimeFromHMS returns QTime(h, m, s, ms).
// Qt: QTime::QTime(int h, int m, int s, int ms = 0).
func NewQTimeFromHMS(h, m, s, ms int) QTime {
	return QTime(C.qtime_from_hms(C.int(h), C.int(m), C.int(s), C.int(ms)))
}

// NewQTimeCurrentTime returns the current time.
// Qt: QTime::currentTime().
func NewQTimeCurrentTime() QTime { return QTime(C.qtime_current_time()) }

// NewQTimeFromString parses str using the given Qt date format.
// Returns a null QTime on failure.
// Qt: QTime::fromString(const QString &string, Qt::DateFormat format).
func NewQTimeFromString(str string, format qdata.DateFormat) QTime {
	var p *C.char
	if len(str) > 0 {
		p = (*C.char)(unsafe.Pointer(unsafe.StringData(str)))
	}
	t := QTime(C.qtime_from_string_utf8(p, C.ptrdiff_t(len(str)), C.int(format)))
	runtime.KeepAlive(str)
	return t
}

// NewQTimeFromStringFormat parses str using the given format string.
// Returns a null QTime on failure.
// Qt: QTime::fromString(const QString &string, const QString &format).
func NewQTimeFromStringFormat(str, format string) QTime {
	var ps, pf *C.char
	if len(str) > 0 {
		ps = (*C.char)(unsafe.Pointer(unsafe.StringData(str)))
	}
	if len(format) > 0 {
		pf = (*C.char)(unsafe.Pointer(unsafe.StringData(format)))
	}
	t := QTime(C.qtime_from_string_format_utf8(
		ps, C.ptrdiff_t(len(str)),
		pf, C.ptrdiff_t(len(format)),
	))
	runtime.KeepAlive(str)
	runtime.KeepAlive(format)
	return t
}

// NewQTimeFromMSecsSinceStartOfDay returns QTime::fromMSecsSinceStartOfDay.
// Qt: QTime::fromMSecsSinceStartOfDay(int msecs).
func NewQTimeFromMSecsSinceStartOfDay(msecs int) QTime {
	return QTime(C.qtime_from_msecs_since_start_of_day(C.int(msecs)))
}

// NewQTimeFromTime returns the time-of-day part of t (millisecond precision).
// Qt: QTime(int h, int m, int s, int ms) with the fields taken from time.Time.
func NewQTimeFromTime(t time.Time) QTime {
	h, m, s := t.Clock()
	return NewQTimeFromHMS(h, m, s, t.Nanosecond()/1_000_000)
}

// IsNull reports whether the time is null.
// Qt: QTime::isNull() const.
func (t QTime) IsNull() bool { return bool(C.qtime_is_null(C.int32_t(t))) }

// IsValid reports whether the time is valid.
// Qt: QTime::isValid() const.
func (t QTime) IsValid() bool { return bool(C.qtime_is_valid(C.int32_t(t))) }

// Hour returns the hour (0..23).
// Qt: QTime::hour() const.
func (t QTime) Hour() int { return int(C.qtime_hour(C.int32_t(t))) }

// Minute returns the minute (0..59).
// Qt: QTime::minute() const.
func (t QTime) Minute() int { return int(C.qtime_minute(C.int32_t(t))) }

// Second returns the second (0..59).
// Qt: QTime::second() const.
func (t QTime) Second() int { return int(C.qtime_second(C.int32_t(t))) }

// Msec returns the millisecond (0..999).
// Qt: QTime::msec() const.
func (t QTime) Msec() int { return int(C.qtime_msec(C.int32_t(t))) }

// MsecsSinceStartOfDay returns the number of milliseconds since midnight.
// Qt: QTime::msecsSinceStartOfDay() const.
func (t QTime) MsecsSinceStartOfDay() int {
	return int(C.qtime_msecs_since_start_of_day(C.int32_t(t)))
}

// AddSecs returns the time shifted by secs seconds (wrapping around midnight).
// Qt: QTime::addSecs(int secs) const.
func (t QTime) AddSecs(secs int) QTime {
	return QTime(C.qtime_add_secs(C.int32_t(t), C.int(secs)))
}

// AddMSecs returns the time shifted by msecs milliseconds (wrapping around midnight).
// Qt: QTime::addMSecs(int msecs) const.
func (t QTime) AddMSecs(msecs int) QTime {
	return QTime(C.qtime_add_msecs(C.int32_t(t), C.int(msecs)))
}

// SecsTo returns the number of seconds from t to other.
// Qt: QTime::secsTo(QTime t) const.
func (t QTime) SecsTo(other QTime) int {
	return int(C.qtime_secs_to(C.int32_t(t), C.int32_t(other)))
}

// MsecsTo returns the number of milliseconds from t to other.
// Qt: QTime::msecsTo(QTime t) const.
func (t QTime) MsecsTo(other QTime) int {
	return int(C.qtime_msecs_to(C.int32_t(t), C.int32_t(other)))
}

// ToString formats the time using the given Qt date format.
// Qt: QTime::toString(Qt::DateFormat f) const.
func (t QTime) ToString(format qdata.DateFormat) *qtext.QString {
	var out unsafe.Pointer
	C.qtime_to_string(C.int32_t(t), C.int(format), &out)
	return qtext.QStringFromUnsafePtr(out)
}

// ToStringFormat formats the time using a format string.
// Qt: QTime::toString(const QString &format) const.
func (t QTime) ToStringFormat(format string) *qtext.QString {
	var p *C.char
	if len(format) > 0 {
		p = (*C.char)(unsafe.Pointer(unsafe.StringData(format)))
	}
	var out unsafe.Pointer
	C.qtime_to_string_format(C.int32_t(t), p, C.ptrdiff_t(len(format)), &out)
	runtime.KeepAlive(format)
	return qtext.QStringFromUnsafePtr(out)
}

// SetHMS changes the time; returns false for an invalid time.
// Qt: QTime::setHMS(int h, int m, int s, int ms = 0).
func (t *QTime) SetHMS(h, m, s, ms int) bool {
	return bool(C.qtime_set_hms((*C.go_qtime_t)(unsafe.Pointer(t)),
		C.int(h), C.int(m), C.int(s), C.int(ms)))
}

// QTimeIsValid mirrors QTime::isValid(h, m, s, ms).
// Qt: QTime::isValid(int h, int m, int s, int ms = 0) static.
func QTimeIsValid(h, m, s, ms int) bool {
	return bool(C.qtime_static_is_valid(C.int(h), C.int(m), C.int(s), C.int(ms)))
}

// ---- time.Time interop --------------------------------------------------

// ToTime returns the time-of-day on 1970-01-01 UTC (epoch date, millisecond
// precision). A null/invalid QTime yields the zero time.Time.
// Qt: QTime::hour()/minute()/second()/msec() combined into a Go time.Time (no
// direct Qt counterpart — QTime has no std::chrono conversion for time.Time).
func (t QTime) ToTime() time.Time {
	if !t.IsValid() {
		return time.Time{}
	}
	return time.Date(1970, 1, 1, t.Hour(), t.Minute(), t.Second(),
		t.Msec()*1_000_000, time.UTC)
}
