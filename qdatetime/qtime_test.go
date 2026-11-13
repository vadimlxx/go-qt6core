package qdatetime

import "testing"

func TestQTimeConstruction(t *testing.T) {
	tt := NewQTimeFromHMS(14, 30, 45, 123)
	if tt.Hour() != 14 || tt.Minute() != 30 || tt.Second() != 45 || tt.Msec() != 123 {
		t.Fatalf("got %d:%d:%d.%d", tt.Hour(), tt.Minute(), tt.Second(), tt.Msec())
	}
	if tt.IsNull() || !tt.IsValid() {
		t.Fatal("valid time reported as null or invalid")
	}
}

func TestQTimeNull(t *testing.T) {
	tt := NewQTime()
	if !tt.IsNull() || tt.IsValid() {
		t.Fatal("NewQTime() must be null and invalid")
	}
}

func TestQTimeZeroValue(t *testing.T) {
	var zero QTime // 0 = midnight
	if zero.IsNull() {
		t.Fatal("Go zero QTime must not be null")
	}
	if !zero.IsValid() {
		t.Fatal("Go zero QTime must be valid")
	}
	if zero.Hour() != 0 || zero.Minute() != 0 || zero.Second() != 0 {
		t.Fatalf("zero QTime = %d:%d:%d, want 00:00:00",
			zero.Hour(), zero.Minute(), zero.Second())
	}
}

func TestQTimeArithmetic(t *testing.T) {
	tt := NewQTimeFromHMS(23, 59, 59, 0)
	tt2 := tt.AddSecs(1)
	if tt2.Hour() != 0 || tt2.Minute() != 0 || tt2.Second() != 0 {
		t.Fatalf("23:59:59 + 1s = %d:%d:%d, want 00:00:00",
			tt2.Hour(), tt2.Minute(), tt2.Second())
	}
}

func TestQTimeSecsTo(t *testing.T) {
	a := NewQTimeFromHMS(10, 0, 0, 0)
	b := NewQTimeFromHMS(10, 0, 30, 0)
	if d := a.SecsTo(b); d != 30 {
		t.Fatalf("SecsTo = %d, want 30", d)
	}
}

func TestQTimeToString(t *testing.T) {
	tt := NewQTimeFromHMS(14, 30, 45, 0)
	s := tt.ToString(ISODate)
	defer s.Close()
	if s.GoString() != "14:30:45" {
		t.Fatalf("ISODate: %q", s.GoString())
	}
}

func TestQTimeToStringFormat(t *testing.T) {
	tt := NewQTimeFromHMS(14, 30, 45, 0)
	s := tt.ToStringFormat("HH-mm-ss")
	defer s.Close()
	if s.GoString() != "14-30-45" {
		t.Fatalf("HH-mm-ss: %q", s.GoString())
	}
}

func TestQTimeSetHMS(t *testing.T) {
	tt := NewQTimeFromHMS(10, 0, 0, 0)
	ok := tt.SetHMS(14, 30, 45, 123)
	if !ok {
		t.Fatal("SetHMS failed")
	}
	if tt.Hour() != 14 || tt.Minute() != 30 || tt.Second() != 45 || tt.Msec() != 123 {
		t.Fatalf("after SetHMS: %d:%d:%d.%d",
			tt.Hour(), tt.Minute(), tt.Second(), tt.Msec())
	}
}

func TestQTimeStatic(t *testing.T) {
	if !QTimeIsValid(0, 0, 0, 0) {
		t.Fatal("00:00:00 must be valid")
	}
	if !QTimeIsValid(23, 59, 59, 999) {
		t.Fatal("23:59:59.999 must be valid")
	}
	if QTimeIsValid(24, 0, 0, 0) {
		t.Fatal("24:00:00 must be invalid")
	}
	if QTimeIsValid(0, 0, 0, 1000) {
		t.Fatal("00:00:00.1000 must be invalid")
	}
}

func TestQTimeFromMSecs(t *testing.T) {
	tt := NewQTimeFromMSecsSinceStartOfDay(3661000) // 01:01:01
	if tt.Hour() != 1 || tt.Minute() != 1 || tt.Second() != 1 {
		t.Fatalf("got %d:%d:%d, want 01:01:01", tt.Hour(), tt.Minute(), tt.Second())
	}
}
