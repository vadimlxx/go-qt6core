/* C ABI for Qt 6.11 QtCore/QUuid. */
#ifndef GO_QT6CORE_QUUID_C_H
#define GO_QT6CORE_QUUID_C_H

#include <stddef.h>
#include <stdint.h>
#include <stdbool.h>

#ifdef __cplusplus
class QUuid;
extern "C" {
#else
typedef struct QUuid QUuid;
#endif

/* --- lifecycle --- */
QUuid* quuid_new_empty(void);
QUuid* quuid_new_from_uint32(uint32_t l, uint16_t w1, uint16_t w2,
                             uint8_t b1, uint8_t b2, uint8_t b3,
                             uint8_t b4, uint8_t b5, uint8_t b6,
                             uint8_t b7, uint8_t b8);
QUuid* quuid_new_from_string(const char* utf8, size_t size);
QUuid* quuid_new_from_bytes(const uint8_t bytes[16], int order);
QUuid* quuid_clone(const QUuid* other);
void   quuid_delete(QUuid* self);
bool   quuid_equal(const QUuid* a, const QUuid* b);
int    quuid_compare(const QUuid* a, const QUuid* b);

/* --- conversion --- */
char* quuid_to_string(const QUuid* self, int format);
char* quuid_to_byte_array(const QUuid* self, int format, size_t* out_size);
void  quuid_to_bytes(const QUuid* self, int order, uint8_t out[16]);
char* quuid_to_rfc4122(const QUuid* self, size_t* out_size);

/* --- properties --- */
int  quuid_variant(const QUuid* self);
int  quuid_version(const QUuid* self);
bool quuid_is_null(const QUuid* self);

/* --- static creators --- */
QUuid* quuid_create_uuid(void);
QUuid* quuid_create_uuid_v3(const QUuid* ns, const uint8_t* data, size_t size);
QUuid* quuid_create_uuid_v5(const QUuid* ns, const uint8_t* data, size_t size);
QUuid* quuid_create_uuid_v7(void);

/* --- static parsers --- */
QUuid* quuid_from_string(const char* utf8, size_t size);
QUuid* quuid_from_bytes(const uint8_t bytes[16], int order);
QUuid* quuid_from_rfc4122(const uint8_t* data, size_t size);

/* --- memory --- */
void quuid_free_string(char* s);

#ifdef __cplusplus
}
#endif
#endif