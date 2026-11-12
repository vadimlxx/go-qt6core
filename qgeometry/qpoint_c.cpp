/* Order: version check first, then real Qt headers.
 *
 * In C++ mode, go_qpoint_t == QPoint and go_qpointf_t == QPointF (see qpoint_c.h). Therefore every function
 * below operates on real Qt values directly, with no fromC/toC helpers and no reinterpret_cast.
 */
#include "../internal/qt6core/check_qt_version.h"

#include <QPoint>

#include "qpoint_c.h"

extern "C" {

/* --- QPoint: construction --- */
go_qpoint_t qpoint_default(void)             { return QPoint();    }
go_qpoint_t qpoint_new(int x, int y)         { return QPoint(x, y); }

/* --- QPoint: accessors --- */
bool qpoint_is_null         (go_qpoint_t p) { return p.isNull();          }
int  qpoint_x               (go_qpoint_t p) { return p.x();               }
int  qpoint_y               (go_qpoint_t p) { return p.y();               }
int  qpoint_manhattan_length(go_qpoint_t p) { return p.manhattanLength(); }

/* --- QPoint: mutating setters ---
 * Pass-by-value + copy-mutate-return: gives Go value semantics while letting C++ call the real setters. */
go_qpoint_t qpoint_set_x(go_qpoint_t p, int x) { p.setX(x); return p; }
go_qpoint_t qpoint_set_y(go_qpoint_t p, int y) { p.setY(y); return p; }

/* --- QPoint: transforms --- */
go_qpoint_t qpoint_transposed(go_qpoint_t p) { return p.transposed(); }

/* --- QPoint: static --- */
int qpoint_dot_product(go_qpoint_t a, go_qpoint_t b) { return QPoint::dotProduct(a, b); }

/* --- QPoint: arithmetic (uses Qt's real operators) --- */
go_qpoint_t qpoint_add       (go_qpoint_t a, go_qpoint_t b) { return a + b; }
go_qpoint_t qpoint_sub       (go_qpoint_t a, go_qpoint_t b) { return a - b; }
go_qpoint_t qpoint_mul_int   (go_qpoint_t p, int f)         { return p * f; }
go_qpoint_t qpoint_mul_float (go_qpoint_t p, double f)      { return p * f; }
go_qpoint_t qpoint_div_float (go_qpoint_t p, double d)      { return p / d; }
go_qpoint_t qpoint_neg       (go_qpoint_t p)                { return -p;   }
go_qpoint_t qpoint_unary_plus(go_qpoint_t p)                { return +p;   }

/* --- QPoint: comparison --- */
bool qpoint_equal(go_qpoint_t a, go_qpoint_t b) { return a == b; }

/* --- QPoint: conversions --- */
go_qpointf_t qpoint_to_pointf(go_qpoint_t p) { return p.toPointF(); }

/* --- QPointF: construction --- */
go_qpointf_t qpointf_default(void)             { return QPointF();    }
go_qpointf_t qpointf_new(double x, double y)   { return QPointF(x, y); }
go_qpointf_t qpointf_from_point(go_qpoint_t p) { return QPointF(p);    }

/* --- QPointF: accessors --- */
bool   qpointf_is_null         (go_qpointf_t p) { return p.isNull();          }
double qpointf_x               (go_qpointf_t p) { return p.x();               }
double qpointf_y               (go_qpointf_t p) { return p.y();               }
double qpointf_manhattan_length(go_qpointf_t p) { return p.manhattanLength(); }

/* --- QPointF: mutating setters --- */
go_qpointf_t qpointf_set_x(go_qpointf_t p, double x) { p.setX(x); return p; }
go_qpointf_t qpointf_set_y(go_qpointf_t p, double y) { p.setY(y); return p; }

/* --- QPointF: transforms --- */
go_qpointf_t qpointf_transposed(go_qpointf_t p) { return p.transposed(); }

/* --- QPointF: static --- */
double qpointf_dot_product(go_qpointf_t a, go_qpointf_t b) { return QPointF::dotProduct(a, b); }

/* --- QPointF: arithmetic --- */
go_qpointf_t qpointf_add       (go_qpointf_t a, go_qpointf_t b) { return a + b; }
go_qpointf_t qpointf_sub       (go_qpointf_t a, go_qpointf_t b) { return a - b; }
go_qpointf_t qpointf_mul_float (go_qpointf_t p, double f)       { return p * f; }
go_qpointf_t qpointf_div_float (go_qpointf_t p, double d)       { return p / d; }
go_qpointf_t qpointf_neg       (go_qpointf_t p)                 { return -p;   }
go_qpointf_t qpointf_unary_plus(go_qpointf_t p)                 { return +p;   }

/* --- QPointF: comparison --- */
bool qpointf_equal      (go_qpointf_t a, go_qpointf_t b) { return a == b; }
bool qpointf_equal_point(go_qpointf_t a, go_qpoint_t b)  { return a == b; }

/* --- QPointF: conversions --- */
go_qpoint_t qpointf_to_point(go_qpointf_t p) { return p.toPoint(); }

} // extern "C"