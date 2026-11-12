package qdatetime

// SystemID mirrors QCalendar::SystemId.
//
// In Qt 6.11, QCalendar::SystemId's non-default constructors —
// SystemId(QCalendar::System) and SystemId(size_t) — are private and callable
// only by QCalendarBackend (via `friend class QCalendarBackend`). Therefore
// external code cannot construct a valid SystemId from C++. The only public
// constructor is the default one, which produces an INVALID id (SIZE_MAX).
//
// Because of this, SystemID is exposed as a pure Go wrapper around the
// underlying size_t (uintptr) value, without any C++-backed construction:
//
//   - The zero value SystemID(0) corresponds to QCalendar::System::Gregorian
//     (enum value 0) — matching SystemId(System::Gregorian).isValid() == true.
//   - The "invalid" value is ^uintptr(0) (SIZE_MAX), matching Qt's default-
//     constructed SystemId() (which has id == ~size_t(0)).
//
// The IsValid/Index methods are pure Go and do not call into C++ — they
// reproduce exactly the logic in QCalendar::SystemId:
//
//   constexpr size_t index()     const noexcept { return id; }
//   constexpr bool   isValid()   const noexcept { return ~id; }
//
// i.e. IsValid() == (id != SIZE_MAX), which is what the Go code below does.
//
// Construction path: use NewQCalendarFromSystemID to obtain a QCalendar whose
// underlying System corresponds to id (for id in [0, System::Last]); for other
// ids, an invalid QCalendar is returned.
type SystemID uintptr

// Index returns the underlying size_t index of the SystemId.
// Mirrors QCalendar::SystemId::index().
func (id SystemID) Index() uintptr { return uintptr(id) }

// IsValid reports whether the SystemId is valid (i.e. not SIZE_MAX).
// Mirrors QCalendar::SystemId::isValid() — which returns `~id`, equivalent to
// `id != SIZE_MAX`.
func (id SystemID) IsValid() bool {
	return uintptr(id) != ^uintptr(0)
}