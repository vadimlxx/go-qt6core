/* C ABI for Qt 6.11 QtCore/QUrlQuery. */
#ifndef GO_QT6CORE_QURLQUERY_C_H
#define GO_QT6CORE_QURLQUERY_C_H

#include <stddef.h>
#include <stdint.h>
#include <stdbool.h>

#ifdef __cplusplus
class QUrlQuery;
class QUrl;
extern "C" {
#else
typedef struct QUrlQuery QUrlQuery;
typedef struct QUrl QUrl;
#endif

/* --- lifecycle --- */
QUrlQuery* qurlquery_new_empty(void);
QUrlQuery* qurlquery_new_from_url(const QUrl* url);
QUrlQuery* qurlquery_new_from_utf8(const char* query, size_t size);
QUrlQuery* qurlquery_new_copy(const QUrlQuery* other);
void       qurlquery_delete(QUrlQuery* self);
void       qurlquery_swap(QUrlQuery* self, QUrlQuery* other);

/* --- state --- */
bool qurlquery_is_empty(const QUrlQuery* self);
bool qurlquery_is_detached(const QUrlQuery* self);
void qurlquery_clear(QUrlQuery* self);

/* --- query string --- */
char* qurlquery_query(const QUrlQuery* self, int encoding, size_t* out_size);
char* qurlquery_to_string(const QUrlQuery* self, int encoding, size_t* out_size);
void  qurlquery_set_query_utf8(QUrlQuery* self, const char* query, size_t size);

/* --- delimiters --- */
void     qurlquery_set_query_delimiters(QUrlQuery* self, uint16_t value, uint16_t pair);
uint16_t qurlquery_query_value_delimiter(const QUrlQuery* self);
uint16_t qurlquery_query_pair_delimiter(const QUrlQuery* self);

/* --- items (flat-buffer, order-preserving on the wire) ---
 *
 * setQueryItems: buffer holds all keys concatenated then all values
 * concatenated; keyOffsets[i]..keyOffsets[i+1] is the range of key i,
 * valOffsets[i]..valOffsets[i+1] the range of value i. Both offset arrays
 * have count+1 entries.
 *
 * queryItems: C allocates *out_keyOffsets and *out_valOffsets (size_t[]
 * of count+1 entries each) and returns a malloc'd buffer holding all keys
 * followed by all values; *out_count receives the number of items.
 * Free the buffer with qurlquery_free_string, both arrays with
 * qurlquery_free_offsets.
 */
void   qurlquery_set_query_items(QUrlQuery* self,
                                 const char* buffer,
                                 const size_t* keyOffsets,
                                 const size_t* valOffsets,
                                 size_t count);
char*  qurlquery_query_items(const QUrlQuery* self, int encoding,
                             size_t** out_keyOffsets, size_t** out_valOffsets,
                             size_t* out_count);

/* --- single-item operations (Go strings) --- */
bool  qurlquery_has_query_item_utf8(const QUrlQuery* self, const char* key, size_t keySize);
void  qurlquery_add_query_item_utf8_utf8(QUrlQuery* self,
                                         const char* key, size_t keySize,
                                         const char* value, size_t valueSize);
void  qurlquery_remove_query_item_utf8(QUrlQuery* self, const char* key, size_t keySize);
char* qurlquery_query_item_value_utf8(const QUrlQuery* self, const char* key, size_t keySize,
                                      int encoding, size_t* out_size);

/* allQueryItemValues: C allocates *out_offsets (size_t[] of count+1 entries)
 * and returns a malloc'd buffer of all values concatenated; *out_count is
 * the number of values. Free buffer with qurlquery_free_string, offsets
 * with qurlquery_free_offsets. */
char* qurlquery_all_query_item_values_utf8(const QUrlQuery* self, const char* key, size_t keySize,
                                           int encoding, size_t** out_offsets, size_t* out_count);

void  qurlquery_remove_all_query_items_utf8(QUrlQuery* self, const char* key, size_t keySize);

/* --- memory --- */
void qurlquery_free_string(char* s);
void qurlquery_free_offsets(size_t* offs);

#ifdef __cplusplus
}
#endif
#endif