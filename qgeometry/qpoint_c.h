/* C ABI for Qt 6.11 QtCore/QPoint and QPointF.
 *
 * QPoint and QPointF are exposed as HANDLES (QPoint* / QPointF*). Each handle owns a
 * heap-allocated C++ QPoint / QPointF. All accessors/setters dereference the handle and go
 * directly to the C++ object — no mirrored Go copy is kept.
 *
 * Memory ownership: every function returning go_qpoint_t / go_qpointf_t allocates a new
 * object; the caller owns it and must eventually qpoint_delete / qpointf_delete.
 *
 * Note on embedded QPoints: QLine, QRect, etc. store QPoint by value in C++. A handle returned
 * from qline_p1() points to a fresh heap copy, NOT to QLine's internal storage. Mutations to
 * that handle do not affect the QLine.
 */
#ifndef GO_QT6CORE_QGEOMETRY_QPOINT_C_H
#define GO_QT6CORE_QGEOMETRY_QPOINT_C_H

#include <stdbool.h>

#ifdef __cplusplus
class QPoint;
class QPointF;
typedef QPoint*  go_qpoint_t;
typedef QPointF* go_qpointf_t;
extern "C" {
#else
typedef struct go_qpoint_handle_s  go_qpoint_handle_s;
typedef struct go_qpointf_handle_s go_qpointf_handle_s;
typedef go_qpoint_handle_s*  go_qpoint_t;
typedef go_qpointf_handle_s* go_qpointf_t;
#endif

/* --- QPoint: lifecycle --- */
go_qpoint_t qpoint_new(int x, int y);
go_qpoint_t qpoint_new_default(void);
go_qpoint_t qpoint_clone(go_qpoint_t p);
void        qpoint_delete(go_qpoint_t p);

/* --- QPoint: accessors --- */
bool qpoint_is_null(go_qpoint_t p);
int  qpoint_x(go_qpoint_t p);
int  qpoint_y(go_qpoint_t p);
int  qpoint_manhattan_length(go_qpoint_t p);

/* --- QPoint: mutating setters (in place) --- */
void qpoint_set_x(go_qpoint_t p, int x);
void qpoint_set_y(go_qpoint_t p, int y);

/* --- QPoint: transforms (return new handles) --- */
go_qpoint_t qpoint_transposed(go_qpoint_t p);

/* --- QPoint: static --- */
int qpoint_dot_product(go_qpoint_t a, go_qpoint_t b);

/* --- QPoint: arithmetic (return new handles) --- */
go_qpoint_t qpoint_add      (go_qpoint_t a, go_qpoint_t b);
go_qpoint_t qpoint_sub      (go_qpoint_t a, go_qpoint_t b);
go_qpoint_t qpoint_mul_int  (go_qpoint_t p, int f);
go_qpoint_t qpoint_mul_float(go_qpoint_t p, double f);
go_qpoint_t qpoint_div_float(go_qpoint_t p, double d);
go_qpoint_t qpoint_neg      (go_qpoint_t p);
go_qpoint_t qpoint_unary_plus(go_qpoint_t p);

/* --- QPoint: in-place assign (mutate self) --- */
void qpoint_add_assign       (go_qpoint_t self, go_qpoint_t other);
void qpoint_sub_assign       (go_qpoint_t self, go_qpoint_t other);
void qpoint_mul_int_assign   (go_qpoint_t self, int f);
void qpoint_mul_float_assign (go_qpoint_t self, double f);
void qpoint_div_float_assign (go_qpoint_t self, double d);

/* --- QPoint: comparison --- */
bool qpoint_equal(go_qpoint_t a, go_qpoint_t b);

/* --- QPoint: conversions --- */
go_qpointf_t qpoint_to_pointf(go_qpoint_t p);

/* --- QPointF: lifecycle --- */
go_qpointf_t qpointf_new(double x, double y);
go_qpointf_t qpointf_new_default(void);
go_qpointf_t qpointf_clone(go_qpointf_t p);
go_qpointf_t qpointf_from_point(go_qpoint_t p);
void         qpointf_delete(go_qpointf_t p);

/* --- QPointF: accessors --- */
bool   qpointf_is_null(go_qpointf_t p);
double qpointf_x(go_qpointf_t p);
double qpointf_y(go_qpointf_t p);
double qpointf_manhattan_length(go_qpointf_t p);

/* --- QPointF: mutating setters --- */
void qpointf_set_x(go_qpointf_t p, double x);
void qpointf_set_y(go_qpointf_t p, double y);

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

/* --- QPointF: in-place assign --- */
void qpointf_add_assign       (go_qpointf_t self, go_qpointf_t other);
void qpointf_sub_assign       (go_qpointf_t self, go_qpointf_t other);
void qpointf_mul_float_assign (go_qpointf_t self, double f);
void qpointf_div_float_assign (go_qpointf_t self, double d);

/* --- QPointF: comparison --- */
bool qpointf_equal      (go_qpointf_t a, go_qpointf_t b);
bool qpointf_equal_point(go_qpointf_t a, go_qpoint_t b);

/* --- QPointF: conversions --- */
go_qpoint_t qpointf_to_point(go_qpointf_t p);

#ifdef __cplusplus
}
#endif

#endif