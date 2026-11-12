/* C ABI for Qt 6.11 QtCore/QMargins and QMarginsF.
 *
 * Both types are Q_RELOCATABLE_TYPE value types. To minimise casts in the
 * C++ implementation, we typedef go_qmargins_t / go_qmarginsf_t directly to
 * QMargins / QMarginsF in C++ mode. In C mode (used by cgo's preamble),
 * they fall back to flat POD structs with an identical memory layout.
 *
 * Safety: QMargins (4×int) and QMarginsF (4×qreal) are trivially copyable,
 * so both representations are ABI-compatible on all supported platforms.
 */
#ifndef GO_QT6CORE_QGEOMETRY_QMARGINS_C_H
#define GO_QT6CORE_QGEOMETRY_QMARGINS_C_H

#include <stdbool.h>

#ifdef __cplusplus
#include <QMargins>
typedef QMargins  go_qmargins_t;
typedef QMarginsF go_qmarginsf_t;
extern "C" {
#else
typedef struct go_qmargins_s {
    int left;
    int top;
    int right;
    int bottom;
} go_qmargins_t;

typedef struct go_qmarginsf_s {
    double left;
    double top;
    double right;
    double bottom;
} go_qmarginsf_t;
#endif

/* --- QMargins: construction --- */
go_qmargins_t qmargins_default(void);
go_qmargins_t qmargins_new(int left, int top, int right, int bottom);

/* --- QMargins: accessors --- */
bool qmargins_is_null(go_qmargins_t m);
int  qmargins_left  (go_qmargins_t m);
int  qmargins_top   (go_qmargins_t m);
int  qmargins_right (go_qmargins_t m);
int  qmargins_bottom(go_qmargins_t m);

/* --- QMargins: mutating setters (return updated value) --- */
go_qmargins_t qmargins_set_left  (go_qmargins_t m, int v);
go_qmargins_t qmargins_set_top   (go_qmargins_t m, int v);
go_qmargins_t qmargins_set_right (go_qmargins_t m, int v);
go_qmargins_t qmargins_set_bottom(go_qmargins_t m, int v);

/* --- QMargins: arithmetic --- */
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

/* --- QMargins: comparison --- */
bool qmargins_equal(go_qmargins_t a, go_qmargins_t b);

/* --- QMargins: conversions --- */
go_qmarginsf_t qmargins_to_marginsf(go_qmargins_t m);

/* --- QMarginsF: construction --- */
go_qmarginsf_t qmarginsf_default(void);
go_qmarginsf_t qmarginsf_new(double left, double top, double right, double bottom);
go_qmarginsf_t qmarginsf_from_margins(go_qmargins_t m);

/* --- QMarginsF: accessors --- */
bool   qmarginsf_is_null(go_qmarginsf_t m);
double qmarginsf_left   (go_qmarginsf_t m);
double qmarginsf_top    (go_qmarginsf_t m);
double qmarginsf_right  (go_qmarginsf_t m);
double qmarginsf_bottom (go_qmarginsf_t m);

/* --- QMarginsF: mutating setters --- */
go_qmarginsf_t qmarginsf_set_left  (go_qmarginsf_t m, double v);
go_qmarginsf_t qmarginsf_set_top   (go_qmarginsf_t m, double v);
go_qmarginsf_t qmarginsf_set_right (go_qmarginsf_t m, double v);
go_qmarginsf_t qmarginsf_set_bottom(go_qmarginsf_t m, double v);

/* --- QMarginsF: arithmetic --- */
go_qmarginsf_t qmarginsf_add       (go_qmarginsf_t a, go_qmarginsf_t b);
go_qmarginsf_t qmarginsf_sub       (go_qmarginsf_t a, go_qmarginsf_t b);
go_qmarginsf_t qmarginsf_add_float (go_qmarginsf_t a, double v);
go_qmarginsf_t qmarginsf_sub_float (go_qmarginsf_t a, double v);
go_qmarginsf_t qmarginsf_mul_float (go_qmarginsf_t a, double f);
go_qmarginsf_t qmarginsf_div_float (go_qmarginsf_t a, double d);
go_qmarginsf_t qmarginsf_or        (go_qmarginsf_t a, go_qmarginsf_t b);
go_qmarginsf_t qmarginsf_neg       (go_qmarginsf_t a);
go_qmarginsf_t qmarginsf_unary_plus(go_qmarginsf_t a);

/* --- QMarginsF: comparison --- */
bool qmarginsf_equal        (go_qmarginsf_t a, go_qmarginsf_t b);
bool qmarginsf_equal_margins(go_qmarginsf_t a, go_qmargins_t b);

/* --- QMarginsF: conversions --- */
go_qmargins_t qmarginsf_to_margins(go_qmarginsf_t m);

#ifdef __cplusplus
}
#endif

#endif