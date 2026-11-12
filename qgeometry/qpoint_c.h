/* C ABI for Qt 6.11 QtCore/QPoint and QPointF.
 *
 * Both types are Q_PRIMITIVE_TYPE value types. To minimise casts in the C++ implementation, we typedef
 * go_qpoint_t / go_qpointf_t directly to QPoint / QPointF in C++ mode. In C mode (used by cgo's preamble),
 * they fall back to flat POD structs with an identical memory layout.
 *
 * Safety: QPoint (2×int = 8 B) and QPointF (2×qreal = 16 B) are trivially copyable, so both representations
 * are ABI-compatible on all supported platforms.
 */
#ifndef GO_QT6CORE_QGEOMETRY_QPOINT_C_H
#define GO_QT6CORE_QGEOMETRY_QPOINT_C_H

#include <stdbool.h>

#ifdef __cplusplus
#include <QPoint>
typedef QPoint  go_qpoint_t;
typedef QPointF go_qpointf_t;
extern "C" {
#else
typedef struct go_qpoint_s {
    int x;
    int y;
} go_qpoint_t;

typedef struct go_qpointf_s {
    double x;
    double y;
} go_qpointf_t;
#endif

/* --- QPoint: construction --- */
go_qpoint_t qpoint_default(void);
go_qpoint_t qpoint_new(int x, int y);

/* --- QPoint: accessors --- */
bool qpoint_is_null         (go_qpoint_t p);
int  qpoint_x               (go_qpoint_t p);
int  qpoint_y               (go_qpoint_t p);
int  qpoint_manhattan_length(go_qpoint_t p);

/* --- QPoint: mutating setters (return updated value) --- */
go_qpoint_t qpoint_set_x(go_qpoint_t p, int x);
go_qpoint_t qpoint_set_y(go_qpoint_t p, int y);

/* --- QPoint: transforms --- */
go_qpoint_t qpoint_transposed(go_qpoint_t p);

/* --- QPoint: static --- */
int qpoint_dot_product(go_qpoint_t a, go_qpoint_t b);

/* --- QPoint: arithmetic --- */
go_qpoint_t qpoint_add      (go_qpoint_t a, go_qpoint_t b);
go_qpoint_t qpoint_sub      (go_qpoint_t a, go_qpoint_t b);
go_qpoint_t qpoint_mul_int  (go_qpoint_t p, int f);
go_qpoint_t qpoint_mul_float(go_qpoint_t p, double f);
go_qpoint_t qpoint_div_float(go_qpoint_t p, double d);
go_qpoint_t qpoint_neg      (go_qpoint_t p);
go_qpoint_t qpoint_unary_plus(go_qpoint_t p);

/* --- QPoint: comparison --- */
bool qpoint_equal(go_qpoint_t a, go_qpoint_t b);

/* --- QPoint: conversions --- */
go_qpointf_t qpoint_to_pointf(go_qpoint_t p);

/* --- QPointF: construction --- */
go_qpointf_t qpointf_default(void);
go_qpointf_t qpointf_new(double x, double y);
go_qpointf_t qpointf_from_point(go_qpoint_t p);

/* --- QPointF: accessors --- */
bool   qpointf_is_null         (go_qpointf_t p);
double qpointf_x               (go_qpointf_t p);
double qpointf_y               (go_qpointf_t p);
double qpointf_manhattan_length(go_qpointf_t p);

/* --- QPointF: mutating setters --- */
go_qpointf_t qpointf_set_x(go_qpointf_t p, double x);
go_qpointf_t qpointf_set_y(go_qpointf_t p, double y);

/* --- QPointF: transforms --- */
go_qpointf_t qpointf_transposed(go_qpointf_t p);

/* --- QPointF: static --- */
double qpointf_dot_product(go_qpointf_t a, go_qpointf_t b);

/* --- QPointF: arithmetic --- */
go_qpointf_t qpointf_add      (go_qpointf_t a, go_qpointf_t b);
go_qpointf_t qpointf_sub      (go_qpointf_t a, go_qpointf_t b);
go_qpointf_t qpointf_mul_float(go_qpointf_t p, double f);
go_qpointf_t qpointf_div_float(go_qpointf_t p, double d);
go_qpointf_t qpointf_neg      (go_qpointf_t p);
go_qpointf_t qpointf_unary_plus(go_qpointf_t p);

/* --- QPointF: comparison --- */
bool qpointf_equal      (go_qpointf_t a, go_qpointf_t b);
bool qpointf_equal_point(go_qpointf_t a, go_qpoint_t b);

/* --- QPointF: conversions --- */
go_qpoint_t qpointf_to_point(go_qpointf_t p);

#ifdef __cplusplus
}
#endif

#endif