/* Order: version check first, then real Qt headers. */
#include "../internal/qt6core/check_qt_version.h"

#include <QPoint>

#include "qpoint_c.h"

extern "C" {

/* --- QPoint: lifecycle --- */
go_qpoint_t qpoint_new(int x, int y)        { return new QPoint(x, y); }
go_qpoint_t qpoint_new_default()            { return new QPoint();    }
go_qpoint_t qpoint_clone(go_qpoint_t p)     { return new QPoint(*p);  }
void        qpoint_delete(go_qpoint_t p)    { delete p;               }

/* --- QPoint: accessors --- */
bool qpoint_is_null(go_qpoint_t p)          { return p->isNull();          }
int  qpoint_x(go_qpoint_t p)                { return p->x();               }
int  qpoint_y(go_qpoint_t p)                { return p->y();               }
int  qpoint_manhattan_length(go_qpoint_t p) { return p->manhattanLength(); }

/* --- QPoint: mutating setters --- */
void qpoint_set_x(go_qpoint_t p, int x)     { p->setX(x); }
void qpoint_set_y(go_qpoint_t p, int y)     { p->setY(y); }

/* --- QPoint: transforms --- */
go_qpoint_t qpoint_transposed(go_qpoint_t p) { return new QPoint(p->transposed()); }

/* --- QPoint: static --- */
int qpoint_dot_product(go_qpoint_t a, go_qpoint_t b) { return QPoint::dotProduct(*a, *b); }

/* --- QPoint: arithmetic --- */
go_qpoint_t qpoint_add       (go_qpoint_t a, go_qpoint_t b) { return new QPoint(*a + *b); }
go_qpoint_t qpoint_sub       (go_qpoint_t a, go_qpoint_t b) { return new QPoint(*a - *b); }
go_qpoint_t qpoint_mul_int   (go_qpoint_t p, int f)         { return new QPoint(*p * f);   }
go_qpoint_t qpoint_mul_float (go_qpoint_t p, double f)      { return new QPoint(*p * f);   }
go_qpoint_t qpoint_div_float (go_qpoint_t p, double d)      { return new QPoint(*p / d);   }
go_qpoint_t qpoint_neg       (go_qpoint_t p)                { return new QPoint(-*p);      }
go_qpoint_t qpoint_unary_plus(go_qpoint_t p)                { return new QPoint(+*p);      }

/* --- QPoint: in-place assign --- */
void qpoint_add_assign       (go_qpoint_t self, go_qpoint_t other) { *self += *other; }
void qpoint_sub_assign       (go_qpoint_t self, go_qpoint_t other) { *self -= *other; }
void qpoint_mul_int_assign   (go_qpoint_t self, int f)             { *self *= f;      }
void qpoint_mul_float_assign (go_qpoint_t self, double f)          { *self *= f;      }
void qpoint_div_float_assign (go_qpoint_t self, double d)          { *self /= d;      }

/* --- QPoint: comparison --- */
bool qpoint_equal(go_qpoint_t a, go_qpoint_t b) { return *a == *b; }

/* --- QPoint: conversions --- */
go_qpointf_t qpoint_to_pointf(go_qpoint_t p) { return new QPointF(*p); }

/* --- QPointF: lifecycle --- */
go_qpointf_t qpointf_new(double x, double y)    { return new QPointF(x, y); }
go_qpointf_t qpointf_new_default()              { return new QPointF();    }
go_qpointf_t qpointf_clone(go_qpointf_t p)      { return new QPointF(*p);  }
go_qpointf_t qpointf_from_point(go_qpoint_t p)  { return new QPointF(*p);  }
void         qpointf_delete(go_qpointf_t p)     { delete p;                }

/* --- QPointF: accessors --- */
bool   qpointf_is_null(go_qpointf_t p)          { return p->isNull();          }
double qpointf_x(go_qpointf_t p)                { return p->x();               }
double qpointf_y(go_qpointf_t p)                { return p->y();               }
double qpointf_manhattan_length(go_qpointf_t p) { return p->manhattanLength(); }

/* --- QPointF: mutating setters --- */
void qpointf_set_x(go_qpointf_t p, double x)    { p->setX(x); }
void qpointf_set_y(go_qpointf_t p, double y)    { p->setY(y); }

/* --- QPointF: transforms --- */
go_qpointf_t qpointf_transposed(go_qpointf_t p) { return new QPointF(p->transposed()); }

/* --- QPointF: static --- */
double qpointf_dot_product(go_qpointf_t a, go_qpointf_t b) { return QPointF::dotProduct(*a, *b); }

/* --- QPointF: arithmetic --- */
go_qpointf_t qpointf_add       (go_qpointf_t a, go_qpointf_t b) { return new QPointF(*a + *b); }
go_qpointf_t qpointf_sub       (go_qpointf_t a, go_qpointf_t b) { return new QPointF(*a - *b); }
go_qpointf_t qpointf_mul_float (go_qpointf_t p, double f)       { return new QPointF(*p * f);   }
go_qpointf_t qpointf_div_float (go_qpointf_t p, double d)       { return new QPointF(*p / d);   }
go_qpointf_t qpointf_neg       (go_qpointf_t p)                 { return new QPointF(-*p);      }
go_qpointf_t qpointf_unary_plus(go_qpointf_t p)                 { return new QPointF(+*p);      }

/* --- QPointF: in-place assign --- */
void qpointf_add_assign       (go_qpointf_t self, go_qpointf_t other) { *self += *other; }
void qpointf_sub_assign       (go_qpointf_t self, go_qpointf_t other) { *self -= *other; }
void qpointf_mul_float_assign (go_qpointf_t self, double f)           { *self *= f;      }
void qpointf_div_float_assign (go_qpointf_t self, double d)           { *self /= d;      }

/* --- QPointF: comparison --- */
bool qpointf_equal      (go_qpointf_t a, go_qpointf_t b) { return *a == *b; }
bool qpointf_equal_point(go_qpointf_t a, go_qpoint_t b)  { return *a == *b; }

/* --- QPointF: conversions --- */
go_qpoint_t qpointf_to_point(go_qpointf_t p) { return new QPoint(p->toPoint()); }

} // extern "C"