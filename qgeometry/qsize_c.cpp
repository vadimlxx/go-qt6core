/* Order: version check first, then real Qt headers.
 *
 * In C++ mode, go_qsize_t == QSize* and go_qsizef_t == QSizeF* (see qsize_c.h).
 * All functions dereference handles directly. No reinterpret_cast.
 */
#include "../internal/qt6core/check_qt_version.h"

#include <QSize>
#include <QMargins>
#include <QtCore/qnamespace.h>

#include "qsize_c.h"

namespace {
inline Qt::AspectRatioMode ARM(int mode) noexcept { return static_cast<Qt::AspectRatioMode>(mode); }
} // namespace

extern "C" {

/* --- QSize: lifecycle --- */
go_qsize_t qsize_new(int w, int h)      { return new QSize(w, h); }
go_qsize_t qsize_new_default()          { return new QSize();     }
go_qsize_t qsize_clone(go_qsize_t s)    { return new QSize(*s);   }
void       qsize_delete(go_qsize_t s)   { delete s;               }

/* --- QSize: accessors --- */
bool qsize_is_null (go_qsize_t s) { return s->isNull();  }
bool qsize_is_empty(go_qsize_t s) { return s->isEmpty(); }
bool qsize_is_valid(go_qsize_t s) { return s->isValid(); }
int  qsize_width   (go_qsize_t s) { return s->width();   }
int  qsize_height  (go_qsize_t s) { return s->height();  }

/* --- QSize: mutating setters (in place) --- */
void qsize_set_width (go_qsize_t self, int w) { self->setWidth(w);  }
void qsize_set_height(go_qsize_t self, int h) { self->setHeight(h); }
void qsize_transpose (go_qsize_t self)        { self->transpose();  }
void qsize_scale     (go_qsize_t self, int w, int h, int mode) {
    self->scale(w, h, ARM(mode));
}
void qsize_scale_size(go_qsize_t self, go_qsize_t other, int mode) {
    self->scale(*other, ARM(mode));
}

/* --- QSize: transforms (new handle) --- */
go_qsize_t qsize_transposed (go_qsize_t s)                         { return new QSize(s->transposed()); }
go_qsize_t qsize_scaled     (go_qsize_t s, int w, int h, int mode) { return new QSize(s->scaled(w, h, ARM(mode))); }
go_qsize_t qsize_scaled_size(go_qsize_t s, go_qsize_t o, int mode) { return new QSize(s->scaled(*o, ARM(mode))); }
go_qsize_t qsize_expanded_to(go_qsize_t s, go_qsize_t o)           { return new QSize(s->expandedTo(*o)); }
go_qsize_t qsize_bounded_to (go_qsize_t s, go_qsize_t o)           { return new QSize(s->boundedTo(*o)); }
go_qsize_t qsize_grown_by   (go_qsize_t s, go_qmargins_t m)        { return new QSize(s->grownBy(*m)); }
go_qsize_t qsize_shrunk_by  (go_qsize_t s, go_qmargins_t m)        { return new QSize(s->shrunkBy(*m)); }

/* --- QSize: arithmetic --- */
go_qsize_t qsize_add      (go_qsize_t a, go_qsize_t b) { return new QSize(*a + *b); }
go_qsize_t qsize_sub      (go_qsize_t a, go_qsize_t b) { return new QSize(*a - *b); }
go_qsize_t qsize_mul_float(go_qsize_t s, double c)     { return new QSize(*s * c);   }
go_qsize_t qsize_div_float(go_qsize_t s, double c)     { return new QSize(*s / c);   }

/* --- QSize: in-place assign --- */
void qsize_add_assign       (go_qsize_t self, go_qsize_t other) { *self += *other; }
void qsize_sub_assign       (go_qsize_t self, go_qsize_t other) { *self -= *other; }
void qsize_mul_float_assign (go_qsize_t self, double c)         { *self *= c;      }
void qsize_div_float_assign (go_qsize_t self, double c)         { *self /= c;      }

/* --- QSize: comparison --- */
bool qsize_equal(go_qsize_t a, go_qsize_t b) { return *a == *b; }

/* --- QSize: conversions --- */
go_qsizef_t qsize_to_sizef(go_qsize_t s) { return new QSizeF(s->toSizeF()); }

/* --- QSizeF: lifecycle --- */
go_qsizef_t qsizef_new(double w, double h)      { return new QSizeF(w, h); }
go_qsizef_t qsizef_new_default()                { return new QSizeF();     }
go_qsizef_t qsizef_from_size(go_qsize_t s)      { return new QSizeF(*s);   }
go_qsizef_t qsizef_clone(go_qsizef_t s)         { return new QSizeF(*s);   }
void        qsizef_delete(go_qsizef_t s)        { delete s;               }

/* --- QSizeF: accessors --- */
bool   qsizef_is_null (go_qsizef_t s) { return s->isNull();  }
bool   qsizef_is_empty(go_qsizef_t s) { return s->isEmpty(); }
bool   qsizef_is_valid(go_qsizef_t s) { return s->isValid(); }
double qsizef_width   (go_qsizef_t s) { return s->width();   }
double qsizef_height  (go_qsizef_t s) { return s->height();  }

/* --- QSizeF: mutating setters (in place) --- */
void qsizef_set_width (go_qsizef_t self, double w) { self->setWidth(w);  }
void qsizef_set_height(go_qsizef_t self, double h) { self->setHeight(h); }
void qsizef_transpose (go_qsizef_t self)           { self->transpose();  }
void qsizef_scale     (go_qsizef_t self, double w, double h, int mode) {
    self->scale(w, h, ARM(mode));
}
void qsizef_scale_size(go_qsizef_t self, go_qsizef_t other, int mode) {
    self->scale(*other, ARM(mode));
}

/* --- QSizeF: transforms (new handle) --- */
go_qsizef_t qsizef_transposed (go_qsizef_t s)                            { return new QSizeF(s->transposed()); }
go_qsizef_t qsizef_scaled     (go_qsizef_t s, double w, double h, int m) { return new QSizeF(s->scaled(w, h, ARM(m))); }
go_qsizef_t qsizef_scaled_size(go_qsizef_t s, go_qsizef_t o, int m)      { return new QSizeF(s->scaled(*o, ARM(m))); }
go_qsizef_t qsizef_expanded_to(go_qsizef_t s, go_qsizef_t o)             { return new QSizeF(s->expandedTo(*o)); }
go_qsizef_t qsizef_bounded_to (go_qsizef_t s, go_qsizef_t o)             { return new QSizeF(s->boundedTo(*o)); }
go_qsizef_t qsizef_grown_by   (go_qsizef_t s, go_qmarginsf_t m)          { return new QSizeF(s->grownBy(*m)); }
go_qsizef_t qsizef_shrunk_by  (go_qsizef_t s, go_qmarginsf_t m)          { return new QSizeF(s->shrunkBy(*m)); }

/* --- QSizeF: arithmetic --- */
go_qsizef_t qsizef_add      (go_qsizef_t a, go_qsizef_t b) { return new QSizeF(*a + *b); }
go_qsizef_t qsizef_sub      (go_qsizef_t a, go_qsizef_t b) { return new QSizeF(*a - *b); }
go_qsizef_t qsizef_mul_float(go_qsizef_t s, double c)      { return new QSizeF(*s * c);   }
go_qsizef_t qsizef_div_float(go_qsizef_t s, double c)      { return new QSizeF(*s / c);   }

/* --- QSizeF: in-place assign --- */
void qsizef_add_assign       (go_qsizef_t self, go_qsizef_t other) { *self += *other; }
void qsizef_sub_assign       (go_qsizef_t self, go_qsizef_t other) { *self -= *other; }
void qsizef_mul_float_assign (go_qsizef_t self, double c)          { *self *= c;      }
void qsizef_div_float_assign (go_qsizef_t self, double c)          { *self /= c;      }

/* --- QSizeF: comparison --- */
bool qsizef_equal     (go_qsizef_t a, go_qsizef_t b) { return *a == *b; }
bool qsizef_equal_size(go_qsizef_t a, go_qsize_t b)  { return *a == *b; }

/* --- QSizeF: conversions --- */
go_qsize_t qsizef_to_size(go_qsizef_t s) { return new QSize(s->toSize()); }

} // extern "C"