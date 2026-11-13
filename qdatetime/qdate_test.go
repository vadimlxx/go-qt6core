package qdatetime

import "testing"

func TestQDateConstruction(t *testing.T) {
	d := NewQDateFromYMD(2024, 1, 15)
	if d.Year() != 2024 || d.Month() != 1 || d.Day() != 15 {
		t.Fatalf("got %d-%d-%d", d.Year(), d.Month(), d.Day())
	}
	if d.IsNull() || !d.IsValid() {
		t.Fatal("valid date reported as null or invalid")
	}
}

func TestQDateNull(t *testing.T) {
	d := NewQDate()
	if !d.IsNull() || d.IsValid() {
		t.Fatal("NewQDate() must be null and invalid")
	}
}

func TestQDateZeroValue(t *testing.T) {
	var zero QDate // 0 = Julian day 0 = 1 January 4713 BCE
	if zero.IsNull() {
		t.Fatal("Go zero QDate must not be null")
	}
	if !zero.IsValid() {
		t.Fatal("Go zero QDate must be valid")
	}
}

func TestQDateArithmetic(t *testing.T) {
	d := NewQDateFromYMD(2024, 1, 31)
	d2 := d.AddMonths(1)
	if d2.Month() != 2 || d2.Day() != 29 {
		t.Fatalf("Jan 31 + 1 month = %d-%d, want 2024-02-29", d2.Month(), d2.Day())
	}
}

func TestQDateAddDays(t *testing.T) {
	d := NewQDateFromYMD(2024, 1, 1)
	d2 := d.AddDays(100)
	if d2.Year() != 2024 || d2.Month() != 4 || d2.Day() != 10 {
		t.Fatalf("2024-01-01 + 100d = %d-%d-%d, want 2024-04-10",
			d2.Year(), d2.Month(), d2.Day())
	}
}

func TestQDateToString(t *testing.T) {
	d := NewQDateFromYMD(2024, 1, 15)
	s := d.ToString(ISODate)
	defer s.Close()
	if s.GoString() != "2024-01-15" {
		t.Fatalf("ISODate: %q", s.GoString())
	}
}

func TestQDateToStringFormat(t *testing.T) {
	d := NewQDateFromYMD(2024, 1, 15)
	s := d.ToStringFormat("dd.MM.yyyy")
	defer s.Close()
	if s.GoString() != "15.01.2024" {
		t.Fatalf("dd.MM.yyyy: %q", s.GoString())
	}
}

func TestQDateSetDate(t *testing.T) {
	d := NewQDateFromYMD(2024, 1, 15)
	ok := d.SetDate(2025, 6, 20)
	if !ok {
		t.Fatal("SetDate failed")
	}
	if d.Year() != 2025 || d.Month() != 6 || d.Day() != 20 {
		t.Fatalf("after SetDate: %d-%d-%d", d.Year(), d.Month(), d.Day())
	}
}

func TestQDateGetDate(t *testing.T) {
	d := NewQDateFromYMD(2024, 1, 15)
	y, m, day := d.GetDate()
	if y != 2024 || m != 1 || day != 15 {
		t.Fatalf("GetDate: %d-%d-%d", y, m, day)
	}
}

func TestQDateStatic(t *testing.T) {
	if !QDateIsValid(2024, 2, 29) {
		t.Fatal("2024-02-29 must be valid")
	}
	if QDateIsValid(2023, 2, 29) {
		t.Fatal("2023-02-29 must be invalid")
	}
	if !QDateIsLeapYear(2024) || QDateIsLeapYear(2023) {
		t.Fatal("leap year check failed")
	}
}

func TestQDateJulianRoundTrip(t *testing.T) {
	d := NewQDateFromYMD(2024, 1, 15)
	jd := d.ToJulianDay()
	d2 := NewQDateFromJulianDay(jd)
	if d2 != d {
		t.Fatalf("julian round-trip: %v != %v", d2, d)
	}
}
