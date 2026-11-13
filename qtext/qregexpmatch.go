package qtext

/*
   QRegularExpressionMatch — обёртка над одноимённым C++-объектом через
   непрозрачный указатель. Хранит указатель C-ABI, а также (в поле
   keepAlive) ссылку на subject, чтобы гарантировать, что буфер UTF-16
   не будет освобождён, пока жив хотя бы один captured-view.

   Методы CapturedView* возвращают QStringView, который держит keepAlive
   на сам матч (newQStringViewWithOwner). Даже если вызывающий не сохранил
   *QRegularExpressionMatch, а только QStringView, сборщик мусора не
   освободит исходный буфер раньше времени.

   Перегрузки "по имени" (CapturedName, HasCapturedName, CapturedStartName,
   ...) принимают Go string и внутри создают временный QString — это
   единственное место в текущей реализации, где происходит выделение под
   имя группы. Для zero-copy варианта нужны *Name-перегрузки с уже
   готовым *QString — они отложены, чтобы не плодить имена.

 Что сознательно не вошло (пока):

   - capturedTexts() — возвращает QStringList (нужна специализация
     контейнера QStringList).
   - Операторы сравнения == и != (Q_DECLARE_EQUALITY_COMPARABLE) — будут
     добавлены вместе с Compare()-хелпером на уровне QRegularExpressionMatch.
   - Range-based for по capturedTexts() — неприменимо в Go.
   - swap() — реализован как метод Swap().
   - Перегрузки captured()/hasCaptured()/... по *QString — API уже
     принимает Go string; для zero-copy overload нужен отдельный
     QStringView-вариант, отложен.
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

// QRegularExpressionMatch wraps a Qt6 QRegularExpressionMatch object.
// The match keeps a reference to the subject passed to MatchString or MatchView
// so that CapturedView remains valid for the lifetime of the match.
type QRegularExpressionMatch struct {
	p         *C.QRegularExpressionMatch
	keepAlive interface{}
}

func wrapQRegularExpressionMatch(p unsafe.Pointer) *QRegularExpressionMatch {
	if p == nil {
		return nil
	}
	m := &QRegularExpressionMatch{p: (*C.QRegularExpressionMatch)(p)}
	runtime.SetFinalizer(m, finalizeQRegularExpressionMatch)
	return m
}

func finalizeQRegularExpressionMatch(m *QRegularExpressionMatch) {
	if m.p != nil {
		C.qregexpmatch_delete(m.p)
		m.p = nil
	}
}

func (m *QRegularExpressionMatch) Close() {
	if m.p != nil {
		C.qregexpmatch_delete(m.p)
		m.p = nil
		runtime.SetFinalizer(m, nil)
	}
}

func NewQRegularExpressionMatch() *QRegularExpressionMatch {
	return wrapQRegularExpressionMatch(unsafe.Pointer(C.qregexpmatch_new_empty()))
}

func (m *QRegularExpressionMatch) Clone() *QRegularExpressionMatch {
	return wrapQRegularExpressionMatch(unsafe.Pointer(C.qregexpmatch_new_copy(m.p)))
}

func (m *QRegularExpressionMatch) RegularExpression() *QRegularExpression {
	return wrapQRegularExpression(unsafe.Pointer(C.qregexpmatch_regular_expression(m.p)))
}

func (m *QRegularExpressionMatch) MatchType() MatchType {
	return MatchType(C.qregexpmatch_match_type(m.p))
}

func (m *QRegularExpressionMatch) MatchOptions() MatchOption {
	return MatchOption(C.qregexpmatch_match_options(m.p))
}

func (m *QRegularExpressionMatch) HasMatch() bool {
	return bool(C.qregexpmatch_has_match(m.p))
}

func (m *QRegularExpressionMatch) HasPartialMatch() bool {
	return bool(C.qregexpmatch_has_partial_match(m.p))
}

func (m *QRegularExpressionMatch) IsValid() bool {
	return bool(C.qregexpmatch_is_valid(m.p))
}

func (m *QRegularExpressionMatch) LastCapturedIndex() int {
	return int(C.qregexpmatch_last_captured_index(m.p))
}

func (m *QRegularExpressionMatch) HasCapturedInt(nth int) bool {
	return bool(C.qregexpmatch_has_captured_int(m.p, C.int(nth)))
}

func (m *QRegularExpressionMatch) HasCapturedName(name string) bool {
	qname := QStringFromString(name)
	defer qname.Close()
	return bool(C.qregexpmatch_has_captured_view(m.p,
		(*C.uint16_t)(qname.Utf16Unsafe()), C.ptrdiff_t(qname.Size())))
}

func (m *QRegularExpressionMatch) CapturedInt(nth int) *QString {
	return wrapQString(unsafe.Pointer(C.qregexpmatch_captured_int(m.p, C.int(nth))))
}

func (m *QRegularExpressionMatch) CapturedName(name string) *QString {
	qname := QStringFromString(name)
	defer qname.Close()
	return wrapQString(unsafe.Pointer(C.qregexpmatch_captured_view(m.p,
		(*C.uint16_t)(qname.Utf16Unsafe()), C.ptrdiff_t(qname.Size()))))
}

func (m *QRegularExpressionMatch) CapturedViewInt(nth int) QStringView {
	var outData *C.uint16_t
	var outSize C.ptrdiff_t
	C.qregexpmatch_captured_view_int(m.p, C.int(nth), &outData, &outSize)
	return newQStringViewWithOwner(unsafe.Pointer(outData), int(outSize), m)
}

func (m *QRegularExpressionMatch) CapturedViewName(name string) QStringView {
	qname := QStringFromString(name)
	defer qname.Close()
	var outData *C.uint16_t
	var outSize C.ptrdiff_t
	C.qregexpmatch_captured_view_view(m.p,
		(*C.uint16_t)(qname.Utf16Unsafe()), C.ptrdiff_t(qname.Size()),
		&outData, &outSize)
	return newQStringViewWithOwner(unsafe.Pointer(outData), int(outSize), m)
}

func (m *QRegularExpressionMatch) CapturedStartInt(nth int) int {
	return int(C.qregexpmatch_captured_start_int(m.p, C.int(nth)))
}

func (m *QRegularExpressionMatch) CapturedLengthInt(nth int) int {
	return int(C.qregexpmatch_captured_length_int(m.p, C.int(nth)))
}

func (m *QRegularExpressionMatch) CapturedEndInt(nth int) int {
	return int(C.qregexpmatch_captured_end_int(m.p, C.int(nth)))
}

func (m *QRegularExpressionMatch) CapturedStartName(name string) int {
	qname := QStringFromString(name)
	defer qname.Close()
	return int(C.qregexpmatch_captured_start_view(m.p,
		(*C.uint16_t)(qname.Utf16Unsafe()), C.ptrdiff_t(qname.Size())))
}

func (m *QRegularExpressionMatch) CapturedLengthName(name string) int {
	qname := QStringFromString(name)
	defer qname.Close()
	return int(C.qregexpmatch_captured_length_view(m.p,
		(*C.uint16_t)(qname.Utf16Unsafe()), C.ptrdiff_t(qname.Size())))
}

func (m *QRegularExpressionMatch) CapturedEndName(name string) int {
	qname := QStringFromString(name)
	defer qname.Close()
	return int(C.qregexpmatch_captured_end_view(m.p,
		(*C.uint16_t)(qname.Utf16Unsafe()), C.ptrdiff_t(qname.Size())))
}

func (m *QRegularExpressionMatch) Swap(other *QRegularExpressionMatch) {
	C.qregexpmatch_swap(m.p, other.p)
}
