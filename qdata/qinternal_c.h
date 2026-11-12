/* C ABI for Qt 6.11 QtCore/QInternal and QKeyCombination.
 *
 * QInternal is a purely static class — no instances. Its static methods and enum values
 * are exposed as free functions.
 *
 * QKeyCombination is a trivially copyable value type (single int). Passed by value.
 */
#ifndef GO_QT6CORE_QDATA_QINTERNAL_C_H
#define GO_QT6CORE_QDATA_QINTERNAL_C_H

#include <stddef.h>
#include <stdint.h>
#include <stdbool.h>

#ifdef __cplusplus
class QInternal;
class QKeyCombination;
typedef QKeyCombination go_qkeycombination_t;
extern "C" {
#else
typedef struct go_qkeycombination_s {
    int32_t combination;
} go_qkeycombination_t;
#endif

/* --- QInternal: static methods ---
 * callback is a C function pointer of type `bool (*)(void **)`.
 * The Go side passes it as unsafe.Pointer obtained from a C function
 * exported via //export; see qinternal_doc.go. */
bool qinternal_register_callback(int cb, void* callback);
bool qinternal_unregister_callback(int cb, void* callback);
bool qinternal_activate_callbacks(int cb, void** parameters);

/* --- QKeyCombination: constructors (return by value) --- */
go_qkeycombination_t qkeycombination_default(void);
go_qkeycombination_t qkeycombination_from_key(int key);
go_qkeycombination_t qkeycombination_from_modifiers_key(int modifiers, int key);
go_qkeycombination_t qkeycombination_from_combined(int combined);

/* --- QKeyCombination: accessors --- */
int qkeycombination_keyboard_modifiers(go_qkeycombination_t self);
int qkeycombination_key(go_qkeycombination_t self);
int qkeycombination_to_combined(go_qkeycombination_t self);

#ifdef __cplusplus
}
#endif
#endif