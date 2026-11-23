/* Order: version check first, then real Qt headers. */
#include "../internal/qt6core/check_qt_version.h"

/* qstringconverter.h declares QStringConverter, QStringEncoder and
 * QStringDecoder (and transitively qstringconverter_base.h, qstring.h,
 * qstringbuilder.h). This is the canonical header for the whole hierarchy, and
 * it is why the bridge for all three classes lives in this single file,
 * mirroring QtCore's qstringconverter.cpp. */
#include <QtCore/qstringconverter.h>

#include <QStringList>

#include <cstdlib>
#include <cstring>
#include <optional>
#include <utility>

#include "qstringconverter_c.h"

/* QSC_Helper exposes the protected members of QStringConverter:
 *   - Interface (protected nested struct)
 *   - iface, state (protected data fields)
 *   - constructors (including the protected ones)
 * The protected destructor is reachable through the derived go_qsc_s below. */
struct QSC_Helper : QStringConverter {
    using QStringConverter::Interface;
    using QStringConverter::iface;
    using QStringConverter::state;

    QSC_Helper() : QStringConverter() {}
    QSC_Helper(QStringConverter::Encoding e, QStringConverter::Flags f) : QStringConverter(e, f) {}
    QSC_Helper(QAnyStringView name, QStringConverter::Flags f) : QStringConverter(name, f) {}
    QSC_Helper(const Interface *i) : QStringConverter(i) {}
};

/* QSE_Helper exposes the protected QStringEncoder(const Interface *) ctor.
 * QStringEncoder itself is mostly public; its base members (iface/state) are
 * reached through the converter ABI, not here. */
struct QSE_Helper : QStringEncoder {
    using Interface = QStringConverter::Interface;

    QSE_Helper() : QStringEncoder() {}
    QSE_Helper(QStringConverter::Encoding e, QStringConverter::Flags f) : QStringEncoder(e, f) {}
    QSE_Helper(QAnyStringView name, QStringConverter::Flags f) : QStringEncoder(name, f) {}
    QSE_Helper(const Interface *i) : QStringEncoder(i) {}
};

/* QSD_Helper exposes the protected QStringDecoder(const Interface *) ctor. */
struct QSD_Helper : QStringDecoder {
    using Interface = QStringConverter::Interface;

    QSD_Helper() : QStringDecoder() {}
    QSD_Helper(QStringConverter::Encoding e, QStringConverter::Flags f) : QStringDecoder(e, f) {}
    QSD_Helper(QAnyStringView name, QStringConverter::Flags f) : QStringDecoder(name, f) {}
    QSD_Helper(const Interface *i) : QStringDecoder(i) {}
    QSD_Helper(QStringDecoder &&decoder) : QStringDecoder(std::move(decoder)) {}
};

/* All conversions between the C ABI (int, ptrdiff_t, void *) and Qt types live
 * in this namespace. The ABI functions below never spell out a cast themselves,
 * which keeps the conversion surface small and auditable: one place per type
 * pair instead of dozens scattered through the file. */
namespace cvt {

inline QStringConverter::Encoding encoding(int v) noexcept {
    return static_cast<QStringConverter::Encoding>(v);
}
inline QStringConverter::Flags flags(int v) noexcept {
    return QStringConverter::Flags(v);
}
inline qsizetype size(ptrdiff_t v) noexcept {
    return static_cast<qsizetype>(v);
}
inline ptrdiff_t ptrdiff(qsizetype v) noexcept {
    return static_cast<ptrdiff_t>(v);
}
inline size_t bytes(qsizetype v) noexcept {
    return static_cast<size_t>(v);
}
inline const char *chars(const void *p) noexcept {
    return static_cast<const char *>(p);
}
inline char *wchars(void *p) noexcept {
    return static_cast<char *>(p);
}
inline const char16_t *utf16(const void *p) noexcept {
    return static_cast<const char16_t *>(p);
}
inline char16_t *utf16Out(void *p) noexcept {
    return static_cast<char16_t *>(p);
}
inline char16_t codeUnit(uint16_t v) noexcept {
    return static_cast<char16_t>(v);
}
inline const QSC_Helper::Interface *iface(const void *p) noexcept {
    return static_cast<const QSC_Helper::Interface *>(p);
}
inline int encodingValue(const std::optional<QStringConverter::Encoding> &v) noexcept {
    return v ? static_cast<int>(*v) : -1;
}
inline int flagValue(QStringConverter::Flags f) noexcept {
    return static_cast<int>(f);
}
inline int64_t int64(qsizetype v) noexcept {
    return static_cast<int64_t>(v);
}
inline int16_t int16(qint16 v) noexcept {
    return static_cast<int16_t>(v);
}
inline int error(QStringConverter::FinalizeResultError e) noexcept {
    return static_cast<int>(e);
}
template <typename T>
inline T *mallocArray(size_t n) noexcept {
    return static_cast<T *>(std::malloc(n));
}

} // namespace cvt

/* Reconstruct a QAnyStringView from (data, size, kind).
 * kind: 0 = Utf8, 1 = Latin1, 2 = Utf16 (matching qanystringview.Kind). */
inline QAnyStringView makeAny(const void *data, ptrdiff_t size, int kind) noexcept {
    switch (kind) {
    case 0:
        return QAnyStringView(cvt::chars(data), cvt::size(size));
    case 1:
        return QAnyStringView(QLatin1StringView(cvt::chars(data), cvt::size(size)));
    case 2:
        return QAnyStringView(QStringView(cvt::utf16(data), cvt::size(size)));
    }
    return QAnyStringView();
}

/* The opaque C handle types are real subclasses of the helper types. This keeps
 * the C ABI unchanged (the go_qs*_t types stay incomplete in the header) while
 * letting every C function below use the object directly: no reinterpret_cast
 * between a handle and a Qt pointer is needed anywhere. */
struct go_qsc_s : QSC_Helper {
    go_qsc_s() : QSC_Helper() {}
    go_qsc_s(QStringConverter::Encoding e, QStringConverter::Flags f) : QSC_Helper(e, f) {}
    go_qsc_s(QAnyStringView name, QStringConverter::Flags f) : QSC_Helper(name, f) {}
    go_qsc_s(const QSC_Helper::Interface *i) : QSC_Helper(i) {}
};

struct go_qse_s : QSE_Helper {
    go_qse_s() : QSE_Helper() {}
    go_qse_s(QStringConverter::Encoding e, QStringConverter::Flags f) : QSE_Helper(e, f) {}
    go_qse_s(QAnyStringView name, QStringConverter::Flags f) : QSE_Helper(name, f) {}
    go_qse_s(const QSC_Helper::Interface *i) : QSE_Helper(i) {}
};

struct go_qsd_s : QSD_Helper {
    go_qsd_s() : QSD_Helper() {}
    go_qsd_s(QStringConverter::Encoding e, QStringConverter::Flags f) : QSD_Helper(e, f) {}
    go_qsd_s(QAnyStringView name, QStringConverter::Flags f) : QSD_Helper(name, f) {}
    go_qsd_s(const QSC_Helper::Interface *i) : QSD_Helper(i) {}
    /* Needed by qstringdecoder_decoder_for_html(), which returns by value. */
    go_qsd_s(QStringDecoder &&decoder) : QSD_Helper(std::move(decoder)) {}
};

extern "C" {

/* ==== QStringConverter ==== */

go_qsc_t *qstringconverter_new_default(void) {
    return new go_qsc_s();
}
go_qsc_t *qstringconverter_new_from_encoding(int encoding, int flags) {
    return new go_qsc_s(cvt::encoding(encoding), cvt::flags(flags));
}
go_qsc_t *qstringconverter_new_from_any(const void *data, ptrdiff_t size, int kind, int flags) {
    return new go_qsc_s(makeAny(data, size, kind), cvt::flags(flags));
}
go_qsc_t *qstringconverter_new_from_iface(const void *iface) {
    return new go_qsc_s(cvt::iface(iface));
}
void qstringconverter_delete(go_qsc_t *self) {
    delete self;
}

/* Qt 6.11 has no QStringConverter::swap(); implement via std::swap on protected
 * fields. State is move-only; std::swap uses move ctor/assign, which is correct. */
void qstringconverter_swap(go_qsc_t *self, go_qsc_t *other) {
    std::swap(self->iface, other->iface);
    std::swap(self->state, other->state);
}

bool qstringconverter_is_valid(const go_qsc_t *self) {
    return self->isValid();
}
const char *qstringconverter_name(const go_qsc_t *self) {
    return self->name();
}
void qstringconverter_reset_state(go_qsc_t *self) {
    self->resetState();
}
bool qstringconverter_has_error(const go_qsc_t *self) {
    return self->hasError();
}

void qstringconverter_get_state(const go_qsc_t *self,
    int *out_flags, int *out_internal,
    int64_t *out_remaining, int64_t *out_invalid, uint32_t out_data[4]) {
    const auto &s = self->state;
    *out_flags     = cvt::flagValue(s.flags);
    *out_internal  = s.internalState;
    *out_remaining = cvt::int64(s.remainingChars);
    *out_invalid   = cvt::int64(s.invalidChars);
    for (int i = 0; i < 4; ++i) out_data[i] = s.state_data[i];
}

void qstringconverter_set_state(go_qsc_t *self,
                                int flags, int internal, int64_t remaining, int64_t invalid,
                                const uint32_t data[4]) {
    auto &target = self->state;
    target.clear();
    target.flags          = cvt::flags(flags);
    target.internalState  = internal;
    target.remainingChars = cvt::size(remaining);
    target.invalidChars   = cvt::size(invalid);
    for (int i = 0; i < 4; ++i) target.state_data[i] = data[i];
    target.clearFn = nullptr;
}

const void *qstringconverter_iface(const go_qsc_t *self) {
    return self->iface;
}
const char *qstringconverter_interface_name(const void *iface) {
    return cvt::iface(iface)->name;
}
/* The Interface callbacks are plain function pointers. Converting them to void *
 * is the only cast the bridge still performs outside cvt, and it cannot be
 * avoided portably: ISO C++ makes function-pointer -> object-pointer conversion
 * conditionally supported. */
void *qstringconverter_interface_to_utf16(const void *iface) {
    return reinterpret_cast<void *>(cvt::iface(iface)->toUtf16);
}
void *qstringconverter_interface_to_utf16_len(const void *iface) {
    return reinterpret_cast<void *>(cvt::iface(iface)->toUtf16Len);
}
void *qstringconverter_interface_from_utf16(const void *iface) {
    return reinterpret_cast<void *>(cvt::iface(iface)->fromUtf16);
}
void *qstringconverter_interface_from_utf16_len(const void *iface) {
    return reinterpret_cast<void *>(cvt::iface(iface)->fromUtf16Len);
}

const char *qstringconverter_name_for_encoding(int encoding) {
    return QStringConverter::nameForEncoding(cvt::encoding(encoding));
}
int qstringconverter_encoding_for_name(const void *data, ptrdiff_t size, int kind) {
    return cvt::encodingValue(QStringConverter::encodingForName(makeAny(data, size, kind)));
}
int qstringconverter_encoding_for_data(const void *data, ptrdiff_t size,
                                       uint16_t expectedFirstChar) {
    return cvt::encodingValue(QStringConverter::encodingForData(
        QByteArrayView(cvt::chars(data), cvt::size(size)),
        cvt::codeUnit(expectedFirstChar)));
}
int qstringconverter_encoding_for_html(const void *data, ptrdiff_t size) {
    return cvt::encodingValue(QStringConverter::encodingForHtml(
        QByteArrayView(cvt::chars(data), cvt::size(size))));
}

char **qstringconverter_available_codecs(ptrdiff_t *out_len) {
    const QStringList codecs = QStringConverter::availableCodecs();
    *out_len = cvt::ptrdiff(codecs.size());
    if (codecs.isEmpty())
        return nullptr;

    auto **arr = cvt::mallocArray<char *>(sizeof(char *) * cvt::bytes(codecs.size()));
    if (!arr) {
        *out_len = 0;
        return nullptr;
    }
    for (qsizetype i = 0; i < codecs.size(); ++i) {
        const QByteArray utf8 = codecs.at(i).toUtf8();
        const size_t n = cvt::bytes(utf8.size()) + 1;
        arr[i] = cvt::mallocArray<char>(n);
        if (!arr[i]) {
            for (qsizetype j = 0; j < i; ++j)
                std::free(arr[j]);
            std::free(arr);
            *out_len = 0;
            return nullptr;
        }
        std::memcpy(arr[i], utf8.constData(), n);
    }
    return arr;
}
void qstringconverter_available_codecs_free(char **names, ptrdiff_t len) {
    if (!names)
        return;
    for (ptrdiff_t i = 0; i < len; ++i)
        std::free(names[i]);
    std::free(names);
}

/* ==== QStringEncoder ==== */

go_qse_t *qstringencoder_new_default(void) {
    return new go_qse_s();
}
go_qse_t *qstringencoder_new_from_encoding(int encoding, int flags) {
    return new go_qse_s(cvt::encoding(encoding), cvt::flags(flags));
}
go_qse_t *qstringencoder_new_from_any(const void *data, ptrdiff_t size, int kind, int flags) {
    return new go_qse_s(makeAny(data, size, kind), cvt::flags(flags));
}
go_qse_t *qstringencoder_new_from_iface(const void *iface) {
    return new go_qse_s(cvt::iface(iface));
}
void qstringencoder_delete(go_qse_t *self) {
    delete self;
}

QByteArray *qstringencoder_encode(go_qse_t *self, const uint16_t *in_data, ptrdiff_t in_size) {
    return new QByteArray(self->encode(QStringView(cvt::utf16(in_data), cvt::size(in_size))));
}

ptrdiff_t qstringencoder_required_space(const go_qse_t *self, ptrdiff_t input_length) {
    return cvt::ptrdiff(self->requiredSpace(cvt::size(input_length)));
}

ptrdiff_t qstringencoder_append_to_buffer(go_qse_t *self,
    void *buf, ptrdiff_t buf_len,
    const uint16_t *in_data, ptrdiff_t in_size) {
    (void)buf_len; /* appendToBuffer writes exactly fromUtf16Len(input) bytes */
    char *out = cvt::wchars(buf);
    char *next = self->appendToBuffer(
        out, QStringView(cvt::utf16(in_data), cvt::size(in_size)));
    return next - out;
}

void qstringencoder_finalize(go_qse_t *self,
    void *buf, ptrdiff_t buf_len,
    ptrdiff_t *out_next, int16_t *out_invalid, int *out_err) {
    char *out = cvt::wchars(buf);
    auto r = self->finalize(out, cvt::size(buf_len));
    *out_next    = (out && r.next) ? r.next - out : 0;
    *out_invalid = cvt::int16(r.invalidChars);
    *out_err     = cvt::error(r.error);
}

/* ==== QStringDecoder ==== */

go_qsd_t *qstringdecoder_new_default(void) {
    return new go_qsd_s();
}
go_qsd_t *qstringdecoder_new_from_encoding(int encoding, int flags) {
    return new go_qsd_s(cvt::encoding(encoding), cvt::flags(flags));
}
go_qsd_t *qstringdecoder_new_from_any(const void *data, ptrdiff_t size, int kind, int flags) {
    return new go_qsd_s(makeAny(data, size, kind), cvt::flags(flags));
}
go_qsd_t *qstringdecoder_new_from_iface(const void *iface) {
    return new go_qsd_s(cvt::iface(iface));
}
void qstringdecoder_delete(go_qsd_t *self) {
    delete self;
}

void qstringdecoder_decode(go_qsd_t *self,
    const char *in_data, ptrdiff_t in_size,
    void **out_qstr) {
    /* decode() returns EncodedData<QByteArrayView>; the implicit conversion to
     * QString calls decodeAsString(), which handles both valid and invalid
     * (iface == nullptr) decoders — for invalid it sets state.invalidChars and
     * returns an empty QString. */
    QString *str = new QString(
        self->decode(QByteArrayView(in_data, cvt::size(in_size))));
    *out_qstr = str;
}

ptrdiff_t qstringdecoder_required_space(const go_qsd_t *self, ptrdiff_t input_length) {
    return cvt::ptrdiff(self->requiredSpace(cvt::size(input_length)));
}

ptrdiff_t qstringdecoder_append_to_buffer(go_qsd_t *self,
    void *buf, ptrdiff_t buf_len,
    const char *in_data, ptrdiff_t in_size) {
    (void)buf_len; /* appendToBuffer writes exactly toUtf16Len(input) units */
    char16_t *out = cvt::utf16Out(buf);
    char16_t *next = self->appendToBuffer(
        out, QByteArrayView(in_data, cvt::size(in_size)));
    return next - out;
}

void qstringdecoder_finalize(go_qsd_t *self,
    void *buf, ptrdiff_t buf_len,
    ptrdiff_t *out_next, int16_t *out_invalid, int *out_err) {
    char16_t *out = cvt::utf16Out(buf);
    auto r = self->finalize(out, cvt::size(buf_len));
    *out_next    = (out && r.next) ? r.next - out : 0;
    *out_invalid = cvt::int16(r.invalidChars);
    *out_err     = cvt::error(r.error);
}

go_qsd_t *qstringdecoder_decoder_for_html(const char *data, ptrdiff_t size) {
    return new go_qsd_s(QStringDecoder::decoderForHtml(
        QByteArrayView(data, cvt::size(size))));
}

} // extern "C"
