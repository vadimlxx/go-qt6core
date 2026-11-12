package qdatetime

import (
	"strings"
	"testing"
	"unsafe"
)

// helpers ------------------------------------------------------------------

func mkCal(t *testing.T) QCalendar {
	t.Helper()
	return NewQCalendar()
}

func mkCalSys(t *testing.T, sys System) QCalendar {
	t.Helper()
	return NewQCalendar(sys)
}

// ---- Construction --------------------------------------------------------

func TestQCalendarDefault(t *testing.T) {
	c := mkCal(t)
	if !c.IsValid() {
		t.Fatal("NewQCalendar() must be valid")
	}
	if !c.IsGregorian() {
		t.Fatal("NewQCalendar() must be Gregorian")
	}
}

func TestQCalendarZeroValueInvalid(t *testing.T) {
	var c QCalendar
	if c.IsValid() {
		t.Fatal("Go zero value QCalendar{} must be invalid (d_ptr == nil)")
	}
	// Methods on an invalid calendar should not crash — they hit the
	// QCalendarBackend null-check inside Qt and return safe defaults.
	if c.IsGregorian() {
		t.Fatal("invalid calendar must not report Gregorian")
	}
	if c.DaysInYear(2024) != 0 {
		t.Fatal("invalid calendar must report DaysInYear() == 0")
	}
}

func TestQCalendarFromSystemGregorian(t *testing.T) {
	c := mkCalSys(t, Gregorian)
	if !c.IsValid() {
		t.Fatal("Gregorian calendar must always be valid")
	}
	if !c.IsGregorian() {
		t.Fatal("Gregorian calendar must report IsGregorian() = true")
	}
}

func TestQCalendarFromSystemJulian(t *testing.T) {
	c := NewQCalendarFromSystem(Julian)
	if !c.IsValid() {
		t.Skip("Julian calendar not available on this Qt build")
	}
	if c.IsGregorian() {
		t.Fatal("Julian must not be Gregorian")
	}
	if !c.IsSolar() {
		t.Fatal("Julian must be solar")
	}
	if !c.IsLeapYear(1900) {
		t.Fatal("Julian year 1900 must be leap (every-4-years rule)")
	}
	if mkCal(t).IsLeapYear(1900) {
		t.Fatal("Gregorian year 1900 must not be leap")
	}
}

func TestQCalendarFromSystemOutOfRange(t *testing.T) {
	// Anything outside [0, System::Last] must produce an invalid calendar
	// rather than UB. System::User = -1 also invalidates.
	cases := []int{-100, -1, 100, 9999}
	for _, s := range cases {
		c := NewQCalendarFromSystem(System(s))
		if c.IsValid() {
			t.Errorf("NewQCalendarFromSystem(%d) must be invalid", s)
		}
	}
}

func TestQCalendarFromSystemID(t *testing.T) {
	c := NewQCalendarFromSystemID(SystemID(Gregorian))
	if !c.IsValid() {
		t.Fatal("calendar from SystemID(Gregorian) must be valid")
	}
	if !c.IsGregorian() {
		t.Fatal("calendar from SystemID(Gregorian) must be Gregorian")
	}
}

func TestQCalendarFromSystemIDOutOfRange(t *testing.T) {
	// SIZE_MAX and anything > System::Last must produce an invalid calendar.
	for _, id := range []SystemID{
		SystemID(^uintptr(0)),       // SIZE_MAX, Qt's invalid SystemId
		SystemID(^uintptr(0) - 1),   // just below
		SystemID(100),
	} {
		c := NewQCalendarFromSystemID(id)
		if c.IsValid() {
			t.Errorf("NewQCalendarFromSystemID(%d) must be invalid", uintptr(id))
		}
	}
}

func TestQCalendarFromAnyUTF8(t *testing.T) {
	name := "gregorian"
	c := NewQCalendarFromAny(
		unsafe.Pointer(unsafe.StringData(name)),
		len(name),
		0, // UTF-8
	)
	if !c.IsValid() {
		t.Skipf("QCalendar(%q) returned invalid on this Qt build", name)
	}
	if !c.IsGregorian() {
		t.Fatal(`QCalendar("gregorian") must be Gregorian`)
	}
}

func TestQCalendarFromAnyUnknownName(t *testing.T) {
	name := "definitely-not-a-real-calendar-name-xyz"
	c := NewQCalendarFromAny(
		unsafe.Pointer(unsafe.StringData(name)),
		len(name),
		0,
	)
	if c.IsValid() {
		t.Fatal("QCalendar(<unknown name>) must be invalid")
	}
}

func TestQCalendarFromAnyEmpty(t *testing.T) {
	c := NewQCalendarFromAny(nil, 0, 0)
	if c.IsValid() {
		t.Fatal("QCalendar(<empty>) must be invalid")
	}
}

// ---- Value semantics -----------------------------------------------------

func TestQCalendarValueSemantics(t *testing.T) {
	a := NewQCalendar()
	b := a // value copy — independent struct

	// Both report the same properties.
	if a.IsGregorian() != b.IsGregorian() {
		t.Fatal("copy must preserve IsGregorian()")
	}
	if a.IsValid() != b.IsValid() {
		t.Fatal("copy must preserve IsValid()")
	}

	// Go has no mutators on QCalendar, so we can only verify that assigning a
	// different value to b does not affect a.
	prev := a.IsValid()
	b = QCalendar{} // zero value
	if a.IsValid() != prev {
		t.Fatal("a.IsValid() changed after b was reassigned to zero value")
	}
	if b.IsValid() {
		t.Fatal("b must be invalid after assignment to zero value")
	}
}

// ---- Date queries --------------------------------------------------------

func TestQCalendarDaysInMonthGregorian(t *testing.T) {
	c := mkCal(t)
	cases := []struct {
		month int
		year  int
		want  int
	}{
		{1, 2024, 31},
		{2, 2024, 29},
		{2, 2023, 28},
		{2, 1900, 28},
		{2, 2000, 29},
		{3, 2024, 31},
		{4, 2024, 30},
		{5, 2024, 31},
		{6, 2024, 30},
		{7, 2024, 31},
		{8, 2024, 31},
		{9, 2024, 30},
		{10, 2024, 31},
		{11, 2024, 30},
		{12, 2024, 31},
	}
	for _, tc := range cases {
		if got := c.DaysInMonth(tc.month, tc.year); got != tc.want {
			t.Errorf("DaysInMonth(%d, %d) = %d, want %d", tc.month, tc.year, got, tc.want)
		}
	}
}

func TestQCalendarDaysInYearGregorian(t *testing.T) {
	c := mkCal(t)
	cases := []struct {
		year int
		want int
	}{
		{2024, 366},
		{2023, 365},
		{2000, 366},
		{1900, 365},
		{1600, 366},
		{1, 365},
	}
	for _, tc := range cases {
		if got := c.DaysInYear(tc.year); got != tc.want {
			t.Errorf("DaysInYear(%d) = %d, want %d", tc.year, got, tc.want)
		}
	}
}

func TestQCalendarMonthsInYearGregorian(t *testing.T) {
	c := mkCal(t)
	for _, y := range []int{1, 1000, 2024, 9999} {
		if got := c.MonthsInYear(y); got != 12 {
			t.Errorf("MonthsInYear(%d) = %d, want 12", y, got)
		}
	}
}

func TestQCalendarIsDateValidGregorian(t *testing.T) {
	c := mkCal(t)
	cases := []struct {
		year, month, day int
		want             bool
	}{
		{2024, 1, 1, true},
		{2024, 2, 29, true},
		{2023, 2, 29, false},
		{2024, 4, 31, false},
		{2024, 12, 31, true},
		{2024, 13, 1, false},
		{2024, 0, 1, false},
		{2024, 1, 0, false},
		{2024, 1, 32, false},
		{2024, -1, 1, false},
	}
	for _, tc := range cases {
		if got := c.IsDateValid(tc.year, tc.month, tc.day); got != tc.want {
			t.Errorf("IsDateValid(%d, %d, %d) = %v, want %v",
				tc.year, tc.month, tc.day, got, tc.want)
		}
	}
}

// ---- Leap years ----------------------------------------------------------

func TestQCalendarIsLeapYearGregorian(t *testing.T) {
	c := mkCal(t)
	cases := []struct {
		year int
		want bool
	}{
		{2024, true}, {2023, false}, {2025, false},
		{2000, true}, {1900, false}, {2100, false}, {2400, true},
	}
	for _, tc := range cases {
		if got := c.IsLeapYear(tc.year); got != tc.want {
			t.Errorf("IsLeapYear(%d) = %v, want %v", tc.year, got, tc.want)
		}
	}
}

// ---- Calendar properties -------------------------------------------------

func TestQCalendarGregorianProperties(t *testing.T) {
	c := mkCal(t)
	if !c.IsGregorian() { t.Error("IsGregorian() = false, want true") }
	if c.IsLunar() { t.Error("IsLunar() = true, want false") }
	if c.IsLuniSolar() { t.Error("IsLuniSolar() = true, want false") }
	if !c.IsSolar() { t.Error("IsSolar() = false, want true") }
	if !c.IsProleptic() { t.Error("IsProleptic() = false, want true") }
	// The proleptic Gregorian calendar has NO year zero: years go
	// ... 2 BCE, 1 BCE, 1 CE, 2 CE ... — see Qt docs for QCalendar.
	if c.HasYearZero() { t.Error("HasYearZero() = true, want false (proleptic Gregorian has no year 0)") }
}

func TestQCalendarGregorianLimits(t *testing.T) {
	c := mkCal(t)
	if got := c.MaximumDaysInMonth(); got != 31 {
		t.Errorf("MaximumDaysInMonth() = %d, want 31", got)
	}
	if got := c.MinimumDaysInMonth(); got != 28 {
		t.Errorf("MinimumDaysInMonth() = %d, want 28", got)
	}
	if got := c.MaximumMonthsInYear(); got != 12 {
		t.Errorf("MaximumMonthsInYear() = %d, want 12", got)
	}
}

func TestQCalendarName(t *testing.T) {
	c := mkCal(t)
	name := mustQString(t, c.Name())
	if name == "" {
		t.Fatal("Name() is empty for Gregorian calendar")
	}
	if !strings.Contains(strings.ToLower(name), "greg") {
		t.Fatalf("Name() = %q, expected to contain \"greg\"", name)
	}
}

func TestQCalendarNameOnZeroValue(t *testing.T) {
	// Invalid calendar must not crash; Name() returns an empty string.
	var c QCalendar
	name := mustQString(t, c.Name())
	if name != "" {
		t.Fatalf("Name() on zero value = %q, want empty", name)
	}
}

// ---- SystemID ------------------------------------------------------------

func TestSystemIDGregorian(t *testing.T) {
	id := SystemID(Gregorian)
	if uintptr(id) != 0 {
		t.Fatalf("SystemID(Gregorian) = %d, want 0", uintptr(id))
	}
	if !id.IsValid() {
		t.Fatal("SystemID(Gregorian) must be valid")
	}
	if id.Index() != 0 {
		t.Fatalf("SystemID(Gregorian).Index() = %d, want 0", id.Index())
	}
}

func TestSystemIDInvalid(t *testing.T) {
	id := SystemID(^uintptr(0))
	if id.IsValid() {
		t.Fatal("SystemID(^uintptr(0)) must be invalid")
	}
}

func TestSystemIDZeroValueIsGregorian(t *testing.T) {
	var id SystemID
	if !id.IsValid() {
		t.Fatal("zero SystemID(0) must be valid (Gregorian)")
	}
}
