/* Order: version check first, then real Qt headers. */
#include "../internal/qt6core/check_qt_version.h"

/* qstringconverter.h declares QStringEncoder, QStringDecoder, QStringConverter, and also
 * transitively includes qstringconverter_base.h, qstring.h, qstringbuilder.h. That is the
 * canonical header for this class. */
#include <QtCore/qstringconverter.h>

#include "qstringdecoder_c.h"

namespace {

/* QSD_Helper exposes the protected QStringDecoder(const Interface*) ctor and the protected
 * destructor inherited from QStringConverter. It also re-exports QStringConverter::Interface
 * to public scope so our C functions can static_cast the opaque Interface pointer.
 *
 * All inherited ctors (default / Encoding+Flags / QAnyStringView+Flags) are public in
 * QStringDecoder — inherited via `using QStringDecoder::QStringDecoder;`. The Interface* ctor
 * is overridden here explicitly as public. */
struct QSD_Helper : QStringDecoder {
    using Interface = QStringConverter::Interface;
    using QStringDecoder::QStringDecoder;

    QSD_Helper() : QStringDecoder() {}
    QSD_Helper(const Interface *i) : QStringDecoder(i) {}

    static void destroy(QStringDecoder *p) { delete static_cast<QSD_Helper *>(p); }
};

inline QStringDecoder       *D(go_qsd_t       *p) { return reinterpret_cast<QStringDecoder       *>(p); }
inline const QStringDecoder *D(const go_qsd_t *p) { return reinterpret_cast<const QStringDecoder *>(p); }
inline go_qsd_t *toD(QStringDecoder *p)           { return reinterpret_cast<go_qsd_t *>(p); }

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
go_qsd_t *qstringdecoder_new_default(void) {
    return toD(new QSD_Helper());
}
go_qsd_t *qstringdecoder_new_from_encoding(int encoding, int flags) {
    return toD(new QSD_Helper(static_cast<QStringConverter::Encoding>(encoding),
                              QStringConverter::Flags(flags)));
}
go_qsd_t *qstringdecoder_new_from_any(const void *data, ptrdiff_t size, int kind, int flags) {
    return toD(new QSD_Helper(makeAny(data, size, kind), QStringConverter::Flags(flags)));
}
go_qsd_t *qstringdecoder_new_from_iface(const void *iface) {
    return toD(new QSD_Helper(static_cast<const QSD_Helper::Interface *>(iface)));
}

/* ---- Lifecycle ---- */
void qstringdecoder_delete(go_qsd_t *self) {
    if (self) QSD_Helper::destroy(D(self));
}

/* ---- Decoding ---- */
void qstringdecoder_decode(go_qsd_t *self,
    const char *in_data, ptrdiff_t in_size,
    void **out_qstr) {
    QByteArrayView in(in_data, static_cast<qsizetype>(in_size));
    /* decode() returns EncodedData<QByteArrayView>; the implicit conversion to QString
     * calls decodeAsString(), which handles both valid and invalid (iface == nullptr)
     * decoders — for invalid it increments state.invalidChars and returns an empty QString. */
    auto *str = new QString(D(self)->decode(in));
    *out_qstr = static_cast<void *>(str);
}

/* ---- Required space ---- */
ptrdiff_t qstringdecoder_required_space(const go_qsd_t *self, ptrdiff_t input_length) {
    return static_cast<ptrdiff_t>(D(self)->requiredSpace(static_cast<qsizetype>(input_length)));
}

/* ---- Append to buffer ---- */
ptrdiff_t qstringdecoder_append_to_buffer(go_qsd_t *self,
    void *buf, ptrdiff_t buf_len,
    const char *in_data, ptrdiff_t in_size) {
    (void)buf_len; /* QStringDecoder::appendToBuffer writes exactly toUtf16Len(input) units */
    char16_t *out = static_cast<char16_t *>(buf);
    QByteArrayView in(in_data, static_cast<qsizetype>(in_size));
    char16_t *next = D(self)->appendToBuffer(out, in);
    return static_cast<ptrdiff_t>(next - out);
}

/* ---- Finalize ---- */
void qstringdecoder_finalize(go_qsd_t *self,
    void *buf, ptrdiff_t buf_len,
    ptrdiff_t *out_next, int16_t *out_invalid, int *out_err) {
    char16_t *out = static_cast<char16_t *>(buf);
    auto r = D(self)->finalize(out, static_cast<qsizetype>(buf_len));
    *out_next    = (out != nullptr && r.next != nullptr)
                     ? static_cast<ptrdiff_t>(r.next - out)
                     : 0;
    *out_invalid = static_cast<int16_t>(r.invalidChars);
    *out_err     = static_cast<int>(r.error);
}

/* ---- decoderForHtml static ---- */
go_qsd_t *qstringdecoder_decoder_for_html(const char *data, ptrdiff_t size) {
    QByteArrayView in(data, static_cast<qsizetype>(size));
    return toD(new QStringDecoder(QStringDecoder::decoderForHtml(in)));
}

} // extern "C"