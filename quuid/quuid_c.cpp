#include "quuid_c.h"

#include <QtCore/QUuid>
#include <QtCore/QString>
#include <QtCore/QByteArray>
#include <QtCore/QByteArrayView>
#include <QtCore/QSysInfo>

#include <cstdlib>
#include <cstring>

namespace {

// copyBytes: returns exactly ba.size() bytes, no NUL terminator. Used for
// QByteArray-returning getters (toByteArray, toRfc4122): Go reads those via
// C.GoBytes(p, n) with an explicit length, so NUL is neither needed nor
// wanted (it would inflate the reported length).
char* copyBytes(const QByteArray& ba, size_t* out_size) {
    const size_t n = static_cast<size_t>(ba.size());
    char* out = static_cast<char*>(std::malloc(n ? n : 1));
    if (!out) {
        if (out_size) *out_size = 0;
        return nullptr;
    }
    if (n) std::memcpy(out, ba.constData(), n);
    if (out_size) *out_size = n;
    return out;
}

// copyUtf8: returns a NUL-terminated UTF-8 buffer. Used for QString-returning
// getters (toString) where the Go side reads with C.GoString(p) — which scans
// until the first '\0'. Must allocate n+1 and write the terminator, otherwise
// C.GoString runs past the buffer into unrelated heap memory.
char* copyUtf8(const QString& s) {
    const QByteArray ba = s.toUtf8();
    const size_t n = static_cast<size_t>(ba.size());
    char* out = static_cast<char*>(std::malloc(n + 1));
    if (!out) return nullptr;
    if (n) std::memcpy(out, ba.constData(), n);
    out[n] = '\0';
    return out;
}

QSysInfo::Endian toEndian(int order) {
    return order == static_cast<int>(QSysInfo::LittleEndian)
               ? QSysInfo::LittleEndian
               : QSysInfo::BigEndian;
}

QString fromUtf8Ptr(const char* p, size_t size) {
    // QByteArrayView handles nullptr with size 0 as a null view, which
    // QString::fromUtf8 turns into a null QString — no qstrlen(nullptr).
    return QString::fromUtf8(QByteArrayView(p, static_cast<qsizetype>(size)));
}

} // namespace

extern "C" {

QUuid* quuid_new_empty() { return new QUuid(); }

QUuid* quuid_new_from_uint32(uint32_t l, uint16_t w1, uint16_t w2,
                             uint8_t b1, uint8_t b2, uint8_t b3,
                             uint8_t b4, uint8_t b5, uint8_t b6,
                             uint8_t b7, uint8_t b8) {
    return new QUuid(l, w1, w2, b1, b2, b3, b4, b5, b6, b7, b8);
}

QUuid* quuid_new_from_string(const char* utf8, size_t size) {
    return new QUuid(QUuid::fromString(fromUtf8Ptr(utf8, size)));
}

QUuid* quuid_new_from_bytes(const uint8_t bytes[16], int order) {
    return new QUuid(QUuid::fromBytes(bytes, toEndian(order)));
}

QUuid* quuid_clone(const QUuid* other) { return new QUuid(*other); }
void   quuid_delete(QUuid* self)       { delete self; }

bool quuid_equal(const QUuid* a, const QUuid* b) { return *a == *b; }
int  quuid_compare(const QUuid* a, const QUuid* b) {
    if (*a < *b) return -1;
    if (*b < *a) return  1;
    return 0;
}

char* quuid_to_string(const QUuid* self, int format) {
    return copyUtf8(self->toString(static_cast<QUuid::StringFormat>(format)));
}

char* quuid_to_byte_array(const QUuid* self, int format, size_t* out_size) {
    return copyBytes(self->toByteArray(static_cast<QUuid::StringFormat>(format)), out_size);
}

void quuid_to_bytes(const QUuid* self, int order, uint8_t out[16]) {
    const auto b = self->toBytes(toEndian(order));
    std::memcpy(out, b.data, 16);
}

char* quuid_to_rfc4122(const QUuid* self, size_t* out_size) {
    return copyBytes(self->toRfc4122(), out_size);
}

int  quuid_variant(const QUuid* self) { return static_cast<int>(self->variant()); }
int  quuid_version(const QUuid* self) { return static_cast<int>(self->version()); }
bool quuid_is_null(const QUuid* self) { return self->isNull(); }

QUuid* quuid_create_uuid() { return new QUuid(QUuid::createUuid()); }

QUuid* quuid_create_uuid_v3(const QUuid* ns, const uint8_t* data, size_t size) {
    return new QUuid(QUuid::createUuidV3(*ns, QByteArrayView(reinterpret_cast<const char*>(data),
                                                             static_cast<qsizetype>(size))));
}

QUuid* quuid_create_uuid_v5(const QUuid* ns, const uint8_t* data, size_t size) {
    return new QUuid(QUuid::createUuidV5(*ns, QByteArrayView(reinterpret_cast<const char*>(data),
                                                             static_cast<qsizetype>(size))));
}

QUuid* quuid_create_uuid_v7() { return new QUuid(QUuid::createUuidV7()); }

QUuid* quuid_from_string(const char* utf8, size_t size) {
    return new QUuid(QUuid::fromString(fromUtf8Ptr(utf8, size)));
}

QUuid* quuid_from_bytes(const uint8_t bytes[16], int order) {
    return new QUuid(QUuid::fromBytes(bytes, toEndian(order)));
}

QUuid* quuid_from_rfc4122(const uint8_t* data, size_t size) {
    return new QUuid(QUuid::fromRfc4122(QByteArrayView(reinterpret_cast<const char*>(data),
                                                       static_cast<qsizetype>(size))));
}

void quuid_free_string(char* s) { std::free(s); }

} // extern "C"