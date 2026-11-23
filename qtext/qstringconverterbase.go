package qtext

/*
QStringConverterBase — базовые типы Qt 6.11 (Flag, Flags, State) и сам базовый
класс QStringConverterBase.

В Qt 6.11 QStringConverterBase почти пуст: в нём объявлены только Flag, Flags,
State и protected-конструктор/деструктор. Всё функциональное API (Encoding,
Interface, protected-поля iface/state, isValid/name/resetState/hasError и
статические помощники) находится в QStringConverter, который наследует
QStringConverterBase (см. qstringconverter.go).

Владение:

  - NewQStringConverterBase создаёт отдельный C++ объект; он освобождается
    Close() или runtime.AddCleanup (когда Go-обёртка становится недостижимой);
  - экземпляр, встроенный в *QStringConverter, — не владеющая проекция на тот
    же C++ объект: его Close() ничего не делает, а владелец удаляет объект.

Оба варианта используют общий qscState с атомарным слотом указателя, поэтому
Close() и cleanup безопасно соревнуются (ровно один забирает объект), а после
закрытия методы видят nil вместо висячего указателя.

Не реализовано:

  - State::clearFn (функциональный указатель) — не экспонируется; SetState
    вызывает clear()/move-присваивание, которые корректно очищают состояние.
  - State::d[2] (union-альтернатива state_data[4]) — те же байты, что State.Data;
    отдельного поля не нужно.
  - QStringConverterBase(QStringConverterBase&&) — Go не имеет move-конструкторов;
    пересоздавайте объект.
  - operator| и operator& для Flags — используйте Flags.Has(flag) и явное
    приведение типов.
*/

/*
#cgo pkg-config: Qt6Core
#cgo CXXFLAGS: -std=c++17 -fPIC -Wno-sfinae-incomplete
#cgo LDFLAGS: -lstdc++
#include <stdlib.h>
#include <stdint.h>
#include <stdbool.h>
#include <stddef.h>
#include "qstringconverterbase_c.h"
*/
import "C"

import (
	"runtime"
	"sync/atomic"
	"unsafe"
)

// qscState holds a C++ QStringConverter/QStringConverterBase pointer in an
// atomic slot so that Close() and the runtime cleanup can race safely: whoever
// swaps the slot to nil owns the object and deletes it.
//
// The same state may be shared by a *QStringConverter and its embedded
// *QStringConverterBase view. In that case only the converter is allowed to
// delete the object (owned == true); the view is a non-owning alias.
//
// *byte (rather than the cgo pointer type) is used because cgo exposes the
// opaque structs as incomplete Go types, while atomic.Pointer[T] requires a
// complete type. The pointer value is unchanged; only the Go-side static type
// differs.
//
// State is a separate allocation because the cleanup function must not capture
// the wrapper itself — otherwise the wrapper would stay reachable forever and
// the cleanup would never run.
//
// Qt: no direct counterpart — internal ownership helper.
type qscState struct {
	p atomic.Pointer[byte]
}

// cpp returns the stored C++ pointer as unsafe.Pointer, or nil if the slot has
// already been claimed (the object is closed). Read-only; does not take
// ownership.
// Qt: no direct counterpart — internal helper.
func (s *qscState) cpp() unsafe.Pointer {
	if bp := s.p.Load(); bp != nil {
		return unsafe.Pointer(bp)
	}
	return nil
}

// claim atomically clears the slot and returns the pointer that was stored
// there (or nil). Exactly one caller receives a non-nil result for a given
// pointer, which makes Close() versus cleanup safe against a double free.
// Qt: no direct counterpart — internal helper.
func (s *qscState) claim() unsafe.Pointer {
	if bp := s.p.Swap(nil); bp != nil {
		return unsafe.Pointer(bp)
	}
	return nil
}

// Flag mirrors QStringConverterBase::Flag.
// Qt: QStringConverterBase::Flag.
type Flag int

const (
	// FlagStateless mirrors Flag::Stateless.
	// Qt: QStringConverterBase::Flag::Stateless.
	FlagStateless Flag = 0x1
	// FlagConvertInvalidToNull mirrors Flag::ConvertInvalidToNull.
	// Qt: QStringConverterBase::Flag::ConvertInvalidToNull.
	FlagConvertInvalidToNull Flag = 0x2
	// FlagWriteBom mirrors Flag::WriteBom.
	// Qt: QStringConverterBase::Flag::WriteBom.
	FlagWriteBom Flag = 0x4
	// FlagConvertInitialBom mirrors Flag::ConvertInitialBom.
	// Qt: QStringConverterBase::Flag::ConvertInitialBom.
	FlagConvertInitialBom Flag = 0x8
	// FlagUsesIcu mirrors Flag::UsesIcu.
	// Qt: QStringConverterBase::Flag::UsesIcu.
	FlagUsesIcu Flag = 0x10
)

// Flags is a bitmask of Flag values (QFlags<Flag>).
// Qt: QStringConverterBase::Flags.
type Flags int

// FlagDefault is the empty Flags value (Qt's Flag::Default = 0). It is typed as
// Flags rather than Flag so it can be passed directly to constructors such as
// NewQStringConverterFromEncoding(e, FlagDefault) without an explicit
// conversion, matching Qt's implicit Flag -> Flags conversion.
// Qt: QStringConverterBase::Flag::Default / Flags(Flag::Default).
const FlagDefault Flags = 0x0

// Has reports whether flag is set in f.
// Qt: QFlags<QStringConverterBase::Flag>::testFlag(Flag) const.
func (f Flags) Has(flag Flag) bool { return int(f)&int(flag) != 0 }

// State mirrors QStringConverterBase::State (simplified).
//
// The C++ State is move-only and contains a function pointer (clearFn); neither
// clearFn nor the union alternative d[2] is exposed in Go. Data holds the same
// bytes as the C++ union of state_data[4] and d[2].
//
// Qt 6.11 State fields: flags (Flags), internalState (int), remainingChars
// (qsizetype), invalidChars (qsizetype), state_data[4] / d[2] (union),
// clearFn (function pointer).
type State struct {
	Flags          Flags
	InternalState  int32
	RemainingChars int64
	InvalidChars   int64
	Data           [4]uint32
}

// QStringConverterBase mirrors Qt 6.11's QStringConverterBase.
//
// In Qt 6.11 the class is almost empty (Flag/Flags/State plus protected
// constructor/destructor); all functionality lives in QStringConverter. The Go
// type exists to reproduce Qt's class hierarchy via embedding and to allow
// explicit construction of the base object.
//
// A base created by NewQStringConverterBase owns its C++ object. A base
// embedded in *QStringConverter is a non-owning view over the converter's
// object; its Close is a no-op.
//
// Qt: QStringConverterBase.
type QStringConverterBase struct {
	state   *qscState
	owned   bool
	cleanup runtime.Cleanup
}

// NewQStringConverterBase creates an empty base object (the protected default
// constructor is exposed).
// Qt: QStringConverterBase::QStringConverterBase() (protected).
func NewQStringConverterBase() *QStringConverterBase {
	c := C.qstringconverterbase_new_default()
	if c == nil {
		return nil
	}
	st := &qscState{}
	st.p.Store((*byte)(unsafe.Pointer(c)))
	b := &QStringConverterBase{state: st, owned: true}
	b.cleanup = runtime.AddCleanup(b, func(st *qscState) {
		if p := st.claim(); p != nil {
			C.qstringconverterbase_delete((*C.go_qscbase_t)(p))
		}
	}, st)
	return b
}

// Close frees the underlying C++ object if this base owns it. For a base
// embedded in *QStringConverter, Close is a no-op (the converter owns the
// object and will delete it).
//
// Close is idempotent and nil-safe. It cancels the registered cleanup; if Close
// races with the cleanup, exactly one of them claims the object, so there is no
// double free.
// Qt: QStringConverterBase::~QStringConverterBase() (protected).
func (b *QStringConverterBase) Close() {
	if b == nil || b.state == nil || !b.owned {
		return
	}
	b.cleanup.Stop()
	if p := b.state.claim(); p != nil {
		C.qstringconverterbase_delete((*C.go_qscbase_t)(p))
	}
}
