/* Order: version check first, then real Qt headers. */
#include "../internal/qt6core/check_qt_version.h"

/* Include the real qstringconverter_base.h directly — NOT <QStringConverterBase>, which
 * pulls qstringconverter_base_deprecated.h and masks the real QStringConverterBase members.
 * In Qt 6.x, this single header declares both QStringConverterBase and QStringConverter. */
#include <QtCore/qstringconverter_base.h>

/* qstring.h is REQUIRED here, not just convenient:
 *  1. It defines QString (complete type). Without it, the compiler's viability check for
 *     QStringConverter(Encoding, Flags) substitutes Encoding into the QAnyStringView
 *     template ctor and fails on `std::is_convertible<Encoding, QString>` with
 *     "invalid use of incomplete type 'class QString'".
 *  2. It provides inline definitions of QAnyStringView(QLatin1StringView),
 *     QAnyStringView(const QString&), QAnyStringView(const QByteArray&) — without them,
 *     using those constructors produces a "used but not defined" warning.
 * qstring.h transitively includes qanystringview.h, qbytearrayview.h, qstringview.h,
 * qlatin1stringview.h, qbytearray.h — so all our other includes come for free. */
#include <QtCore/qstring.h>

#include <utility>

#include "qstringconverter_c.h"

namespace {

/* QSC_Helper exposes protected members of QStringConverter:
 *   - Interface (protected nested struct)
 *   - iface, state (protected data fields)
 *   - constructors
 *   - destructor (for destroy) */
struct QSC_Helper : QStringConverter {
    using QStringConverter::Interface;
    using QStringConverter::iface;
    using QStringConverter::state;

    QSC_Helper() : QStringConverter() {}
    QSC_Helper(QStringConverter::Encoding e, QStringConverter::Flags f) : QStringConverter(e, f) {}
    QSC_Helper(QAnyStringView name, QStringConverter::Flags f) : QStringConverter(name, f) {}
    QSC_Helper(const Interface *i) : QStringConverter(i) {}

    static void destroy(QStringConverter *p) { delete static_cast<QSC_Helper *>(p); }
};

/* QSB_Helper is only needed because ~QStringConverterBase() is protected. */
struct QSB_Helper : QStringConverterBase {
    static void destroy(QStringConverterBase *p) { delete static_cast<QSB_Helper *>(p); }
};

inline QStringConverter       *C (go_qsc_t       *p) { return reinterpret_cast<QStringConverter       *>(p); }
inline const QStringConverter *C (const go_qsc_t *p) { return reinterpret_cast<const QStringConverter *>(p); }
inline QStringConverterBase   *B (go_qscbase_t   *p) { return reinterpret_cast<QStringConverterBase   *>(p); }
inline go_qsc_t     *toC(QStringConverter     *p) { return reinterpret_cast<go_qsc_t     *>(p); }
inline go_qscbase_t *toB(QStringConverterBase *p) { return reinterpret_cast<go_qscbase_t *>(p); }
inline QSC_Helper       *CA(go_qsc_t       *p) { return reinterpret_cast<QSC_Helper       *>(p); }
inline const QSC_Helper *CA(const go_qsc_t *p) { return reinterpret_cast<const QSC_Helper *>(p); }

/* Reconstruct a QAnyStringView from (data, size, kind).
 * kind: 0 = Utf8, 1 = Latin1, 2 = Utf16 (matching qanystringview.Kind). */
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

/* ---- QStringConverterBase ---- */
go_qscbase_t *qstringconverterbase_new_default(void) {
    return toB(new QSB_Helper());
}
void qstringconverterbase_delete(go_qscbase_t *self) {
    if (self) QSB_Helper::destroy(B(self));
}

/* ---- QStringConverter ctors ---- */
go_qsc_t *qstringconverter_new_default(void) {
    return toC(new QSC_Helper());
}
go_qsc_t *qstringconverter_new_from_encoding(int encoding, int flags) {
    return toC(new QSC_Helper(static_cast<QStringConverter::Encoding>(encoding),
                              QStringConverter::Flags(flags)));
}
go_qsc_t *qstringconverter_new_from_any(const void *data, ptrdiff_t size, int kind, int flags) {
    return toC(new QSC_Helper(makeAny(data, size, kind), QStringConverter::Flags(flags)));
}
go_qsc_t *qstringconverter_new_from_iface(const void *iface) {
    return toC(new QSC_Helper(static_cast<const QSC_Helper::Interface *>(iface)));
}
void qstringconverter_delete(go_qsc_t *self) {
    if (self) QSC_Helper::destroy(C(self));
}

/* Qt 6.11 has no QStringConverter::swap(); implement via std::swap on protected fields.
 * State is move-only; std::swap uses move ctor/assign, which is correct. */
void qstringconverter_swap(go_qsc_t *self, go_qsc_t *other) {
    auto *a = CA(self);
    auto *b = CA(other);
    std::swap(a->iface, b->iface);
    std::swap(a->state, b->state);
}

/* ---- Public methods ---- */
bool qstringconverter_is_valid(const go_qsc_t *self) {
    return C(self)->isValid();
}
const char *qstringconverter_name(const go_qsc_t *self) {
    return C(self)->name();
}
void qstringconverter_reset_state(go_qsc_t *self) {
    C(self)->resetState();
}
bool qstringconverter_has_error(const go_qsc_t *self) {
    return C(self)->hasError();
}

/* ---- Protected state accessors ---- */
void qstringconverter_get_state(const go_qsc_t *self,
    int *out_flags, int *out_internal,
    int64_t *out_remaining, int64_t *out_invalid, uint32_t out_data[4]) {
    const auto &s = CA(self)->state;
    *out_flags     = static_cast<int>(s.flags);
    *out_internal  = s.internalState;
    *out_remaining = static_cast<int64_t>(s.remainingChars);
    *out_invalid   = static_cast<int64_t>(s.invalidChars);
    for (int i = 0; i < 4; ++i) out_data[i] = s.state_data[i];
}

void qstringconverter_set_state(go_qsc_t *self,
								int flags, int internal, int64_t remaining, int64_t invalid, const uint32_t data[4]) {
	auto &target = CA(self)->state;
	target.clear();
	target.flags          = QStringConverter::Flags(flags);
	target.internalState  = internal;
	target.remainingChars = static_cast<qsizetype>(remaining);
	target.invalidChars   = static_cast<qsizetype>(invalid);
	for (int i = 0; i < 4; ++i) target.state_data[i] = data[i];
	target.clearFn = nullptr;
}

/* ---- Protected nested Interface ---- */
const void *qstringconverter_iface(const go_qsc_t *self) {
    return static_cast<const void *>(CA(self)->iface);
}
const char *qstringconverter_interface_name(const void *iface) {
    return static_cast<const QSC_Helper::Interface *>(iface)->name;
}
void *qstringconverter_interface_to_utf16(const void *iface) {
    return reinterpret_cast<void *>(static_cast<const QSC_Helper::Interface *>(iface)->toUtf16);
}
void *qstringconverter_interface_to_utf16_len(const void *iface) {
    return reinterpret_cast<void *>(static_cast<const QSC_Helper::Interface *>(iface)->toUtf16Len);
}
void *qstringconverter_interface_from_utf16(const void *iface) {
    return reinterpret_cast<void *>(static_cast<const QSC_Helper::Interface *>(iface)->fromUtf16);
}
void *qstringconverter_interface_from_utf16_len(const void *iface) {
    return reinterpret_cast<void *>(static_cast<const QSC_Helper::Interface *>(iface)->fromUtf16Len);
}

/* ---- Static helpers ---- */
const char *qstringconverter_name_for_encoding(int encoding) {
    return QStringConverter::nameForEncoding(static_cast<QStringConverter::Encoding>(encoding));
}
int qstringconverter_encoding_for_name(const void *data, ptrdiff_t size, int kind) {
    auto opt = QStringConverter::encodingForName(makeAny(data, size, kind));
    if (!opt) return -1;
    return static_cast<int>(*opt);
}
int qstringconverter_encoding_for_data(const void *data, ptrdiff_t size,
                                       uint16_t expectedFirstChar) {
    auto opt = QStringConverter::encodingForData(
        QByteArrayView(static_cast<const char *>(data), static_cast<qsizetype>(size)),
        static_cast<char16_t>(expectedFirstChar));
    if (!opt) return -1;
    return static_cast<int>(*opt);
}
int qstringconverter_encoding_for_html(const void *data, ptrdiff_t size) {
    auto opt = QStringConverter::encodingForHtml(
        QByteArrayView(static_cast<const char *>(data), static_cast<qsizetype>(size)));
    if (!opt) return -1;
    return static_cast<int>(*opt);
}

} // extern "C"
