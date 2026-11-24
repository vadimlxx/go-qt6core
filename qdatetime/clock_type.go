package qdatetime

// ClockType mirrors QElapsedTimer::ClockType.
//
// The clock type reports which underlying clock source QElapsedTimer uses
// on the current platform.
//
// Qt: QElapsedTimer::ClockType.
type ClockType int

const (
    // SystemTime — QElapsedTimer::SystemTime.
    SystemTime ClockType = 0

    // MonotonicClock — QElapsedTimer::MonotonicClock.
    MonotonicClock ClockType = 1

    // TickCounter is deprecated in Qt and mapped internally to
    // PerformanceCounter on supported platforms.
    //
    // Deprecated: QElapsedTimer::TickCounter. Use PerformanceCounter.
    TickCounter ClockType = 2

    // MachAbsoluteTime — QElapsedTimer::MachAbsoluteTime (Darwin).
    MachAbsoluteTime ClockType = 3

    // PerformanceCounter — QElapsedTimer::PerformanceCounter.
    PerformanceCounter ClockType = 4
)