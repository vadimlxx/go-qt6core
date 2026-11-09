/* Order: version check first, then real Qt headers. */
#include "../internal/qt6core/check_qt_version.h"

#include <QChar>   /* brings QLatin1Char and QChar */

#include "qlatin1char_c.h"

namespace {
inline QLatin1Char C(uint8_t u) noexcept {
    return QLatin1Char(static_cast<char>(u));
}
} // namespace

extern "C" {

uint8_t qlatin1char_to_latin1(uint8_t self) {
    return static_cast<uint8_t>(C(self).toLatin1());
}

uint16_t qlatin1char_unicode(uint8_t self) {
    return static_cast<uint16_t>(C(self).unicode());
}

/* QLatin1Char vs QLatin1Char: generated operator== from
 * Q_DECLARE_STRONGLY_ORDERED_LITERAL_TYPE(QLatin1Char). */
bool qlatin1char_compares_equal(uint8_t a, uint8_t b) {
    return C(a) == C(b);
}

/* QLatin1Char vs QLatin1Char ordering: uchar(lhs) vs uchar(rhs)
 * (unsigned) via Qt's compareThreeWay, called implicitly by
 * the generated operator<. */
int qlatin1char_compare(uint8_t a, uint8_t b) {
    if (C(a) == C(b)) return 0;
    return C(a) < C(b) ? -1 : +1;
}

/* QLatin1Char vs raw char: bit-equality — same as uint8_t equality. */
bool qlatin1char_compares_equal_byte(uint8_t self, uint8_t other) {
    return C(self) == static_cast<char>(other);
}

/* QLatin1Char vs raw char ordering: unsigned (uchar) comparison. */
int qlatin1char_compare_byte(uint8_t self, uint8_t other) {
    const char o = static_cast<char>(other);
    if (C(self) == o) return 0;
    return C(self) < o ? -1 : +1;
}

} // extern "C"