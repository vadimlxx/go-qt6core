/* C ABI for Qt 6.11 QtCore/QSize and QSizeF.
 *
 * QSize / QSizeF are exposed as HANDLES (QSize* / QSizeF*). Each handle owns a heap-allocated
 * C++ QSize / QSizeF. All accessors/setters go directly to the C++ object.
 *
 * Memory: every function returning go_qsize_t / go_qsizef_t allocates; caller owns it and must
 * delete via qsize_delete / qsizef_delete.
 */
#ifndef GO_QT6CORE_QGEOMETRY_QSIZE_C_H
#define GO_QT6CORE_QGEOMETRY_QSIZE_C_H

#include <stdbool.h>
#include "qmargins_c.h"

#ifdef __cplusplus
class QSize;
class QSizeF;
typedef QSize*  go_qsize_t;
typedef QSizeF* go_qsizef_t;
extern "C" {
#else
typedef struct go_qsize_handle_s  go_qsize_handle_s;
typedef struct go_qsizef_handle_s go_qsizef_handle_s;
typedef go_qsize_handle_s*  go_qsize_t;
typedef go_qsizef_handle_s* go_qsizef_t;
#endif

/* --- QSize: lifecycle --- */
go_qsize_t qsize_new(int w, int h);
go_qsize_t qsize_new_default(void);
go_qsize_t qsize_clone(go_qsize_t s);
void       qsize_delete(go_qsize_t s);

/* --- QSize: accessors --- */
bool qsize_is_null (go_qsize_t s);
bool qsize_is_empty(go_qsize_t s);
bool qsize_is_valid(go_qsize_t s);
int  qsize_width    (go_qsize_t s);
int  qsize_height   (go_qsize_t s);

/* --- QSize: mutating setters (in place, void) --- */
void qsize_set_width (go_qsize_t self, int w);
void qsize_set_height(go_qsize_t self, int h);
void qsize_transpose (go_qsize_t self);
void qsize_scale     (go_qsize_t self, int w, int h, int mode);
void qsize_scale_size(go_qsize_t self, go_qsize_t other, int mode);

/* --- QSize: transforms (return new handle) --- */
go_qsize_t qsize_transposed (go_qsize_t s);
go_qsize_t qsize_scaled     (go_qsize_t s, int w, int h, int mode);
go_qsize_t qsize_scaled_size(go_qsize_t s, go_qsize_t other, int mode);
go_qsize_t qsize_expanded_to(go_qsize_t s, go_qsize_t other);
go_qsize_t qsize_bounded_to (go_qsize_t s, go_qsize_t other);
go_qsize_t qsize_grown_by   (go_qsize_t s, go_qmargins_t m);
go_qsize_t qsize_shrunk_by  (go_qsize_t s, go_qmargins_t m);

/* --- QSize: arithmetic (return new handles) --- */
go_qsize_t qsize_add      (go_qsize_t a, go_qsize_t b);
go_qsize_t qsize_sub      (go_qsize_t a, go_qsize_t b);
go_qsize_t qsize_mul_float(go_qsize_t s, double c);
go_qsize_t qsize_div_float(go_qsize_t s, double c);

/* --- QSize: in-place assign (mutate self) --- */
void qsize_add_assign       (go_qsize_t self, go_qsize_t other);
void qsize_sub_assign       (go_qsize_t self, go_qsize_t other);
void qsize_mul_float_assign (go_qsize_t self, double c);
void qsize_div_float_assign (go_qsize_t self, double c);

/* --- QSize: comparison --- */
bool qsize_equal(go_qsize_t a, go_qsize_t b);

/* --- QSize: conversions --- */
go_qsizef_t qsize_to_sizef(go_qsize_t s);

/* --- QSizeF: lifecycle --- */
go_qsizef_t qsizef_new(double w, double h);
go_qsizef_t qsizef_new_default(void);
go_qsizef_t qsizef_from_size(go_qsize_t s);
go_qsizef_t qsizef_clone(go_qsizef_t s);
void        qsizef_delete(go_qsizef_t s);

/* --- QSizeF: accessors --- */
bool   qsizef_is_null (go_qsizef_t s);
bool   qsizef_is_empty(go_qsizef_t s);
bool   qsizef_is_valid(go_qsizef_t s);
double qsizef_width   (go_qsizef_t s);
double qsizef_height  (go_qsizef_t s);

/* --- QSizeF: mutating setters (in place, void) --- */
void qsizef_set_width (go_qsizef_t self, double w);
void qsizef_set_height(go_qsizef_t self, double h);
void qsizef_transpose (go_qsizef_t self);
void qsizef_scale     (go_qsizef_t self, double w, double h, int mode);
void qsizef_scale_size(go_qsizef_t self, go_qsizef_t other, int mode);

/* --- QSizeF: transforms (return new handle) --- */
go_qsizef_t qsizef_transposed (go_qsizef_t s);
go_qsizef_t qsizef_scaled     (go_qsizef_t s, double w, double h, int mode);
go_qsizef_t qsizef_scaled_size(go_qsizef_t s, go_qsizef_t other, int mode);
go_qsizef_t qsizef_expanded_to(go_qsizef_t s, go_qsizef_t other);
go_qsizef_t qsizef_bounded_to (go_qsizef_t s, go_qsizef_t other);
go_qsizef_t qsizef_grown_by   (go_qsizef_t s, go_qmarginsf_t m);
go_qsizef_t qsizef_shrunk_by  (go_qsizef_t s, go_qmarginsf_t m);

/* --- QSizeF: arithmetic (return new handles) --- */
go_qsizef_t qsizef_add      (go_qsizef_t a, go_qsizef_t b);
go_qsizef_t qsizef_sub      (go_qsizef_t a, go_qsizef_t b);
go_qsizef_t qsizef_mul_float(go_qsizef_t s, double c);
go_qsizef_t qsizef_div_float(go_qsizef_t s, double c);

/* --- QSizeF: in-place assign (mutate self) --- */
void qsizef_add_assign       (go_qsizef_t self, go_qsizef_t other);
void qsizef_sub_assign       (go_qsizef_t self, go_qsizef_t other);
void qsizef_mul_float_assign (go_qsizef_t self, double c);
void qsizef_div_float_assign (go_qsizef_t self, double c);

/* --- QSizeF: comparison --- */
bool qsizef_equal     (go_qsizef_t a, go_qsizef_t b);
bool qsizef_equal_size(go_qsizef_t a, go_qsize_t b);

/* --- QSizeF: conversions --- */
go_qsize_t qsizef_to_size(go_qsizef_t s);

#ifdef __cplusplus
}
#endif

#endif