#include "qurlquery_c.h"

#include <QtCore/QUrlQuery>
#include <QtCore/QUrl>
#include <QtCore/QString>
#include <QtCore/QStringList>
#include <QtCore/QByteArray>
#include <QtCore/QByteArrayView>
#include <QtCore/QList>
#include <QtCore/QPair>

#include <cstdlib>
#include <cstring>
#include <utility>

namespace {

char* copyUtf8(const QString& s, size_t* out_size) {
    const QByteArray ba = s.toUtf8();
    const size_t n = static_cast<size_t>(ba.size());
    char* out = static_cast<char*>(std::malloc(n + 1));
    if (!out) { if (out_size) *out_size = 0; return nullptr; }
    if (n) std::memcpy(out, ba.constData(), n);
    out[n] = '\0';
    if (out_size) *out_size = n;
    return out;
}

QString fromUtf8Ptr(const char* p, size_t size) {
    return QString::fromUtf8(QByteArrayView(p, static_cast<qsizetype>(size)));
}

} // namespace

extern "C" {

QUrlQuery* qurlquery_new_empty() { return new QUrlQuery(); }

QUrlQuery* qurlquery_new_from_url(const QUrl* url) { return new QUrlQuery(*url); }

QUrlQuery* qurlquery_new_from_utf8(const char* query, size_t size) {
    return new QUrlQuery(fromUtf8Ptr(query, size));
}

QUrlQuery* qurlquery_new_copy(const QUrlQuery* other) { return new QUrlQuery(*other); }
void qurlquery_delete(QUrlQuery* self) { delete self; }
void qurlquery_swap(QUrlQuery* self, QUrlQuery* other) { self->swap(*other); }

bool qurlquery_is_empty(const QUrlQuery* self) { return self->isEmpty(); }
bool qurlquery_is_detached(const QUrlQuery* self) { return self->isDetached(); }
void qurlquery_clear(QUrlQuery* self) { self->clear(); }

char* qurlquery_query(const QUrlQuery* self, int encoding, size_t* out_size) {
    return copyUtf8(self->query(QUrl::ComponentFormattingOptions(encoding)), out_size);
}
char* qurlquery_to_string(const QUrlQuery* self, int encoding, size_t* out_size) {
    return copyUtf8(self->toString(QUrl::ComponentFormattingOptions(encoding)), out_size);
}
void qurlquery_set_query_utf8(QUrlQuery* self, const char* query, size_t size) {
    self->setQuery(fromUtf8Ptr(query, size));
}

void qurlquery_set_query_delimiters(QUrlQuery* self, uint16_t value, uint16_t pair) {
    self->setQueryDelimiters(QChar(value), QChar(pair));
}
uint16_t qurlquery_query_value_delimiter(const QUrlQuery* self) {
    return self->queryValueDelimiter().unicode();
}
uint16_t qurlquery_query_pair_delimiter(const QUrlQuery* self) {
    return self->queryPairDelimiter().unicode();
}

void qurlquery_set_query_items(QUrlQuery* self,
                               const char* buffer,
                               const size_t* keyOffsets,
                               const size_t* valOffsets,
                               size_t count) {
    QList<std::pair<QString, QString>> list;
    list.reserve(static_cast<qsizetype>(count));
    for (size_t i = 0; i < count; ++i) {
        const size_t klo = keyOffsets[i];
        const size_t khi = keyOffsets[i + 1];
        const size_t vlo = valOffsets[i];
        const size_t vhi = valOffsets[i + 1];
        const char* kp = (buffer && khi > klo) ? (buffer + klo) : nullptr;
        const char* vp = (buffer && vhi > vlo) ? (buffer + vlo) : nullptr;
        list.append(std::make_pair(
            QString::fromUtf8(QByteArrayView(kp, static_cast<qsizetype>(khi - klo))),
            QString::fromUtf8(QByteArrayView(vp, static_cast<qsizetype>(vhi - vlo)))));
    }
    self->setQueryItems(list);
}

char* qurlquery_query_items(const QUrlQuery* self, int encoding,
                            size_t** out_keyOffsets, size_t** out_valOffsets,
                            size_t* out_count) {
    const auto list = self->queryItems(QUrl::ComponentFormattingOptions(encoding));
    const size_t n = static_cast<size_t>(list.size());

    // Concatenate all keys (as UTF-8), then all values (as UTF-8).
    QByteArray keys, values;
    for (size_t i = 0; i < n; ++i) {
        keys.append(list[static_cast<qsizetype>(i)].first.toUtf8());
        values.append(list[static_cast<qsizetype>(i)].second.toUtf8());
    }
    const size_t keysLen = static_cast<size_t>(keys.size());
    const size_t total = keysLen + static_cast<size_t>(values.size());

    char* buf = static_cast<char*>(std::malloc(total ? total : 1));
    size_t* ko = static_cast<size_t*>(std::malloc((n + 1) * sizeof(size_t)));
    size_t* vo = static_cast<size_t*>(std::malloc((n + 1) * sizeof(size_t)));
    if (!buf || !ko || !vo) {
        std::free(buf);
        std::free(ko);
        std::free(vo);
        *out_keyOffsets = nullptr;
        *out_valOffsets = nullptr;
        *out_count = 0;
        return nullptr;
    }
    if (keysLen) std::memcpy(buf, keys.constData(), keysLen);
    if (!values.isEmpty())
        std::memcpy(buf + keysLen, values.constData(), static_cast<size_t>(values.size()));

    size_t k = 0, v = keysLen;
    for (size_t i = 0; i < n; ++i) {
        ko[i] = k;
        vo[i] = v;
        k += static_cast<size_t>(list[static_cast<qsizetype>(i)].first.toUtf8().size());
        v += static_cast<size_t>(list[static_cast<qsizetype>(i)].second.toUtf8().size());
    }
    ko[n] = keysLen;
    vo[n] = total;

    *out_keyOffsets = ko;
    *out_valOffsets = vo;
    *out_count = n;
    return buf;
}

bool qurlquery_has_query_item_utf8(const QUrlQuery* self, const char* key, size_t keySize) {
    return self->hasQueryItem(fromUtf8Ptr(key, keySize));
}
void qurlquery_add_query_item_utf8_utf8(QUrlQuery* self,
                                        const char* key, size_t keySize,
                                        const char* value, size_t valueSize) {
    self->addQueryItem(fromUtf8Ptr(key, keySize), fromUtf8Ptr(value, valueSize));
}
void qurlquery_remove_query_item_utf8(QUrlQuery* self, const char* key, size_t keySize) {
    self->removeQueryItem(fromUtf8Ptr(key, keySize));
}
char* qurlquery_query_item_value_utf8(const QUrlQuery* self, const char* key, size_t keySize,
                                      int encoding, size_t* out_size) {
    return copyUtf8(self->queryItemValue(fromUtf8Ptr(key, keySize),
                                         QUrl::ComponentFormattingOptions(encoding)), out_size);
}

char* qurlquery_all_query_item_values_utf8(const QUrlQuery* self, const char* key, size_t keySize,
                                           int encoding, size_t** out_offsets, size_t* out_count) {
    const QStringList list = self->allQueryItemValues(fromUtf8Ptr(key, keySize),
                                                      QUrl::ComponentFormattingOptions(encoding));
    const size_t n = static_cast<size_t>(list.size());

    QByteArray joined;
    for (size_t i = 0; i < n; ++i)
        joined.append(list[static_cast<qsizetype>(i)].toUtf8());
    const size_t total = static_cast<size_t>(joined.size());

    char* buf = static_cast<char*>(std::malloc(total ? total : 1));
    size_t* offs = static_cast<size_t*>(std::malloc((n + 1) * sizeof(size_t)));
    if (!buf || !offs) {
        std::free(buf);
        std::free(offs);
        *out_offsets = nullptr;
        *out_count = 0;
        return nullptr;
    }
    if (total) std::memcpy(buf, joined.constData(), total);

    size_t pos = 0;
    offs[0] = 0;
    for (size_t i = 0; i < n; ++i) {
        pos += static_cast<size_t>(list[static_cast<qsizetype>(i)].toUtf8().size());
        offs[i + 1] = pos;
    }

    *out_offsets = offs;
    *out_count = n;
    return buf;
}

void qurlquery_remove_all_query_items_utf8(QUrlQuery* self, const char* key, size_t keySize) {
    self->removeAllQueryItems(fromUtf8Ptr(key, keySize));
}

void qurlquery_free_string(char* s) { std::free(s); }
void qurlquery_free_offsets(size_t* offs) { std::free(offs); }

} // extern "C"