/* C ABI for Qt 6.11 QtCore/QCalendar.
 *
 * QCalendar is a trivially copyable value type: a single pointer (d_ptr). We pass it by
 * value everywhere — no heap allocation, no lifecycle management.
 *
 * In C++ mode go_qcalendar_t is typedef'd directly to QCalendar; in C mode it falls back
 * to a single-pointer POD struct with the same layout (8 bytes on LP64, 4 on ILP32).
 *
 * String outputs use void** out-parameter; the Go side wraps the result via
 * qtext.QStringFromUnsafePtr.
 */
#ifndef GO_QT6CORE_QDATETIME_QCALENDAR_C_H
#define GO_QT6CORE_QDATETIME_QCALENDAR_C_H

#include <stddef.h>
#include <stdint.h>
#include <stdbool.h>

#ifdef __cplusplus
#include <QCalendar>
typedef QCalendar go_qcalendar_t;
extern "C" {
#else
typedef struct go_qcalendar_s {
    const void* d_ptr;
} go_qcalendar_t;
#endif

/* --- Construction (return by value, no allocation) --- */
go_qcalendar_t qcalendar_default(void);
go_qcalendar_t qcalendar_from_system(int system);
go_qcalendar_t qcalendar_from_systemid(uintptr_t id);
go_qcalendar_t qcalendar_from_any(const void* data, ptrdiff_t size, int tag);

/* --- Validity --- */
bool qcalendar_is_valid(go_qcalendar_t self);

/* --- Date queries --- */
int  qcalendar_days_in_month(go_qcalendar_t self, int month, int year);
int  qcalendar_days_in_year(go_qcalendar_t self, int year);
int  qcalendar_months_in_year(go_qcalendar_t self, int year);
bool qcalendar_is_date_valid(go_qcalendar_t self, int year, int month, int day);

/* --- Leap years --- */
bool qcalendar_is_leap_year(go_qcalendar_t self, int year);

/* --- Calendar properties --- */
bool qcalendar_is_gregorian(go_qcalendar_t self);
bool qcalendar_is_lunar(go_qcalendar_t self);
bool qcalendar_is_luni_solar(go_qcalendar_t self);
bool qcalendar_is_solar(go_qcalendar_t self);
bool qcalendar_is_proleptic(go_qcalendar_t self);
bool qcalendar_has_year_zero(go_qcalendar_t self);
int  qcalendar_maximum_days_in_month(go_qcalendar_t self);
int  qcalendar_minimum_days_in_month(go_qcalendar_t self);
int  qcalendar_maximum_months_in_year(go_qcalendar_t self);
void qcalendar_name(go_qcalendar_t self, void** out_qstr);

#ifdef __cplusplus
}
#endif
#endif