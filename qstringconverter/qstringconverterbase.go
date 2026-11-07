package qstringconverter

// Package qstringconverter: базовые типы (Flag, Flags, State) и QStringConverterBase.

/*
 В Qt 6.11 QStringConverterBase почти пуст — в нём объявлены только Flag, Flags, State
 и protected-конструкторы/деструктор. Всё функциональное API (Encoding, Interface,
 protected iface/state, isValid/name/resetState/hasError, статические помощники)
 находится в QStringConverter, наследующем QStringConverterBase.

 не реализовано:
   - State::clearFn (функциональный указатель) — не экспонируется. State — move-only
     в C++; SetState использует move-assignment, который корректно вызывает clearFn
     для старого состояния (если кодировка этого требует).
   - State::d[2] (union-альтернатива state_data[4]) — те же байты, что State.Data;
     отдельного поля не нужно.
   - QStringConverterBase(QStringConverterBase&&) — Go не имеет move-конструкторов;
     пересоздавайте объект.
   - Операторы | и & для Flags — используйте Flags.Has(flag).
*/

/*
#cgo pkg-config: Qt6Core
#cgo CXXFLAGS: -std=c++17 -fPIC -Wno-sfinae-incomplete
#cgo LDFLAGS: -lstdc++
#include <stdlib.h>
#include <stdint.h>
#include <stdbool.h>
#include <stddef.h>
#include "qstringconverter_c.h"
*/
import "C"

import "runtime"

// Flag mirrors QStringConverterBase::Flag.
type Flag int

const (
	FlagStateless            Flag = 0x1
	FlagConvertInvalidToNull Flag = 0x2
	FlagWriteBom             Flag = 0x4
	FlagConvertInitialBom    Flag = 0x8
	FlagUsesIcu              Flag = 0x10
)

// Flags is a bitmask of Flag values (mirrors QStringConverterBase::Flags = QFlags<Flag>).
type Flags int

// FlagDefault is the empty Flags value (Qt's Flag::Default = 0). It is typed as Flags
// (rather than Flag) so it can be passed directly to constructors like
// NewQStringConverterFromEncoding(e, FlagDefault) without an explicit conversion —
// matching Qt's implicit Flag → Flags conversion.
const FlagDefault Flags = 0x0

// Has reports whether flag is set in f.
func (f Flags) Has(flag Flag) bool { return int(f)&int(flag) != 0 }

// State mirrors QStringConverterBase::State (simplified).
//
// The C++ State is move-only and contains a function pointer (clearFn); neither clearFn
// nor the union alternative d[2] is exposed in Go. Data holds the same bytes as the C++
// union of state_data[4] and d[2].
//
// Qt 6.11 State fields: flags (Flags), internalState (int), remainingChars (qsizetype),
// invalidChars (qsizetype), state_data[4] / d[2] (union), clearFn (fn ptr).
type State struct {
	Flags          Flags
	InternalState  int32    // C++ int
	RemainingChars int64    // C++ qsizetype (ptrdiff_t; int64 is a safe superset for 32-bit too)
	InvalidChars   int64    // C++ qsizetype
	Data           [4]uint32
}

// QStringConverterBase mirrors Qt 6.11's QStringConverterBase.
//
// In Qt 6.11 this class is almost empty (Flag/Flags/State + protected ctor/dtor); all
// functionality lives in QStringConverter. The Go type exists mainly to be embedded in
// *QStringConverter (matching Qt's class hierarchy) and to allow explicit construction.
type QStringConverterBase struct {
	p         *C.go_qscbase_t
	ownedByGo bool
}

// NewQStringConverterBase creates an empty base (protected default ctor exposed).
func NewQStringConverterBase() *QStringConverterBase {
	c := C.qstringconverterbase_new_default()
	if c == nil {
		return nil
	}
	b := &QStringConverterBase{p: c, ownedByGo: true}
	runtime.SetFinalizer(b, finalizeQStringConverterBase)
	return b
}

func finalizeQStringConverterBase(b *QStringConverterBase) {
	if b.p != nil && b.ownedByGo {
		C.qstringconverterbase_delete(b.p)
		b.p = nil
	}
}

// Close frees the underlying C++ object if this base owns it.
// For a base embedded in *QStringConverter, Close is a no-op (the parent owns the object).
func (b *QStringConverterBase) Close() {
	if b.p != nil && b.ownedByGo {
		C.qstringconverterbase_delete(b.p)
		b.p = nil
		runtime.SetFinalizer(b, nil)
	}
}
