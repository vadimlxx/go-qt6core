/* Order matters: version check first, then real Qt headers. */
#include "check_qt_version.h"

#include "bridge.h"

extern "C" {

void qt6core_get_version(int *major, int *minor, int *patch) {
    *major = (QT_VERSION >> 16) & 0xFF;
    *minor = (QT_VERSION >>  8) & 0xFF;
    *patch =  QT_VERSION        & 0xFF;
}

} // extern "C"

