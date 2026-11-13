package qdata

/*
 QKeyCombination — биндинг Qt 6.11 QtCore/QKeyCombination как VALUE-тип.

 QKeyCombination хранит одно поле int combination — упакованную пару
 (Qt::KeyboardModifiers | Qt::Key). Размер 4 байта, trivially copyable.
 В Go хранится как плоская структура из тех же 4 байт. Присваивание создаёт
 независимую копию. Указателей, финализаторов и Close() нет.

 ВАЖНО: Go-нулевое значение QKeyCombination{} соответствует Qt Key_unknown (0).
 Qt-конструктор по умолчанию QKeyCombination() тоже даёт Key_unknown —
 то есть Go zero value и Qt default эквивалентны (в отличие от QCalendar).

 Соглашения:
   - Конструкторы: NewQKeyCombination(), NewQKeyCombinationFromKey(key),
     NewQKeyCombinationFromModifiersKey(mods, key), QKeyCombinationFromCombined(n).
   - Геттеры: KeyboardModifiers(), Key(), ToCombined().
   - Пакетный комбинатор: Combine(mods, key) — аналог operator| в C++.
   - Сравнение: Equal(other) (comparesEqual по полю combination).

 Что не реализовано (см. qinternal_doc.go):

   - QKeyCombination::operator< — удалён в Qt (bool operator<(QKeyCombination) = delete).
   - Deprecated operator int() — Qt 6 удалил неявную конверсию в int.
   - Deprecated operator+ — Qt 6 заменил на operator|; в Go используйте Combine.
   - Свободные функции Qt::operator|(...) — заменены методом Combine.

 Protected-секций у QKeyCombination нет.
*/

/*
#cgo pkg-config: Qt6Core
#cgo CXXFLAGS: -std=c++17 -fPIC -Wno-sfinae-incomplete
#cgo LDFLAGS: -lstdc++
#include <stdint.h>
#include "qinternal_c.h"
*/
import "C"

// QKeyCombination is a value type mirroring Qt6 QKeyCombination (single int).
// The Go zero value equals Qt's default-constructed Key_unknown.
type QKeyCombination struct {
	v C.go_qkeycombination_t
}

// NewQKeyCombination returns QKeyCombination() — Key_unknown with no modifiers.
func NewQKeyCombination(arg ...any) QKeyCombination {
	if len(arg) == 0 || arg[0] == nil {
		return QKeyCombination{v: C.qkeycombination_default()}
	}

	if len(arg) == 1 {
		key, ok := arg[0].(Key)
		if ok {
			return NewQKeyCombinationFromKey(key)
		}
	}

	mods, ok1 := arg[0].(KeyboardModifiers)
	key,  ok2 := arg[1].(Key)
	if ok1 && ok2 {
		return NewQKeyCombinationFromModifiersKey(mods, key)
	}

	panic("NewQKeyCombination unknown arguments")
}

// NewQKeyCombinationFromKey returns QKeyCombination(key).
func NewQKeyCombinationFromKey(key Key) QKeyCombination {
	return QKeyCombination{v: C.qkeycombination_from_key(C.int(key))}
}

// NewQKeyCombinationFromModifiersKey returns QKeyCombination(modifiers, key).
func NewQKeyCombinationFromModifiersKey(mods KeyboardModifiers, key Key) QKeyCombination {
	return QKeyCombination{v: C.qkeycombination_from_modifiers_key(C.int(mods), C.int(key))}
}

// QKeyCombinationFromCombined returns QKeyCombination::fromCombined(combined).
func QKeyCombinationFromCombined(combined int) QKeyCombination {
	return QKeyCombination{v: C.qkeycombination_from_combined(C.int(combined))}
}

// Combine is the Go equivalent of C++ operator|(mods, key).
func Combine(mods KeyboardModifiers, key Key) QKeyCombination {
	return NewQKeyCombinationFromModifiersKey(mods, key)
}

// KeyboardModifiers returns the keyboard modifiers part of the combination.
func (q QKeyCombination) KeyboardModifiers() KeyboardModifiers {
	return KeyboardModifiers(C.qkeycombination_keyboard_modifiers(q.v))
}

// Key returns the key part of the combination.
func (q QKeyCombination) Key() Key {
	return Key(C.qkeycombination_key(q.v))
}

// ToCombined returns the raw int representation.
func (q QKeyCombination) ToCombined() int {
	return int(C.qkeycombination_to_combined(q.v))
}

// Equal reports whether q and other hold the same raw combination.
func (q QKeyCombination) Equal(other QKeyCombination) bool {
	return q.ToCombined() == other.ToCombined()
}
