package qurl

/*
   QUrlQuery — обёртка над одноимённым C++-классом. Владение — финализатор
   и явный Close(). Строки передаются парами (ptr, size), без C.CString.
   Все QString-варианты в Go конвертируют входной *qtext.QString в string
   через GoString() (поле .p в qtext.QString неэкспортировано); QString
   C ABI-перегрузки не нужны и не объявляются.

   QList<std::pair<QString, QString>> передаётся через плоский буфер +
   offsets (без Go-указателей на Go-память — правило cgo). Для возврата
   данных C сам выделяет offsets через malloc; Go освобождает их через
   qurlquery_free_offsets. Это устраняет необходимость угадывать
   количество элементов заранее.

   Константы QueryValueDelimiterDefault / QueryPairDelimiterDefault — uint16,
   соответствуют constexpr char16_t из qurlquery.h.
*/

/*
#cgo pkg-config: Qt6Core
#cgo CXXFLAGS: -std=c++17 -fPIC -Wno-sfinae-incomplete
#cgo LDFLAGS: -lstdc++
#include <stdint.h>
#include <stdbool.h>
#include <stddef.h>
#include "qurlquery_c.h"
*/
import "C"

import (
	"runtime"
	"unsafe"

	"github.com/vadimlxx/go-qt6core/qtext"
)

// QueryValueDelimiterDefault is QUrlQuery::defaultQueryValueDelimiter() (u'=').
const QueryValueDelimiterDefault uint16 = '='

// QueryPairDelimiterDefault is QUrlQuery::defaultQueryPairDelimiter() (u'&').
const QueryPairDelimiterDefault uint16 = '&'

// QUrlQuery wraps a Qt6 QUrlQuery object.
type QUrlQuery struct {
	p *C.QUrlQuery
}

func wrapQUrlQuery(p unsafe.Pointer) *QUrlQuery {
	if p == nil {
		return nil
	}
	q := &QUrlQuery{p: (*C.QUrlQuery)(p)}
	runtime.SetFinalizer(q, finalizeQUrlQuery)
	return q
}

func finalizeQUrlQuery(q *QUrlQuery) {
	if q.p != nil {
		C.qurlquery_delete(q.p)
		q.p = nil
	}
}

// Close releases the underlying C++ object immediately and disables the finalizer.
// Qt: QUrlQuery::~QUrlQuery().
func (q *QUrlQuery) Close() {
	if q.p != nil {
		C.qurlquery_delete(q.p)
		q.p = nil
		runtime.SetFinalizer(q, nil)
	}
}

// ---- construction --------------------------------------------------------

// NewQUrlQuery returns an empty QUrlQuery.
// Qt: QUrlQuery().
func NewQUrlQuery(arg ...any) *QUrlQuery {
	if len(arg) == 0 || arg[0] == nil {
		return wrapQUrlQuery(unsafe.Pointer(C.qurlquery_new_empty()))
	}

	switch v := arg[0].(type) {
	case *QUrl:
		return QUrlQueryFromUrl(v)
	case *qtext.QString:
		return QUrlQueryFromQString(v)
	case string:
		return QUrlQueryFromString(v)
	}

	panic("NewQUrlQuery invalid argument")
}

// QUrlQueryFromUrl builds a query from the query component of url.
// Qt: QUrlQuery(const QUrl&).
func QUrlQueryFromUrl(url *QUrl) *QUrlQuery {
	if url == nil {
		return NewQUrlQuery()
	}
	return wrapQUrlQuery(unsafe.Pointer(C.qurlquery_new_from_url(url.p)))
}

// QUrlQueryFromQString parses a query string (QString).
// Qt: QUrlQuery(const QString&) — with QString converted to UTF-8 internally.
func QUrlQueryFromQString(s *qtext.QString) *QUrlQuery {
	if s == nil {
		return NewQUrlQuery()
	}
	return QUrlQueryFromString(s.GoString())
}

// QUrlQueryFromString parses a query string (Go string, UTF-8).
// Qt: QUrlQuery(const QString&) — with QString::fromUtf8.
func QUrlQueryFromString(s string) *QUrlQuery {
	var p *C.char
	if len(s) > 0 {
		p = (*C.char)(unsafe.Pointer(unsafe.StringData(s)))
	}
	q := wrapQUrlQuery(unsafe.Pointer(C.qurlquery_new_from_utf8(p, C.size_t(len(s)))))
	runtime.KeepAlive(s)
	return q
}

// Clone returns a deep copy of the query.
// Qt: QUrlQuery(const QUrlQuery&) — copy constructor.
func (q *QUrlQuery) Clone() *QUrlQuery {
	return wrapQUrlQuery(unsafe.Pointer(C.qurlquery_new_copy(q.p)))
}

// Swap exchanges the contents of q and other.
// Qt: QUrlQuery::swap(QUrlQuery&).
func (q *QUrlQuery) Swap(other *QUrlQuery) {
	C.qurlquery_swap(q.p, other.p)
	runtime.KeepAlive(other)
}

// ---- state ---------------------------------------------------------------

// IsEmpty reports whether the query has no items.
// Qt: QUrlQuery::isEmpty() const.
func (q *QUrlQuery) IsEmpty() bool { return bool(C.qurlquery_is_empty(q.p)) }

// IsDetached reports whether the query is detached from any shared copy.
// Qt: QUrlQuery::isDetached() const.
func (q *QUrlQuery) IsDetached() bool { return bool(C.qurlquery_is_detached(q.p)) }

// Clear removes all query items.
// Qt: QUrlQuery::clear().
func (q *QUrlQuery) Clear() { C.qurlquery_clear(q.p) }

// ---- query string --------------------------------------------------------

// QueryString returns the query as a Go string.
// Qt: QUrlQuery::query(ComponentFormattingOptions) const → QString (as string).
func (q *QUrlQuery) QueryString(encoding ComponentFormattingOptions) string {
	var n C.size_t
	p := C.qurlquery_query(q.p, C.int(encoding), &n)
	return qurlQueryGoStringN(p, n)
}

// Query returns the query as a QString.
// Qt: QUrlQuery::query(ComponentFormattingOptions) const → QString.
func (q *QUrlQuery) Query(encoding ComponentFormattingOptions) *qtext.QString {
	return qtext.QStringFromString(q.QueryString(encoding))
}

// ToString returns the query as a Go string (alias of QueryString).
// Qt: QUrlQuery::toString(ComponentFormattingOptions) const → QString (as string).
func (q *QUrlQuery) ToString(encoding ComponentFormattingOptions) string {
	var n C.size_t
	p := C.qurlquery_to_string(q.p, C.int(encoding), &n)
	return qurlQueryGoStringN(p, n)
}

// ToQString returns the query as a QString (alias of Query).
// Qt: QUrlQuery::toString(ComponentFormattingOptions) const → QString.
func (q *QUrlQuery) ToQString(encoding ComponentFormattingOptions) *qtext.QString {
	return qtext.QStringFromString(q.ToString(encoding))
}

// SetQueryQString replaces the query with the given QString.
// Qt: QUrlQuery::setQuery(const QString&).
func (q *QUrlQuery) SetQueryQString(s *qtext.QString) {
	if s == nil {
		q.SetQueryString("")
		return
	}
	q.SetQueryString(s.GoString())
}

// SetQueryString replaces the query with the given Go string (UTF-8).
// Qt: QUrlQuery::setQuery(const QString&) — with QString::fromUtf8.
func (q *QUrlQuery) SetQueryString(s string) {
	var p *C.char
	if len(s) > 0 {
		p = (*C.char)(unsafe.Pointer(unsafe.StringData(s)))
	}
	C.qurlquery_set_query_utf8(q.p, p, C.size_t(len(s)))
	runtime.KeepAlive(s)
}

// ---- delimiters ----------------------------------------------------------

// SetQueryDelimiters sets the value and pair delimiters.
// Qt: QUrlQuery::setQueryDelimiters(QChar, QChar).
func (q *QUrlQuery) SetQueryDelimiters(valueDelim, pairDelim qtext.QChar) {
	C.qurlquery_set_query_delimiters(q.p, C.uint16_t(valueDelim.Unicode()),
		C.uint16_t(pairDelim.Unicode()))
}

// QueryValueDelimiter returns the value delimiter.
// Qt: QUrlQuery::queryValueDelimiter() const → QChar.
func (q *QUrlQuery) QueryValueDelimiter() qtext.QChar {
	return qtext.NewQChar(uint16(C.qurlquery_query_value_delimiter(q.p)))
}

// QueryPairDelimiter returns the pair delimiter.
// Qt: QUrlQuery::queryPairDelimiter() const → QChar.
func (q *QUrlQuery) QueryPairDelimiter() qtext.QChar {
	return qtext.NewQChar(uint16(C.qurlquery_query_pair_delimiter(q.p)))
}

// ---- items (map) ---------------------------------------------------------

// SetQueryItems replaces all items with the given map.
// The order of map iteration is not preserved; if multiple items share a key,
// the last one encountered wins (mirroring repeated addQueryItem).
// Qt: QUrlQuery::setQueryItems(const QList<std::pair<QString,QString>>&).
func (q *QUrlQuery) SetQueryItems(items map[string]string) {
	if len(items) == 0 {
		q.Clear()
		return
	}

	keys := make([]string, 0, len(items))
	values := make([]string, 0, len(items))
	for k, v := range items {
		keys = append(keys, k)
		values = append(values, v)
	}

	var totalKeys, totalVals int
	for _, k := range keys { totalKeys += len(k) }
	for _, v := range values { totalVals += len(v) }
	total := totalKeys + totalVals

	var data []byte
	if total > 0 {
		data = make([]byte, 0, total)
		for _, k := range keys { data = append(data, k...) }
		for _, v := range values { data = append(data, v...) }
	}

	keyOffsets := make([]C.size_t, len(keys)+1)
	keyOffsets[0] = 0
	for i, k := range keys {
		keyOffsets[i+1] = keyOffsets[i] + C.size_t(len(k))
	}
	valOffsets := make([]C.size_t, len(values)+1)
	valOffsets[0] = C.size_t(totalKeys)
	for i, v := range values {
		valOffsets[i+1] = valOffsets[i] + C.size_t(len(v))
	}

	var dataPtr *C.char
	if len(data) > 0 {
		dataPtr = (*C.char)(unsafe.Pointer(&data[0]))
	}
	C.qurlquery_set_query_items(q.p, dataPtr,
		(*C.size_t)(unsafe.Pointer(&keyOffsets[0])),
		(*C.size_t)(unsafe.Pointer(&valOffsets[0])),
		C.size_t(len(keys)))
	runtime.KeepAlive(data)
	runtime.KeepAlive(keyOffsets)
	runtime.KeepAlive(valOffsets)
}

// QueryItems returns all query items as a map.
// The order is not preserved; if multiple items share a key, the last one in
// the internal list wins.
// Qt: QUrlQuery::queryItems(ComponentFormattingOptions) const
//     → QList<std::pair<QString,QString>> (returned here as map).
func (q *QUrlQuery) QueryItems(encoding ComponentFormattingOptions) map[string]string {
	var ko, vo *C.size_t
	var count C.size_t
	buf := C.qurlquery_query_items(q.p, C.int(encoding), &ko, &vo, &count)
	if buf == nil || count == 0 {
		return nil
	}
	defer C.qurlquery_free_string(buf)
	if ko != nil {
		defer C.qurlquery_free_offsets(ko)
	}
	if vo != nil {
		defer C.qurlquery_free_offsets(vo)
	}

	n := int(count)
	keyOffsets := unsafe.Slice(ko, n+1)
	valOffsets := unsafe.Slice(vo, n+1)
	data := C.GoBytes(unsafe.Pointer(buf), C.int(valOffsets[n]))

	result := make(map[string]string, n)
	for i := 0; i < n; i++ {
		kLo, kHi := int(keyOffsets[i]), int(keyOffsets[i+1])
		vLo, vHi := int(valOffsets[i]), int(valOffsets[i+1])
		result[string(data[kLo:kHi])] = string(data[vLo:vHi])
	}
	return result
}

// ---- single-item operations ---------------------------------------------

// HasQueryItemQString reports whether key is present (QString variant).
// Qt: QUrlQuery::hasQueryItem(const QString&) const.
func (q *QUrlQuery) HasQueryItemQString(key *qtext.QString) bool {
	if key == nil {
		return false
	}
	return q.HasQueryItem(key.GoString())
}

// HasQueryItem reports whether key is present (Go string, UTF-8).
// Qt: QUrlQuery::hasQueryItem(const QString&) const — with QString::fromUtf8.
func (q *QUrlQuery) HasQueryItem(key string) bool {
	var p *C.char
	if len(key) > 0 { p = (*C.char)(unsafe.Pointer(unsafe.StringData(key))) }
	r := bool(C.qurlquery_has_query_item_utf8(q.p, p, C.size_t(len(key))))
	runtime.KeepAlive(key)
	return r
}

// AddQueryItemQString appends an item (key, value) using Qt strings.
// Qt: QUrlQuery::addQueryItem(const QString&, const QString&).
func (q *QUrlQuery) AddQueryItemQString(key, value *qtext.QString) {
	k, v := "", ""
	if key != nil { k = key.GoString() }
	if value != nil { v = value.GoString() }
	q.AddQueryItem(k, v)
}

// AddQueryItem appends an item (key, value) using Go strings (UTF-8).
// Qt: QUrlQuery::addQueryItem(const QString&, const QString&) — with QString::fromUtf8.
func (q *QUrlQuery) AddQueryItem(key, value string) {
	var kp, vp *C.char
	if len(key) > 0 { kp = (*C.char)(unsafe.Pointer(unsafe.StringData(key))) }
	if len(value) > 0 { vp = (*C.char)(unsafe.Pointer(unsafe.StringData(value))) }
	C.qurlquery_add_query_item_utf8_utf8(q.p, kp, C.size_t(len(key)), vp, C.size_t(len(value)))
	runtime.KeepAlive(key)
	runtime.KeepAlive(value)
}

// RemoveQueryItemQString removes the first item whose key matches.
// Qt: QUrlQuery::removeQueryItem(const QString&).
func (q *QUrlQuery) RemoveQueryItemQString(key *qtext.QString) {
	if key == nil {
		return
	}
	q.RemoveQueryItem(key.GoString())
}

// RemoveQueryItem removes the first item whose key matches (Go string).
// Qt: QUrlQuery::removeQueryItem(const QString&) — with QString::fromUtf8.
func (q *QUrlQuery) RemoveQueryItem(key string) {
	var p *C.char
	if len(key) > 0 { p = (*C.char)(unsafe.Pointer(unsafe.StringData(key))) }
	C.qurlquery_remove_query_item_utf8(q.p, p, C.size_t(len(key)))
	runtime.KeepAlive(key)
}

// QueryItemValueQString returns the value of the first item whose key matches.
// Qt: QUrlQuery::queryItemValue(const QString&, ComponentFormattingOptions) const.
func (q *QUrlQuery) QueryItemValueQString(key *qtext.QString, encoding ComponentFormattingOptions) string {
	if key == nil {
		return ""
	}
	return q.QueryItemValue(key.GoString(), encoding)
}

// QueryItemValue returns the value of the first item whose key matches (Go string).
// Qt: QUrlQuery::queryItemValue(const QString&, ComponentFormattingOptions) const.
func (q *QUrlQuery) QueryItemValue(key string, encoding ComponentFormattingOptions) string {
	var p *C.char
	if len(key) > 0 { p = (*C.char)(unsafe.Pointer(unsafe.StringData(key))) }
	var n C.size_t
	out := C.qurlquery_query_item_value_utf8(q.p, p, C.size_t(len(key)), C.int(encoding), &n)
	runtime.KeepAlive(key)
	return qurlQueryGoStringN(out, n)
}

// AllQueryItemValuesQString returns all values for the given key (QString variant).
// Qt: QUrlQuery::allQueryItemValues(const QString&, ComponentFormattingOptions) const
//     → QStringList (returned here as []string).
func (q *QUrlQuery) AllQueryItemValuesQString(key *qtext.QString, encoding ComponentFormattingOptions) []string {
	if key == nil {
		return nil
	}
	return q.AllQueryItemValues(key.GoString(), encoding)
}

// AllQueryItemValues returns all values for the given key (Go string).
// Qt: QUrlQuery::allQueryItemValues(const QString&, ComponentFormattingOptions) const.
func (q *QUrlQuery) AllQueryItemValues(key string, encoding ComponentFormattingOptions) []string {
	var p *C.char
	if len(key) > 0 { p = (*C.char)(unsafe.Pointer(unsafe.StringData(key))) }
	var offs *C.size_t
	var count C.size_t
	buf := C.qurlquery_all_query_item_values_utf8(q.p, p, C.size_t(len(key)), C.int(encoding), &offs, &count)
	runtime.KeepAlive(key)
	if buf == nil || count == 0 {
		return nil
	}
	defer C.qurlquery_free_string(buf)
	if offs != nil {
		defer C.qurlquery_free_offsets(offs)
	}

	n := int(count)
	offsets := unsafe.Slice(offs, n+1)
	data := C.GoBytes(unsafe.Pointer(buf), C.int(offsets[n]))

	result := make([]string, n)
	for i := 0; i < n; i++ {
		lo, hi := int(offsets[i]), int(offsets[i+1])
		result[i] = string(data[lo:hi])
	}
	return result
}

// RemoveAllQueryItemsQString removes all items whose key matches.
// Qt: QUrlQuery::removeAllQueryItems(const QString&).
func (q *QUrlQuery) RemoveAllQueryItemsQString(key *qtext.QString) {
	if key == nil {
		return
	}
	q.RemoveAllQueryItems(key.GoString())
}

// RemoveAllQueryItems removes all items whose key matches (Go string).
// Qt: QUrlQuery::removeAllQueryItems(const QString&).
func (q *QUrlQuery) RemoveAllQueryItems(key string) {
	var p *C.char
	if len(key) > 0 { p = (*C.char)(unsafe.Pointer(unsafe.StringData(key))) }
	C.qurlquery_remove_all_query_items_utf8(q.p, p, C.size_t(len(key)))
	runtime.KeepAlive(key)
}

// ---- helpers -------------------------------------------------------------

// qurlQueryGoStringN converts a heap-allocated NUL-terminated C string to a
// Go string and releases the C allocation via qurlquery_free_string.
// Named distinctly from qurl.go's goStringN to avoid redeclaration.
func qurlQueryGoStringN(p *C.char, n C.size_t) string {
	if p == nil || n == 0 {
		return ""
	}
	defer C.qurlquery_free_string(p)
	return C.GoStringN(p, C.int(n))
}
