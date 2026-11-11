/* Order: version check first, then real Qt headers.
 *
 * In C++ mode, go_qmargins_t == QMargins* and go_qmarginsf_t == QMarginsF* (see qmargins_c.h).
 * All functions dereference handles directly. No reinterpret_cast.
 */
#include "../internal/qt6core/check_qt_version.h"

#include <QMargins>

#include "qmargins_c.h"

extern "C" {

/* --- QMargins: lifecycle --- */
go_qmargins_t qmargins_new(int l, int t, int r, int b) { return new QMargins(l, t, r, b); }
go_qmargins_t qmargins_new_default()                   { return new QMargins(); }
go_qmargins_t qmargins_clone(go_qmargins_t m)          { return new QMargins(*m); }
void          qmargins_delete(go_qmargins_t m)         { delete m; }

/* --- QMargins: accessors --- */
bool qmargins_is_null(go_qmargins_t m) { return m->isNull(); }
int  qmargins_left  (go_qmargins_t m)  { return m->left();   }
int  qmargins_top   (go_qmargins_t m)  { return m->top();    }
int  qmargins_right (go_qmargins_t m)  { return m->right();  }
int  qmargins_bottom(go_qmargins_t m)  { return m->bottom(); }

/* --- QMargins: mutating setters --- */
void qmargins_set_left  (go_qmargins_t m, int v) { m->setLeft(v);   }
void qmargins_set_top   (go_qmargins_t m, int v) { m->setTop(v);    }
void qmargins_set_right (go_qmargins_t m, int v) { m->setRight(v);  }
void qmargins_set_bottom(go_qmargins_t m, int v) { m->setBottom(v); }

/* --- QMargins: arithmetic --- */
go_qmargins_t qmargins_add       (go_qmargins_t a, go_qmargins_t b) { return new QMargins(*a + *b); }
go_qmargins_t qmargins_sub       (go_qmargins_t a, go_qmargins_t b) { return new QMargins(*a - *b); }
go_qmargins_t qmargins_add_int   (go_qmargins_t a, int v)           { return new QMargins(*a + v);   }
go_qmargins_t qmargins_sub_int   (go_qmargins_t a, int v)           { return new QMargins(*a - v);   }
go_qmargins_t qmargins_mul_int   (go_qmargins_t a, int f)           { return new QMargins(*a * f);   }
go_qmargins_t qmargins_div_int   (go_qmargins_t a, int d)           { return new QMargins(*a / d);   }
go_qmargins_t qmargins_mul_float (go_qmargins_t a, double f)        { return new QMargins(*a * f);   }
go_qmargins_t qmargins_div_float (go_qmargins_t a, double d)        { return new QMargins(*a / d);   }
go_qmargins_t qmargins_or        (go_qmargins_t a, go_qmargins_t b) { return new QMargins(*a | *b); }
go_qmargins_t qmargins_neg       (go_qmargins_t a)                  { return new QMargins(-*a);      }
go_qmargins_t qmargins_unary_plus(go_qmargins_t a)                  { return new QMargins(+*a);      }

/* --- QMargins: in-place assign --- */
void qmargins_add_assign       (go_qmargins_t self, go_qmargins_t other) { *self += *other; }
void qmargins_sub_assign       (go_qmargins_t self, go_qmargins_t other) { *self -= *other; }
void qmargins_add_int_assign   (go_qmargins_t self, int v)               { *self += v;      }
void qmargins_sub_int_assign   (go_qmargins_t self, int v)               { *self -= v;      }
void qmargins_mul_int_assign   (go_qmargins_t self, int f)               { *self *= f;      }
void qmargins_div_int_assign   (go_qmargins_t self, int d)               { *self /= d;      }
void qmargins_mul_float_assign (go_qmargins_t self, double f)            { *self *= f;      }
void qmargins_div_float_assign (go_qmargins_t self, double d)            { *self /= d;      }

/* --- QMargins: comparison --- */
bool qmargins_equal(go_qmargins_t a, go_qmargins_t b) { return *a == *b; }

/* --- QMargins: conversions --- */
go_qmarginsf_t qmargins_to_marginsf(go_qmargins_t m) { return new QMarginsF(m->toMarginsF()); }

/* --- QMarginsF: lifecycle --- */
go_qmarginsf_t qmarginsf_new(double l, double t, double r, double b) { return new QMarginsF(l, t, r, b); }
go_qmarginsf_t qmarginsf_new_default()                               { return new QMarginsF(); }
go_qmarginsf_t qmarginsf_from_margins(go_qmargins_t m)               { return new QMarginsF(*m); }
go_qmarginsf_t qmarginsf_clone(go_qmarginsf_t m)                     { return new QMarginsF(*m); }
void           qmarginsf_delete(go_qmarginsf_t m)                    { delete m; }

/* --- QMarginsF: accessors --- */
bool   qmarginsf_is_null(go_qmarginsf_t m) { return m->isNull();  }
double qmarginsf_left   (go_qmarginsf_t m) { return m->left();    }
double qmarginsf_top    (go_qmarginsf_t m) { return m->top();     }
double qmarginsf_right  (go_qmarginsf_t m) { return m->right();   }
double qmarginsf_bottom (go_qmarginsf_t m) { return m->bottom();  }

/* --- QMarginsF: mutating setters --- */
void qmarginsf_set_left  (go_qmarginsf_t m, double v) { m->setLeft(v);   }
void qmarginsf_set_top   (go_qmarginsf_t m, double v) { m->setTop(v);    }
void qmarginsf_set_right (go_qmarginsf_t m, double v) { m->setRight(v);  }
void qmarginsf_set_bottom(go_qmarginsf_t m, double v) { m->setBottom(v); }

/* --- QMarginsF: arithmetic --- */
go_qmarginsf_t qmarginsf_add       (go_qmarginsf_t a, go_qmarginsf_t b) { return new QMarginsF(*a + *b); }
go_qmarginsf_t qmarginsf_sub       (go_qmarginsf_t a, go_qmarginsf_t b) { return new QMarginsF(*a - *b); }
go_qmarginsf_t qmarginsf_add_float (go_qmarginsf_t a, double v)         { return new QMarginsF(*a + v);   }
go_qmarginsf_t qmarginsf_sub_float (go_qmarginsf_t a, double v)         { return new QMarginsF(*a - v);   }
go_qmarginsf_t qmarginsf_mul_float (go_qmarginsf_t a, double f)         { return new QMarginsF(*a * f);   }
go_qmarginsf_t qmarginsf_div_float (go_qmarginsf_t a, double d)         { return new QMarginsF(*a / d);   }
go_qmarginsf_t qmarginsf_or        (go_qmarginsf_t a, go_qmarginsf_t b) { return new QMarginsF(*a | *b); }
go_qmarginsf_t qmarginsf_neg       (go_qmarginsf_t a)                   { return new QMarginsF(-*a);      }
go_qmarginsf_t qmarginsf_unary_plus(go_qmarginsf_t a)                   { return new QMarginsF(+*a);      }

/* --- QMarginsF: in-place assign --- */
void qmarginsf_add_assign       (go_qmarginsf_t self, go_qmarginsf_t other) { *self += *other; }
void qmarginsf_sub_assign       (go_qmarginsf_t self, go_qmarginsf_t other) { *self -= *other; }
void qmarginsf_add_float_assign (go_qmarginsf_t self, double v)             { *self += v;      }
void qmarginsf_sub_float_assign (go_qmarginsf_t self, double v)             { *self -= v;      }
void qmarginsf_mul_float_assign (go_qmarginsf_t self, double f)             { *self *= f;      }
void qmarginsf_div_float_assign (go_qmarginsf_t self, double d)             { *self /= d;      }

/* --- QMarginsF: comparison --- */
bool qmarginsf_equal        (go_qmarginsf_t a, go_qmarginsf_t b) { return *a == *b; }
bool qmarginsf_equal_margins(go_qmarginsf_t a, go_qmargins_t b)  { return *a == *b; }

/* --- QMarginsF: conversions --- */
go_qmargins_t qmarginsf_to_margins(go_qmarginsf_t m) { return new QMargins(m->toMargins()); }

} // extern "C"