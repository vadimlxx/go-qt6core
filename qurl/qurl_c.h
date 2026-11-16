/* C ABI for Qt 6.11 QtCore/QUrl. */
#ifndef GO_QT6CORE_QURL_C_H
#define GO_QT6CORE_QURL_C_H

#include <stddef.h>
#include <stdint.h>
#include <stdbool.h>

#ifdef __cplusplus
class QUrl;
class QString;
class QUrlQuery;
extern "C" {
#else
typedef struct QUrl QUrl;
typedef struct QString QString;
typedef struct QUrlQuery QUrlQuery;
#endif

/* --- lifecycle --- */
QUrl* qurl_new_empty(void);
QUrl* qurl_new_from_qstring(const QString* url, int mode);
QUrl* qurl_new_from_utf8(const char* url, size_t size, int mode);
QUrl* qurl_new_copy(const QUrl* other);
void  qurl_delete(QUrl* self);
void  qurl_swap(QUrl* self, QUrl* other);

/* --- setters / getters (QString in/out) --- */
void   qurl_set_url_qstring(QUrl* self, const QString* url, int mode);
void   qurl_set_url_utf8(QUrl* self, const char* url, size_t size, int mode);

char*  qurl_url(const QUrl* self, int options, size_t* out_size);
char*  qurl_to_string(const QUrl* self, int options, size_t* out_size);
char*  qurl_to_display_string(const QUrl* self, int options, size_t* out_size);

QUrl*  qurl_adjusted(const QUrl* self, int options);

char*  qurl_to_encoded(const QUrl* self, int options, size_t* out_size);
QUrl*  qurl_from_encoded(const uint8_t* data, size_t size, int mode);

QUrl*  qurl_from_user_input_utf8(const char* input, size_t input_size,
                                 const char* workdir, size_t workdir_size,
                                 int options);
QUrl*  qurl_from_user_input_qstring(const QString* input, const QString* workdir, int options);

/* --- validity --- */
bool  qurl_is_valid(const QUrl* self);
char* qurl_error_string(const QUrl* self);
bool  qurl_is_empty(const QUrl* self);
void  qurl_clear(QUrl* self);

/* --- components --- */
void  qurl_set_scheme_qstring(QUrl* self, const QString* scheme);
void  qurl_set_scheme_utf8(QUrl* self, const char* scheme, size_t size);
char* qurl_scheme(const QUrl* self);

void  qurl_set_authority_qstring(QUrl* self, const QString* auth, int mode);
void  qurl_set_authority_utf8(QUrl* self, const char* auth, size_t size, int mode);
char* qurl_authority(const QUrl* self, int options);

void  qurl_set_user_info_qstring(QUrl* self, const QString* ui, int mode);
void  qurl_set_user_info_utf8(QUrl* self, const char* ui, size_t size, int mode);
char* qurl_user_info(const QUrl* self, int options);

void  qurl_set_user_name_qstring(QUrl* self, const QString* un, int mode);
void  qurl_set_user_name_utf8(QUrl* self, const char* un, size_t size, int mode);
char* qurl_user_name(const QUrl* self, int options);

void  qurl_set_password_qstring(QUrl* self, const QString* pw, int mode);
void  qurl_set_password_utf8(QUrl* self, const char* pw, size_t size, int mode);
char* qurl_password(const QUrl* self, int options);

void  qurl_set_host_qstring(QUrl* self, const QString* host, int mode);
void  qurl_set_host_utf8(QUrl* self, const char* host, size_t size, int mode);
char* qurl_host(const QUrl* self, int options);

void  qurl_set_port(QUrl* self, int port);
int   qurl_port(const QUrl* self, int defaultPort);

void  qurl_set_path_qstring(QUrl* self, const QString* path, int mode);
void  qurl_set_path_utf8(QUrl* self, const char* path, size_t size, int mode);
char* qurl_path(const QUrl* self, int options);
char* qurl_file_name(const QUrl* self, int options);

bool  qurl_has_query(const QUrl* self);
void  qurl_set_query_qstring(QUrl* self, const QString* query, int mode);
void  qurl_set_query_utf8(QUrl* self, const char* query, size_t size, int mode);
char* qurl_query(const QUrl* self, int options);

/* QUrlQuery interop.
 * Qt: QUrl::setQuery(const QUrlQuery&).
 * Takes ownership semantics: the pointer's lifetime is not transferred; the
 * QUrlQuery may be reused or closed immediately after the call. */
void  qurl_set_query_qurlquery(QUrl* self, const QUrlQuery* query);

bool  qurl_has_fragment(const QUrl* self);
char* qurl_fragment(const QUrl* self, int options);
void  qurl_set_fragment_qstring(QUrl* self, const QString* frag, int mode);
void  qurl_set_fragment_utf8(QUrl* self, const char* frag, size_t size, int mode);

/* --- relations --- */
QUrl* qurl_resolved(const QUrl* self, const QUrl* relative);
bool  qurl_is_relative(const QUrl* self);
bool  qurl_is_parent_of(const QUrl* self, const QUrl* url);

bool  qurl_is_local_file(const QUrl* self);
QUrl* qurl_from_local_file_qstring(const QString* path);
QUrl* qurl_from_local_file_utf8(const char* path, size_t size);
char* qurl_to_local_file(const QUrl* self);

void  qurl_detach(QUrl* self);
bool  qurl_is_detached(const QUrl* self);

bool  qurl_matches(const QUrl* self, const QUrl* other, int options);

/* --- static helpers --- */
char* qurl_from_percent_encoding(const uint8_t* ba, size_t size, size_t* out_size);
char* qurl_to_percent_encoding(const char* utf8, size_t size,
                               const uint8_t* exclude, size_t exclude_size,
                               const uint8_t* include, size_t include_size,
                               size_t* out_size);

char* qurl_from_ace(const uint8_t* ba, size_t size);
char* qurl_to_ace(const char* utf8, size_t size, size_t* out_size);

char* qurl_to_string_list(const QUrl** urls, size_t count, int options, size_t* out_size);
/* from_string_list: flat-buffer variant. `data` holds the concatenation of all
   input strings; `offsets` has count+1 entries; string i spans
   data[offsets[i] .. offsets[i+1]). This avoids passing a Go array of Go
   pointers (which cgo forbids). */
QUrl** qurl_from_string_list(const char* data, const size_t* offsets, size_t count, int mode,
                             size_t* out_count);

/* --- memory --- */
void qurl_free_string(char* s);
void qurl_free_string_array(QUrl** arr);

#ifdef __cplusplus
}
#endif
#endif