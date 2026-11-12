#include "../internal/qt6core/check_qt_version.h"

#include <QInternal>
#include <QKeyCombination>
#include <QtCore/qnamespace.h>

#include "qinternal_c.h"

extern "C" {

/* --- QInternal: static methods --- */
bool qinternal_register_callback(int cb, void* callback) {
    return QInternal::registerCallback(static_cast<QInternal::Callback>(cb),
                                       reinterpret_cast<qInternalCallback>(callback));
}

bool qinternal_unregister_callback(int cb, void* callback) {
    return QInternal::unregisterCallback(static_cast<QInternal::Callback>(cb),
                                         reinterpret_cast<qInternalCallback>(callback));
}

bool qinternal_activate_callbacks(int cb, void** parameters) {
    return QInternal::activateCallbacks(static_cast<QInternal::Callback>(cb), parameters);
}

/* --- QKeyCombination: constructors --- */
go_qkeycombination_t qkeycombination_default() {
    return QKeyCombination();
}

go_qkeycombination_t qkeycombination_from_key(int key) {
    return QKeyCombination(static_cast<Qt::Key>(key));
}

go_qkeycombination_t qkeycombination_from_modifiers_key(int modifiers, int key) {
    return QKeyCombination(Qt::KeyboardModifiers(modifiers), static_cast<Qt::Key>(key));
}

go_qkeycombination_t qkeycombination_from_combined(int combined) {
    return QKeyCombination::fromCombined(combined);
}

/* --- QKeyCombination: accessors --- */
int qkeycombination_keyboard_modifiers(go_qkeycombination_t self) {
    return self.keyboardModifiers().toInt();
}

int qkeycombination_key(go_qkeycombination_t self) {
    return static_cast<int>(self.key());
}

int qkeycombination_to_combined(go_qkeycombination_t self) {
    return self.toCombined();
}

} // extern "C"