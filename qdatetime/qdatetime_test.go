package qdatetime

import "testing"

func TestQDateTimeConstruction(t *testing.T) {
	d := NewQDateTimeFromDateTime(
		NewQDateFromYMD(2024, 1, 15),
		NewQTimeFromHMS(12, 0, 0, 0),
	)
	defer d.Close()
	if !d.IsValid() || d.IsNull() {
		t.Fatal("valid datetime reported as invalid or null")
	}
	if d.Date().Year() != 2024 || d.Date().Month() != 1 || d.Date().Day() != 15 {
		t.Fatalf("date mismatch: %d-%d-%d",
			d.Date().Year(), d.Date().Month(), d.Date().Day())
	}
	if d.Time().Hour() != 12 || d.Time().Minute() != 0 {
		t.Fatalf("time mismatch: %d:%d", d.Time().Hour(), d.Time().Minute())
	}
}

func TestQDateTimeNull(t *testing.T) {
	d := NewQDateTime()
	defer d.Close()
	if !d.IsNull() || d.IsValid() {
		t.Fatal("NewQDateTime() must be null and invalid")
	}
}

func TestQDateTimeArithmetic(t *testing.T) {
	d := NewQDateTimeFromDateTime(
		NewQDateFromYMD(2024, 1, 31),
		NewQTimeFromHMS(12, 0, 0, 0),
	)
	defer d.Close()

	d2 := d.AddMonths(1)
	defer d2.Close()
	if d2.Date().Month() != 2 || d2.Date().Day() != 29 {
		t.Fatalf("Jan 31 + 1 month = %d-%d, want 2024-02-29",
			d2.Date().Month(), d2.Date().Day())
	}
}

func TestQDateTimeAddSecs(t *testing.T) {
	d := NewQDateTimeFromDateTime(
		NewQDateFromYMD(2024, 1, 1),
		NewQTimeFromHMS(0, 0, 0, 0),
	)
	defer d.Close()

	d2 := d.AddSecs(3600)
	defer d2.Close()
	if d2.Time().Hour() != 1 || d2.Time().Minute() != 0 {
		t.Fatalf("+1h = %d:%d, want 01:00", d2.Time().Hour(), d2.Time().Minute())
	}
}

func TestQDateTimeMsecsRoundTrip(t *testing.T) {
	// 2024-01-15T12:30:00Z
	const ms int64 = 1705321800000

	d := NewQDateTimeFromMSecsSinceEpoch(ms)
	defer d.Close()
	if d.ToMSecsSinceEpoch() != ms {
		t.Fatalf("msecs round-trip: %d != %d", d.ToMSecsSinceEpoch(), ms)
	}
}

func TestQDateTimeToString(t *testing.T) {
	d := NewQDateTimeFromDateTime(
		NewQDateFromYMD(2024, 1, 15),
		NewQTimeFromHMS(12, 30, 0, 0),
	)
	defer d.Close()

	s := d.ToString(ISODate)
	defer s.Close()
	if s.GoString() != "2024-01-15T12:30:00" {
		t.Fatalf("ISODate: %q", s.GoString())
	}
}

func TestQDateTimeToStringFormat(t *testing.T) {
	d := NewQDateTimeFromDateTime(
		NewQDateFromYMD(2024, 1, 15),
		NewQTimeFromHMS(12, 30, 45, 0),
	)
	defer d.Close()

	s := d.ToStringFormat("yyyy-MM-dd HH:mm:ss")
	defer s.Close()
	if s.GoString() != "2024-01-15 12:30:45" {
		t.Fatalf("custom format: %q", s.GoString())
	}
}

func TestQDateTimeCurrent(t *testing.T) {
	d := NewQDateTimeCurrentDateTime()
	defer d.Close()
	if d.IsNull() || !d.IsValid() {
		t.Fatal("currentDateTime must be valid and not null")
	}
	if d.ToMSecsSinceEpoch() <= 0 {
		t.Fatal("ToMSecsSinceEpoch should be positive")
	}
}

func TestQDateTimeEquals(t *testing.T) {
	d1 := NewQDateTimeFromDateTime(
		NewQDateFromYMD(2024, 1, 15),
		NewQTimeFromHMS(12, 0, 0, 0),
	)
	defer d1.Close()
	d2 := NewQDateTimeFromDateTime(
		NewQDateFromYMD(2024, 1, 15),
		NewQTimeFromHMS(12, 0, 0, 0),
	)
	defer d2.Close()

	if !d1.Equals(d2) {
		t.Fatal("equal datetimes must compare as equal")
	}
}

func TestQDateTimeSetDate(t *testing.T) {
	d := NewQDateTimeFromDateTime(
		NewQDateFromYMD(2024, 1, 15),
		NewQTimeFromHMS(12, 0, 0, 0),
	)
	defer d.Close()

	d.SetDate(NewQDateFromYMD(2025, 6, 20))
	if d.Date().Year() != 2025 || d.Date().Month() != 6 {
		t.Fatalf("after SetDate: %d-%d",
			d.Date().Year(), d.Date().Month())
	}
	// Time unchanged.
	if d.Time().Hour() != 12 {
		t.Fatalf("time should remain 12:00, got %d", d.Time().Hour())
	}
}
