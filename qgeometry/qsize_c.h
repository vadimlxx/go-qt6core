/* C ABI for Qt 6.11 QtCore/QSize and QSizeF.
 *
 * Both types are Q_RELOCATABLE_TYPE value types. To minimise casts in the C++ implementation, we typedef
 * go_qsize_t / go_qsizef_t directly to QSize / QSizeF in C++ mode. In C mode (used by cgo's preamble),
 * they fall back to flat POD structs with an identical memory layout.
 *
 * Safety: QSize (2×int = 8 B) and QSizeF (2×qreal = 16 B) are trivially copyable, so both representations
 * are ABI-compatible on all supported platforms.
 */
#ifndef GO_QT6CORE_QGEOMETRY_QSIZE_C_H
#define GO_QT6CORE_QGEOMETRY_QSIZE_C_H

#include <stdbool.h>
#include "qmargins_c.h"

#ifdef __cplusplus
#include <QSize>
typedef QSize  go_qsize_t;
typedef QSizeF go_qsizef_t;
extern "C" {
#else
typedef struct go_qsize_s {
    int width;
    int height;
} go_qsize_t;

typedef struct go_qsizef_s {
    double width;
    double height;
} go_qsizef_t;
#endif

/* --- QSize: construction --- */
go_qsize_t qsize_default(void);
go_qsize_t qsize_new(int w, int h);

/* --- QSize: accessors --- */
bool qsize_is_null (go_qsize_t s);
bool qsize_is_empty(go_qsize_t s);
bool qsize_is_valid(go_qsize_t s);
int  qsize_width    (go_qsize_t s);
int  qsize_height   (go_qsize_t s);

/* --- QSize: mutating setters (return updated value) --- */
go_qsize_t qsize_set_width (go_qsize_t s, int w);
go_qsize_t qsize_set_height(go_qsize_t s, int h);

/* --- QSize: transforms (return new QSize) --- */
go_qsize_t qsize_transposed (go_qsize_t s);
go_qsize_t qsize_scaled     (go_qsize_t s, int w, int h, int mode);
go_qsize_t qsize_scaled_size(go_qsize_t s, go_qsize_t other, int mode);
go_qsize_t qsize_expanded_to(go_qsize_t s, go_qsize_t other);
go_qsize_t qsize_bounded_to (go_qsize_t s, go_qsize_t other);
go_qsize_t qsize_grown_by   (go_qsize_t s, go_qmargins_t m);
go_qsize_t qsize_shrunk_by  (go_qsize_t s, go_qmargins_t m);

/* --- QSize: arithmetic --- */
go_qsize_t qsize_add      (go_qsize_t a, go_qsize_t b);
go_qsize_t qsize_sub      (go_qsize_t a, go_qsize_t b);
go_qsize_t qsize_mul_float(go_qsize_t s, double c);
go_qsize_t qsize_div_float(go_qsize_t s, double c);

/* --- QSize: comparison --- */
bool qsize_equal(go_qsize_t a, go_qsize_t b);

/* --- QSize: conversions --- */
go_qsizef_t qsize_to_sizef(go_qsize_t s);

/* --- QSizeF: construction --- */
go_qsizef_t qsizef_default(void);
go_qsizef_t qsizef_new(double w, double h);
go_qsizef_t qsizef_from_size(go_qsize_t s);

/* --- QSizeF: accessors --- */
bool   qsizef_is_null (go_qsizef_t s);
bool   qsizef_is_empty(go_qsizef_t s);
bool   qsizef_is_valid(go_qsizef_t s);
double qsizef_width   (go_qsizef_t s);
double qsizef_height  (go_qsizef_t s);

/* --- QSizeF: mutating setters --- */
go_qsizef_t qsizef_set_width (go_qsizef_t s, double w);
go_qsizef_t qsizef_set_height(go_qsizef_t s, double h);

/* --- QSizeF: transforms --- */
go_qsizef_t qsizef_transposed (go_qsizef_t s);
go_qsizef_t qsizef_scaled     (go_qsizef_t s, double w, double h, int mode);
go_qsizef_t qsizef_scaled_size(go_qsizef_t s, go_qsizef_t other, int mode);
go_qsizef_t qsizef_expanded_to(go_qsizef_t s, go_qsizef_t other);
go_qsizef_t qsizef_bounded_to (go_qsizef_t s, go_qsizef_t other);
go_qsizef_t qsizef_grown_by   (go_qsizef_t s, go_qmarginsf_t m);
go_qsizef_t qsizef_shrunk_by  (go_qsizef_t s, go_qmarginsf_t m);

/* --- QSizeF: arithmetic --- */
go_qsizef_t qsizef_add      (go_qsizef_t a, go_qsizef_t b);
go_qsizef_t qsizef_sub      (go_qsizef_t a, go_qsizef_t b);
go_qsizef_t qsizef_mul_float(go_qsizef_t s, double c);
go_qsizef_t qsizef_div_float(go_qsizef_t s, double c);

/* --- QSizeF: comparison --- */
bool qsizef_equal     (go_qsizef_t a, go_qsizef_t b);
bool qsizef_equal_size(go_qsizef_t a, go_qsize_t b);

/* --- QSizeF: conversions --- */
go_qsize_t qsizef_to_size(go_qsizef_t s);

#ifdef __cplusplus
}
#endif

#endif