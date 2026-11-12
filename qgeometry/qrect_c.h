/* C ABI for Qt 6.11 QtCore/QRect and QRectF.
 *
 * Both types are Q_RELOCATABLE_TYPE value types. To minimise casts in the C++ implementation,
 * we typedef go_qrect_t / go_qrectf_t directly to QRect / QRectF in C++ mode. In C mode
 * (used by cgo's preamble), they fall back to flat POD structs with an identical memory layout:
 *   - QRect  stores (x1, y1, x2, y2) — inclusive left/top/right/bottom. 16 B, align 4.
 *   - QRectF stores (xp, yp, w, h)   — left/top/width/height.              32 B, align 8.
 *
 * Safety: both are trivially copyable, so the two representations are ABI-compatible.
 */
#ifndef GO_QT6CORE_QGEOMETRY_QRECT_C_H
#define GO_QT6CORE_QGEOMETRY_QRECT_C_H

#include <stdbool.h>
#include "qpoint_c.h"
#include "qsize_c.h"
#include "qmargins_c.h"

#ifdef __cplusplus
#include <QRect>
typedef QRect  go_qrect_t;
typedef QRectF go_qrectf_t;
extern "C" {
#else
typedef struct go_qrect_s {
    int x1;
    int y1;
    int x2;
    int y2;
} go_qrect_t;

typedef struct go_qrectf_s {
    double xp;
    double yp;
    double w;
    double h;
} go_qrectf_t;
#endif

/* --- QRect: construction --- */
go_qrect_t qrect_default(void);
go_qrect_t qrect_new(int left, int top, int width, int height);
go_qrect_t qrect_from_points(go_qpoint_t topLeft, go_qpoint_t bottomRight);
go_qrect_t qrect_from_point_size(go_qpoint_t topLeft, go_qsize_t size);

/* --- QRect: accessors (values) --- */
bool qrect_is_null (go_qrect_t r);
bool qrect_is_empty(go_qrect_t r);
bool qrect_is_valid(go_qrect_t r);
int  qrect_left    (go_qrect_t r);
int  qrect_top     (go_qrect_t r);
int  qrect_right   (go_qrect_t r);
int  qrect_bottom  (go_qrect_t r);
int  qrect_x       (go_qrect_t r);
int  qrect_y       (go_qrect_t r);
int  qrect_width   (go_qrect_t r);
int  qrect_height  (go_qrect_t r);

/* --- QRect: accessors (values, returns a packed value) --- */
go_qsize_t  qrect_size        (go_qrect_t r);
go_qpoint_t qrect_top_left    (go_qrect_t r);
go_qpoint_t qrect_top_right   (go_qrect_t r);
go_qpoint_t qrect_bottom_left (go_qrect_t r);
go_qpoint_t qrect_bottom_right(go_qrect_t r);
go_qpoint_t qrect_center      (go_qrect_t r);

/* --- QRect: transforms (return new value) --- */
go_qrect_t qrect_normalized (go_qrect_t r);
go_qrect_t qrect_transposed (go_qrect_t r);
go_qrect_t qrect_translated (go_qrect_t r, int dx, int dy);
go_qrect_t qrect_translated_point(go_qrect_t r, go_qpoint_t p);
go_qrect_t qrect_adjusted   (go_qrect_t r, int dx1, int dy1, int dx2, int dy2);

/* --- QRect: mutating setters (return updated value) --- */
go_qrect_t qrect_set_left  (go_qrect_t r, int pos);
go_qrect_t qrect_set_top   (go_qrect_t r, int pos);
go_qrect_t qrect_set_right (go_qrect_t r, int pos);
go_qrect_t qrect_set_bottom(go_qrect_t r, int pos);
go_qrect_t qrect_set_x     (go_qrect_t r, int pos);
go_qrect_t qrect_set_y     (go_qrect_t r, int pos);
go_qrect_t qrect_set_width (go_qrect_t r, int w);
go_qrect_t qrect_set_height(go_qrect_t r, int h);
go_qrect_t qrect_set_size  (go_qrect_t r, go_qsize_t s);
go_qrect_t qrect_set_top_left    (go_qrect_t r, go_qpoint_t p);
go_qrect_t qrect_set_top_right   (go_qrect_t r, go_qpoint_t p);
go_qrect_t qrect_set_bottom_left (go_qrect_t r, go_qpoint_t p);
go_qrect_t qrect_set_bottom_right(go_qrect_t r, go_qpoint_t p);
go_qrect_t qrect_set_rect  (go_qrect_t r, int x, int y, int w, int h);
go_qrect_t qrect_set_coords(go_qrect_t r, int x1, int y1, int x2, int y2);

/* --- QRect: movers (return updated value) --- */
go_qrect_t qrect_move_left  (go_qrect_t r, int pos);
go_qrect_t qrect_move_top   (go_qrect_t r, int pos);
go_qrect_t qrect_move_right (go_qrect_t r, int pos);
go_qrect_t qrect_move_bottom(go_qrect_t r, int pos);
go_qrect_t qrect_move_top_left    (go_qrect_t r, go_qpoint_t p);
go_qrect_t qrect_move_top_right   (go_qrect_t r, go_qpoint_t p);
go_qrect_t qrect_move_bottom_left (go_qrect_t r, go_qpoint_t p);
go_qrect_t qrect_move_bottom_right(go_qrect_t r, go_qpoint_t p);
go_qrect_t qrect_move_center(go_qrect_t r, go_qpoint_t p);
go_qrect_t qrect_move_to    (go_qrect_t r, int x, int y);
go_qrect_t qrect_move_to_point(go_qrect_t r, go_qpoint_t p);

/* --- QRect: translate (return updated value) --- */
go_qrect_t qrect_translate      (go_qrect_t r, int dx, int dy);
go_qrect_t qrect_translate_point(go_qrect_t r, go_qpoint_t p);

/* --- QRect: adjust (return updated value) --- */
go_qrect_t qrect_adjust(go_qrect_t r, int dx1, int dy1, int dx2, int dy2);

/* --- QRect: boolean ops --- */
bool qrect_contains_rect (go_qrect_t r, go_qrect_t other, bool proper);
bool qrect_contains_point(go_qrect_t r, go_qpoint_t p, bool proper);
bool qrect_contains_xy   (go_qrect_t r, int x, int y);
bool qrect_intersects    (go_qrect_t r, go_qrect_t other);
bool qrect_equal         (go_qrect_t a, go_qrect_t b);

/* --- QRect: unions/intersections (new values) --- */
go_qrect_t qrect_united     (go_qrect_t r, go_qrect_t other);
go_qrect_t qrect_intersected(go_qrect_t r, go_qrect_t other);

/* --- QRect: margins (new values) --- */
go_qrect_t qrect_margins_added  (go_qrect_t r, go_qmargins_t m);
go_qrect_t qrect_margins_removed(go_qrect_t r, go_qmargins_t m);

/* --- QRect: static --- */
go_qrect_t qrect_span(go_qpoint_t p1, go_qpoint_t p2);

/* --- QRect: conversions --- */
go_qrectf_t qrect_to_rectf(go_qrect_t r);

/* --- QRectF: construction --- */
go_qrectf_t qrectf_default(void);
go_qrectf_t qrectf_new(double left, double top, double width, double height);
go_qrectf_t qrectf_from_points(go_qpointf_t topLeft, go_qpointf_t bottomRight);
go_qrectf_t qrectf_from_point_size(go_qpointf_t topLeft, go_qsizef_t size);
go_qrectf_t qrectf_from_rect(go_qrect_t r);

/* --- QRectF: accessors (values) --- */
bool   qrectf_is_null (go_qrectf_t r);
bool   qrectf_is_empty(go_qrectf_t r);
bool   qrectf_is_valid(go_qrectf_t r);
double qrectf_left    (go_qrectf_t r);
double qrectf_top     (go_qrectf_t r);
double qrectf_right   (go_qrectf_t r);
double qrectf_bottom  (go_qrectf_t r);
double qrectf_x       (go_qrectf_t r);
double qrectf_y       (go_qrectf_t r);
double qrectf_width   (go_qrectf_t r);
double qrectf_height  (go_qrectf_t r);

/* --- QRectF: accessors (values, returns a packed value) --- */
go_qsizef_t  qrectf_size        (go_qrectf_t r);
go_qpointf_t qrectf_top_left    (go_qrectf_t r);
go_qpointf_t qrectf_top_right   (go_qrectf_t r);
go_qpointf_t qrectf_bottom_left (go_qrectf_t r);
go_qpointf_t qrectf_bottom_right(go_qrectf_t r);
go_qpointf_t qrectf_center      (go_qrectf_t r);

/* --- QRectF: transforms (return new value) --- */
go_qrectf_t qrectf_normalized (go_qrectf_t r);
go_qrectf_t qrectf_transposed (go_qrectf_t r);
go_qrectf_t qrectf_translated (go_qrectf_t r, double dx, double dy);
go_qrectf_t qrectf_translated_point(go_qrectf_t r, go_qpointf_t p);
go_qrectf_t qrectf_adjusted   (go_qrectf_t r, double dx1, double dy1, double dx2, double dy2);

/* --- QRectF: mutating setters (return updated value) --- */
go_qrectf_t qrectf_set_left  (go_qrectf_t r, double pos);
go_qrectf_t qrectf_set_top   (go_qrectf_t r, double pos);
go_qrectf_t qrectf_set_right (go_qrectf_t r, double pos);
go_qrectf_t qrectf_set_bottom(go_qrectf_t r, double pos);
go_qrectf_t qrectf_set_x     (go_qrectf_t r, double pos);
go_qrectf_t qrectf_set_y     (go_qrectf_t r, double pos);
go_qrectf_t qrectf_set_width (go_qrectf_t r, double w);
go_qrectf_t qrectf_set_height(go_qrectf_t r, double h);
go_qrectf_t qrectf_set_size  (go_qrectf_t r, go_qsizef_t s);
go_qrectf_t qrectf_set_top_left    (go_qrectf_t r, go_qpointf_t p);
go_qrectf_t qrectf_set_top_right   (go_qrectf_t r, go_qpointf_t p);
go_qrectf_t qrectf_set_bottom_left (go_qrectf_t r, go_qpointf_t p);
go_qrectf_t qrectf_set_bottom_right(go_qrectf_t r, go_qpointf_t p);
go_qrectf_t qrectf_set_rect  (go_qrectf_t r, double x, double y, double w, double h);
go_qrectf_t qrectf_set_coords(go_qrectf_t r, double x1, double y1, double x2, double y2);

/* --- QRectF: movers (return updated value) --- */
go_qrectf_t qrectf_move_left  (go_qrectf_t r, double pos);
go_qrectf_t qrectf_move_top   (go_qrectf_t r, double pos);
go_qrectf_t qrectf_move_right (go_qrectf_t r, double pos);
go_qrectf_t qrectf_move_bottom(go_qrectf_t r, double pos);
go_qrectf_t qrectf_move_top_left    (go_qrectf_t r, go_qpointf_t p);
go_qrectf_t qrectf_move_top_right   (go_qrectf_t r, go_qpointf_t p);
go_qrectf_t qrectf_move_bottom_left (go_qrectf_t r, go_qpointf_t p);
go_qrectf_t qrectf_move_bottom_right(go_qrectf_t r, go_qpointf_t p);
go_qrectf_t qrectf_move_center(go_qrectf_t r, go_qpointf_t p);
go_qrectf_t qrectf_move_to    (go_qrectf_t r, double x, double y);
go_qrectf_t qrectf_move_to_point(go_qrectf_t r, go_qpointf_t p);

/* --- QRectF: translate (return updated value) --- */
go_qrectf_t qrectf_translate      (go_qrectf_t r, double dx, double dy);
go_qrectf_t qrectf_translate_point(go_qrectf_t r, go_qpointf_t p);

/* --- QRectF: adjust (return updated value) --- */
go_qrectf_t qrectf_adjust(go_qrectf_t r, double dx1, double dy1, double dx2, double dy2);

/* --- QRectF: boolean ops --- */
bool qrectf_contains_rect (go_qrectf_t r, go_qrectf_t other);
bool qrectf_contains_point(go_qrectf_t r, go_qpointf_t p);
bool qrectf_contains_xy   (go_qrectf_t r, double x, double y);
bool qrectf_intersects    (go_qrectf_t r, go_qrectf_t other);
bool qrectf_equal         (go_qrectf_t a, go_qrectf_t b);
bool qrectf_equal_rect    (go_qrectf_t a, go_qrect_t b);

/* --- QRectF: unions/intersections (new values) --- */
go_qrectf_t qrectf_united     (go_qrectf_t r, go_qrectf_t other);
go_qrectf_t qrectf_intersected(go_qrectf_t r, go_qrectf_t other);

/* --- QRectF: margins (new values) --- */
go_qrectf_t qrectf_margins_added  (go_qrectf_t r, go_qmarginsf_t m);
go_qrectf_t qrectf_margins_removed(go_qrectf_t r, go_qmarginsf_t m);

/* --- QRectF: conversions --- */
go_qrect_t qrectf_to_rect        (go_qrectf_t r);
go_qrect_t qrectf_to_aligned_rect(go_qrectf_t r);

#ifdef __cplusplus
}
#endif

#endif