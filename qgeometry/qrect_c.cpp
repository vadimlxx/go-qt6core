/* Order: version check first, then real Qt headers.
 *
 * In C++ mode, go_qrect_t == QRect and go_qrectf_t == QRectF (see qrect_c.h).
 * All functions take/return them by value. No casts, no heap allocation.
 */
#include "../internal/qt6core/check_qt_version.h"

#include <QRect>

#include "qrect_c.h"

extern "C" {

/* --- QRect: construction --- */
go_qrect_t qrect_default()                              { return QRect(); }
go_qrect_t qrect_new(int l, int t, int w, int h)        { return QRect(l, t, w, h); }
go_qrect_t qrect_from_points(go_qpoint_t tl, go_qpoint_t br) { return QRect(tl, br); }
go_qrect_t qrect_from_point_size(go_qpoint_t tl, go_qsize_t sz) { return QRect(tl, sz); }

/* --- QRect: accessors --- */
bool qrect_is_null (go_qrect_t r) { return r.isNull();  }
bool qrect_is_empty(go_qrect_t r) { return r.isEmpty(); }
bool qrect_is_valid(go_qrect_t r) { return r.isValid(); }
int  qrect_left    (go_qrect_t r) { return r.left();    }
int  qrect_top     (go_qrect_t r) { return r.top();     }
int  qrect_right   (go_qrect_t r) { return r.right();   }
int  qrect_bottom  (go_qrect_t r) { return r.bottom();  }
int  qrect_x       (go_qrect_t r) { return r.x();       }
int  qrect_y       (go_qrect_t r) { return r.y();       }
int  qrect_width   (go_qrect_t r) { return r.width();   }
int  qrect_height  (go_qrect_t r) { return r.height();  }

/* --- QRect: accessors (packed values) --- */
go_qsize_t  qrect_size        (go_qrect_t r) { return r.size();        }
go_qpoint_t qrect_top_left    (go_qrect_t r) { return r.topLeft();     }
go_qpoint_t qrect_top_right   (go_qrect_t r) { return r.topRight();    }
go_qpoint_t qrect_bottom_left (go_qrect_t r) { return r.bottomLeft();  }
go_qpoint_t qrect_bottom_right(go_qrect_t r) { return r.bottomRight(); }
go_qpoint_t qrect_center      (go_qrect_t r) { return r.center();      }

/* --- QRect: transforms --- */
go_qrect_t qrect_normalized (go_qrect_t r)                       { return r.normalized();         }
go_qrect_t qrect_transposed (go_qrect_t r)                       { return r.transposed();         }
go_qrect_t qrect_translated (go_qrect_t r, int dx, int dy)       { return r.translated(dx, dy);   }
go_qrect_t qrect_translated_point(go_qrect_t r, go_qpoint_t p)   { return r.translated(p);        }
go_qrect_t qrect_adjusted   (go_qrect_t r, int dx1, int dy1, int dx2, int dy2) {
    return r.adjusted(dx1, dy1, dx2, dy2);
}

/* --- QRect: mutating setters --- */
go_qrect_t qrect_set_left  (go_qrect_t r, int pos) { r.setLeft(pos);   return r; }
go_qrect_t qrect_set_top   (go_qrect_t r, int pos) { r.setTop(pos);    return r; }
go_qrect_t qrect_set_right (go_qrect_t r, int pos) { r.setRight(pos);  return r; }
go_qrect_t qrect_set_bottom(go_qrect_t r, int pos) { r.setBottom(pos); return r; }
go_qrect_t qrect_set_x     (go_qrect_t r, int pos) { r.setX(pos);      return r; }
go_qrect_t qrect_set_y     (go_qrect_t r, int pos) { r.setY(pos);      return r; }
go_qrect_t qrect_set_width (go_qrect_t r, int w)   { r.setWidth(w);    return r; }
go_qrect_t qrect_set_height(go_qrect_t r, int h)   { r.setHeight(h);   return r; }
go_qrect_t qrect_set_size  (go_qrect_t r, go_qsize_t s)  { r.setSize(s);  return r; }
go_qrect_t qrect_set_top_left    (go_qrect_t r, go_qpoint_t p) { r.setTopLeft(p);     return r; }
go_qrect_t qrect_set_top_right   (go_qrect_t r, go_qpoint_t p) { r.setTopRight(p);    return r; }
go_qrect_t qrect_set_bottom_left (go_qrect_t r, go_qpoint_t p) { r.setBottomLeft(p);  return r; }
go_qrect_t qrect_set_bottom_right(go_qrect_t r, go_qpoint_t p) { r.setBottomRight(p); return r; }
go_qrect_t qrect_set_rect  (go_qrect_t r, int x, int y, int w, int h)   { r.setRect(x, y, w, h);   return r; }
go_qrect_t qrect_set_coords(go_qrect_t r, int x1, int y1, int x2, int y2) { r.setCoords(x1, y1, x2, y2); return r; }

/* --- QRect: movers --- */
go_qrect_t qrect_move_left  (go_qrect_t r, int pos) { r.moveLeft(pos);   return r; }
go_qrect_t qrect_move_top   (go_qrect_t r, int pos) { r.moveTop(pos);    return r; }
go_qrect_t qrect_move_right (go_qrect_t r, int pos) { r.moveRight(pos);  return r; }
go_qrect_t qrect_move_bottom(go_qrect_t r, int pos) { r.moveBottom(pos); return r; }
go_qrect_t qrect_move_top_left    (go_qrect_t r, go_qpoint_t p) { r.moveTopLeft(p);     return r; }
go_qrect_t qrect_move_top_right   (go_qrect_t r, go_qpoint_t p) { r.moveTopRight(p);    return r; }
go_qrect_t qrect_move_bottom_left (go_qrect_t r, go_qpoint_t p) { r.moveBottomLeft(p);  return r; }
go_qrect_t qrect_move_bottom_right(go_qrect_t r, go_qpoint_t p) { r.moveBottomRight(p); return r; }
go_qrect_t qrect_move_center(go_qrect_t r, go_qpoint_t p)       { r.moveCenter(p);      return r; }
go_qrect_t qrect_move_to    (go_qrect_t r, int x, int y)        { r.moveTo(x, y);       return r; }
go_qrect_t qrect_move_to_point(go_qrect_t r, go_qpoint_t p)     { r.moveTo(p);          return r; }

/* --- QRect: translate --- */
go_qrect_t qrect_translate      (go_qrect_t r, int dx, int dy) { r.translate(dx, dy); return r; }
go_qrect_t qrect_translate_point(go_qrect_t r, go_qpoint_t p)  { r.translate(p);      return r; }

/* --- QRect: adjust --- */
go_qrect_t qrect_adjust(go_qrect_t r, int dx1, int dy1, int dx2, int dy2) {
    r.adjust(dx1, dy1, dx2, dy2);
    return r;
}

/* --- QRect: boolean ops --- */
bool qrect_contains_rect (go_qrect_t r, go_qrect_t o, bool proper) { return r.contains(o, proper); }
bool qrect_contains_point(go_qrect_t r, go_qpoint_t p, bool proper) { return r.contains(p, proper); }
bool qrect_contains_xy   (go_qrect_t r, int x, int y)               { return r.contains(x, y);     }
bool qrect_intersects    (go_qrect_t r, go_qrect_t o)               { return r.intersects(o);      }
bool qrect_equal         (go_qrect_t a, go_qrect_t b)               { return a == b;                }

/* --- QRect: unions/intersections --- */
go_qrect_t qrect_united     (go_qrect_t r, go_qrect_t o) { return r.united(o);      }
go_qrect_t qrect_intersected(go_qrect_t r, go_qrect_t o) { return r.intersected(o); }

/* --- QRect: margins --- */
go_qrect_t qrect_margins_added  (go_qrect_t r, go_qmargins_t m) { return r.marginsAdded(m);   }
go_qrect_t qrect_margins_removed(go_qrect_t r, go_qmargins_t m) { return r.marginsRemoved(m); }

/* --- QRect: static --- */
go_qrect_t qrect_span(go_qpoint_t p1, go_qpoint_t p2) { return QRect::span(p1, p2); }

/* --- QRect: conversions --- */
go_qrectf_t qrect_to_rectf(go_qrect_t r) { return r.toRectF(); }

/* --- QRectF: construction --- */
go_qrectf_t qrectf_default()                                { return QRectF(); }
go_qrectf_t qrectf_new(double l, double t, double w, double h) { return QRectF(l, t, w, h); }
go_qrectf_t qrectf_from_points(go_qpointf_t tl, go_qpointf_t br) { return QRectF(tl, br); }
go_qrectf_t qrectf_from_point_size(go_qpointf_t tl, go_qsizef_t sz) { return QRectF(tl, sz); }
go_qrectf_t qrectf_from_rect(go_qrect_t r)                  { return QRectF(r); }

/* --- QRectF: accessors --- */
bool   qrectf_is_null (go_qrectf_t r) { return r.isNull();  }
bool   qrectf_is_empty(go_qrectf_t r) { return r.isEmpty(); }
bool   qrectf_is_valid(go_qrectf_t r) { return r.isValid(); }
double qrectf_left    (go_qrectf_t r) { return r.left();    }
double qrectf_top     (go_qrectf_t r) { return r.top();     }
double qrectf_right   (go_qrectf_t r) { return r.right();   }
double qrectf_bottom  (go_qrectf_t r) { return r.bottom();  }
double qrectf_x       (go_qrectf_t r) { return r.x();       }
double qrectf_y       (go_qrectf_t r) { return r.y();       }
double qrectf_width   (go_qrectf_t r) { return r.width();   }
double qrectf_height  (go_qrectf_t r) { return r.height();  }

/* --- QRectF: accessors (packed values) --- */
go_qsizef_t  qrectf_size        (go_qrectf_t r) { return r.size();        }
go_qpointf_t qrectf_top_left    (go_qrectf_t r) { return r.topLeft();     }
go_qpointf_t qrectf_top_right   (go_qrectf_t r) { return r.topRight();    }
go_qpointf_t qrectf_bottom_left (go_qrectf_t r) { return r.bottomLeft();  }
go_qpointf_t qrectf_bottom_right(go_qrectf_t r) { return r.bottomRight(); }
go_qpointf_t qrectf_center      (go_qrectf_t r) { return r.center();      }

/* --- QRectF: transforms --- */
go_qrectf_t qrectf_normalized (go_qrectf_t r)                          { return r.normalized();       }
go_qrectf_t qrectf_transposed (go_qrectf_t r)                          { return r.transposed();       }
go_qrectf_t qrectf_translated (go_qrectf_t r, double dx, double dy)    { return r.translated(dx, dy); }
go_qrectf_t qrectf_translated_point(go_qrectf_t r, go_qpointf_t p)     { return r.translated(p);      }
go_qrectf_t qrectf_adjusted   (go_qrectf_t r, double dx1, double dy1, double dx2, double dy2) {
    return r.adjusted(dx1, dy1, dx2, dy2);
}

/* --- QRectF: mutating setters --- */
go_qrectf_t qrectf_set_left  (go_qrectf_t r, double pos) { r.setLeft(pos);   return r; }
go_qrectf_t qrectf_set_top   (go_qrectf_t r, double pos) { r.setTop(pos);    return r; }
go_qrectf_t qrectf_set_right (go_qrectf_t r, double pos) { r.setRight(pos);  return r; }
go_qrectf_t qrectf_set_bottom(go_qrectf_t r, double pos) { r.setBottom(pos); return r; }
go_qrectf_t qrectf_set_x     (go_qrectf_t r, double pos) { r.setX(pos);      return r; }
go_qrectf_t qrectf_set_y     (go_qrectf_t r, double pos) { r.setY(pos);      return r; }
go_qrectf_t qrectf_set_width (go_qrectf_t r, double w)   { r.setWidth(w);    return r; }
go_qrectf_t qrectf_set_height(go_qrectf_t r, double h)   { r.setHeight(h);   return r; }
go_qrectf_t qrectf_set_size  (go_qrectf_t r, go_qsizef_t s)  { r.setSize(s);  return r; }
go_qrectf_t qrectf_set_top_left    (go_qrectf_t r, go_qpointf_t p) { r.setTopLeft(p);     return r; }
go_qrectf_t qrectf_set_top_right   (go_qrectf_t r, go_qpointf_t p) { r.setTopRight(p);    return r; }
go_qrectf_t qrectf_set_bottom_left (go_qrectf_t r, go_qpointf_t p) { r.setBottomLeft(p);  return r; }
go_qrectf_t qrectf_set_bottom_right(go_qrectf_t r, go_qpointf_t p) { r.setBottomRight(p); return r; }
go_qrectf_t qrectf_set_rect  (go_qrectf_t r, double x, double y, double w, double h) {
    r.setRect(x, y, w, h);
    return r;
}
go_qrectf_t qrectf_set_coords(go_qrectf_t r, double x1, double y1, double x2, double y2) {
    r.setCoords(x1, y1, x2, y2);
    return r;
}

/* --- QRectF: movers --- */
go_qrectf_t qrectf_move_left  (go_qrectf_t r, double pos) { r.moveLeft(pos);   return r; }
go_qrectf_t qrectf_move_top   (go_qrectf_t r, double pos) { r.moveTop(pos);    return r; }
go_qrectf_t qrectf_move_right (go_qrectf_t r, double pos) { r.moveRight(pos);  return r; }
go_qrectf_t qrectf_move_bottom(go_qrectf_t r, double pos) { r.moveBottom(pos); return r; }
go_qrectf_t qrectf_move_top_left    (go_qrectf_t r, go_qpointf_t p) { r.moveTopLeft(p);     return r; }
go_qrectf_t qrectf_move_top_right   (go_qrectf_t r, go_qpointf_t p) { r.moveTopRight(p);    return r; }
go_qrectf_t qrectf_move_bottom_left (go_qrectf_t r, go_qpointf_t p) { r.moveBottomLeft(p);  return r; }
go_qrectf_t qrectf_move_bottom_right(go_qrectf_t r, go_qpointf_t p) { r.moveBottomRight(p); return r; }
go_qrectf_t qrectf_move_center(go_qrectf_t r, go_qpointf_t p)       { r.moveCenter(p);      return r; }
go_qrectf_t qrectf_move_to    (go_qrectf_t r, double x, double y)   { r.moveTo(x, y);       return r; }
go_qrectf_t qrectf_move_to_point(go_qrectf_t r, go_qpointf_t p)     { r.moveTo(p);          return r; }

/* --- QRectF: translate --- */
go_qrectf_t qrectf_translate      (go_qrectf_t r, double dx, double dy) { r.translate(dx, dy); return r; }
go_qrectf_t qrectf_translate_point(go_qrectf_t r, go_qpointf_t p)       { r.translate(p);      return r; }

/* --- QRectF: adjust --- */
go_qrectf_t qrectf_adjust(go_qrectf_t r, double dx1, double dy1, double dx2, double dy2) {
    r.adjust(dx1, dy1, dx2, dy2);
    return r;
}

/* --- QRectF: boolean ops --- */
bool qrectf_contains_rect (go_qrectf_t r, go_qrectf_t o) { return r.contains(o); }
bool qrectf_contains_point(go_qrectf_t r, go_qpointf_t p){ return r.contains(p); }
bool qrectf_contains_xy   (go_qrectf_t r, double x, double y) { return r.contains(x, y); }
bool qrectf_intersects    (go_qrectf_t r, go_qrectf_t o) { return r.intersects(o); }
bool qrectf_equal         (go_qrectf_t a, go_qrectf_t b) { return a == b;         }
bool qrectf_equal_rect    (go_qrectf_t a, go_qrect_t b)  { return a == b;         }

/* --- QRectF: unions/intersections --- */
go_qrectf_t qrectf_united     (go_qrectf_t r, go_qrectf_t o) { return r.united(o);      }
go_qrectf_t qrectf_intersected(go_qrectf_t r, go_qrectf_t o) { return r.intersected(o); }

/* --- QRectF: margins --- */
go_qrectf_t qrectf_margins_added  (go_qrectf_t r, go_qmarginsf_t m) { return r.marginsAdded(m);   }
go_qrectf_t qrectf_margins_removed(go_qrectf_t r, go_qmarginsf_t m) { return r.marginsRemoved(m); }

/* --- QRectF: conversions --- */
go_qrect_t qrectf_to_rect        (go_qrectf_t r) { return r.toRect();        }
go_qrect_t qrectf_to_aligned_rect(go_qrectf_t r) { return r.toAlignedRect(); }

} // extern "C"