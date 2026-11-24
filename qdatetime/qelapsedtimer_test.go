package qdatetime

import (
	"testing"
	"time"
)

func TestNewQElapsedTimerInvalid(t *testing.T) {
	et := NewQElapsedTimer()
	if et.IsValid() {
		t.Fatal("NewQElapsedTimer() must be invalid")
	}
}

func TestQElapsedTimerZeroValueIsNotQtDefault(t *testing.T) {
	// Go zero value of QElapsedTimer{} is {t1:0, t2:0}. Qt's default-
	// constructed QElapsedTimer() is {INT64_MIN, INT64_MIN}. Qt's isValid()
	// checks only t1 != INT64_MIN, so:
	//   - NewQElapsedTimer().IsValid() == false  (Qt default)
	//   - QElapsedTimer{}.IsValid()      == true   (Go zero value)
	//
	// This is documented on the QElapsedTimer type. The test asserts the
	// actual contract, not a hypothetical "zero value is invalid" one.
	var et QElapsedTimer
	if et.Equal(NewQElapsedTimer()) {
		t.Fatal("Go zero value must differ from Qt default QElapsedTimer()")
	}

	// Go zero value reports IsValid()==true (t1==0 != INT64_MIN), but the
	// timer was never started. Invalidate() brings it to Qt's invalid state.
	et.Invalidate()
	if et.IsValid() {
		t.Fatal("after Invalidate(), IsValid must be false")
	}
	if !et.Equal(NewQElapsedTimer()) {
		t.Fatal("after Invalidate(), value must equal Qt default QElapsedTimer()")
	}
}

func TestQElapsedTimerStart(t *testing.T) {
	var et QElapsedTimer
	et.Start()
	if !et.IsValid() {
		t.Fatal("after Start(), IsValid must be true")
	}
	if et.Elapsed() < 0 {
		t.Fatalf("Elapsed() = %d, must be non-negative", et.Elapsed())
	}
	if et.NsecsElapsed() < 0 {
		t.Fatalf("NsecsElapsed() = %d, must be non-negative", et.NsecsElapsed())
	}
}

func TestQElapsedTimerRestart(t *testing.T) {
	var et QElapsedTimer
	et.Start()
	time.Sleep(10 * time.Millisecond)
	r := et.Restart()
	if r < 5 || r > 1000 {
		t.Fatalf("Restart() = %d ms, expected roughly 10", r)
	}
	if !et.IsValid() {
		t.Fatal("after Restart(), IsValid must remain true")
	}
}

func TestQElapsedTimerInvalidate(t *testing.T) {
	var et QElapsedTimer
	et.Start()
	if !et.IsValid() {
		t.Fatal("precondition: timer must be valid")
	}
	et.Invalidate()
	if et.IsValid() {
		t.Fatal("after Invalidate(), IsValid must be false")
	}
}

func TestQElapsedTimerHasExpired(t *testing.T) {
	var et QElapsedTimer
	et.Start()
	time.Sleep(10 * time.Millisecond)

	if !et.HasExpired(1) {
		t.Fatalf("HasExpired(1) = false, expected true (10ms elapsed)")
	}
	if et.HasExpired(60_000) {
		t.Fatalf("HasExpired(60000) = true, expected false")
	}
}

func TestQElapsedTimerMSecsSinceReference(t *testing.T) {
	var et QElapsedTimer
	et.Start()
	// msecsSinceReference uses the same reference point as Start(). After
	// Start() it should be a large number (time since system boot / epoch).
	// It can be 0 only if we happened to call it exactly at the reference
	// point, which is extremely unlikely.
	if et.MSecsSinceReference() == 0 {
		t.Log("MSecsSinceReference() = 0 (unusual but not impossible)")
	}
}

func TestQElapsedTimerMSecsTo(t *testing.T) {
	var a, b QElapsedTimer
	a.Start()
	time.Sleep(10 * time.Millisecond)
	b.Start()

	// Qt's msecsTo(other) returns NEGATIVE if `other` was started BEFORE
	// `this`, and POSITIVE if `other` was started AFTER `this`.
	//
	// Here: a started first, b started second. Therefore:
	//   b.MSecsTo(a) — other=a started before this=b  → negative
	//   a.MSecsTo(b) — other=b started after  this=a  → positive
	if got := b.MSecsTo(a); got >= 0 {
		t.Fatalf("b.MSecsTo(a) = %d, expected negative (a started before b)", got)
	}
	if got := a.MSecsTo(b); got <= 0 {
		t.Fatalf("a.MSecsTo(b) = %d, expected positive (b started after a)", got)
	}
}

func TestQElapsedTimerSecsTo(t *testing.T) {
	var a, b QElapsedTimer
	a.Start()
	time.Sleep(10 * time.Millisecond)
	b.Start()

	// Same sign convention as MSecsTo: negative when other is earlier,
	// positive when other is later. 10ms is < 1s, so magnitudes are 0; we
	// only check the sign.
	if got := b.SecsTo(a); got > 0 {
		t.Fatalf("b.SecsTo(a) = %d, expected non-positive (a started before b)", got)
	}
	if got := a.SecsTo(b); got < 0 {
		t.Fatalf("a.SecsTo(b) = %d, expected non-negative (b started after a)", got)
	}
}

func TestQElapsedTimerEqual(t *testing.T) {
	var a QElapsedTimer
	a.Start()
	// Both started ~now; their (t1, t2) may differ slightly, so compare
	// a timer to a copy of itself.
	c := a
	if !a.Equal(c) {
		t.Fatal("a.Equal(copy of a) must be true")
	}
}

func TestQElapsedTimerLess(t *testing.T) {
	var a, b QElapsedTimer
	a.Start()
	time.Sleep(2 * time.Millisecond)
	b.Start()

	// Qt's operator< compares only t1. a started earlier → a.t1 < b.t1.
	if !a.Less(b) {
		t.Fatal("a (started earlier) must be Less than b")
	}
	if b.Less(a) {
		t.Fatal("b (started later) must not be Less than a")
	}
	if a.Less(a) {
		t.Fatal("a must not be Less than itself")
	}
}

func TestQElapsedTimerClockType(t *testing.T) {
	ct := QElapsedTimerClockType()
	switch ct {
	case SystemTime, MonotonicClock, TickCounter, MachAbsoluteTime, PerformanceCounter:
		// valid enum value
	default:
		t.Fatalf("QElapsedTimerClockType() returned unknown value %d", int(ct))
	}
}

func TestQElapsedTimerIsMonotonic(t *testing.T) {
	// On any production platform the clock should be monotonic; we only
	// verify the call does not panic and returns a bool.
	_ = QElapsedTimerIsMonotonic()
}

func TestQElapsedTimerNilReceiver(t *testing.T) {
	var et *QElapsedTimer
	// These calls must not panic.
	et.Start()
	et.Invalidate()
	if got := et.Restart(); got != 0 {
		t.Fatalf("nil Restart() = %d, want 0", got)
	}
}

func TestQElapsedTimerCopyIndependence(t *testing.T) {
	var a QElapsedTimer
	a.Start()
	b := a // copy

	time.Sleep(5 * time.Millisecond)
	b.Start() // mutate the copy only

	// a and b should now differ (b started later); Qt's Less uses t1 only.
	if !a.Less(b) {
		t.Fatal("after copying and re-starting b, a must be Less than b")
	}
	if a.Equal(b) {
		t.Fatal("after copying and re-starting b, a must not Equal b")
	}
}
