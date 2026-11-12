/* Order: version check first, then real Qt headers.
 *
 * In C++ mode, go_qmargins_t  == QMargins  and  go_qmarginsf_t == QMarginsF
 * (see qmargins_c.h). Therefore every function below operates on real Qt
 * values directly, with no fromC/toC helpers and no reinterpret_cast.
 */
#include "../internal/qt6core/check_qt_version.h"

#include <QMargins>

#include "qmargins_c.h"

extern "C" {

/* --- QMargins: construction --- */
go_qmargins_t qmargins_default(void)                   { return QMargins(); }
go_qmargins_t qmargins_new(int l, int t, int r, int b) { return QMargins(l, t, r, b); }

/* --- QMargins: accessors --- */
bool qmargins_is_null(go_qmargins_t m) { return m.isNull(); }
int  qmargins_left  (go_qmargins_t m)  { return m.left();   }
int  qmargins_top   (go_qmargins_t m)  { return m.top();    }
int  qmargins_right (go_qmargins_t m)  { return m.right();  }
int  qmargins_bottom(go_qmargins_t m)  { return m.bottom(); }

/* --- QMargins: mutating setters ---
 * Pass-by-value + copy-mutate-return: gives Go value semantics while letting
 * C++ call the real QMargins setters. */
go_qmargins_t qmargins_set_left  (go_qmargins_t m, int v) { m.setLeft(v);   return m; }
go_qmargins_t qmargins_set_top   (go_qmargins_t m, int v) { m.setTop(v);    return m; }
go_qmargins_t qmargins_set_right (go_qmargins_t m, int v) { m.setRight(v);  return m; }
go_qmargins_t qmargins_set_bottom(go_qmargins_t m, int v) { m.setBottom(v); return m; }

/* --- QMargins: arithmetic (uses Qt's real operators) --- */
go_qmargins_t qmargins_add       (go_qmargins_t a, go_qmargins_t b) { return a + b; }
go_qmargins_t qmargins_sub       (go_qmargins_t a, go_qmargins_t b) { return a - b; }
go_qmargins_t qmargins_add_int   (go_qmargins_t a, int v)           { return a + v; }
go_qmargins_t qmargins_sub_int   (go_qmargins_t a, int v)           { return a - v; }
go_qmargins_t qmargins_mul_int   (go_qmargins_t a, int f)           { return a * f; }
go_qmargins_t qmargins_div_int   (go_qmargins_t a, int d)           { return a / d; }
go_qmargins_t qmargins_mul_float (go_qmargins_t a, double f)        { return a * f; }
go_qmargins_t qmargins_div_float (go_qmargins_t a, double d)        { return a / d; }
go_qmargins_t qmargins_or        (go_qmargins_t a, go_qmargins_t b) { return a | b; }
go_qmargins_t qmargins_neg       (go_qmargins_t a)                  { return -a; }
go_qmargins_t qmargins_unary_plus(go_qmargins_t a)                  { return +a; }

/* --- QMargins: comparison --- */
bool qmargins_equal(go_qmargins_t a, go_qmargins_t b) { return a == b; }

/* --- QMargins: conversions --- */
go_qmarginsf_t qmargins_to_marginsf(go_qmargins_t m) { return m.toMarginsF(); }

/* --- QMarginsF: construction --- */
go_qmarginsf_t qmarginsf_default(void) { return QMarginsF(); }
go_qmarginsf_t qmarginsf_new(double l, double t, double r, double b) { return QMarginsF(l, t, r, b); }
go_qmarginsf_t qmarginsf_from_margins(go_qmargins_t m) { return QMarginsF(m); }

/* --- QMarginsF: accessors --- */
bool   qmarginsf_is_null(go_qmarginsf_t m) { return m.isNull(); }
double qmarginsf_left   (go_qmarginsf_t m) { return m.left();   }
double qmarginsf_top    (go_qmarginsf_t m) { return m.top();    }
double qmarginsf_right  (go_qmarginsf_t m) { return m.right();  }
double qmarginsf_bottom (go_qmarginsf_t m) { return m.bottom(); }

/* --- QMarginsF: mutating setters --- */
go_qmarginsf_t qmarginsf_set_left  (go_qmarginsf_t m, double v) { m.setLeft(v);   return m; }
go_qmarginsf_t qmarginsf_set_top   (go_qmarginsf_t m, double v) { m.setTop(v);    return m; }
go_qmarginsf_t qmarginsf_set_right (go_qmarginsf_t m, double v) { m.setRight(v);  return m; }
go_qmarginsf_t qmarginsf_set_bottom(go_qmarginsf_t m, double v) { m.setBottom(v); return m; }

/* --- QMarginsF: arithmetic --- */
go_qmarginsf_t qmarginsf_add       (go_qmarginsf_t a, go_qmarginsf_t b) { return a + b; }
go_qmarginsf_t qmarginsf_sub       (go_qmarginsf_t a, go_qmarginsf_t b) { return a - b; }
go_qmarginsf_t qmarginsf_add_float (go_qmarginsf_t a, double v)         { return a + v; }
go_qmarginsf_t qmarginsf_sub_float (go_qmarginsf_t a, double v)         { return a - v; }
go_qmarginsf_t qmarginsf_mul_float (go_qmarginsf_t a, double f)         { return a * f; }
go_qmarginsf_t qmarginsf_div_float (go_qmarginsf_t a, double d)         { return a / d; }
go_qmarginsf_t qmarginsf_or        (go_qmarginsf_t a, go_qmarginsf_t b) { return a | b; }
go_qmarginsf_t qmarginsf_neg       (go_qmarginsf_t a)                   { return -a; }
go_qmarginsf_t qmarginsf_unary_plus(go_qmarginsf_t a)                   { return +a; }

/* --- QMarginsF: comparison --- */
bool qmarginsf_equal        (go_qmarginsf_t a, go_qmarginsf_t b) { return a == b; }
bool qmarginsf_equal_margins(go_qmarginsf_t a, go_qmargins_t b)  { return a == b; }

/* --- QMarginsF: conversions --- */
go_qmargins_t qmarginsf_to_margins(go_qmarginsf_t m) { return m.toMargins(); }

} // extern "C"