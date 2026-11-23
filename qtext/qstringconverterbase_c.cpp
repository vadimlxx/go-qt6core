/* Order: version check first, then real Qt headers. */
#include "../internal/qt6core/check_qt_version.h"

#include <QtCore/qstringconverter_base.h>

#include "qstringconverterbase_c.h"

/* QSB_Helper makes the protected QStringConverterBase constructor and
 * destructor accessible to the go_qscbase_s subclass below. */
struct QSB_Helper : QStringConverterBase {};

/* The opaque handle is a real subclass: no cast between go_qscbase_t* and
 * QStringConverterBase* is needed anywhere. */
struct go_qscbase_s : QSB_Helper {};

extern "C" {

go_qscbase_t *qstringconverterbase_new_default(void) {
    return new go_qscbase_s();
}
void qstringconverterbase_delete(go_qscbase_t *self) {
    delete self;
}

} // extern "C"
