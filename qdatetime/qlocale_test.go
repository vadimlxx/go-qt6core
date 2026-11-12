package qdatetime

import (
	"testing"

	"github.com/vadimlxx/go-qt6core/qtext"
)

// helpers ------------------------------------------------------------------

// mustQString converts a *qtext.QString to a Go string and closes it.
func mustQString(t *testing.T, s *qtext.QString) string {
	t.Helper()
	if s == nil {
		t.Fatal("nil *QString")
	}
	v := s.GoString()
	s.Close()
	return v
}

func mkLocale(t *testing.T, name string) *QLocale {
	t.Helper()
	l := NewQLocaleFromName(name)
	t.Cleanup(l.Close)
	return l
}

// ---- Construction --------------------------------------------------------

func TestQLocaleDefault(t *testing.T) {
	l := NewQLocale()
	t.Cleanup(l.Close)
	if l == nil {
		t.Fatal("NewQLocale returned nil")
	}
	// Default QLocale() must be valid and report some language.
	if l.Language() == AnyLanguage {
		// Not an error per se, but the default should be a real language on any
		// normal system. Log only, do not fail — CI environments can be minimal.
		t.Logf("default locale reported AnyLanguage")
	}
}

func TestQLocaleFromName(t *testing.T) {
	cases := []string{"en_US", "ru_RU", "de_DE", "fr_FR", "ja_JP", "C"}
	for _, name := range cases {
		l := mkLocale(t, name)
		if l == nil {
			t.Fatalf("NewQLocaleFromName(%q) = nil", name)
		}
		// Name() must round-trip for well-known locales.
		got := mustQString(t, l.Name())
		if got == "" {
			t.Errorf("locale %q: Name() is empty", name)
		}
	}
}

func TestQLocaleFromLanguageCountry(t *testing.T) {
	l := NewQLocaleFromLanguageCountry(English, UnitedStates)
	t.Cleanup(l.Close)

	if l.Language() != English {
		t.Fatalf("Language() = %v, want English", l.Language())
	}
	if l.Country() != UnitedStates {
		t.Fatalf("Country() = %v, want UnitedStates", l.Country())
	}
}

func TestQLocaleClone(t *testing.T) {
	orig := mkLocale(t, "en_US")
	clone := orig.Clone()
	t.Cleanup(clone.Close)

	if clone.Language() != orig.Language() {
		t.Fatalf("clone.Language() = %v, want %v", clone.Language(), orig.Language())
	}
	if clone.Country() != orig.Country() {
		t.Fatalf("clone.Country() = %v, want %v", clone.Country(), orig.Country())
	}
	// Both handles point to different C++ objects — closing clone must not
	// affect orig (orig.Close() is called by cleanup).
	if mustQString(t, orig.Name()) == "" {
		t.Fatal("orig.Name() empty after clone creation")
	}
}

// ---- Accessors -----------------------------------------------------------

func TestQLocaleAccessors(t *testing.T) {
	l := mkLocale(t, "en_US")

	if got := l.Language(); got != English {
		t.Errorf("Language() = %v, want English", got)
	}
	if got := l.Script(); got != LatinScript {
		t.Errorf("Script() = %v, want LatinScript", got)
	}
	if got := l.Country(); got != UnitedStates {
		t.Errorf("Country() = %v, want UnitedStates", got)
	}
	// en_US uses imperial (US) measurement system.
	if got := l.MeasurementSystem(); got != ImperialUSSystem {
		t.Errorf("MeasurementSystem() = %v, want ImperialUSSystem", got)
	}
	// English is left-to-right. Qt::LeftToRight == 0.
	if got := l.TextDirection(); got != 0 {
		t.Errorf("TextDirection() = %d, want 0 (LeftToRight)", got)
	}
}

func TestQLocaleAccessorsRussian(t *testing.T) {
	l := mkLocale(t, "ru_RU")

	if got := l.Language(); got != Russian {
		t.Errorf("Language() = %v, want Russian", got)
	}
	if got := l.Country(); got != Russia {
		t.Errorf("Country() = %v, want Russia", got)
	}
	if got := l.MeasurementSystem(); got != MetricSystem {
		t.Errorf("MeasurementSystem() = %v, want MetricSystem", got)
	}
}

// ---- String representation ----------------------------------------------

func TestQLocaleName(t *testing.T) {
	l := mkLocale(t, "en_US")
	got := mustQString(t, l.Name())
	if got != "en_US" {
		t.Fatalf("Name() = %q, want %q", got, "en_US")
	}
}

func TestQLocaleBcp47Name(t *testing.T) {
	l := mkLocale(t, "en_US")

	// Verify the locale was fully parsed (territory = US). On minimal Qt
	// builds without complete CLDR territory data, "en_US" silently falls
	// back to language-only "en" + AnyTerritory.
	if l.Country() != UnitedStates {
		t.Skipf("en_US not fully loaded on this system (country = %v); "+
		"skipping BCP-47 name test", l.Country())
	}

	got := mustQString(t, l.Bcp47Name())

	// Both "en-US" (full BCP-47 tag) and "en" (US is the CLDR likely-subtags
	// default region for English) are valid representations of US English.
	// Qt 6.11 may return either depending on the build and CLDR data.
	if got != "en-US" && got != "en" {
		t.Fatalf("Bcp47Name() = %q, want %q or %q", got, "en-US", "en")
	}
}
func TestQLocaleNativeLanguageName(t *testing.T) {
	l := mkLocale(t, "ru_RU")
	got := mustQString(t, l.NativeLanguageName())
	// Russian locale should report "русский" or similar; we only check non-empty
	// because exact wording is CLDR-version dependent.
	if got == "" {
		t.Fatal("NativeLanguageName() is empty for ru_RU")
	}
}

func TestQLocaleNativeTerritoryName(t *testing.T) {
	l := mkLocale(t, "ru_RU")
	got := mustQString(t, l.NativeTerritoryName())
	if got == "" {
		t.Fatal("NativeTerritoryName() is empty for ru_RU")
	}
}

// ---- Number formatting ---------------------------------------------------

func TestQLocaleToStringInt64(t *testing.T) {
	en := mkLocale(t, "en_US")
	if got := mustQString(t, en.ToStringInt64(1234567)); got != "1,234,567" {
		t.Fatalf("en_US.ToStringInt64(1234567) = %q, want %q", got, "1,234,567")
	}

	de := mkLocale(t, "de_DE")
	if got := mustQString(t, de.ToStringInt64(1234567)); got != "1.234.567" {
		t.Fatalf("de_DE.ToStringInt64(1234567) = %q, want %q", got, "1.234.567")
	}

	// Negative numbers.
	if got := mustQString(t, en.ToStringInt64(-42)); got != "-42" {
		t.Fatalf("en_US.ToStringInt64(-42) = %q, want %q", got, "-42")
	}
}

func TestQLocaleToStringUInt64(t *testing.T) {
	en := mkLocale(t, "en_US")
	if got := mustQString(t, en.ToStringUInt64(9876543210)); got != "9,876,543,210" {
		t.Fatalf("en_US.ToStringUInt64(...) = %q, want %q", got, "9,876,543,210")
	}
}

func TestQLocaleToStringDouble(t *testing.T) {
	en := mkLocale(t, "en_US")

	// Default Qt number format 'g'.
	if got := mustQString(t, en.ToStringDouble(3.14, 'g', 6)); got != "3.14" {
		t.Fatalf("ToStringDouble(3.14,'g',6) = %q, want %q", got, "3.14")
	}

	// Fixed format with precision.
	if got := mustQString(t, en.ToStringDouble(3.14159, 'f', 2)); got != "3.14" {
		t.Fatalf("ToStringDouble(3.14159,'f',2) = %q, want %q", got, "3.14")
	}

	// Scientific.
	if got := mustQString(t, en.ToStringDouble(12345.678, 'e', 2)); got != "1.23e+04" {
		t.Fatalf("ToStringDouble(12345.678,'e',2) = %q, want %q", got, "1.23e+04")
	}

	// German locale uses comma as decimal separator.
	de := mkLocale(t, "de_DE")
	if got := mustQString(t, de.ToStringDouble(3.14, 'f', 2)); got != "3,14" {
		t.Fatalf("de_DE.ToStringDouble(3.14,'f',2) = %q, want %q", got, "3,14")
	}
}

// ---- Number parsing ------------------------------------------------------

func TestQLocaleToInt64(t *testing.T) {
	en := mkLocale(t, "en_US")

	v, ok := en.ToInt64("1234")
	if !ok || v != 1234 {
		t.Fatalf("ToInt64(\"1234\") = (%d, %v), want (1234, true)", v, ok)
	}

	v, ok = en.ToInt64("not a number")
	if ok {
		t.Fatalf("ToInt64(\"not a number\") ok=true, want false")
	}
	_ = v

	// Empty string.
	_, ok = en.ToInt64("")
	if ok {
		t.Fatal("ToInt64(\"\") ok=true, want false")
	}
}

func TestQLocaleToDouble(t *testing.T) {
	en := mkLocale(t, "en_US")

	v, ok := en.ToDouble("3.14")
	if !ok || v < 3.13 || v > 3.15 {
		t.Fatalf("en_US.ToDouble(\"3.14\") = (%v, %v)", v, ok)
	}

	de := mkLocale(t, "de_DE")
	v, ok = de.ToDouble("3,14")
	if !ok || v < 3.13 || v > 3.15 {
		t.Fatalf("de_DE.ToDouble(\"3,14\") = (%v, %v)", v, ok)
	}

	_, ok = en.ToDouble("garbage")
	if ok {
		t.Fatal("ToDouble(\"garbage\") ok=true, want false")
	}
}

// ---- Locale-specific strings --------------------------------------------

func TestQLocaleDecimalPoint(t *testing.T) {
	en := mkLocale(t, "en_US")
	if got := mustQString(t, en.DecimalPoint()); got != "." {
		t.Fatalf("en_US.DecimalPoint() = %q, want %q", got, ".")
	}

	de := mkLocale(t, "de_DE")
	if got := mustQString(t, de.DecimalPoint()); got != "," {
		t.Fatalf("de_DE.DecimalPoint() = %q, want %q", got, ",")
	}
}

func TestQLocaleGroupSeparator(t *testing.T) {
	en := mkLocale(t, "en_US")
	if got := mustQString(t, en.GroupSeparator()); got != "," {
		t.Fatalf("en_US.GroupSeparator() = %q, want %q", got, ",")
	}

	de := mkLocale(t, "de_DE")
	if got := mustQString(t, de.GroupSeparator()); got != "." {
		t.Fatalf("de_DE.GroupSeparator() = %q, want %q", got, ".")
	}
}

func TestQLocalePercent(t *testing.T) {
	en := mkLocale(t, "en_US")
	if got := mustQString(t, en.Percent()); got != "%" {
		t.Fatalf("Percent() = %q, want %q", got, "%")
	}
}

func TestQLocaleZeroDigit(t *testing.T) {
	en := mkLocale(t, "en_US")
	if got := mustQString(t, en.ZeroDigit()); got != "0" {
		t.Fatalf("en_US.ZeroDigit() = %q, want %q", got, "0")
	}

	// Arabic-Indic digits (Egypt locale uses ٠ U+0660 as zero).
	ar := mkLocale(t, "ar_EG")
	got := mustQString(t, ar.ZeroDigit())
	if got == "" || got == "0" {
		t.Fatalf("ar_EG.ZeroDigit() = %q, expected non-ASCII digit", got)
	}
}

func TestQLocaleSigns(t *testing.T) {
	en := mkLocale(t, "en_US")
	if got := mustQString(t, en.NegativeSign()); got != "-" {
		t.Fatalf("NegativeSign() = %q, want %q", got, "-")
	}
	if got := mustQString(t, en.PositiveSign()); got != "+" {
		t.Fatalf("PositiveSign() = %q, want %q", got, "+")
	}
}

func TestQLocaleExponential(t *testing.T) {
	en := mkLocale(t, "en_US")
	got := mustQString(t, en.Exponential())
	if got == "" {
		t.Fatal("Exponential() is empty")
	}
}

// ---- Month / day names ---------------------------------------------------

func TestQLocaleMonthName(t *testing.T) {
	en := mkLocale(t, "en_US")
	if got := mustQString(t, en.MonthName(1, LongFormat)); got != "January" {
		t.Fatalf("MonthName(1, LongFormat) = %q, want %q", got, "January")
	}
	if got := mustQString(t, en.MonthName(12, LongFormat)); got != "December" {
		t.Fatalf("MonthName(12, LongFormat) = %q, want %q", got, "December")
	}
	if got := mustQString(t, en.MonthName(1, ShortFormat)); got != "Jan" {
		t.Fatalf("MonthName(1, ShortFormat) = %q, want %q", got, "Jan")
	}
}

func TestQLocaleStandaloneMonthName(t *testing.T) {
	en := mkLocale(t, "en_US")
	got := mustQString(t, en.StandaloneMonthName(1, LongFormat))
	if got != "January" {
		t.Fatalf("StandaloneMonthName(1, LongFormat) = %q, want %q", got, "January")
	}
}

func TestQLocaleDayName(t *testing.T) {
	en := mkLocale(t, "en_US")
	// Qt::Monday == 1.
	if got := mustQString(t, en.DayName(1, LongFormat)); got != "Monday" {
		t.Fatalf("DayName(1, LongFormat) = %q, want %q", got, "Monday")
	}
	if got := mustQString(t, en.DayName(7, LongFormat)); got != "Sunday" {
		t.Fatalf("DayName(7, LongFormat) = %q, want %q", got, "Sunday")
	}
}

func TestQLocaleStandaloneDayName(t *testing.T) {
	en := mkLocale(t, "en_US")
	got := mustQString(t, en.StandaloneDayName(1, LongFormat))
	if got != "Monday" {
		t.Fatalf("StandaloneDayName(1, LongFormat) = %q, want %q", got, "Monday")
	}
}

// ---- AM/PM ---------------------------------------------------------------

func TestQLocaleAmPm(t *testing.T) {
	en := mkLocale(t, "en_US")
	if got := mustQString(t, en.AmText()); got != "AM" {
		t.Fatalf("AmText() = %q, want %q", got, "AM")
	}
	if got := mustQString(t, en.PmText()); got != "PM" {
		t.Fatalf("PmText() = %q, want %q", got, "PM")
	}
}

// ---- Case conversion -----------------------------------------------------

func TestQLocaleToUpper(t *testing.T) {
	en := mkLocale(t, "en_US")
	if got := mustQString(t, en.ToUpper("hello world")); got != "HELLO WORLD" {
		t.Fatalf("ToUpper = %q, want %q", got, "HELLO WORLD")
	}

	// Turkish: 'i' → 'İ' (U+0130), not 'I'. This requires the Turkish CLDR
	// data to be present on the system. On minimal CI images the locale silently
	// falls back to C, so we skip rather than fail if tr_TR did not load.
	tr := mkLocale(t, "tr_TR")
	if tr.Language() != Turkish {
		t.Skipf("tr_TR not available on this system (loaded language = %v); "+
		"skipping Turkish-specific casing test", tr.Language())
	}

	got := mustQString(t, tr.ToUpper("i"))
	if got != "İ" {
		t.Fatalf("tr_TR.ToUpper(\"i\") = %q (first rune U+%04X), want %q (U+0130)",
				 got, firstRune(got), "İ")
	}
}

func TestQLocaleToLower(t *testing.T) {
	en := mkLocale(t, "en_US")
	if got := mustQString(t, en.ToLower("HELLO WORLD")); got != "hello world" {
		t.Fatalf("ToLower = %q, want %q", got, "hello world")
	}
}

// firstRune returns the first rune of s as an int (for test error messages).
func firstRune(s string) int {
	for _, r := range s {
		return int(r)
	}
	return 0
}

// ---- Static methods ------------------------------------------------------

func TestCLocale(t *testing.T) {
	l := CLocale()
	t.Cleanup(l.Close)
	if l.Language() != C {
		t.Fatalf("CLocale().Language() = %v, want C", l.Language())
	}
}

func TestSystemLocale(t *testing.T) {
	l := SystemLocale()
	t.Cleanup(l.Close)
	if l == nil {
		t.Fatal("SystemLocale() = nil")
	}
}

func TestSetDefaultLocale(t *testing.T) {
	// Save the current default (best-effort: we can't easily snapshot it,
	// but we restore to system locale at the end).
	defer SetDefaultLocale(SystemLocale())

	ru := mkLocale(t, "ru_RU")
	SetDefaultLocale(ru)

	// NewQLocale() must return the new default.
	check := NewQLocale()
	t.Cleanup(check.Close)

	if check.Language() != Russian {
		t.Fatalf("after SetDefaultLocale(ru_RU), default.Language() = %v, want Russian",
			check.Language())
	}
}
