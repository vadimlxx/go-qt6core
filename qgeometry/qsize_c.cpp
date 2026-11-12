/* Order: version check first, then real Qt headers.
 *
 * In C++ mode, go_qsize_t == QSize and go_qsizef_t == QSizeF (see qsize_c.h). Therefore every function
 * below operates on real Qt values directly, with no fromC/toC helpers and no reinterpret_cast.
 */
#include "../internal/qt6core/check_qt_version.h"

#include <QSize>
#include <QtCore/qnamespace.h>

#include "qsize_c.h"

namespace {
inline Qt::AspectRatioMode ARM(int mode) noexcept { return static_cast<Qt::AspectRatioMode>(mode); }
} // namespace

extern "C" {

/* --- QSize: construction --- */
go_qsize_t qsize_default(void)     { return QSize(); }
go_qsize_t qsize_new(int w, int h) { return QSize(w, h); }

/* --- QSize: accessors --- */
bool qsize_is_null (go_qsize_t s) { return s.isNull();  }
bool qsize_is_empty(go_qsize_t s) { return s.isEmpty(); }
bool qsize_is_valid(go_qsize_t s) { return s.isValid(); }
int  qsize_width   (go_qsize_t s) { return s.width();   }
int  qsize_height  (go_qsize_t s) { return s.height();  }

/* --- QSize: mutating setters ---
 * Pass-by-value + copy-mutate-return: gives Go value semantics while letting C++ call the real setters. */
go_qsize_t qsize_set_width (go_qsize_t s, int w) { s.setWidth(w);  return s; }
go_qsize_t qsize_set_height(go_qsize_t s, int h) { s.setHeight(h); return s; }

/* --- QSize: transforms (delegate to Qt's real methods) --- */
go_qsize_t qsize_transposed (go_qsize_t s)                         { return s.transposed(); }
go_qsize_t qsize_scaled     (go_qsize_t s, int w, int h, int mode) { return s.scaled(w, h, ARM(mode)); }
go_qsize_t qsize_scaled_size(go_qsize_t s, go_qsize_t o, int mode) { return s.scaled(o, ARM(mode)); }
go_qsize_t qsize_expanded_to(go_qsize_t s, go_qsize_t o)           { return s.expandedTo(o); }
go_qsize_t qsize_bounded_to (go_qsize_t s, go_qsize_t o)           { return s.boundedTo(o); }
go_qsize_t qsize_grown_by   (go_qsize_t s, go_qmargins_t m)        { return s.grownBy(m); }
go_qsize_t qsize_shrunk_by  (go_qsize_t s, go_qmargins_t m)        { return s.shrunkBy(m); }

/* --- QSize: arithmetic (uses Qt's real operators) --- */
go_qsize_t qsize_add      (go_qsize_t a, go_qsize_t b) { return a + b; }
go_qsize_t qsize_sub      (go_qsize_t a, go_qsize_t b) { return a - b; }
go_qsize_t qsize_mul_float(go_qsize_t s, double c)     { return s * c; }
go_qsize_t qsize_div_float(go_qsize_t s, double c)     { return s / c; }

/* --- QSize: comparison --- */
bool qsize_equal(go_qsize_t a, go_qsize_t b) { return a == b; }

/* --- QSize: conversions --- */
go_qsizef_t qsize_to_sizef(go_qsize_t s) { return s.toSizeF(); }

/* --- QSizeF: construction --- */
go_qsizef_t qsizef_default(void)           { return QSizeF(); }
go_qsizef_t qsizef_new(double w, double h) { return QSizeF(w, h); }
go_qsizef_t qsizef_from_size(go_qsize_t s) { return QSizeF(s); }

/* --- QSizeF: accessors --- */
bool   qsizef_is_null (go_qsizef_t s) { return s.isNull();  }
bool   qsizef_is_empty(go_qsizef_t s) { return s.isEmpty(); }
bool   qsizef_is_valid(go_qsizef_t s) { return s.isValid(); }
double qsizef_width   (go_qsizef_t s) { return s.width();   }
double qsizef_height  (go_qsizef_t s) { return s.height();  }

/* --- QSizeF: mutating setters --- */
go_qsizef_t qsizef_set_width (go_qsizef_t s, double w) { s.setWidth(w);  return s; }
go_qsizef_t qsizef_set_height(go_qsizef_t s, double h) { s.setHeight(h); return s; }

/* --- QSizeF: transforms --- */
go_qsizef_t qsizef_transposed (go_qsizef_t s)                            { return s.transposed(); }
go_qsizef_t qsizef_scaled     (go_qsizef_t s, double w, double h, int m) { return s.scaled(w, h, ARM(m)); }
go_qsizef_t qsizef_scaled_size(go_qsizef_t s, go_qsizef_t o, int m)      { return s.scaled(o, ARM(m)); }
go_qsizef_t qsizef_expanded_to(go_qsizef_t s, go_qsizef_t o)             { return s.expandedTo(o); }
go_qsizef_t qsizef_bounded_to (go_qsizef_t s, go_qsizef_t o)             { return s.boundedTo(o); }
go_qsizef_t qsizef_grown_by   (go_qsizef_t s, go_qmarginsf_t m)          { return s.grownBy(m); }
go_qsizef_t qsizef_shrunk_by  (go_qsizef_t s, go_qmarginsf_t m)          { return s.shrunkBy(m); }

/* --- QSizeF: arithmetic --- */
go_qsizef_t qsizef_add      (go_qsizef_t a, go_qsizef_t b) { return a + b; }
go_qsizef_t qsizef_sub      (go_qsizef_t a, go_qsizef_t b) { return a - b; }
go_qsizef_t qsizef_mul_float(go_qsizef_t s, double c)      { return s * c; }
go_qsizef_t qsizef_div_float(go_qsizef_t s, double c)      { return s / c; }

/* --- QSizeF: comparison --- */
bool qsizef_equal     (go_qsizef_t a, go_qsizef_t b) { return a == b; }
bool qsizef_equal_size(go_qsizef_t a, go_qsize_t b)  { return a == b; }

/* --- QSizeF: conversions --- */
go_qsize_t qsizef_to_size(go_qsizef_t s) { return s.toSize(); }

} // extern "C"