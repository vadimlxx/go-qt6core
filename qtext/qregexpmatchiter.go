package qtext

/*
   QRegularExpressionMatchIterator — обёртка над одноимённым C++-объектом.
   Внутренне Qt хранит копию subject, поэтому отдельный keepAlive
   не требуется: итератор владеет всем необходимым для возврата
   CapturedView-ов из последовательных next().

   Метод Next() возвращает новый *QRegularExpressionMatch; его
   владение — за вызывающим (закрывается финализатором или Close()).
   PeekNext() не мутирует итератор.

   Swap() реализован как метод, чтобы не засорять пакет одноимённой
   free-функцией (в Qt это member + ADL-функция, в Go — только метод).

 Что сознательно не вошло (пока):

   - Range-based for: qregularexpression.h определяет
     QRegularExpressionMatchIteratorRangeBasedForIterator с
     Q_DECLARE_EQUALITY_COMPARABLE. В Go это не идиоматично;
     используйте HasNext() / Next().
   - swap() free-функция — покрыт методом Swap().
   - Операторы сравнения == и != для итератора — не имеют
     осмысленной семантики в Go (пользователь итерирует вручную).
*/

/*
#cgo pkg-config: Qt6Core
#cgo CXXFLAGS: -std=c++17 -fPIC -Wno-sfinae-incomplete
#cgo LDFLAGS: -lstdc++
#include <stdint.h>
#include <stdbool.h>
#include <stddef.h>
#include "qregexp_c.h"
*/
import "C"

import (
	"runtime"
	"unsafe"
)

// QRegularExpressionMatchIterator wraps a Qt6 QRegularExpressionMatchIterator.
type QRegularExpressionMatchIterator struct {
	p *C.QRegularExpressionMatchIterator
}

func wrapQRegularExpressionMatchIterator(p unsafe.Pointer) *QRegularExpressionMatchIterator {
	if p == nil {
		return nil
	}
	it := &QRegularExpressionMatchIterator{p: (*C.QRegularExpressionMatchIterator)(p)}
	runtime.SetFinalizer(it, finalizeQRegularExpressionMatchIterator)
	return it
}

func finalizeQRegularExpressionMatchIterator(it *QRegularExpressionMatchIterator) {
	if it.p != nil {
		C.qregexpmatchiter_delete(it.p)
		it.p = nil
	}
}

func (it *QRegularExpressionMatchIterator) Close() {
	if it.p != nil {
		C.qregexpmatchiter_delete(it.p)
		it.p = nil
		runtime.SetFinalizer(it, nil)
	}
}

func NewQRegularExpressionMatchIterator() *QRegularExpressionMatchIterator {
	return wrapQRegularExpressionMatchIterator(unsafe.Pointer(C.qregexpmatchiter_new_empty()))
}

func (it *QRegularExpressionMatchIterator) Clone() *QRegularExpressionMatchIterator {
	return wrapQRegularExpressionMatchIterator(unsafe.Pointer(C.qregexpmatchiter_new_copy(it.p)))
}

func (it *QRegularExpressionMatchIterator) IsValid() bool {
	return bool(C.qregexpmatchiter_is_valid(it.p))
}

func (it *QRegularExpressionMatchIterator) HasNext() bool {
	return bool(C.qregexpmatchiter_has_next(it.p))
}

func (it *QRegularExpressionMatchIterator) Next() *QRegularExpressionMatch {
	return wrapQRegularExpressionMatch(unsafe.Pointer(C.qregexpmatchiter_next(it.p)))
}

func (it *QRegularExpressionMatchIterator) PeekNext() *QRegularExpressionMatch {
	return wrapQRegularExpressionMatch(unsafe.Pointer(C.qregexpmatchiter_peek_next(it.p)))
}

func (it *QRegularExpressionMatchIterator) RegularExpression() *QRegularExpression {
	return wrapQRegularExpression(unsafe.Pointer(C.qregexpmatchiter_regular_expression(it.p)))
}

func (it *QRegularExpressionMatchIterator) MatchType() MatchType {
	return MatchType(C.qregexpmatchiter_match_type(it.p))
}

func (it *QRegularExpressionMatchIterator) MatchOptions() MatchOption {
	return MatchOption(C.qregexpmatchiter_match_options(it.p))
}

func (it *QRegularExpressionMatchIterator) Swap(other *QRegularExpressionMatchIterator) {
	C.qregexpmatchiter_swap(it.p, other.p)
}
