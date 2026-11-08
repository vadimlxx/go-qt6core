/* Order: version check first, then real Qt headers. */
#include "../internal/qt6core/check_qt_version.h"

#include <QtCore/qbitarray.h>
#include <QtCore/qsysinfo.h>

#include "qbitarray_c.h"

namespace {

inline QBitArray       *A(go_qbitarray_t       *p) { return reinterpret_cast<QBitArray       *>(p); }
inline const QBitArray *A(const go_qbitarray_t *p) { return reinterpret_cast<const QBitArray *>(p); }
inline go_qbitarray_t  *toA(QBitArray *p)          { return reinterpret_cast<go_qbitarray_t *>(p); }

} // namespace

extern "C" {

/* ---- Constructors ---- */
go_qbitarray_t *qbitarray_new_empty(void) {
    return toA(new QBitArray());
}
go_qbitarray_t *qbitarray_new_sized(ptrdiff_t size, bool val) {
    return toA(new QBitArray(static_cast<qsizetype>(size), val));
}
go_qbitarray_t *qbitarray_new_copy(const go_qbitarray_t *other) {
    return toA(new QBitArray(*A(other)));
}

/* ---- Lifecycle ---- */
void qbitarray_delete(go_qbitarray_t *self) {
    if (self) delete A(self);
}
void qbitarray_assign(go_qbitarray_t *self, const go_qbitarray_t *other) {
    if (self && other) *A(self) = *A(other);
}
void qbitarray_swap(go_qbitarray_t *self, go_qbitarray_t *other) {
    if (self && other) A(self)->swap(*A(other));
}

/* ---- Size / status ---- */
ptrdiff_t qbitarray_size(const go_qbitarray_t *self) {
    return static_cast<ptrdiff_t>(A(self)->size());
}
ptrdiff_t qbitarray_count_set(const go_qbitarray_t *self, bool on) {
    return static_cast<ptrdiff_t>(A(self)->count(on));
}
bool qbitarray_is_empty(const go_qbitarray_t *self) {
    return A(self)->isEmpty();
}
bool qbitarray_is_null(const go_qbitarray_t *self) {
    return A(self)->isNull();
}

/* ---- Modifiers ---- */
void qbitarray_resize(go_qbitarray_t *self, ptrdiff_t size) {
    A(self)->resize(static_cast<qsizetype>(size));
}
void qbitarray_detach(go_qbitarray_t *self) {
    A(self)->detach();
}
bool qbitarray_is_detached(const go_qbitarray_t *self) {
    return A(self)->isDetached();
}
void qbitarray_clear(go_qbitarray_t *self) {
    A(self)->clear();
}

/* ---- Bit operations ---- */
bool qbitarray_test_bit(const go_qbitarray_t *self, ptrdiff_t i) {
    return A(self)->testBit(static_cast<qsizetype>(i));
}
void qbitarray_set_bit(go_qbitarray_t *self, ptrdiff_t i) {
    A(self)->setBit(static_cast<qsizetype>(i));
}
void qbitarray_set_bit_val(go_qbitarray_t *self, ptrdiff_t i, bool val) {
    A(self)->setBit(static_cast<qsizetype>(i), val);
}
void qbitarray_clear_bit(go_qbitarray_t *self, ptrdiff_t i) {
    A(self)->clearBit(static_cast<qsizetype>(i));
}
bool qbitarray_toggle_bit(go_qbitarray_t *self, ptrdiff_t i) {
    return A(self)->toggleBit(static_cast<qsizetype>(i));
}

/* ---- Bitwise assign ---- */
void qbitarray_and_assign(go_qbitarray_t *self, const go_qbitarray_t *other) {
    *A(self) &= *A(other);
}
void qbitarray_or_assign(go_qbitarray_t *self, const go_qbitarray_t *other) {
    *A(self) |= *A(other);
}
void qbitarray_xor_assign(go_qbitarray_t *self, const go_qbitarray_t *other) {
    *A(self) ^= *A(other);
}

/* ---- Bitwise non-assign ---- */
go_qbitarray_t *qbitarray_op_and(const go_qbitarray_t *a, const go_qbitarray_t *b) {
    return toA(new QBitArray(*A(a) & *A(b)));
}
go_qbitarray_t *qbitarray_op_or(const go_qbitarray_t *a, const go_qbitarray_t *b) {
    return toA(new QBitArray(*A(a) | *A(b)));
}
go_qbitarray_t *qbitarray_op_xor(const go_qbitarray_t *a, const go_qbitarray_t *b) {
    return toA(new QBitArray(*A(a) ^ *A(b)));
}
go_qbitarray_t *qbitarray_inverted(const go_qbitarray_t *a) {
    return toA(new QBitArray(~*A(a)));
}

/* ---- Fill / truncate ---- */
void qbitarray_fill_all(go_qbitarray_t *self, bool val) {
    A(self)->fill(val);
}
void qbitarray_fill_sized(go_qbitarray_t *self, bool val, ptrdiff_t size) {
    A(self)->fill(val, static_cast<qsizetype>(size));
}
void qbitarray_fill_range(go_qbitarray_t *self, bool val, ptrdiff_t first, ptrdiff_t last) {
    A(self)->fill(val, static_cast<qsizetype>(first), static_cast<qsizetype>(last));
}
void qbitarray_truncate(go_qbitarray_t *self, ptrdiff_t pos) {
    A(self)->truncate(static_cast<qsizetype>(pos));
}

/* ---- Raw bits ---- */
const char *qbitarray_bits(const go_qbitarray_t *self, ptrdiff_t *out_bytes) {
    const QBitArray *a = A(self);
    if (a->isEmpty()) {
        *out_bytes = 0;
        return nullptr;
    }
    *out_bytes = static_cast<ptrdiff_t>((a->size() + 7) / 8);
    return a->bits();
}
go_qbitarray_t *qbitarray_from_bits(const char *data, ptrdiff_t bit_len) {
    return toA(new QBitArray(QBitArray::fromBits(data, static_cast<qsizetype>(bit_len))));
}
uint32_t qbitarray_to_uint32(const go_qbitarray_t *self, bool littleEndian, bool *ok) {
    const QSysInfo::Endian e = littleEndian ? QSysInfo::LittleEndian : QSysInfo::BigEndian;
    return static_cast<uint32_t>(A(self)->toUInt32(e, ok));
}

/* ---- Comparison ---- */
bool qbitarray_equal(const go_qbitarray_t *self, const go_qbitarray_t *other) {
    return *A(self) == *A(other);
}

} // extern "C"