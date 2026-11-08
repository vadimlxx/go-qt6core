/* Order: version check first, then real Qt headers. */
#include "../internal/qt6core/check_qt_version.h"

/* qstringconverter.h declares QStringEncoder, QStringDecoder, QStringConverter, and also
 * transitively includes qstringconverter_base.h, qstring.h, qstringbuilder.h. That is the
 * canonical header for this class. */
#include <QtCore/qstringconverter.h>

#include "qstringencoder_c.h"

namespace {

/* QSE_Helper exposes the protected QStringEncoder(const Interface*) ctor and the protected
 * destructor inherited from QStringConverter. It also re-exports QStringConverter::Interface
 * to public scope so our C functions can static_cast the opaque Interface pointer.
 *
 * All inherited ctors (default / Encoding+Flags / QAnyStringView+Flags) are public in
 * QStringEncoder — inherited via `using QStringEncoder::QStringEncoder;`. The Interface* ctor
 * is overridden here explicitly as public. */
struct QSE_Helper : QStringEncoder {
    using Interface = QStringConverter::Interface;
    using QStringEncoder::QStringEncoder;

    QSE_Helper() : QStringEncoder() {}
    QSE_Helper(const Interface *i) : QStringEncoder(i) {}

    static void destroy(QStringEncoder *p) { delete static_cast<QSE_Helper *>(p); }
};

inline QStringEncoder       *C(go_qse_t       *p) { return reinterpret_cast<QStringEncoder       *>(p); }
inline const QStringEncoder *C(const go_qse_t *p) { return reinterpret_cast<const QStringEncoder *>(p); }
inline go_qse_t *toC(QStringEncoder *p)           { return reinterpret_cast<go_qse_t *>(p); }

inline QAnyStringView makeAny(const void *data, ptrdiff_t size, int kind) noexcept {
    switch (kind) {
    case 0:
        return QAnyStringView(static_cast<const char *>(data), static_cast<qsizetype>(size));
    case 1:
        return QAnyStringView(QLatin1StringView(static_cast<const char *>(data),
                                                 static_cast<qsizetype>(size)));
    case 2:
        return QAnyStringView(QStringView(static_cast<const char16_t *>(data),
                                          static_cast<qsizetype>(size)));
    }
    return QAnyStringView();
}

} // namespace

extern "C" {

/* ---- Constructors ---- */
go_qse_t *qstringencoder_new_default(void) {
    return toC(new QSE_Helper());
}
go_qse_t *qstringencoder_new_from_encoding(int encoding, int flags) {
    return toC(new QSE_Helper(static_cast<QStringConverter::Encoding>(encoding),
                              QStringConverter::Flags(flags)));
}
go_qse_t *qstringencoder_new_from_any(const void *data, ptrdiff_t size, int kind, int flags) {
    return toC(new QSE_Helper(makeAny(data, size, kind), QStringConverter::Flags(flags)));
}
go_qse_t *qstringencoder_new_from_iface(const void *iface) {
    return toC(new QSE_Helper(static_cast<const QSE_Helper::Interface *>(iface)));
}

/* ---- Lifecycle ---- */
void qstringencoder_delete(go_qse_t *self) {
    if (self) QSE_Helper::destroy(C(self));
}

/* ---- Encoding ---- */
void qstringencoder_encode(go_qse_t *self,
    const uint16_t *in_data, ptrdiff_t in_size,
    void **out_qba, ptrdiff_t *out_size) {
    QStringView in(reinterpret_cast<const char16_t *>(in_data),
                   static_cast<qsizetype>(in_size));
    /* encode() returns DecodedData<QStringView>; the implicit conversion to QByteArray
     * calls encodeAsByteArray(), which handles both valid and invalid (iface == nullptr)
     * encoders — for invalid it increments state.invalidChars and returns an empty QByteArray. */
    auto *ba = new QByteArray(C(self)->encode(in));
    *out_qba = static_cast<void *>(ba);
    *out_size = static_cast<ptrdiff_t>(ba->size());
}

/* ---- Required space ---- */
ptrdiff_t qstringencoder_required_space(const go_qse_t *self, ptrdiff_t input_length) {
    return static_cast<ptrdiff_t>(C(self)->requiredSpace(static_cast<qsizetype>(input_length)));
}

/* ---- Append to buffer ---- */
ptrdiff_t qstringencoder_append_to_buffer(go_qse_t *self,
    void *buf, ptrdiff_t buf_len,
    const uint16_t *in_data, ptrdiff_t in_size) {
    (void)buf_len; /* QStringEncoder::appendToBuffer writes at most fromUtf16Len(input) bytes */
    char *out = static_cast<char *>(buf);
    QStringView in(reinterpret_cast<const char16_t *>(in_data),
                   static_cast<qsizetype>(in_size));
    char *next = C(self)->appendToBuffer(out, in);
    return static_cast<ptrdiff_t>(next - out);
}

/* ---- Finalize ---- */
void qstringencoder_finalize(go_qse_t *self,
    void *buf, ptrdiff_t buf_len,
    ptrdiff_t *out_next, int16_t *out_invalid, int *out_err) {
    char *out = static_cast<char *>(buf);
    auto r = C(self)->finalize(out, static_cast<qsizetype>(buf_len));
    *out_next    = static_cast<ptrdiff_t>(r.next - out);
    *out_invalid = static_cast<int16_t>(r.invalidChars);
    *out_err     = static_cast<int>(r.error);
}

} // extern "C"