/* C ABI for Qt 6.11 QtCore/QLine and QLineF.
 *
 * Both types are Q_PRIMITIVE_TYPE value types. To minimise casts in the C++ implementation,
 * we typedef go_qline_t / go_qlinef_t directly to QLine / QLineF in C++ mode. In C mode
 * (used by cgo's preamble), they fall back to flat POD structs with an identical memory
 * layout: QLine is 2×QPoint = 4×int, QLineF is 2×QPointF = 4×qreal.
 *
 * Safety: QLine (16 B) and QLineF (32 B) are trivially copyable, so both representations
 * are ABI-compatible on all supported platforms.
 */
#ifndef GO_QT6CORE_QGEOMETRY_QLINE_C_H
#define GO_QT6CORE_QGEOMETRY_QLINE_C_H

#include <stdbool.h>
#include "qpoint_c.h"

#ifdef __cplusplus
#include <QLine>
typedef QLine  go_qline_t;
typedef QLineF go_qlinef_t;
extern "C" {
#else
typedef struct go_qline_s {
    int x1;
    int y1;
    int x2;
    int y2;
} go_qline_t;

typedef struct go_qlinef_s {
    double x1;
    double y1;
    double x2;
    double y2;
} go_qlinef_t;
#endif

/* --- QLine: construction --- */
go_qline_t qline_default(void);
go_qline_t qline_new(int x1, int y1, int x2, int y2);
go_qline_t qline_from_points(go_qpoint_t p1, go_qpoint_t p2);

/* --- QLine: accessors (values; P1/P2/Center return go_qpoint_t by value) --- */
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

/* --- QLine: mutating setters (return updated value) --- */
go_qline_t qline_translate(go_qline_t l, go_qpoint_t p);
go_qline_t qline_translate_xy(go_qline_t l, int dx, int dy);
go_qline_t qline_set_p1(go_qline_t l, go_qpoint_t p);
go_qline_t qline_set_p2(go_qline_t l, go_qpoint_t p);
go_qline_t qline_set_points(go_qline_t l, go_qpoint_t p1, go_qpoint_t p2);
go_qline_t qline_set_line(go_qline_t l, int x1, int y1, int x2, int y2);

/* --- QLine: transforms (return new value) --- */
go_qline_t qline_translated(go_qline_t l, go_qpoint_t p);
go_qline_t qline_translated_xy(go_qline_t l, int dx, int dy);

/* --- QLine: comparison --- */
bool qline_equal(go_qline_t a, go_qline_t b);

/* --- QLine: conversions --- */
go_qlinef_t qline_to_linef(go_qline_t l);

/* --- QLineF: construction --- */
go_qlinef_t qlinef_default(void);
go_qlinef_t qlinef_new(double x1, double y1, double x2, double y2);
go_qlinef_t qlinef_from_points(go_qpointf_t p1, go_qpointf_t p2);
go_qlinef_t qlinef_from_line(go_qline_t l);
go_qlinef_t qlinef_from_polar(double length, double angle);

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

/* --- QLineF: mutating setters (return updated value) --- */
go_qlinef_t qlinef_translate(go_qlinef_t l, go_qpointf_t p);
go_qlinef_t qlinef_translate_xy(go_qlinef_t l, double dx, double dy);
go_qlinef_t qlinef_set_p1(go_qlinef_t l, go_qpointf_t p);
go_qlinef_t qlinef_set_p2(go_qlinef_t l, go_qpointf_t p);
go_qlinef_t qlinef_set_points(go_qlinef_t l, go_qpointf_t p1, go_qpointf_t p2);
go_qlinef_t qlinef_set_line(go_qlinef_t l, double x1, double y1, double x2, double y2);
go_qlinef_t qlinef_set_length(go_qlinef_t l, double len);
go_qlinef_t qlinef_set_angle(go_qlinef_t l, double angle);

/* --- QLineF: transforms (return new value) --- */
go_qlinef_t qlinef_translated(go_qlinef_t l, go_qpointf_t p);
go_qlinef_t qlinef_translated_xy(go_qlinef_t l, double dx, double dy);
go_qlinef_t qlinef_unit_vector(go_qlinef_t l);
go_qlinef_t qlinef_normal_vector(go_qlinef_t l);

/* --- QLineF: intersection ---
 * If intersection exists (Bounded or Unbounded), writes the point to *out_point and
 * returns the intersection type. If no intersection, *out_point is left untouched and
 * NoIntersection is returned. */
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