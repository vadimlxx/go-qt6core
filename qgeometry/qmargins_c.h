/* C ABI for Qt 6.11 QtCore/QMargins and QMarginsF.
 *
 * QMargins / QMarginsF are exposed as HANDLES (QMargins* / QMarginsF*). Each handle owns a
 * heap-allocated C++ QMargins / QMarginsF. All accessors/setters go directly to the C++ object.
 *
 * Memory: every function returning go_qmargins_t / go_qmarginsf_t allocates; caller owns it and
 * must delete via qmargins_delete / qmarginsf_delete.
 */
#ifndef GO_QT6CORE_QGEOMETRY_QMARGINS_C_H
#define GO_QT6CORE_QGEOMETRY_QMARGINS_C_H

#include <stdbool.h>

#ifdef __cplusplus
class QMargins;
class QMarginsF;
typedef QMargins*  go_qmargins_t;
typedef QMarginsF* go_qmarginsf_t;
extern "C" {
#else
typedef struct go_qmargins_handle_s  go_qmargins_handle_s;
typedef struct go_qmarginsf_handle_s go_qmarginsf_handle_s;
typedef go_qmargins_handle_s*  go_qmargins_t;
typedef go_qmarginsf_handle_s* go_qmarginsf_t;
#endif

/* --- QMargins: lifecycle --- */
go_qmargins_t qmargins_new(int left, int top, int right, int bottom);
go_qmargins_t qmargins_new_default(void);
go_qmargins_t qmargins_clone(go_qmargins_t m);
void          qmargins_delete(go_qmargins_t m);

/* --- QMargins: accessors --- */
bool qmargins_is_null(go_qmargins_t m);
int  qmargins_left  (go_qmargins_t m);
int  qmargins_top   (go_qmargins_t m);
int  qmargins_right (go_qmargins_t m);
int  qmargins_bottom(go_qmargins_t m);

/* --- QMargins: mutating setters (in place, void) --- */
void qmargins_set_left  (go_qmargins_t self, int v);
void qmargins_set_top   (go_qmargins_t self, int v);
void qmargins_set_right (go_qmargins_t self, int v);
void qmargins_set_bottom(go_qmargins_t self, int v);

/* --- QMargins: arithmetic (return new handle) --- */
go_qmargins_t qmargins_add       (go_qmargins_t a, go_qmargins_t b);
go_qmargins_t qmargins_sub       (go_qmargins_t a, go_qmargins_t b);
go_qmargins_t qmargins_add_int   (go_qmargins_t a, int v);
go_qmargins_t qmargins_sub_int   (go_qmargins_t a, int v);
go_qmargins_t qmargins_mul_int   (go_qmargins_t a, int f);
go_qmargins_t qmargins_div_int   (go_qmargins_t a, int d);
go_qmargins_t qmargins_mul_float (go_qmargins_t a, double f);
go_qmargins_t qmargins_div_float (go_qmargins_t a, double d);
go_qmargins_t qmargins_or        (go_qmargins_t a, go_qmargins_t b);
go_qmargins_t qmargins_neg       (go_qmargins_t a);
go_qmargins_t qmargins_unary_plus(go_qmargins_t a);

/* --- QMargins: in-place assign (mutate self) --- */
void qmargins_add_assign       (go_qmargins_t self, go_qmargins_t other);
void qmargins_sub_assign       (go_qmargins_t self, go_qmargins_t other);
void qmargins_add_int_assign   (go_qmargins_t self, int v);
void qmargins_sub_int_assign   (go_qmargins_t self, int v);
void qmargins_mul_int_assign   (go_qmargins_t self, int f);
void qmargins_div_int_assign   (go_qmargins_t self, int d);
void qmargins_mul_float_assign (go_qmargins_t self, double f);
void qmargins_div_float_assign (go_qmargins_t self, double d);

/* --- QMargins: comparison --- */
bool qmargins_equal(go_qmargins_t a, go_qmargins_t b);

/* --- QMargins: conversions --- */
go_qmarginsf_t qmargins_to_marginsf(go_qmargins_t m);

/* --- QMarginsF: lifecycle --- */
go_qmarginsf_t qmarginsf_new(double left, double top, double right, double bottom);
go_qmarginsf_t qmarginsf_new_default(void);
go_qmarginsf_t qmarginsf_from_margins(go_qmargins_t m);
go_qmarginsf_t qmarginsf_clone(go_qmarginsf_t m);
void           qmarginsf_delete(go_qmarginsf_t m);

/* --- QMarginsF: accessors --- */
bool   qmarginsf_is_null(go_qmarginsf_t m);
double qmarginsf_left   (go_qmarginsf_t m);
double qmarginsf_top    (go_qmarginsf_t m);
double qmarginsf_right  (go_qmarginsf_t m);
double qmarginsf_bottom (go_qmarginsf_t m);

/* --- QMarginsF: mutating setters (in place, void) --- */
void qmarginsf_set_left  (go_qmarginsf_t self, double v);
void qmarginsf_set_top   (go_qmarginsf_t self, double v);
void qmarginsf_set_right (go_qmarginsf_t self, double v);
void qmarginsf_set_bottom(go_qmarginsf_t self, double v);

/* --- QMarginsF: arithmetic (return new handle) --- */
go_qmarginsf_t qmarginsf_add       (go_qmarginsf_t a, go_qmarginsf_t b);
go_qmarginsf_t qmarginsf_sub       (go_qmarginsf_t a, go_qmarginsf_t b);
go_qmarginsf_t qmarginsf_add_float (go_qmarginsf_t a, double v);
go_qmarginsf_t qmarginsf_sub_float (go_qmarginsf_t a, double v);
go_qmarginsf_t qmarginsf_mul_float (go_qmarginsf_t a, double f);
go_qmarginsf_t qmarginsf_div_float (go_qmarginsf_t a, double d);
go_qmarginsf_t qmarginsf_or        (go_qmarginsf_t a, go_qmarginsf_t b);
go_qmarginsf_t qmarginsf_neg       (go_qmarginsf_t a);
go_qmarginsf_t qmarginsf_unary_plus(go_qmarginsf_t a);

/* --- QMarginsF: in-place assign (mutate self) --- */
void qmarginsf_add_assign       (go_qmarginsf_t self, go_qmarginsf_t other);
void qmarginsf_sub_assign       (go_qmarginsf_t self, go_qmarginsf_t other);
void qmarginsf_add_float_assign (go_qmarginsf_t self, double v);
void qmarginsf_sub_float_assign (go_qmarginsf_t self, double v);
void qmarginsf_mul_float_assign (go_qmarginsf_t self, double f);
void qmarginsf_div_float_assign (go_qmarginsf_t self, double d);

/* --- QMarginsF: comparison --- */
bool qmarginsf_equal        (go_qmarginsf_t a, go_qmarginsf_t b);
bool qmarginsf_equal_margins(go_qmarginsf_t a, go_qmargins_t b);

/* --- QMarginsF: conversions --- */
go_qmargins_t qmarginsf_to_margins(go_qmarginsf_t m);

#ifdef __cplusplus
}
#endif

#endif