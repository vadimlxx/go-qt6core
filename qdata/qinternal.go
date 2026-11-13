package qdata

/*
 QInternal — биндинг Qt 6.11 QtCore/QInternal как набора статических функций.

 QInternal — чисто статический класс в Qt: у него нет экземпляров, только статические
 методы и enum'ы. В Go он представлен:
   - типом QInternal (пустая структура) с методами RegisterCallback,
     UnregisterCallback, ActivateCallbacks;
   - пакетными функциями тех же имён (алиасы для удобства).

 Статические методы RegisterCallback/UnregisterCallback/ActivateCallbacks принимают
 qInternalCallback = bool (*)(void**). В Go невозможно создать указатель на Go-функцию
 для передачи в C. Поэтому эти методы предназначены для использования совместно с
 C-функциями, экспортированными через //export в C-части пользовательского проекта.
 Сигнатура: func RegisterCallback(cb Callback, fn unsafe.Pointer) bool.

 Enum'ы QInternal (Callback, PaintDeviceFlags, RelayoutType, DockPosition, CallMode,
 SaveStateRule) вынесены в отдельные файлы и доступны как Go-константы.

 Что не реализовано (см. qinternal_doc.go):

   - QInternal::callFunction(InternalFunction, void**) — InternalFunction enum
     не выставлен (в qnamespace.h его нет; он определён в qglobal.cpp).
   - QInternal::registerCallback с проверкой типа функции — Go-типобезопасность
     не пересекается с C-callback'ами; пользователь отвечает за корректность.

 Protected-секций у QInternal нет — все методы публично-статические.
*/

/*
#cgo pkg-config: Qt6Core
#cgo CXXFLAGS: -std=c++17 -fPIC -Wno-sfinae-incomplete
#cgo LDFLAGS: -lstdc++
#include <stdint.h>
#include "qinternal_c.h"
*/
import "C"

import "unsafe"

// QInternal is a purely static class; the Go type exists only to host methods.
type QInternal struct{}

func NewQInternal() *QInternal {
	return &QInternal{}
}

// RegisterCallback registers a C callback for the given event. Mirrors
// QInternal::registerCallback(Callback, qInternalCallback).
//
// fn must be an unsafe.Pointer to a C function with signature bool(*)(void**),
// obtained from a //export'ed Go function or from C directly. Returns false if
// the callback slot is invalid or full.
func (QInternal) RegisterCallback(cb Callback, fn unsafe.Pointer) bool {
	return bool(C.qinternal_register_callback(C.int(cb), fn))
}

// UnregisterCallback removes a previously registered callback.
func (QInternal) UnregisterCallback(cb Callback, fn unsafe.Pointer) bool {
	return bool(C.qinternal_unregister_callback(C.int(cb), fn))
}

// ActivateCallbacks invokes all registered callbacks for cb with the given
// parameters array. Returns true if at least one callback ran successfully.
func (QInternal) ActivateCallbacks(cb Callback, params unsafe.Pointer) bool {
	return bool(C.qinternal_activate_callbacks(C.int(cb), (*unsafe.Pointer)(params)))
}

// Package-level aliases — QInternal is static, so these forward to the type methods.
func RegisterCallback(cb Callback, fn unsafe.Pointer) bool {
	return QInternal{}.RegisterCallback(cb, fn)
}

func UnregisterCallback(cb Callback, fn unsafe.Pointer) bool {
	return QInternal{}.UnregisterCallback(cb, fn)
}

func ActivateCallbacks(cb Callback, params unsafe.Pointer) bool {
	return QInternal{}.ActivateCallbacks(cb, params)
}
