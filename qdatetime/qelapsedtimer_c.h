/* C ABI for Qt 6.11 QtCore/QElapsedTimer.
 *
 * QElapsedTimer is a trivially copyable value type: two qint64 fields
 * (t1, t2), no heap allocation, no destructor. We pass it across the C ABI
 * by value (for queries) or by pointer (for mutators), using a struct with
 * the same memory layout as QElapsedTimer.
 */
#ifndef GO_QT6CORE_QDATETIME_QELAPSEDTIMER_C_H
#define GO_QT6CORE_QDATETIME_QELAPSEDTIMER_C_H

#include <stdint.h>
#include <stdbool.h>
#include <stddef.h>

#ifdef __cplusplus
extern "C" {
#endif

typedef struct go_qelapsedtimer_s {
    int64_t t1;
    int64_t t2;
} go_qelapsedtimer_t;

/* ---- Static ---- */
int  qelapsedtimer_clock_type(void);   /* static QElapsedTimer::clockType() */
bool qelapsedtimer_is_monotonic(void); /* static QElapsedTimer::isMonotonic() */

/* ---- Construction ---- */
go_qelapsedtimer_t qelapsedtimer_new(void); /* QElapsedTimer() — invalid */

/* ---- Mutators ---- */
void    qelapsedtimer_start(go_qelapsedtimer_t* self);      /* QElapsedTimer::start() */
int64_t qelapsedtimer_restart(go_qelapsedtimer_t* self);    /* QElapsedTimer::restart() */
void    qelapsedtimer_invalidate(go_qelapsedtimer_t* self); /* QElapsedTimer::invalidate() */

/* ---- Queries ---- */
bool    qelapsedtimer_is_valid(go_qelapsedtimer_t self);                     /* QElapsedTimer::isValid() */
int64_t qelapsedtimer_nsecs_elapsed(go_qelapsedtimer_t self);                /* QElapsedTimer::nsecsElapsed() */
int64_t qelapsedtimer_elapsed(go_qelapsedtimer_t self);                      /* QElapsedTimer::elapsed() */
bool    qelapsedtimer_has_expired(go_qelapsedtimer_t self, int64_t timeout); /* QElapsedTimer::hasExpired(qint64) */
int64_t qelapsedtimer_msecs_since_reference(go_qelapsedtimer_t self);        /* QElapsedTimer::msecsSinceReference() */
int64_t qelapsedtimer_msecs_to(go_qelapsedtimer_t self, go_qelapsedtimer_t other); /* QElapsedTimer::msecsTo(const QElapsedTimer&) */
int64_t qelapsedtimer_secs_to(go_qelapsedtimer_t self, go_qelapsedtimer_t other);  /* QElapsedTimer::secsTo(const QElapsedTimer&) */

/* ---- Comparison ---- */
bool qelapsedtimer_equals(go_qelapsedtimer_t self, go_qelapsedtimer_t other); /* QElapsedTimer::operator== */
bool qelapsedtimer_less(go_qelapsedtimer_t self, go_qelapsedtimer_t other);   /* QElapsedTimer::operator< */

#ifdef __cplusplus
}
#endif
#endif /* GO_QT6CORE_QDATETIME_QELAPSEDTIMER_C_H */