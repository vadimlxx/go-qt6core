/* Order: version check first, then real Qt headers.
 *
 * In C++ mode, go_qline_t == QLine and go_qlinef_t == QLineF (see qline_c.h).
 * All functions take/return them by value. No casts, no heap allocation.
 */
#include "../internal/qt6core/check_qt_version.h"

#include <QLine>

#include "qline_c.h"

extern "C" {

/* --- QLine: construction --- */
go_qline_t qline_default()                              { return QLine(); }
go_qline_t qline_new(int x1, int y1, int x2, int y2)    { return QLine(x1, y1, x2, y2); }
go_qline_t qline_from_points(go_qpoint_t p1, go_qpoint_t p2) { return QLine(p1, p2); }

/* --- QLine: accessors --- */
bool qline_is_null(go_qline_t l)          { return l.isNull(); }
go_qpoint_t qline_p1(go_qline_t l)        { return l.p1(); }
go_qpoint_t qline_p2(go_qline_t l)        { return l.p2(); }
int  qline_x1(go_qline_t l)               { return l.x1(); }
int  qline_y1(go_qline_t l)               { return l.y1(); }
int  qline_x2(go_qline_t l)               { return l.x2(); }
int  qline_y2(go_qline_t l)               { return l.y2(); }
int  qline_dx(go_qline_t l)               { return l.dx(); }
int  qline_dy(go_qline_t l)               { return l.dy(); }
go_qpoint_t qline_center(go_qline_t l)    { return l.center(); }

/* --- QLine: mutating setters (return updated value) --- */
go_qline_t qline_translate(go_qline_t l, go_qpoint_t p)  { l.translate(p); return l; }
go_qline_t qline_translate_xy(go_qline_t l, int dx, int dy) { l.translate(dx, dy); return l; }
go_qline_t qline_set_p1(go_qline_t l, go_qpoint_t p)     { l.setP1(p); return l; }
go_qline_t qline_set_p2(go_qline_t l, go_qpoint_t p)     { l.setP2(p); return l; }
go_qline_t qline_set_points(go_qline_t l, go_qpoint_t p1, go_qpoint_t p2) {
    l.setPoints(p1, p2);
    return l;
}
go_qline_t qline_set_line(go_qline_t l, int x1, int y1, int x2, int y2) {
    l.setLine(x1, y1, x2, y2);
    return l;
}

/* --- QLine: transforms --- */
go_qline_t qline_translated(go_qline_t l, go_qpoint_t p)  { return l.translated(p); }
go_qline_t qline_translated_xy(go_qline_t l, int dx, int dy) { return l.translated(dx, dy); }

/* --- QLine: comparison --- */
bool qline_equal(go_qline_t a, go_qline_t b) { return a == b; }

/* --- QLine: conversions --- */
go_qlinef_t qline_to_linef(go_qline_t l) { return l.toLineF(); }

/* --- QLineF: construction --- */
go_qlinef_t qlinef_default()                                     { return QLineF(); }
go_qlinef_t qlinef_new(double x1, double y1, double x2, double y2) { return QLineF(x1, y1, x2, y2); }
go_qlinef_t qlinef_from_points(go_qpointf_t p1, go_qpointf_t p2) { return QLineF(p1, p2); }
go_qlinef_t qlinef_from_line(go_qline_t l)                       { return QLineF(l); }
go_qlinef_t qlinef_from_polar(double length, double angle)       { return QLineF::fromPolar(length, angle); }

/* --- QLineF: accessors --- */
bool qlinef_is_null(go_qlinef_t l)               { return l.isNull(); }
go_qpointf_t qlinef_p1(go_qlinef_t l)            { return l.p1(); }
go_qpointf_t qlinef_p2(go_qlinef_t l)            { return l.p2(); }
double qlinef_x1(go_qlinef_t l)                  { return l.x1(); }
double qlinef_y1(go_qlinef_t l)                  { return l.y1(); }
double qlinef_x2(go_qlinef_t l)                  { return l.x2(); }
double qlinef_y2(go_qlinef_t l)                  { return l.y2(); }
double qlinef_dx(go_qlinef_t l)                  { return l.dx(); }
double qlinef_dy(go_qlinef_t l)                  { return l.dy(); }
double qlinef_length(go_qlinef_t l)              { return l.length(); }
double qlinef_angle(go_qlinef_t l)               { return l.angle(); }
double qlinef_angle_to(go_qlinef_t a, go_qlinef_t b) { return a.angleTo(b); }
go_qpointf_t qlinef_center(go_qlinef_t l)        { return l.center(); }
go_qpointf_t qlinef_point_at(go_qlinef_t l, double t) { return l.pointAt(t); }

/* --- QLineF: mutating setters --- */
go_qlinef_t qlinef_translate(go_qlinef_t l, go_qpointf_t p)  { l.translate(p); return l; }
go_qlinef_t qlinef_translate_xy(go_qlinef_t l, double dx, double dy) { l.translate(dx, dy); return l; }
go_qlinef_t qlinef_set_p1(go_qlinef_t l, go_qpointf_t p)     { l.setP1(p); return l; }
go_qlinef_t qlinef_set_p2(go_qlinef_t l, go_qpointf_t p)     { l.setP2(p); return l; }
go_qlinef_t qlinef_set_points(go_qlinef_t l, go_qpointf_t p1, go_qpointf_t p2) {
    l.setPoints(p1, p2);
    return l;
}
go_qlinef_t qlinef_set_line(go_qlinef_t l, double x1, double y1, double x2, double y2) {
    l.setLine(x1, y1, x2, y2);
    return l;
}
go_qlinef_t qlinef_set_length(go_qlinef_t l, double len)     { l.setLength(len); return l; }
go_qlinef_t qlinef_set_angle(go_qlinef_t l, double angle)    { l.setAngle(angle); return l; }

/* --- QLineF: transforms --- */
go_qlinef_t qlinef_translated(go_qlinef_t l, go_qpointf_t p)  { return l.translated(p); }
go_qlinef_t qlinef_translated_xy(go_qlinef_t l, double dx, double dy) { return l.translated(dx, dy); }
go_qlinef_t qlinef_unit_vector(go_qlinef_t l)                 { return l.unitVector(); }
go_qlinef_t qlinef_normal_vector(go_qlinef_t l)               { return l.normalVector(); }

/* --- QLineF: intersection --- */
int qlinef_intersects(go_qlinef_t a, go_qlinef_t b, go_qpointf_t* out_point) {
    if (!out_point)
        return static_cast<int>(a.intersects(b));
    QPointF pt;
    auto r = a.intersects(b, &pt);
    if (r != QLineF::NoIntersection)
        *out_point = pt;
    return static_cast<int>(r);
}

/* --- QLineF: comparison --- */
bool qlinef_equal(go_qlinef_t a, go_qlinef_t b)      { return a == b; }
bool qlinef_equal_line(go_qlinef_t a, go_qline_t b)  { return a == b; }

/* --- QLineF: conversions --- */
go_qline_t qlinef_to_line(go_qlinef_t l) { return l.toLine(); }

} // extern "C"