package qdatetime

import (
	"strings"
	"testing"
)

// helpers ------------------------------------------------------------------

func mkTZ(t *testing.T, id string) *QTimeZone {
	t.Helper()
	z := NewQTimeZoneFromIANA(id)
	t.Cleanup(z.Close)
	return z
}

// ---- Construction --------------------------------------------------------

func TestQTimeZoneNull(t *testing.T) {
	z := NewQTimeZone()
	defer z.Close()
	if z.IsValid() {
		t.Fatal("NewQTimeZone() must be invalid")
	}
}

func TestQTimeZoneUTC(t *testing.T) {
	z := NewQTimeZoneUTC()
	defer z.Close()
	if !z.IsValid() {
		t.Fatal("UTC zone must be valid")
	}
	if z.TimeSpec() != UTC {
		t.Fatalf("TimeSpec = %v, want UTC", z.TimeSpec())
	}
	if !z.IsUTCOrFixedOffset() {
		t.Fatal("UTC zone must report IsUTCOrFixedOffset")
	}
}

func TestQTimeZoneLocal(t *testing.T) {
	z := NewQTimeZoneLocal()
	defer z.Close()
	if !z.IsValid() {
		t.Fatal("LocalTime zone must be valid")
	}
	if z.TimeSpec() != LocalTime {
		t.Fatalf("TimeSpec = %v, want LocalTime", z.TimeSpec())
	}
}

func TestQTimeZoneFromOffset(t *testing.T) {
	// +05:30
	z := NewQTimeZoneFromOffset(5*3600 + 30*60)
	defer z.Close()
	if !z.IsValid() {
		t.Fatal("offset zone must be valid")
	}
	if z.TimeSpec() != OffsetFromUTC {
		t.Fatalf("TimeSpec = %v, want OffsetFromUTC", z.TimeSpec())
	}
	if got := z.FixedSecondsAheadOfUtc(); got != 5*3600+30*60 {
		t.Fatalf("FixedSecondsAheadOfUtc = %d, want %d", got, 5*3600+30*60)
	}
}

func TestQTimeZoneFromIANA(t *testing.T) {
	z := mkTZ(t, "Europe/Berlin")
	if !z.IsValid() {
		t.Skip("Europe/Berlin not available on this system")
	}
	id := z.Id()
	defer id.Close()
	if got := id.String(); got != "Europe/Berlin" {
		t.Fatalf("Id = %q, want %q", got, "Europe/Berlin")
	}
	if z.TimeSpec() != TimeZone {
		t.Fatalf("TimeSpec = %v, want TimeZone", z.TimeSpec())
	}
}

func TestQTimeZoneFromIANAInvalid(t *testing.T) {
	z := mkTZ(t, "Invalid/Zone/Name")
	if z.IsValid() {
		t.Fatal("invalid IANA id must produce an invalid zone")
	}
}

func TestQTimeZoneSystem(t *testing.T) {
	z := SystemTimeZone()
	defer z.Close()
	if !z.IsValid() {
		t.Fatal("system zone must be valid")
	}
}

// ---- Static --------------------------------------------------------------

func TestQTimeZoneSystemTimeZoneId(t *testing.T) {
	id := SystemTimeZoneId()
	defer id.Close()
	if id.IsEmpty() {
		t.Fatal("systemTimeZoneId must not be empty")
	}
}

func TestQTimeZoneIsTimeZoneIdAvailable(t *testing.T) {
	if !IsTimeZoneIdAvailable("Europe/Berlin") {
		t.Skip("Europe/Berlin not available on this system")
	}
	if !IsTimeZoneIdAvailable("UTC") {
		t.Fatal("UTC must always be available")
	}
	if IsTimeZoneIdAvailable("Invalid/Zone/Name") {
		t.Fatal("invalid id must report unavailable")
	}
}

// ---- Offset queries ------------------------------------------------------

func TestQTimeZoneOffsets(t *testing.T) {
	z := mkTZ(t, "Europe/Berlin")
	if !z.IsValid() {
		t.Skip("Europe/Berlin not available on this system")
	}

	// 2024-07-01 12:00 UTC — Berlin is in CEST (UTC+2), DST active.
	summer := NewQDateTimeFromDateTime(
		NewQDateFromYMD(2024, 7, 1),
		NewQTimeFromHMS(12, 0, 0, 0),
	)
	defer summer.Close()
	// Set the time spec to UTC by using fromMSecsSinceEpoch.
	summerUTC := NewQDateTimeFromMSecsSinceEpoch(summer.ToMSecsSinceEpoch())
	defer summerUTC.Close()

	if got := z.OffsetFromUtc(summerUTC); got != 2*3600 {
		t.Fatalf("summer OffsetFromUtc = %d, want %d", got, 2*3600)
	}
	if got := z.StandardTimeOffset(summerUTC); got != 1*3600 {
		t.Fatalf("summer StandardTimeOffset = %d, want %d", got, 1*3600)
	}
	if got := z.DaylightTimeOffset(summerUTC); got != 1*3600 {
		t.Fatalf("summer DaylightTimeOffset = %d, want %d", got, 1*3600)
	}
	if !z.IsDaylightTime(summerUTC) {
		t.Fatal("July in Berlin must be daylight time")
	}

	// 2024-01-01 12:00 UTC — Berlin is in CET (UTC+1), no DST.
	winter := NewQDateTimeFromMSecsSinceEpoch(1704110400000)
	defer winter.Close()
	if got := z.OffsetFromUtc(winter); got != 1*3600 {
		t.Fatalf("winter OffsetFromUtc = %d, want %d", got, 1*3600)
	}
	if z.IsDaylightTime(winter) {
		t.Fatal("January in Berlin must not be daylight time")
	}
}

func TestQTimeZoneHasDaylightTime(t *testing.T) {
	berlin := mkTZ(t, "Europe/Berlin")
	if berlin.IsValid() && !berlin.HasDaylightTime() {
		t.Fatal("Europe/Berlin must have daylight time")
	}
	utc := NewQTimeZoneUTC()
	defer utc.Close()
	if utc.HasDaylightTime() {
		t.Fatal("UTC must not have daylight time")
	}
}

func TestQTimeZoneHasTransitions(t *testing.T) {
	berlin := mkTZ(t, "Europe/Berlin")
	if berlin.IsValid() && !berlin.HasTransitions() {
		t.Fatal("Europe/Berlin must have transitions")
	}
	utc := NewQTimeZoneUTC()
	defer utc.Close()
	if utc.HasTransitions() {
		t.Fatal("UTC must not have transitions")
	}
}

// ---- Display names -------------------------------------------------------

func TestQTimeZoneDisplayNameAt(t *testing.T) {
	z := mkTZ(t, "Europe/Berlin")
	if !z.IsValid() {
		t.Skip("Europe/Berlin not available on this system")
	}
	at := NewQDateTimeFromMSecsSinceEpoch(1704110400000) // 2024-01-01 12:00 UTC
	defer at.Close()

	name := z.DisplayNameAt(at, LongName, nil)
	defer name.Close()
	if name.IsEmpty() {
		t.Fatal("display name must not be empty")
	}
}

func TestQTimeZoneDisplayNameType(t *testing.T) {
	z := mkTZ(t, "Europe/Berlin")
	if !z.IsValid() {
		t.Skip("Europe/Berlin not available on this system")
	}
	name := z.DisplayNameType(StandardTime, LongName, nil)
	defer name.Close()
	if name.IsEmpty() {
		t.Fatal("display name must not be empty")
	}
}

func TestQTimeZoneAbbreviation(t *testing.T) {
	z := mkTZ(t, "Europe/Berlin")
	if !z.IsValid() {
		t.Skip("Europe/Berlin not available on this system")
	}
	at := NewQDateTimeFromMSecsSinceEpoch(1704110400000)
	defer at.Close()
	abbr := z.Abbreviation(at)
	defer abbr.Close()
	if abbr.IsEmpty() {
		t.Fatal("abbreviation must not be empty")
	}
}

// ---- Transitions ---------------------------------------------------------

func TestQTimeZoneNextTransition(t *testing.T) {
	z := mkTZ(t, "Europe/Berlin")
	if !z.IsValid() || !z.HasTransitions() {
		t.Skip("Europe/Berlin not available or has no transitions")
	}
	// 2024-01-01 12:00 UTC — next transition should be the spring DST change.
	start := NewQDateTimeFromMSecsSinceEpoch(1704110400000)
	defer start.Close()

	offset, stdOffset, dstOffset, atUtc, abbr := z.NextTransition(start)
	defer abbr.Close()
	if atUtc != nil {
		defer atUtc.Close()
	}
	if offset == 0 && stdOffset == 0 && dstOffset == 0 {
		t.Fatal("NextTransition returned all-zero offsets")
	}
	if atUtc == nil || !atUtc.IsValid() {
		t.Fatal("NextTransition atUtc must be valid")
	}
	if abbr.IsEmpty() {
		t.Fatal("NextTransition abbreviation must not be empty")
	}
}

func TestQTimeZonePreviousTransition(t *testing.T) {
	z := mkTZ(t, "Europe/Berlin")
	if !z.IsValid() || !z.HasTransitions() {
		t.Skip("Europe/Berlin not available or has no transitions")
	}
	start := NewQDateTimeFromMSecsSinceEpoch(1704110400000)
	defer start.Close()

	offset, stdOffset, dstOffset, atUtc, abbr := z.PreviousTransition(start)
	defer abbr.Close()
	if atUtc != nil {
		defer atUtc.Close()
	}
	if offset == 0 && stdOffset == 0 && dstOffset == 0 {
		t.Fatal("PreviousTransition returned all-zero offsets")
	}
	if atUtc == nil || !atUtc.IsValid() {
		t.Fatal("PreviousTransition atUtc must be valid")
	}
	if abbr.IsEmpty() {
		t.Fatal("PreviousTransition abbreviation must not be empty")
	}
}

// ---- Conversions ---------------------------------------------------------

func TestQTimeZoneAsBackendZone(t *testing.T) {
	z := mkTZ(t, "Europe/Berlin")
	if !z.IsValid() {
		t.Skip("Europe/Berlin not available on this system")
	}
	backend := z.AsBackendZone()
	defer backend.Close()
	if !backend.IsValid() {
		t.Fatal("backend zone must be valid")
	}
	if backend.TimeSpec() != TimeZone {
		t.Fatalf("backend TimeSpec = %v, want TimeZone", backend.TimeSpec())
	}
}

func TestQTimeZoneTerritory(t *testing.T) {
	z := mkTZ(t, "Europe/Berlin")
	if !z.IsValid() {
		t.Skip("Europe/Berlin not available on this system")
	}
	terr := z.Territory()
	// Germany is the territory for Europe/Berlin.
	if terr != Germany {
		t.Fatalf("Territory = %v, want Germany", terr)
	}
}

func TestQTimeZoneIDContains(t *testing.T) {
	z := mkTZ(t, "America/New_York")
	if !z.IsValid() {
		t.Skip("America/New_York not available on this system")
	}
	id := z.Id()
	defer id.Close()
	if !strings.Contains(id.String(), "New_York") {
		t.Fatalf("Id = %q, want it to contain \"New_York\"", id.String())
	}
}