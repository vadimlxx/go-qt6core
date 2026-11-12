// Package qdata provides Go bindings for Qt 6.11 QtCore data types.
//
// Currently supported:
//   - QLocale       (handle-based; COW via QSharedDataPointer)
//   - QCalendar     (value-based; trivially copyable)
//   - QInternal     (static-only; enum + static methods)
//   - QKeyCombination (value-based; trivially copyable)
//
// # QInternal
//
// QInternal is a purely static utility class in Qt: no instances, only static methods
// and enums. In Go it is exposed as:
//
//	type QInternal struct{}   // method host, never instantiated
//	func (QInternal) RegisterCallback(cb Callback, fn unsafe.Pointer) bool
//	func (QInternal) UnregisterCallback(cb Callback, fn unsafe.Pointer) bool
//	func (QInternal) ActivateCallbacks(cb Callback, params unsafe.Pointer) bool
//
// plus package-level aliases RegisterCallback / UnregisterCallback / ActivateCallbacks.
//
// The callback parameter is qInternalCallback = bool(*)(void**). Since Go cannot hand
// a Go function pointer to C, callers must supply an unsafe.Pointer obtained from a
// C function exported via //export in their own cgo file. Example:
//
//	//export myEventNotify
//	func myEventNotify(params unsafe.Pointer) C.bool { ... }
//
//	qdata.RegisterCallback(qdata.EventNotifyCallback, unsafe.Pointer(C.myEventNotify))
//
// QInternal enums (Callback, PaintDeviceFlags, RelayoutType, DockPosition, CallMode,
// SaveStateRule) live in their own files (callback.go, paint_device_flags.go, …).
//
// # Not yet implemented (QInternal)
//
//   - QInternal::callFunction(InternalFunction, void**) — InternalFunction enum is
//     defined only in qglobal.cpp, not in public headers, so it cannot be bound
//     without a hand-maintained table.
//   - QInternal::Callback currently covers only EventNotifyCallback and LastCallback,
//     matching the Qt 6.11 header. Earlier Qt versions had ConnectCallback /
//     DisconnectCallback / AdoptCurrentThread / RefAdoptedThread — these were removed
//     in Qt 6 and are intentionally not exposed.
//
// # QKeyCombination
//
// Value-based, trivially copyable (single int). Constructors and accessors:
//
//	NewQKeyCombination()                                 — QKeyCombination()
//	NewQKeyCombinationFromKey(Key)                       — QKeyCombination(Qt::Key)
//	NewQKeyCombinationFromModifiersKey(KeyboardModifiers, Key)
//	QKeyCombinationFromCombined(int)                     — QKeyCombination::fromCombined
//	(qkc).KeyboardModifiers() KeyboardModifiers
//	(qkc).Key() Key
//	(qkc).ToCombined() int
//
// Package-level combinators (Go has no operator overloading):
//
//	Combine(mods KeyboardModifiers, key Key) QKeyCombination
//
// Requires Qt >= 6.11, C++17, pkg-config Qt6Core.
package qdata
