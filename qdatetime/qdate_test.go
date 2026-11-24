package qdatetime

import (
	"testing"
	"time"

	"github.com/vadimlxx/go-qt6core/qdata"
)

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
	s := d.ToString(qdata.ISODate)
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

func TestQDateFromStringFormatCal(t *testing.T) {
	d := NewQDateFromStringFormatCal("15.01.2024", "dd.MM.yyyy", NewQCalendar())
	if d.Year() != 2024 || d.Month() != 1 || d.Day() != 15 {
		t.Fatalf("NewQDateFromStringFormatCal: %d-%d-%d", d.Year(), d.Month(), d.Day())
	}
}

func TestQDateFromTime(t *testing.T) {
	ts := time.Date(2024, time.February, 29, 13, 45, 59, 0, time.UTC)
	d := NewQDateFromTime(ts)
	if d.Year() != 2024 || d.Month() != 2 || d.Day() != 29 {
		t.Fatalf("NewQDateFromTime: %d-%d-%d", d.Year(), d.Month(), d.Day())
	}
}

func TestQDateToTime(t *testing.T) {
	d := NewQDateFromYMD(2024, 1, 15)
	got := d.ToTime()
	if got.Year() != 2024 || got.Month() != time.January || got.Day() != 15 {
		t.Fatalf("ToTime date: %v", got)
	}
	if got.Hour() != 0 || got.Minute() != 0 || got.Second() != 0 || got.Nanosecond() != 0 {
		t.Fatalf("ToTime time-of-day: %v", got)
	}

	if !NewQDate().ToTime().IsZero() {
		t.Fatal("null QDate must yield the zero time.Time")
	}
}

func TestQDateStartOfDay(t *testing.T) {
	d := NewQDateFromYMD(2024, 1, 15)

	start := d.StartOfDay()
	if start == nil || !start.IsValid() {
		t.Fatal("StartOfDay returned an invalid QDateTime")
	}
	defer start.Close()
	if start.Date().Year() != 2024 || start.Date().Month() != 1 || start.Date().Day() != 15 {
		t.Fatalf("StartOfDay date: %v", start.Date())
	}
	if start.Time().Hour() != 0 || start.Time().Minute() != 0 ||
		start.Time().Second() != 0 || start.Time().Msec() != 0 {
		t.Fatalf("StartOfDay time: %v", start.Time())
	}

	end := d.EndOfDay()
	if end == nil || !end.IsValid() {
		t.Fatal("EndOfDay returned an invalid QDateTime")
	}
	defer end.Close()
	if end.Time().Hour() != 23 || end.Time().Minute() != 59 ||
		end.Time().Second() != 59 || end.Time().Msec() != 999 {
		t.Fatalf("EndOfDay time: %v", end.Time())
	}
}

func TestQDateStartOfDaySpecAndZone(t *testing.T) {
	d := NewQDateFromYMD(2024, 1, 15)

	utcStart := d.StartOfDaySpec(qdata.UTC, 0)
	if utcStart == nil || !utcStart.IsValid() {
		t.Fatal("StartOfDaySpec(UTC) returned an invalid QDateTime")
	}
	defer utcStart.Close()
	if utcStart.Time().Hour() != 0 || utcStart.OffsetFromUtc() != 0 {
		t.Fatalf("StartOfDaySpec(UTC): time=%v offset=%d",
			utcStart.Time(), utcStart.OffsetFromUtc())
	}

	tz := NewQTimeZoneFromOffset(2 * 3600)
	defer tz.Close()
	zoned := d.StartOfDayZone(tz)
	if zoned == nil || !zoned.IsValid() {
		t.Fatal("StartOfDayZone returned an invalid QDateTime")
	}
	defer zoned.Close()
	if zoned.Time().Hour() != 0 || zoned.OffsetFromUtc() != 2*3600 {
		t.Fatalf("StartOfDayZone(+02:00): time=%v offset=%d",
			zoned.Time(), zoned.OffsetFromUtc())
	}
}

func TestQDateFromStringBaseYear(t *testing.T) {
	d := NewQDateFromStringBaseYear("15.01.03", "dd.MM.yy", 1900)
	if d.Year() != 1903 || d.Month() != 1 || d.Day() != 15 {
		t.Fatalf("NewQDateFromStringBaseYear: %d-%d-%d", d.Year(), d.Month(), d.Day())
	}

	d2 := NewQDateFromStringBaseYearCal("15.01.03", "dd.MM.yy", 1900, NewQCalendar())
	if d2.Year() != 1903 || d2.Month() != 1 || d2.Day() != 15 {
		t.Fatalf("NewQDateFromStringBaseYearCal: %d-%d-%d", d2.Year(), d2.Month(), d2.Day())
	}
}
