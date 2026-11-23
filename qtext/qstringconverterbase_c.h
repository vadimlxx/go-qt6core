/* C ABI for Qt 6.11 QtCore/QStringConverterBase.
 *
 * QStringConverterBase is almost empty in Qt 6.11 (Flag/Flags/State plus a
 * protected constructor/destructor), so this bridge only exposes construction
 * and destruction. It is a separate translation unit from QStringConverter,
 * mirroring QtCore's split between qstringconverter_base.h and
 * qstringconverter.h.
 *
 * The go_qscbase_t handle is opaque to C. In qstringconverterbase_c.cpp it is
 * defined as a real subclass of a helper adapter over QStringConverterBase, so
 * the two functions below delete the object directly and need no casts.
 *
 * Ownership: qstringconverterbase_new_default() returns a handle owned by the
 * caller; release it with qstringconverterbase_delete().
 */
#ifndef GO_QT6CORE_QSTRINGCONVERTERBASE_C_H
#define GO_QT6CORE_QSTRINGCONVERTERBASE_C_H

#include <stddef.h>

#ifdef __cplusplus
extern "C" {
#endif

typedef struct go_qscbase_s go_qscbase_t;

go_qscbase_t *qstringconverterbase_new_default(void);
void          qstringconverterbase_delete(go_qscbase_t *self);

#ifdef __cplusplus
}
#endif

#endif
