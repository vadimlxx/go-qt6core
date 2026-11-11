/* C ABI for Qt 6.11 QtCore/QLine and QLineF.
 *
 * QLine / QLineF are exposed as HANDLES (QLine* / QLineF*). Each handle owns a heap-allocated
 * C++ QLine / QLineF. Accessors/setters go directly to the C++ object. QPoint / QPointF I/O
 * uses the handle types from qpoint_c.h: inputs are dereferenced (*p), outputs are freshly
 * allocated (`new QPoint(...)`).
 *
 * Memory: every function returning go_qline_t / go_qlinef_t / go_qpoint_t / go_qpointf_t
 * allocates; caller owns it and must delete via qline_delete / qlinef_delete / qpoint_delete /
 * qpointf_delete.
 */
#ifndef GO_QT6CORE_QGEOMETRY_QLINE_C_H
#define GO_QT6CORE_QGEOMETRY_QLINE_C_H

#include <stdbool.h>
#include "qpoint_c.h"

#ifdef __cplusplus
class QLine;
class QLineF;
typedef QLine*  go_qline_t;
typedef QLineF* go_qlinef_t;
extern "C" {
#else
typedef struct go_qline_handle_s  go_qline_handle_s;
typedef struct go_qlinef_handle_s go_qlinef_handle_s;
typedef go_qline_handle_s*  go_qline_t;
typedef go_qlinef_handle_s* go_qlinef_t;
#endif

/* --- QLine: lifecycle --- */
go_qline_t qline_new(int x1, int y1, int x2, int y2);
go_qline_t qline_new_default(void);
go_qline_t qline_from_points(go_qpoint_t p1, go_qpoint_t p2);
go_qline_t qline_clone(go_qline_t l);
void       qline_delete(go_qline_t l);

/* --- QLine: accessors (P1/P2/Center allocate a fresh QPoint) --- */
bool        qline_is_null(go_qline_t l);
go_qpoint_t qline_p1(go_qline_t l);
go_qpoint_t qline_p2(go_qline_t l);
int         qline_x1(go_qline_t l);
int         qline_y1(go_qline_t l);
int         qline_x2(go_qline_t l);
int         qline_y2(go_qline_t l);
int         qline_dx(go_qline_t l);
int         qline_dy(go_qline_t l);
go_qpoint_t qline_center(go_qline_t l);

/* --- QLine: mutating (in place, void) --- */
void qline_translate(go_qline_t self, go_qpoint_t p);
void qline_translate_xy(go_qline_t self, int dx, int dy);
void qline_set_p1(go_qline_t self, go_qpoint_t p);
void qline_set_p2(go_qline_t self, go_qpoint_t p);
void qline_set_points(go_qline_t self, go_qpoint_t p1, go_qpoint_t p2);
void qline_set_line(go_qline_t self, int x1, int y1, int x2, int y2);

/* --- QLine: transforms (return new handle) --- */
go_qline_t qline_translated(go_qline_t l, go_qpoint_t p);
go_qline_t qline_translated_xy(go_qline_t l, int dx, int dy);

/* --- QLine: comparison --- */
bool qline_equal(go_qline_t a, go_qline_t b);

/* --- QLine: conversions --- */
go_qlinef_t qline_to_linef(go_qline_t l);

/* --- QLineF: lifecycle --- */
go_qlinef_t qlinef_new(double x1, double y1, double x2, double y2);
go_qlinef_t qlinef_new_default(void);
go_qlinef_t qlinef_from_points(go_qpointf_t p1, go_qpointf_t p2);
go_qlinef_t qlinef_from_line(go_qline_t l);
go_qlinef_t qlinef_from_polar(double length, double angle);
go_qlinef_t qlinef_clone(go_qlinef_t l);
void        qlinef_delete(go_qlinef_t l);

/* --- QLineF: accessors --- */
bool         qlinef_is_null(go_qlinef_t l);
go_qpointf_t qlinef_p1(go_qlinef_t l);
go_qpointf_t qlinef_p2(go_qlinef_t l);
double       qlinef_x1(go_qlinef_t l);
double       qlinef_y1(go_qlinef_t l);
double       qlinef_x2(go_qlinef_t l);
double       qlinef_y2(go_qlinef_t l);
double       qlinef_dx(go_qlinef_t l);
double       qlinef_dy(go_qlinef_t l);
double       qlinef_length(go_qlinef_t l);
double       qlinef_angle(go_qlinef_t l);
double       qlinef_angle_to(go_qlinef_t a, go_qlinef_t b);
go_qpointf_t qlinef_center(go_qlinef_t l);
go_qpointf_t qlinef_point_at(go_qlinef_t l, double t);

/* --- QLineF: mutating (in place, void) --- */
void qlinef_translate(go_qlinef_t self, go_qpointf_t p);
void qlinef_translate_xy(go_qlinef_t self, double dx, double dy);
void qlinef_set_p1(go_qlinef_t self, go_qpointf_t p);
void qlinef_set_p2(go_qlinef_t self, go_qpointf_t p);
void qlinef_set_points(go_qlinef_t self, go_qpointf_t p1, go_qpointf_t p2);
void qlinef_set_line(go_qlinef_t self, double x1, double y1, double x2, double y2);
void qlinef_set_length(go_qlinef_t self, double len);
void qlinef_set_angle(go_qlinef_t self, double angle);

/* --- QLineF: transforms (return new handle) --- */
go_qlinef_t qlinef_translated(go_qlinef_t l, go_qpointf_t p);
go_qlinef_t qlinef_translated_xy(go_qlinef_t l, double dx, double dy);
go_qlinef_t qlinef_unit_vector(go_qlinef_t l);
go_qlinef_t qlinef_normal_vector(go_qlinef_t l);

/* --- QLineF: intersection ---
 * out_point may be NULL. If intersection exists (Bounded or Unbounded), *out_point is
 * allocated (new QPointF). If no intersection, *out_point is set to NULL. */
int qlinef_intersects(go_qlinef_t a, go_qlinef_t b, go_qpointf_t* out_point);

/* --- QLineF: comparison --- */
bool qlinef_equal(go_qlinef_t a, go_qlinef_t b);
bool qlinef_equal_line(go_qlinef_t a, go_qline_t b);

/* --- QLineF: conversions --- */
go_qline_t qlinef_to_line(go_qlinef_t l);

#ifdef __cplusplus
}
#endif

#endif