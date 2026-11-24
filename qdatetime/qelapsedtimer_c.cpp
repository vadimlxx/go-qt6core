#include "../internal/qt6core/check_qt_version.h"

#include <QElapsedTimer>
#include <cstring>

#include "qelapsedtimer_c.h"

namespace {

/* QElapsedTimer is trivially copyable: two qint64 fields (t1, t2).
 * go_qelapsedtimer_t has the same layout, so memcpy between them is
 * well-defined in both directions and avoids any reinterpret_cast.
 */
static_assert(sizeof(QElapsedTimer) == sizeof(go_qelapsedtimer_t),
              "QElapsedTimer and go_qelapsedtimer_t must have identical size");
static_assert(alignof(QElapsedTimer) == alignof(go_qelapsedtimer_t),
              "QElapsedTimer and go_qelapsedtimer_t must have identical alignment");

inline QElapsedTimer toQt(const go_qelapsedtimer_t& g) {
    QElapsedTimer t;
    std::memcpy(&t, &g, sizeof(t));
    return t;
}

inline void fromQt(const QElapsedTimer& t, go_qelapsedtimer_t* g) {
    std::memcpy(g, &t, sizeof(t));
}

} // namespace

extern "C" {

/* ---- Static ---- */

int qelapsedtimer_clock_type(void) {
    return static_cast<int>(QElapsedTimer::clockType());
}
bool qelapsedtimer_is_monotonic(void) {
    return QElapsedTimer::isMonotonic();
}

/* ---- Construction ---- */

go_qelapsedtimer_t qelapsedtimer_new(void) {
    go_qelapsedtimer_t g;
    QElapsedTimer t;   /* default-constructed: invalid (t1 = t2 = INT64_MIN) */
    fromQt(t, &g);
    return g;
}

/* ---- Mutators ---- */

void qelapsedtimer_start(go_qelapsedtimer_t* self) {
    if (!self) return;
    QElapsedTimer t = toQt(*self);
    t.start();
    fromQt(t, self);
}

int64_t qelapsedtimer_restart(go_qelapsedtimer_t* self) {
    if (!self) return 0;
    QElapsedTimer t = toQt(*self);
    const qint64 r = t.restart();
    fromQt(t, self);
    return r;
}

void qelapsedtimer_invalidate(go_qelapsedtimer_t* self) {
    if (!self) return;
    QElapsedTimer t = toQt(*self);
    t.invalidate();
    fromQt(t, self);
}

/* ---- Queries ---- */

bool qelapsedtimer_is_valid(go_qelapsedtimer_t self) {
    return toQt(self).isValid();
}
int64_t qelapsedtimer_nsecs_elapsed(go_qelapsedtimer_t self) {
    return toQt(self).nsecsElapsed();
}
int64_t qelapsedtimer_elapsed(go_qelapsedtimer_t self) {
    return toQt(self).elapsed();
}
bool qelapsedtimer_has_expired(go_qelapsedtimer_t self, int64_t timeout) {
    return toQt(self).hasExpired(static_cast<qint64>(timeout));
}
int64_t qelapsedtimer_msecs_since_reference(go_qelapsedtimer_t self) {
    return toQt(self).msecsSinceReference();
}
int64_t qelapsedtimer_msecs_to(go_qelapsedtimer_t self, go_qelapsedtimer_t other) {
    return toQt(self).msecsTo(toQt(other));
}
int64_t qelapsedtimer_secs_to(go_qelapsedtimer_t self, go_qelapsedtimer_t other) {
    return toQt(self).secsTo(toQt(other));
}

/* ---- Comparison ---- */

bool qelapsedtimer_equals(go_qelapsedtimer_t self, go_qelapsedtimer_t other) {
    return toQt(self) == toQt(other);
}
bool qelapsedtimer_less(go_qelapsedtimer_t self, go_qelapsedtimer_t other) {
    return toQt(self) < toQt(other);
}

} // extern "C"