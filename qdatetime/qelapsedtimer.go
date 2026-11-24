package qdatetime

/*
	QElapsedTimer — биндинг Qt 6.11 QtCore/QElapsedTimer как VALUE-тип.

	# Семантика

	QElapsedTimer — value-тип в Qt: два qint64 поля (t1, t2), trivially
	copyable, без heap, без деструктора, без COW. В Go хранится как
	struct{t1, t2 int64} с идентичным layout. Присваивание Go-переменной
	создаёт независимую копию.

	# Нет Close / AddCleanup

	QElapsedTimer не владеет никакими ресурсами. Close() не нужен,
	runtime.AddCleanup не регистрируется.

	# Нулевое значение Go ≠ Qt-default

	ВАЖНО: Go-нулевое значение QElapsedTimer{t1:0, t2:0} НЕ равно
	Qt-конструктору QElapsedTimer() (t1 == t2 == INT64_MIN).

	- NewQElapsedTimer() → t1 == t2 == INT64_MIN → IsValid() == false
	- QElapsedTimer{}    → t1 == t2 == 0         → IsValid() == true,
							но таймер НЕ запущен: Elapsed() / NsecsElapsed()
							вернут бессмысленные значения (msecs от нуля
							reference-точки).

	Это сознательное следствие value-семантики Go (нулевое значение struct —
	это все нули). Та же ситуация у QDate{} / QTime{} в этом пакете.

	Правила:
	- всегда создавайте таймер через NewQElapsedTimer() и вызывайте Start()
		перед чтением; либо
	- вызывайте Invalidate() на существующем значении, чтобы привести его
		к Qt-невалидному состоянию.

	# Соглашения по именам

	- Конструктор: NewQElapsedTimer() — QElapsedTimer().
	- Статические: QElapsedTimerClockType() (ClockType) и
		QElapsedTimerIsMonotonic() — QElapsedTimer::clockType() /
		QElapsedTimer::isMonotonic(). Префикс "QElapsedTimer" нужен, чтобы
		имена не конфликтовали с одноимёнными методами других классов пакета.
	- Мутаторы: Start, Restart, Invalidate.
	- Геттеры: IsValid, NsecsElapsed, Elapsed, HasExpired,
		MSecsSinceReference, MSecsTo, SecsTo.
	- Сравнение: Equal, Less.

	# Знак MSecsTo / SecsTo

	Qt: this.MSecsTo(other) (this.SecsTo(other)) возвращает ОТРИЦАТЕЛЬНОЕ
	значение, если other запущен РАНЬШЕ this, и ПОЛОЖИТЕЛЬНОЕ — если other
	запущен ПОЗЖЕ this.

	- a.Start(); sleep; b.Start()
	- b.MSecsTo(a) < 0   (a запущен раньше b)
	- a.MSecsTo(b) > 0   (b запущен позже a)

	# Что не реализовано

	- Duration / durationElapsed() / durationTo() — std::chrono; покрыто
		через NsecsElapsed().
	- operator> / operator<= / operator>= — выводятся из Equal/Less:
		A >  B  ⇔  B.Less(A)
		A <= B  ⇔  !(B.Less(A))
		A >= B  ⇔  !(A.Less(B))
	- QDataStream / QDebug операторы.
	- qHash(QElapsedTimer).
*/

/*
#cgo pkg-config: Qt6Core
#cgo CXXFLAGS: -std=c++17 -fPIC
#cgo LDFLAGS: -lstdc++
#include <stdint.h>
#include <stdbool.h>
#include <stddef.h>
#include "qelapsedtimer_c.h"
*/
import "C"

import "unsafe"

// QElapsedTimer mirrors Qt6 QElapsedTimer as a value type.
//
// The Go zero value is QElapsedTimer{t1:0, t2:0}. This is NOT equal to Qt's
// default-constructed QElapsedTimer() (which has t1 == t2 == INT64_MIN and
// reports IsValid() == false). Because Qt's isValid() checks only
// t1 != INT64_MIN, the Go zero value reports IsValid() == true even though
// the timer was never started — reading Elapsed()/NsecsElapsed() on it
// yields meaningless values.
//
// The contract is therefore: always create the timer with
// NewQElapsedTimer() (Qt-invalid), then call Start() before reading, or
// call Invalidate() on an existing value to bring it to Qt's invalid state.
type QElapsedTimer struct {
    t1 int64
    t2 int64
}

// toC converts the Go value to its C representation field by field, so
// the cgo boundary needs no pointer cast for by-value calls.
func (t QElapsedTimer) toC() C.go_qelapsedtimer_t {
    return C.go_qelapsedtimer_t{
        t1: C.int64_t(t.t1),
        t2: C.int64_t(t.t2),
    }
}

// fromC converts a C struct into the Go value field by field.
func fromC(c C.go_qelapsedtimer_t) QElapsedTimer {
    return QElapsedTimer{t1: int64(c.t1), t2: int64(c.t2)}
}

// cptr returns a pointer to the value for by-pointer C calls.
//
// Safe because QElapsedTimer has the same memory layout as
// C.go_qelapsedtimer_t (two 8-byte fields). No runtime.KeepAlive is
// needed: QElapsedTimer is a value type and C does not hold any pointer
// to Go memory beyond the call duration.
func (t *QElapsedTimer) cptr() *C.go_qelapsedtimer_t {
    return (*C.go_qelapsedtimer_t)(unsafe.Pointer(t))
}

// ---- Construction --------------------------------------------------------

// NewQElapsedTimer returns a default-constructed QElapsedTimer().
//
// Qt: QElapsedTimer() — invalid (t1 == t2 == INT64_MIN); call Start() first.
func NewQElapsedTimer() QElapsedTimer {
    return fromC(C.qelapsedtimer_new())
}

// ---- Static --------------------------------------------------------------

// QElapsedTimerClockType returns the clock type used by QElapsedTimer on
// this platform.
//
// Qt: static QElapsedTimer::clockType().
func QElapsedTimerClockType() ClockType {
    return ClockType(C.qelapsedtimer_clock_type())
}

// QElapsedTimerIsMonotonic reports whether the platform's clock is monotonic.
//
// Qt: static QElapsedTimer::isMonotonic().
func QElapsedTimerIsMonotonic() bool {
    return bool(C.qelapsedtimer_is_monotonic())
}

// ---- Mutators ------------------------------------------------------------

// Start resets the timer and starts counting.
//
// Qt: QElapsedTimer::start().
func (t *QElapsedTimer) Start() {
    if t == nil {
        return
    }
    C.qelapsedtimer_start(t.cptr())
}

// Restart resets the timer and returns the elapsed milliseconds since the
// previous Start() (or since construction, if never started).
//
// Qt: QElapsedTimer::restart().
func (t *QElapsedTimer) Restart() int64 {
    if t == nil {
        return 0
    }
    return int64(C.qelapsedtimer_restart(t.cptr()))
}

// Invalidate marks the timer as invalid (equivalent to NewQElapsedTimer()).
//
// Qt: QElapsedTimer::invalidate().
func (t *QElapsedTimer) Invalidate() {
    if t == nil {
        return
    }
    C.qelapsedtimer_invalidate(t.cptr())
}

// ---- Queries -------------------------------------------------------------

// IsValid reports whether the timer has been started (and not invalidated).
//
// Qt: QElapsedTimer::isValid().
//
// Note: the Go zero value QElapsedTimer{} reports IsValid() == true because
// its t1 == 0 differs from Qt's invalid sentinel INT64_MIN — even though the
// timer was never started. See the type comment on QElapsedTimer.
func (t QElapsedTimer) IsValid() bool {
    return bool(C.qelapsedtimer_is_valid(t.toC()))
}

// NsecsElapsed returns the nanoseconds elapsed since the last Start().
//
// Qt: QElapsedTimer::nsecsElapsed().
func (t QElapsedTimer) NsecsElapsed() int64 {
    return int64(C.qelapsedtimer_nsecs_elapsed(t.toC()))
}

// Elapsed returns the milliseconds elapsed since the last Start().
//
// Qt: QElapsedTimer::elapsed().
func (t QElapsedTimer) Elapsed() int64 {
    return int64(C.qelapsedtimer_elapsed(t.toC()))
}

// HasExpired reports whether Elapsed() is at least timeout milliseconds.
//
// Qt: QElapsedTimer::hasExpired(qint64 timeout).
func (t QElapsedTimer) HasExpired(timeout int64) bool {
    return bool(C.qelapsedtimer_has_expired(t.toC(), C.int64_t(timeout)))
}

// MSecsSinceReference returns milliseconds since the platform-specific
// reference point (the same as Start() uses).
//
// Qt: QElapsedTimer::msecsSinceReference().
func (t QElapsedTimer) MSecsSinceReference() int64 {
    return int64(C.qelapsedtimer_msecs_since_reference(t.toC()))
}

// MSecsTo returns the milliseconds from t to other.
//
// Sign convention (matches Qt): the result is NEGATIVE if other was started
// BEFORE t, and POSITIVE if other was started AFTER t.
//
//   - a.Start(); ...; b.Start()   (a started first)
//   - b.MSecsTo(a) < 0            (a was started before b)
//   - a.MSecsTo(b) > 0            (b was started after a)
//
// Qt: QElapsedTimer::msecsTo(const QElapsedTimer& other).
func (t QElapsedTimer) MSecsTo(other QElapsedTimer) int64 {
    return int64(C.qelapsedtimer_msecs_to(t.toC(), other.toC()))
}

// SecsTo returns the seconds from t to other. Same sign convention as
// MSecsTo: negative if other started before t, positive if after.
//
// Qt: QElapsedTimer::secsTo(const QElapsedTimer& other).
func (t QElapsedTimer) SecsTo(other QElapsedTimer) int64 {
    return int64(C.qelapsedtimer_secs_to(t.toC(), other.toC()))
}

// ---- Comparison ----------------------------------------------------------

// Equal reports whether t and other hold the same (t1, t2) pair.
//
// Qt: QElapsedTimer::operator== (via comparesEqual).
func (t QElapsedTimer) Equal(other QElapsedTimer) bool {
    return bool(C.qelapsedtimer_equals(t.toC(), other.toC()))
}

// Less reports whether t < other, using the same ordering as Qt
// (Qt's compareThreeWay compares only t1).
//
// Qt: QElapsedTimer::operator< (uses t1 only).
func (t QElapsedTimer) Less(other QElapsedTimer) bool {
    return bool(C.qelapsedtimer_less(t.toC(), other.toC()))
}
