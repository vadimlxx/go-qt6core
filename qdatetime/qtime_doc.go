// Package qdatetime provides the QTime binding: a value-type wrapper around
// Qt 6.11 QtCore/QTime.
//
// QTime stores milliseconds since the start of the day as int32: -1 means null,
// 0..86399999 is a valid time. In Go it is QTime (int32): trivially copyable,
// no pointers, no destructor, no Close.
//
// # Null vs zero
//
//	QTime{}       msecs = 0   valid 00:00:00.000
//	NewQTime()    null QTime  mirrors Qt QTime() (msecs = -1)
//
// # Type mapping
//
//	Qt                                  Go
//	QTime                               qdatetime.QTime
//	Qt::DateFormat                      qdata.DateFormat
//	QString                             *qtext.QString
//	time.Time                           time.Time
//
// # Constructors
//
//	NewQTime()                          QTime()
//	NewQTimeFromHMS(h, m, s, ms)        QTime(h, m, s, ms)
//	NewQTimeCurrentTime()               QTime::currentTime()
//	NewQTimeFromString(s, format)       QTime::fromString(s, Qt::DateFormat)
//	NewQTimeFromStringFormat(s, format) QTime::fromString(s, QString)
//	NewQTimeFromMSecsSinceStartOfDay(n) QTime::fromMSecsSinceStartOfDay(n)
//	NewQTimeFromTime(time.Time)         fields taken from time.Time
//
// # Accessors and arithmetic
//
//	IsNull, IsValid, Hour, Minute, Second, Msec, MsecsSinceStartOfDay,
//	AddSecs, AddMSecs, SecsTo, MsecsTo
//
// # Strings and setters
//
//	ToString, ToStringFormat, SetHMS, QTimeIsValid
//
// # time.Time
//
//	ToTime() time.Time
//
// ToTime returns the time-of-day on 1970-01-01 UTC (millisecond precision); a
// null/invalid time yields the zero time.Time.
//
// # Not implemented yet
//
//   - fromStdTimePoint / toStdTimePoint and other std::chrono interop.
//   - QDataStream / qHash / QDebug integration.
//   - Comparison operators: use MsecsSinceStartOfDay / MsecsTo.
//
// The C bridge is shared by QDate, QTime and QDateTime in qdatetime_c.h /
// qdatetime_c.cpp, mirroring QtCore's qdatetime.h / qdatetime.cpp.
//
// This package is a subpackage of github.com/vadimlxx/go-qt6core.
package qdatetime
