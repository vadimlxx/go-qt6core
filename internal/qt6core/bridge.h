/* Internal C ABI for runtime diagnostics.
 *
 * After the qtext unification (QString + QByteArray + QChar in one package
 * with direct C bindings), the previous QString/QByteArray bridge has been
 * removed. This header now exposes only the runtime version query.
 *
 * check_qt_version.h in the same directory is still the single source of
 * truth for the minimum Qt version and is #included by every *_c.cpp
 * of the qtext package.
 */
#ifndef GO_QT6CORE_BRIDGE_H
#define GO_QT6CORE_BRIDGE_H

#ifdef __cplusplus
extern "C" {
#endif

/* Writes the linked Qt version to *major / *minor / *patch. */
void qt6core_get_version(int *major, int *minor, int *patch);

#ifdef __cplusplus
}
#endif

#endif
