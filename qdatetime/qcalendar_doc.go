// Package qdatetime provides Go bindings for Qt 6.11 QtCore date/time-related types.
//
// Currently supported:
//   - QLocale  (handle-based; COW via QSharedDataPointer)
//   - QCalendar (VALUE-based; trivially copyable struct with a single pointer)
//
// # Value semantics (QCalendar)
//
// QCalendar is a value type in Qt: a single pointer to an internal backend, 8 bytes on
// LP64 (4 on ILP32), trivially copyable. In Go it is stored as a flat struct of the same
// bytes. Assignment and function-passing copy the value; mutating one copy does not affect
// another. There are no pointers, no finalizers and no Close().
//
// This mirrors how QPoint / QPointF / QSize / QSizeF / QMargins / QMarginsF / QLine /
// QLineF / QRect / QRectF are represented in package qgeometry — value types with full
// copies, no lifecycle management.
//
// # Go zero value vs Qt default
//
// The Go zero value QCalendar{} is NOT a valid calendar:
//
//	QCalendar{}       → d_ptr == nil → IsValid() == false
//	NewQCalendar()    → Qt's QCalendar() → valid, Gregorian
//
// In Qt, QCalendar() is an explicit constructor that returns a valid Gregorian calendar;
// there is no way to get an invalid QCalendar through the default constructor. To obtain
// an invalid calendar, pass an unknown name to NewQCalendarFromAny.
//
// # Construction
//
//	NewQCalendar()                            — QCalendar(), Gregorian
//	NewQCalendarFromSystem(System)            — QCalendar(System)
//	NewQCalendarFromSystemID(SystemID)        — QCalendar(SystemId)
//	NewQCalendarFromAny(data, size, tag)      — QCalendar(QAnyStringView)
//
// For System values outside [0, System::Last] and for SystemID values outside the same
// range, the result is an invalid calendar.
//
// # Dependencies
//
//   - qtext.QString — for Name().
//   - qtext.QStringView / QAnyStringView (future) — for constructors that take views.
//
// # Not yet implemented
//
// Methods requiring QDate / QDateTime / QTime (will be added together with those classes):
//   - DateFromParts(year, month, day) QDate
//   - DateFromParts(YearMonthDay parts) QDate
//   - MatchCenturyToWeekday(YearMonthDay parts, int dow) QDate
//   - PartsFromDate(QDate date) YearMonthDay
//   - DayOfWeek(QDate date) int
//   - MonthName(QLocale locale, int month, int year, QLocale::FormatType format) QString
//   - StandaloneMonthName(...)
//   - WeekDayName(QLocale locale, int day, QLocale::FormatType format) QString
//   - StandaloneWeekDayName(...)
//   - DateTimeToString(QStringView format, QDateTime datetime, QDate dateOnly, QTime timeOnly,
//                      QLocale locale) QString
//
// Methods requiring QStringList (will be added together with QStringList in qtext):
//   - AvailableCalendars() QStringList
//
// Struct YearMonthDay (not yet exposed as a Go type):
//   - DateFromParts(YearMonthDay parts)
//   - PartsFromDate(QDate date)
//   - MatchCenturyToWeekday(YearMonthDay parts, int dow)
//
// # QCalendar::SystemId
//
// SystemId is exposed as a Go type SystemID (uintptr) with methods Index() and IsValid().
// Since QCalendar::SystemId's non-default constructors are private in Qt 6.11 (friend
// QCalendarBackend only), SystemID is a pure Go wrapper: IsValid() == (id != SIZE_MAX),
// matching Qt's QCalendar::SystemId::isValid() (which returns `~id`).
//
// Requires Qt >= 6.11, C++17, pkg-config Qt6Core.
package qdatetime
