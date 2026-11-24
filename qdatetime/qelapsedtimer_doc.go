// QElapsedTimer is a Go binding for Qt 6.11 QtCore's QElapsedTimer.
//
// QElapsedTimer is a monotonic millisecond/nanosecond stopwatch. Typical
// usage is:
//
//	et := qdatetime.NewQElapsedTimer()
//	et.Start()
//	// ... work ...
//	if et.Elapsed() > 100 {
//	    // more than 100 ms elapsed
//	}
//
// # Type semantics
//
// QElapsedTimer is a trivially-copyable value type in Qt: two qint64 fields
// (t1, t2) with no heap allocation and no destructor. In Go it is stored as
// struct{ t1, t2 int64 } with the same memory layout. Assignment and function
// passing copy the value; there are no pointers, no finalizers, and no Close().
//
// # Go zero value vs Qt default
//
// The Go zero value QElapsedTimer{} (t1 == t2 == 0) is NOT equal to Qt's
// default-constructed QElapsedTimer() (t1 == t2 == INT64_MIN).
//
//	NewQElapsedTimer()   → t1 == t2 == INT64_MIN   → IsValid() == false
//	QElapsedTimer{}      → t1 == t2 == 0           → IsValid() == true,
//	                        but the timer was NEVER started: Elapsed() and
//	                        NsecsElapsed() return meaningless values
//	                        (ms since the reference point, not since Start).
//
// This is a consequence of Go's zero-value semantics for structs: it cannot
// be made equal to Qt's sentinel INT64_MIN. The same situation exists for
// QDate{} (valid date 4713 BCE) and QTime{} (valid midnight) in this package.
//
// Contract:
//   - Always create the timer via NewQElapsedTimer(), then call Start()
//     before reading.
//   - Or call Invalidate() on an existing value to bring it to Qt's invalid
//     state.
//   - Never read Elapsed()/NsecsElapsed() on a value that was not Start()ed.
//
// # Lifecycle
//
// Because QElapsedTimer owns no resources, there is no Close() and no
// runtime.AddCleanup. Unlike QString / QDateTime / QLocale, no cleanup is
// registered.
//
// # Construction
//
//   - NewQElapsedTimer()   — QElapsedTimer() (invalid; call Start first)
//
// # Static methods
//
//   - QElapsedTimerClockType() ClockType  — static QElapsedTimer::clockType()
//   - QElapsedTimerIsMonotonic() bool     — static QElapsedTimer::isMonotonic()
//
// # Mutators
//
//   - Start()           — QElapsedTimer::start()
//   - Restart() int64   — QElapsedTimer::restart()
//   - Invalidate()      — QElapsedTimer::invalidate()
//
// # Queries
//
//   - IsValid() bool              — QElapsedTimer::isValid()
//   - NsecsElapsed() int64        — QElapsedTimer::nsecsElapsed()
//   - Elapsed() int64             — QElapsedTimer::elapsed()
//   - HasExpired(timeout int64)   — QElapsedTimer::hasExpired(qint64)
//   - MSecsSinceReference() int64 — QElapsedTimer::msecsSinceReference()
//   - MSecsTo(other) int64        — QElapsedTimer::msecsTo(const QElapsedTimer&)
//   - SecsTo(other) int64         — QElapsedTimer::secsTo(const QElapsedTimer&)
//
// # Sign convention of MSecsTo / SecsTo
//
// Matches Qt: this.MSecsTo(other) (this.SecsTo(other)) returns NEGATIVE if
// other was started BEFORE this, and POSITIVE if other was started AFTER this.
//
//   - a.Start(); sleep; b.Start()
//   - b.MSecsTo(a) < 0    (a was started before b)
//   - a.MSecsTo(b) > 0    (b was started after a)
//
// # Comparison
//
//   - Equal(other) bool — QElapsedTimer::operator==
//   - Less(other) bool  — QElapsedTimer::operator<
//
// # Enum types
//
// ClockType is exposed as a Go type in clock_type.go:
//
//	type ClockType int
//	const (
//	    SystemTime         ClockType = 0
//	    MonotonicClock     ClockType = 1
//	    TickCounter        ClockType = 2  // deprecated; maps to PerformanceCounter
//	    MachAbsoluteTime   ClockType = 3
//	    PerformanceCounter ClockType = 4
//	)
//
// # Not yet implemented
//
//   - Duration / durationElapsed() / durationTo() — std::chrono types; covered
//     by NsecsElapsed().
//   - operator> / operator<= / operator>= — derive from Equal/Less:
//       A >  B  ⇔  B.Less(A)
//       A <= B  ⇔  !(B.Less(A))
//       A >= B  ⇔  !(A.Less(B))
//   - QDataStream / QDebug operators.
//   - qHash(QElapsedTimer).
//
// # Requirements
//
// Qt >= 6.11, C++17, pkg-config Qt6Core.
//
// This package is a subpackage of github.com/vadimlxx/go-qt6core.
package qdatetime
