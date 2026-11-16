#include "qurl_c.h"

#include <QtCore/QUrl>
#include <QtCore/QUrlQuery>
#include <QtCore/QString>
#include <QtCore/QStringList>
#include <QtCore/QByteArray>
#include <QtCore/QByteArrayView>
#include <QtCore/QList>

#include <cstdlib>
#include <cstring>

namespace {

char* copyUtf8(const QString& s) {
    const QByteArray ba = s.toUtf8();
    const size_t n = static_cast<size_t>(ba.size());
    char* out = static_cast<char*>(std::malloc(n + 1));
    if (!out) return nullptr;
    if (n) std::memcpy(out, ba.constData(), n);
    out[n] = '\0';
    return out;
}

char* copyBytes(const QByteArray& ba, size_t* out_size) {
    const size_t n = static_cast<size_t>(ba.size());
    char* out = static_cast<char*>(std::malloc(n ? n : 1));
    if (!out) { if (out_size) *out_size = 0; return nullptr; }
    if (n) std::memcpy(out, ba.constData(), n);
    if (out_size) *out_size = n;
    return out;
}

QString fromUtf8Ptr(const char* p, size_t size) {
    return QString::fromUtf8(QByteArrayView(p, static_cast<qsizetype>(size)));
}

QUrl* newUrl(const QString& s, int mode) {
    return new QUrl(s, static_cast<QUrl::ParsingMode>(mode));
}

} // namespace

extern "C" {

QUrl* qurl_new_empty() { return new QUrl(); }

QUrl* qurl_new_from_qstring(const QString* url, int mode) {
    return newUrl(*url, mode);
}

QUrl* qurl_new_from_utf8(const char* url, size_t size, int mode) {
    return newUrl(fromUtf8Ptr(url, size), mode);
}

QUrl* qurl_new_copy(const QUrl* other) { return new QUrl(*other); }
void  qurl_delete(QUrl* self) { delete self; }
void  qurl_swap(QUrl* self, QUrl* other) { self->swap(*other); }

void qurl_set_url_qstring(QUrl* self, const QString* url, int mode) {
    self->setUrl(*url, static_cast<QUrl::ParsingMode>(mode));
}
void qurl_set_url_utf8(QUrl* self, const char* url, size_t size, int mode) {
    self->setUrl(fromUtf8Ptr(url, size), static_cast<QUrl::ParsingMode>(mode));
}

char* qurl_url(const QUrl* self, int options, size_t* out_size) {
    return copyBytes(self->url(QUrl::FormattingOptions(options)).toUtf8(), out_size);
}
char* qurl_to_string(const QUrl* self, int options, size_t* out_size) {
    return copyBytes(self->toString(QUrl::FormattingOptions(options)).toUtf8(), out_size);
}
char* qurl_to_display_string(const QUrl* self, int options, size_t* out_size) {
    return copyBytes(self->toDisplayString(QUrl::FormattingOptions(options)).toUtf8(), out_size);
}

QUrl* qurl_adjusted(const QUrl* self, int options) {
    return new QUrl(self->adjusted(QUrl::FormattingOptions(options)));
}

char* qurl_to_encoded(const QUrl* self, int options, size_t* out_size) {
    return copyBytes(self->toEncoded(QUrl::FormattingOptions(options)), out_size);
}
QUrl* qurl_from_encoded(const uint8_t* data, size_t size, int mode) {
    return new QUrl(QUrl::fromEncoded(QByteArrayView(reinterpret_cast<const char*>(data),
                                                     static_cast<qsizetype>(size)),
                                      static_cast<QUrl::ParsingMode>(mode)));
}

QUrl* qurl_from_user_input_utf8(const char* input, size_t input_size,
                                const char* workdir, size_t workdir_size,
                                int options) {
    return new QUrl(QUrl::fromUserInput(fromUtf8Ptr(input, input_size),
                                        fromUtf8Ptr(workdir, workdir_size),
                                        QUrl::UserInputResolutionOptions(options)));
}
QUrl* qurl_from_user_input_qstring(const QString* input, const QString* workdir, int options) {
    return new QUrl(QUrl::fromUserInput(*input, *workdir,
                                        QUrl::UserInputResolutionOptions(options)));
}

bool qurl_is_valid(const QUrl* self) { return self->isValid(); }
char* qurl_error_string(const QUrl* self) { return copyUtf8(self->errorString()); }
bool qurl_is_empty(const QUrl* self) { return self->isEmpty(); }
void qurl_clear(QUrl* self) { self->clear(); }

void qurl_set_scheme_qstring(QUrl* self, const QString* scheme) { self->setScheme(*scheme); }
void qurl_set_scheme_utf8(QUrl* self, const char* scheme, size_t size) {
    self->setScheme(fromUtf8Ptr(scheme, size));
}
char* qurl_scheme(const QUrl* self) { return copyUtf8(self->scheme()); }

void qurl_set_authority_qstring(QUrl* self, const QString* auth, int mode) {
    self->setAuthority(*auth, static_cast<QUrl::ParsingMode>(mode));
}
void qurl_set_authority_utf8(QUrl* self, const char* auth, size_t size, int mode) {
    self->setAuthority(fromUtf8Ptr(auth, size), static_cast<QUrl::ParsingMode>(mode));
}
char* qurl_authority(const QUrl* self, int options) {
    return copyUtf8(self->authority(QUrl::ComponentFormattingOptions(options)));
}

void qurl_set_user_info_qstring(QUrl* self, const QString* ui, int mode) {
    self->setUserInfo(*ui, static_cast<QUrl::ParsingMode>(mode));
}
void qurl_set_user_info_utf8(QUrl* self, const char* ui, size_t size, int mode) {
    self->setUserInfo(fromUtf8Ptr(ui, size), static_cast<QUrl::ParsingMode>(mode));
}
char* qurl_user_info(const QUrl* self, int options) {
    return copyUtf8(self->userInfo(QUrl::ComponentFormattingOptions(options)));
}

void qurl_set_user_name_qstring(QUrl* self, const QString* un, int mode) {
    self->setUserName(*un, static_cast<QUrl::ParsingMode>(mode));
}
void qurl_set_user_name_utf8(QUrl* self, const char* un, size_t size, int mode) {
    self->setUserName(fromUtf8Ptr(un, size), static_cast<QUrl::ParsingMode>(mode));
}
char* qurl_user_name(const QUrl* self, int options) {
    return copyUtf8(self->userName(QUrl::ComponentFormattingOptions(options)));
}

void qurl_set_password_qstring(QUrl* self, const QString* pw, int mode) {
    self->setPassword(*pw, static_cast<QUrl::ParsingMode>(mode));
}
void qurl_set_password_utf8(QUrl* self, const char* pw, size_t size, int mode) {
    self->setPassword(fromUtf8Ptr(pw, size), static_cast<QUrl::ParsingMode>(mode));
}
char* qurl_password(const QUrl* self, int options) {
    return copyUtf8(self->password(QUrl::ComponentFormattingOptions(options)));
}

void qurl_set_host_qstring(QUrl* self, const QString* host, int mode) {
    self->setHost(*host, static_cast<QUrl::ParsingMode>(mode));
}
void qurl_set_host_utf8(QUrl* self, const char* host, size_t size, int mode) {
    self->setHost(fromUtf8Ptr(host, size), static_cast<QUrl::ParsingMode>(mode));
}
char* qurl_host(const QUrl* self, int options) {
    return copyUtf8(self->host(QUrl::ComponentFormattingOptions(options)));
}

void qurl_set_port(QUrl* self, int port) { self->setPort(port); }
int  qurl_port(const QUrl* self, int defaultPort) { return self->port(defaultPort); }

void qurl_set_path_qstring(QUrl* self, const QString* path, int mode) {
    self->setPath(*path, static_cast<QUrl::ParsingMode>(mode));
}
void qurl_set_path_utf8(QUrl* self, const char* path, size_t size, int mode) {
    self->setPath(fromUtf8Ptr(path, size), static_cast<QUrl::ParsingMode>(mode));
}
char* qurl_path(const QUrl* self, int options) {
    return copyUtf8(self->path(QUrl::ComponentFormattingOptions(options)));
}
char* qurl_file_name(const QUrl* self, int options) {
    return copyUtf8(self->fileName(QUrl::ComponentFormattingOptions(options)));
}

bool qurl_has_query(const QUrl* self) { return self->hasQuery(); }
void qurl_set_query_qstring(QUrl* self, const QString* query, int mode) {
    self->setQuery(*query, static_cast<QUrl::ParsingMode>(mode));
}
void qurl_set_query_utf8(QUrl* self, const char* query, size_t size, int mode) {
    self->setQuery(fromUtf8Ptr(query, size), static_cast<QUrl::ParsingMode>(mode));
}
char* qurl_query(const QUrl* self, int options) {
    return copyUtf8(self->query(QUrl::ComponentFormattingOptions(options)));
}

void qurl_set_query_qurlquery(QUrl* self, const QUrlQuery* query) {
    self->setQuery(*query);
}

bool qurl_has_fragment(const QUrl* self) { return self->hasFragment(); }
char* qurl_fragment(const QUrl* self, int options) {
    return copyUtf8(self->fragment(QUrl::ComponentFormattingOptions(options)));
}
void qurl_set_fragment_qstring(QUrl* self, const QString* frag, int mode) {
    self->setFragment(*frag, static_cast<QUrl::ParsingMode>(mode));
}
void qurl_set_fragment_utf8(QUrl* self, const char* frag, size_t size, int mode) {
    self->setFragment(fromUtf8Ptr(frag, size), static_cast<QUrl::ParsingMode>(mode));
}

QUrl* qurl_resolved(const QUrl* self, const QUrl* relative) {
    return new QUrl(self->resolved(*relative));
}
bool qurl_is_relative(const QUrl* self) { return self->isRelative(); }
bool qurl_is_parent_of(const QUrl* self, const QUrl* url) { return self->isParentOf(*url); }

bool qurl_is_local_file(const QUrl* self) { return self->isLocalFile(); }
QUrl* qurl_from_local_file_qstring(const QString* path) {
    return new QUrl(QUrl::fromLocalFile(*path));
}
QUrl* qurl_from_local_file_utf8(const char* path, size_t size) {
    return new QUrl(QUrl::fromLocalFile(fromUtf8Ptr(path, size)));
}
char* qurl_to_local_file(const QUrl* self) { return copyUtf8(self->toLocalFile()); }

void qurl_detach(QUrl* self) { self->detach(); }
bool qurl_is_detached(const QUrl* self) { return self->isDetached(); }

bool qurl_matches(const QUrl* self, const QUrl* other, int options) {
    return self->matches(*other, QUrl::FormattingOptions(options));
}

char* qurl_from_percent_encoding(const uint8_t* ba, size_t size, size_t* out_size) {
    const QString s = QUrl::fromPercentEncoding(
        QByteArray(reinterpret_cast<const char*>(ba), static_cast<qsizetype>(size)));
    return copyBytes(s.toUtf8(), out_size);
}

char* qurl_to_percent_encoding(const char* utf8, size_t size,
                               const uint8_t* exclude, size_t exclude_size,
                               const uint8_t* include, size_t include_size,
                               size_t* out_size) {
    const QByteArray ex(reinterpret_cast<const char*>(exclude), static_cast<qsizetype>(exclude_size));
    const QByteArray in(reinterpret_cast<const char*>(include), static_cast<qsizetype>(include_size));
    return copyBytes(QUrl::toPercentEncoding(fromUtf8Ptr(utf8, size), ex, in), out_size);
}

char* qurl_from_ace(const uint8_t* ba, size_t size) {
    return copyUtf8(QUrl::fromAce(
        QByteArray(reinterpret_cast<const char*>(ba), static_cast<qsizetype>(size))));
}
char* qurl_to_ace(const char* utf8, size_t size, size_t* out_size) {
    return copyBytes(QUrl::toAce(fromUtf8Ptr(utf8, size)), out_size);
}

char* qurl_to_string_list(const QUrl** urls, size_t count, int options, size_t* out_size) {
    QList<QUrl> list;
    list.reserve(static_cast<qsizetype>(count));
    for (size_t i = 0; i < count; ++i)
        list.append(*urls[i]);
    return copyBytes(QUrl::toStringList(list, QUrl::FormattingOptions(options))
                         .join(QLatin1Char('\n')).toUtf8(),
                     out_size);
}

QUrl** qurl_from_string_list(const char* data, const size_t* offsets, size_t count, int mode,
                             size_t* out_count) {
    QStringList sl;
    sl.reserve(static_cast<qsizetype>(count));
    for (size_t i = 0; i < count; ++i) {
        const size_t lo = offsets[i];
        const size_t hi = offsets[i + 1];
        const char* p = (data && hi > lo) ? (data + lo) : nullptr;
        sl.append(QString::fromUtf8(QByteArrayView(p, static_cast<qsizetype>(hi - lo))));
    }
    const QList<QUrl> urls = QUrl::fromStringList(sl, static_cast<QUrl::ParsingMode>(mode));
    const size_t n = static_cast<size_t>(urls.size());
    auto** arr = new QUrl*[n ? n : 1];
    for (size_t i = 0; i < n; ++i)
        arr[i] = new QUrl(urls[static_cast<qsizetype>(i)]);
    *out_count = n;
    return arr;
}

void qurl_free_string(char* s) { std::free(s); }
void qurl_free_string_array(QUrl** arr) { delete[] arr; }

} // extern "C"