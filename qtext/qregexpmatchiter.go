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
	"sync/atomic"
	"unsafe"
)

type qregexpmatchiterState struct {
	p atomic.Pointer[byte]
}

func (s *qregexpmatchiterState) cpp() *C.QRegularExpressionMatchIterator {
	if bp := s.p.Load(); bp != nil {
		return (*C.QRegularExpressionMatchIterator)(unsafe.Pointer(bp))
	}
	return nil
}

func (s *qregexpmatchiterState) claim() *C.QRegularExpressionMatchIterator {
	if bp := s.p.Swap(nil); bp != nil {
		return (*C.QRegularExpressionMatchIterator)(unsafe.Pointer(bp))
	}
	return nil
}

// QRegularExpressionMatchIterator wraps a Qt6 QRegularExpressionMatchIterator.
type QRegularExpressionMatchIterator struct {
	// p *C.QRegularExpressionMatchIterator
	state   *qregexpmatchiterState
	cleanup runtime.Cleanup
}

func wrapQRegularExpressionMatchIterator(p unsafe.Pointer) *QRegularExpressionMatchIterator {
	if p == nil {
		return nil
	}

	st := &qregexpmatchiterState{}
	st.p.Store((*byte)(p))
	q := &QRegularExpressionMatchIterator{state: st}
	q.cleanup = runtime.AddCleanup(q, func(st *qregexpmatchiterState) {
		if cpp := st.claim(); cpp != nil {
			C.qregexpmatchiter_delete(cpp)
		}
	}, st)
	return q
}

func (q *QRegularExpressionMatchIterator) handle() *C.QRegularExpressionMatchIterator {
	if cpp, ok := q.tryHandle(); ok {
		return cpp
	}
	panic("qregexpmatchiter: use after Close (or nil receiver)")
}

func (q *QRegularExpressionMatchIterator) tryHandle() (*C.QRegularExpressionMatchIterator, bool) {
	if q == nil || q.state == nil {
		return nil, false
	}
	if cpp := q.state.cpp(); cpp != nil {
		return cpp, true
	}
	return nil, false
}

func (q *QRegularExpressionMatchIterator) Close() {
	if q == nil || q.state == nil {
		return
	}
	q.cleanup.Stop()
	if cpp := q.state.claim(); cpp != nil {
		C.qregexpmatchiter_delete(cpp)
	}
}

func NewQRegularExpressionMatchIterator() *QRegularExpressionMatchIterator {
	return wrapQRegularExpressionMatchIterator(unsafe.Pointer(C.qregexpmatchiter_new_empty()))
}

func (it *QRegularExpressionMatchIterator) Clone() *QRegularExpressionMatchIterator {
	return wrapQRegularExpressionMatchIterator(unsafe.Pointer(C.qregexpmatchiter_new_copy(it.handle())))
}

func (it *QRegularExpressionMatchIterator) IsValid() bool {
	return bool(C.qregexpmatchiter_is_valid(it.handle()))
}

func (it *QRegularExpressionMatchIterator) HasNext() bool {
	return bool(C.qregexpmatchiter_has_next(it.handle()))
}

func (it *QRegularExpressionMatchIterator) Next() *QRegularExpressionMatch {
	return wrapQRegularExpressionMatch(unsafe.Pointer(C.qregexpmatchiter_next(it.handle())))
}

func (it *QRegularExpressionMatchIterator) PeekNext() *QRegularExpressionMatch {
	return wrapQRegularExpressionMatch(unsafe.Pointer(C.qregexpmatchiter_peek_next(it.handle())))
}

func (it *QRegularExpressionMatchIterator) RegularExpression() *QRegularExpression {
	return wrapQRegularExpression(unsafe.Pointer(C.qregexpmatchiter_regular_expression(it.handle())))
}

func (it *QRegularExpressionMatchIterator) MatchType() MatchType {
	return MatchType(C.qregexpmatchiter_match_type(it.handle()))
}

func (it *QRegularExpressionMatchIterator) MatchOptions() MatchOption {
	return MatchOption(C.qregexpmatchiter_match_options(it.handle()))
}

func (it *QRegularExpressionMatchIterator) Swap(other *QRegularExpressionMatchIterator) {
	C.qregexpmatchiter_swap(it.handle(), other.handle())
}
