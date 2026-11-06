/* Order: version check first, then real Qt headers. */
#include "../internal/qt6core/check_qt_version.h"

#include <QAnyStringView>
#include <QString>            /* for QAnyStringView::toString() and ::visit() definitions */
#include <QStringView>
#include <QLatin1StringView>
#include <QUtf8StringView>

#include <type_traits>

#include "qanystringview_c.h"

namespace {

inline QAnyStringView makeAny(const void *data, ptrdiff_t sz, int tag) noexcept {
    switch (tag) {
    case 0:
        return QAnyStringView(QUtf8StringView(static_cast<const char *>(data), sz));
    case 1:
        return QAnyStringView(QLatin1StringView(static_cast<const char *>(data), sz));
    case 2:
        return QAnyStringView(QStringView(static_cast<const char16_t *>(data), sz));
    default:
        return QAnyStringView();
    }
}

/* Writes the (data, size, tag) triple of a QAnyStringView through the visitor. */
inline void writeOut(const QAnyStringView &v,
                     const void **out_data, ptrdiff_t *out_size, int *out_tag) {
    *out_data = nullptr;
    *out_size = 0;
    *out_tag = 0;
    v.visit([out_data, out_size, out_tag](auto sv) {
        using T = std::decay_t<decltype(sv)>;
        *out_data = static_cast<const void *>(sv.data());
        *out_size = static_cast<ptrdiff_t>(sv.size());
        if constexpr (std::is_same_v<T, QStringView>)
            *out_tag = 2;
        else if constexpr (std::is_same_v<T, QLatin1StringView>)
            *out_tag = 1;
        else
            *out_tag = 0;
    });
}

inline Qt::CaseSensitivity CS(int cs) noexcept {
    return cs == 0 ? Qt::CaseSensitive : Qt::CaseInsensitive;
}

} // namespace

extern "C" {

/* --- default (empty) constructor --- */
void qanystringview_default(const void **out_data, ptrdiff_t *out_size, int *out_tag) {
    writeOut(QAnyStringView(), out_data, out_size, out_tag);
}

/* --- introspection --- */
ptrdiff_t qanystringview_size(const void *data, ptrdiff_t size, int tag) {
    return static_cast<ptrdiff_t>(makeAny(data, size, tag).size());
}
bool qanystringview_is_null(const void *data, ptrdiff_t size, int tag) {
    return makeAny(data, size, tag).isNull();
}
bool qanystringview_is_empty(const void *data, ptrdiff_t size, int tag) {
    return makeAny(data, size, tag).isEmpty();
}
bool qanystringview_is_utf8(const void *data, ptrdiff_t size, int tag) {
    return makeAny(data, size, tag).visit([](auto sv) {
        using T = std::decay_t<decltype(sv)>;
        return std::is_same_v<T, QUtf8StringView>;
    });
}
bool qanystringview_is_latin1(const void *data, ptrdiff_t size, int tag) {
    return makeAny(data, size, tag).visit([](auto sv) {
        using T = std::decay_t<decltype(sv)>;
        return std::is_same_v<T, QLatin1StringView>;
    });
}
bool qanystringview_is_utf16(const void *data, ptrdiff_t size, int tag) {
    return makeAny(data, size, tag).visit([](auto sv) {
        using T = std::decay_t<decltype(sv)>;
        return std::is_same_v<T, QStringView>;
    });
}

/* --- element access --- */
uint16_t qanystringview_front(const void *data, ptrdiff_t size, int tag) {
    return static_cast<uint16_t>(makeAny(data, size, tag).front().unicode());
}
uint16_t qanystringview_back(const void *data, ptrdiff_t size, int tag) {
    return static_cast<uint16_t>(makeAny(data, size, tag).back().unicode());
}

/* --- slicing --- */
void qanystringview_left(const void *data, ptrdiff_t size, int tag, ptrdiff_t n,
                         const void **out_data, ptrdiff_t *out_size, int *out_tag) {
    writeOut(makeAny(data, size, tag).left(static_cast<qsizetype>(n)), out_data, out_size, out_tag);
}
void qanystringview_right(const void *data, ptrdiff_t size, int tag, ptrdiff_t n,
                          const void **out_data, ptrdiff_t *out_size, int *out_tag) {
    writeOut(makeAny(data, size, tag).right(static_cast<qsizetype>(n)), out_data, out_size, out_tag);
}
void qanystringview_mid(const void *data, ptrdiff_t size, int tag, ptrdiff_t pos, ptrdiff_t n,
                        const void **out_data, ptrdiff_t *out_size, int *out_tag) {
    writeOut(makeAny(data, size, tag).mid(static_cast<qsizetype>(pos), static_cast<qsizetype>(n)),
             out_data, out_size, out_tag);
}
void qanystringview_first(const void *data, ptrdiff_t size, int tag, ptrdiff_t n,
                          const void **out_data, ptrdiff_t *out_size, int *out_tag) {
    writeOut(makeAny(data, size, tag).first(static_cast<qsizetype>(n)), out_data, out_size, out_tag);
}
void qanystringview_last(const void *data, ptrdiff_t size, int tag, ptrdiff_t n,
                         const void **out_data, ptrdiff_t *out_size, int *out_tag) {
    writeOut(makeAny(data, size, tag).last(static_cast<qsizetype>(n)), out_data, out_size, out_tag);
}
void qanystringview_sliced(const void *data, ptrdiff_t size, int tag, ptrdiff_t pos,
                           const void **out_data, ptrdiff_t *out_size, int *out_tag) {
    writeOut(makeAny(data, size, tag).sliced(static_cast<qsizetype>(pos)), out_data, out_size, out_tag);
}
void qanystringview_sliced_n(const void *data, ptrdiff_t size, int tag, ptrdiff_t pos, ptrdiff_t n,
                             const void **out_data, ptrdiff_t *out_size, int *out_tag) {
    writeOut(makeAny(data, size, tag).sliced(static_cast<qsizetype>(pos), static_cast<qsizetype>(n)),
             out_data, out_size, out_tag);
}
void qanystringview_chopped(const void *data, ptrdiff_t size, int tag, ptrdiff_t n,
                            const void **out_data, ptrdiff_t *out_size, int *out_tag) {
    writeOut(makeAny(data, size, tag).chopped(static_cast<qsizetype>(n)), out_data, out_size, out_tag);
}

/* --- dispatch to concrete type --- */
void qanystringview_as_qstringview(const void *data, ptrdiff_t size, int tag,
                                   const uint16_t **out_data, ptrdiff_t *out_size) {
    *out_data = nullptr;
    *out_size = 0;
    makeAny(data, size, tag).visit([out_data, out_size](auto sv) {
        using T = std::decay_t<decltype(sv)>;
        if constexpr (std::is_same_v<T, QStringView>) {
            *out_data = reinterpret_cast<const uint16_t *>(sv.data());
            *out_size = static_cast<ptrdiff_t>(sv.size());
        }
    });
}
void qanystringview_as_qlatin1stringview(const void *data, ptrdiff_t size, int tag,
                                         const char **out_data, ptrdiff_t *out_size) {
    *out_data = nullptr;
    *out_size = 0;
    makeAny(data, size, tag).visit([out_data, out_size](auto sv) {
        using T = std::decay_t<decltype(sv)>;
        if constexpr (std::is_same_v<T, QLatin1StringView>) {
            *out_data = sv.data();
            *out_size = static_cast<ptrdiff_t>(sv.size());
        }
    });
}
void qanystringview_as_qutf8stringview(const void *data, ptrdiff_t size, int tag,
                                       const char **out_data, ptrdiff_t *out_size) {
    *out_data = nullptr;
    *out_size = 0;
    makeAny(data, size, tag).visit([out_data, out_size](auto sv) {
        using T = std::decay_t<decltype(sv)>;
        if constexpr (std::is_same_v<T, QUtf8StringView>) {
            *out_data = reinterpret_cast<const char *>(sv.data());
            *out_size = static_cast<ptrdiff_t>(sv.size());
        }
    });
}

/* --- comparison --- */
int qanystringview_compare(const void *ldata, ptrdiff_t lsize, int ltag,
                           const void *rdata, ptrdiff_t rsize, int rtag, int cs) {
    return QAnyStringView::compare(makeAny(ldata, lsize, ltag), makeAny(rdata, rsize, rtag), CS(cs));
}
bool qanystringview_equal(const void *ldata, ptrdiff_t lsize, int ltag,
                          const void *rdata, ptrdiff_t rsize, int rtag) {
    return QAnyStringView::equal(makeAny(ldata, lsize, ltag), makeAny(rdata, rsize, rtag));
}

/* --- conversion to QString --- */
void qanystringview_to_qstring(const void *data, ptrdiff_t size, int tag, void **out_qstring) {
    *out_qstring = static_cast<void *>(new QString(makeAny(data, size, tag).toString()));
}

} // extern "C"
