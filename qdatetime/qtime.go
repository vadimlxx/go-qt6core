package qdatetime

/*
 QTime — биндинг Qt 6.11 QtCore/QTime как VALUE-тип.

 QTime хранит миллисекунды с начала дня (int32, 4 байта), trivially copyable, без COW.
 В Go хранится как int32. Присваивание создаёт независимую копию.
 Указателей, финализаторов и Close() нет.

 ВАЖНО:
   - NewQTime() возвращает null-время (Qt QTime() → mds = -1).
   - Go-нулевое значение QTime{} = 0 — это полночь (00:00:00.000),
     отличная от null. Для null используйте NewQTime(); для «нулевого» времени — QTime{}.
   - Внутреннее представление: -1 = null, 0..86399999 = msecs since start of day.

 Соглашения:
   - Конструкторы: NewQTime, NewQTimeFromHMS, NewQTimeCurrentTime, NewQTimeFromString,
     NewQTimeFromStringFormat, NewQTimeFromMSecsSinceStartOfDay.
   - Геттеры: Hour, Minute, Second, Msec, MsecsSinceStartOfDay, IsNull, IsValid.
   - Арифметика: AddSecs, AddMSecs, SecsTo, MsecsTo.
   - Строки: ToString, ToStringFormat → *qtext.QString.
   - Сеттеры: SetHMS.

 Что не реализовано (см. qdatetime_doc.go):
   - QDataStream операторы, qHash.
*/

/*
#cgo pkg-config: Qt6Core
#cgo CXXFLAGS: -std=c++17 -fPIC -Wno-sfinae-incomplete
#cgo LDFLAGS: -lstdc++
#include "qtime_c.h"
*/
import "C"

import (
	"runtime"
	"unsafe"

	"github.com/vadimlxx/go-qt6core/qtext"
)

// QTime is a value type mirroring Qt6 QTime (int32 msecs; -1 = null).
type QTime int32

// NewQTime returns a null QTime (mirrors Qt QTime()).
func NewQTime() QTime { return QTime(C.qtime_default()) }

// NewQTimeFromHMS returns QTime(h, m, s, ms).
func NewQTimeFromHMS(h, m, s, ms int) QTime {
	return QTime(C.qtime_from_hms(C.int(h), C.int(m), C.int(s), C.int(ms)))
}

// NewQTimeCurrentTime returns QTime::currentTime().
func NewQTimeCurrentTime() QTime { return QTime(C.qtime_current_time()) }

// NewQTimeFromString returns QTime::fromString(str, format).
func NewQTimeFromString(str string, format DateFormat) QTime {
	var p *C.char
	if len(str) > 0 {
		p = (*C.char)(unsafe.Pointer(unsafe.StringData(str)))
	}
	t := QTime(C.qtime_from_string_utf8(p, C.ptrdiff_t(len(str)), C.int(format)))
	runtime.KeepAlive(str)
	return t
}

// NewQTimeFromStringFormat returns QTime::fromString(str, formatQString).
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

// NewQTimeFromMSecsSinceStartOfDay returns QTime::fromMSecsSinceStartOfDay(msecs).
// Negative msecs produces a null QTime.
func NewQTimeFromMSecsSinceStartOfDay(msecs int) QTime {
	return QTime(C.qtime_from_msecs_since_start_of_day(C.int(msecs)))
}

// ---- Accessors -----------------------------------------------------------

func (t QTime) IsNull() bool  { return bool(C.qtime_is_null(C.int32_t(t))) }
func (t QTime) IsValid() bool { return bool(C.qtime_is_valid(C.int32_t(t))) }
func (t QTime) Hour() int     { return int(C.qtime_hour(C.int32_t(t))) }
func (t QTime) Minute() int   { return int(C.qtime_minute(C.int32_t(t))) }
func (t QTime) Second() int   { return int(C.qtime_second(C.int32_t(t))) }
func (t QTime) Msec() int     { return int(C.qtime_msec(C.int32_t(t))) }
func (t QTime) MsecsSinceStartOfDay() int {
	return int(C.qtime_msecs_since_start_of_day(C.int32_t(t)))
}

// ---- Arithmetic ----------------------------------------------------------

func (t QTime) AddSecs(secs int) QTime {
	return QTime(C.qtime_add_secs(C.int32_t(t), C.int(secs)))
}
func (t QTime) AddMSecs(msecs int) QTime {
	return QTime(C.qtime_add_msecs(C.int32_t(t), C.int(msecs)))
}
func (t QTime) SecsTo(other QTime) int {
	return int(C.qtime_secs_to(C.int32_t(t), C.int32_t(other)))
}
func (t QTime) MsecsTo(other QTime) int {
	return int(C.qtime_msecs_to(C.int32_t(t), C.int32_t(other)))
}

// ---- Strings -------------------------------------------------------------

func (t QTime) ToString(format DateFormat) *qtext.QString {
	var out unsafe.Pointer
	C.qtime_to_string(C.int32_t(t), C.int(format), &out)
	return qtext.QStringFromUnsafePtr(out)
}
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

// ---- Setter --------------------------------------------------------------

func (t *QTime) SetHMS(h, m, s, ms int) bool {
	return bool(C.qtime_set_hms((*C.go_qtime_t)(unsafe.Pointer(t)),
		C.int(h), C.int(m), C.int(s), C.int(ms)))
}

// ---- Static --------------------------------------------------------------

// QTimeIsValid mirrors QTime::isValid(h, m, s, ms).
func QTimeIsValid(h, m, s, ms int) bool {
	return bool(C.qtime_static_is_valid(C.int(h), C.int(m), C.int(s), C.int(ms)))
}
