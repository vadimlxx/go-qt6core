package qurl

/*
   QUrl в Qt 6.11 — value-type с implicit sharing (QUrlPrivate* внутри).
   В Go — непрозрачный указатель на heap-аллоцированный QUrl; владение —
   финализатор и явный Close().

   Все флаговые типы (UrlFormattingOption, ComponentFormattingOption,
   ComponentFormattingOptions, FormattingOptions) — алиасы int, повторяют
   семантику QUrlTwoFlags (принимает флаги обоих enum-ов без приведения).

   Для каждой функции, использующей QString, предоставлены две формы:
     - Qt-версия (bare name), принимающая/возвращающая *qtext.QString;
     - Go-версия с суффиксом String, принимающая/возвращающая string.
   Для QByteArray — суффикс Bytes; для QStringList — суффикс Strings.

   QUrlQuery interop (см. qurlquery.go):
     - SetQueryQUrlQuery(q) — QUrl::setQuery(const QUrlQuery&);
     - QUrlQueryFromUrl(url) — QUrlQuery(const QUrl&) — уже реализован
       в qurlquery.go.

   FromStringListStrings передаёт входные строки как плоский буфер
   (склейка в []byte + массив offsets), а не как массив Go-указателей на
   Go-память — это требование cgo (запрещено передавать Go pointer to
   unpinned Go pointer).

   Что сознательно не вошло:

   - toCFURL / fromCFURL, toNSURL / fromNSURL — только Darwin.
   - idnWhitelist() / setIdnWhitelist(QStringList) — TODO.
   - QDataStream / QDebug операторы — не идиоматичны в Go.
   - qHash(QUrl) — не критично.
   - operator<, operator==, operator!= — заменены на Compare/Equal.
   - QUrlTwoFlags-операторы — Go использует int + побитовое ИЛИ.
*/

/*
#cgo pkg-config: Qt6Core
#cgo CXXFLAGS: -std=c++17 -fPIC -Wno-sfinae-incomplete
#cgo LDFLAGS: -lstdc++
#include <stdint.h>
#include <stdbool.h>
#include <stddef.h>
#include "qurl_c.h"
*/
import "C"

import (
	"runtime"
	"strings"
	"sync/atomic"
	"unsafe"

	"github.com/vadimlxx/go-qt6core/qtext"
)

type qurlState struct {
	p atomic.Pointer[byte]
}

func (s *qurlState) cpp() *C.QUrl {
	if bp := s.p.Load(); bp != nil {
		return (*C.QUrl)(unsafe.Pointer(bp))
	}
	return nil
}

func (s *qurlState) claim() *C.QUrl {
	if bp := s.p.Swap(nil); bp != nil {
		return (*C.QUrl)(unsafe.Pointer(bp))
	}
	return nil
}

// QUrl wraps a Qt6 QUrl object.
type QUrl struct {
	// p *C.QUrl
	state   *qurlState
	cleanup runtime.Cleanup
}

func wrapQUrl(p unsafe.Pointer) *QUrl {
	if p == nil {
		return nil
	}

	st := &qurlState{}
	st.p.Store((*byte)(p))
	q := &QUrl{state: st}
	q.cleanup = runtime.AddCleanup(q, func(st *qurlState) {
		if cpp := st.claim(); cpp != nil {
			C.qurl_delete(cpp)
		}
	}, st)

	return q
}

func (q *QUrl) handle() *C.QUrl {
	if cpp, ok := q.tryHandle(); ok {
		return cpp
	}
	panic("qurl: use after Close (or nil receiver)")
}

func (q *QUrl) tryHandle() (*C.QUrl, bool) {
	if q == nil || q.state == nil {
		return nil, false
	}
	if cpp := q.state.cpp(); cpp != nil {
		return cpp, true
	}
	return nil, false
}

// Close releases the underlying C++ object immediately and disables the finalizer.
// Qt: QUrl::~QUrl().
func (q *QUrl) Close() {
	if q == nil || q.state == nil {
		return
	}

	q.cleanup.Stop()
	if cpp := q.state.claim(); cpp != nil {
		C.qurl_delete(cpp)
	}
}

// ---- construction --------------------------------------------------------

// QUrlFromQString parses a URL from a QString. nil maps to the empty QUrl.
// Qt: QUrl(const QString&, ParsingMode).
func QUrlFromQString(s *qtext.QString, mode ParsingMode) *QUrl {
	if s == nil {
		return NewQUrl()
	}
	return QUrlFromString(s.GoString(), mode)
}

// QUrlFromString parses a URL from a Go string (UTF-8).
// Qt: QUrl(const QString&, ParsingMode) — with QString::fromUtf8.
func QUrlFromString(s string, mode ParsingMode) *QUrl {
	var p *C.char
	if len(s) > 0 {
		p = (*C.char)(unsafe.Pointer(unsafe.StringData(s)))
	}
	q := wrapQUrl(unsafe.Pointer(C.qurl_new_from_utf8(p, C.size_t(len(s)), C.int(mode))))
	runtime.KeepAlive(s)
	return q
}

// NewQUrl is a convenience dispatcher for the most common constructors:
//   - NewQUrl()             → QUrl()
//   - NewQUrl(*QUrl)        → QUrl(const QUrl&)  (clone)
//   - NewQUrl(string, m)    → QUrl(const QString&, ParsingMode) [fromUtf8]
//   - NewQUrl(*qtext.QString, m) → QUrl(const QString&, ParsingMode)
func NewQUrl(arg ...any) *QUrl {
	if len(arg) == 0 || arg[0] == nil {
		return wrapQUrl(unsafe.Pointer(C.qurl_new_empty()))
	}

	if len(arg) == 1 {
		q, ok := arg[0].(*QUrl)
		if ok {
			return q.Clone()
		}
	}

	if len(arg) == 2 {
		mode, ok := arg[1].(ParsingMode)
		if ok {
			switch s := arg[0].(type) {
			case string:
				return QUrlFromString(s, mode)
			case *qtext.QString:
				return QUrlFromQString(s, mode)
			case *QUrl:
				return s.Clone()
			}
		}
	}

	panic("NewQUrl invalid argument")
}

// Clone returns a deep copy of the URL.
// Qt: QUrl(const QUrl&) — copy constructor.
func (q *QUrl) Clone() *QUrl {
	return wrapQUrl(unsafe.Pointer(C.qurl_new_copy(q.handle())))
}

// Swap exchanges the contents of q and other.
// Qt: QUrl::swap(QUrl&).
func (q *QUrl) Swap(other *QUrl) {
	C.qurl_swap(q.handle(), other.handle())
	runtime.KeepAlive(other)
}

// ---- setters -------------------------------------------------------------

// SetUrl re-parses the URL from a QString.
// Qt: QUrl::setUrl(const QString&, ParsingMode).
func (q *QUrl) SetUrl(url *qtext.QString, mode ParsingMode) {
	if url == nil {
		q.SetUrlString("", mode)
		return
	}
	q.SetUrlString(url.GoString(), mode)
}

// SetUrlString re-parses the URL from a Go string (UTF-8).
// Qt: QUrl::setUrl(const QString&, ParsingMode) — with QString::fromUtf8.
func (q *QUrl) SetUrlString(s string, mode ParsingMode) {
	var p *C.char
	if len(s) > 0 {
		p = (*C.char)(unsafe.Pointer(unsafe.StringData(s)))
	}
	C.qurl_set_url_utf8(q.handle(), p, C.size_t(len(s)), C.int(mode))
	runtime.KeepAlive(s)
}

// ---- string representation ----------------------------------------------

// URLString returns the URL as a Go string with the given formatting options.
// Qt: QUrl::url(FormattingOptions) → QString (returned here as string).
func (q *QUrl) URLString(options FormattingOptions) string {
	var n C.size_t
	p := C.qurl_url(q.handle(), C.int(options), &n)
	return goStringN(p, n)
}

// URL returns the URL as a QString.
// Qt: QUrl::url(FormattingOptions) → QString.
func (q *QUrl) URL(options FormattingOptions) *qtext.QString {
	return qtext.NewQString(q.URLString(options))
}

// ToString returns the URL as a Go string (alias of URLString).
// Qt: QUrl::toString(FormattingOptions) → QString (returned here as string).
func (q *QUrl) ToString(options FormattingOptions) string {
	var n C.size_t
	p := C.qurl_to_string(q.handle(), C.int(options), &n)
	return goStringN(p, n)
}

// ToQString returns the URL as a QString (alias of URL).
// Qt: QUrl::toString(FormattingOptions) → QString.
func (q *QUrl) ToQString(options FormattingOptions) *qtext.QString {
	return qtext.NewQString(q.ToString(options))
}

// ToDisplayString returns the URL as a display-friendly Go string.
// Qt: QUrl::toDisplayString(FormattingOptions) → QString (returned here as string).
func (q *QUrl) ToDisplayString(options FormattingOptions) string {
	var n C.size_t
	p := C.qurl_to_display_string(q.handle(), C.int(options), &n)
	return goStringN(p, n)
}

// ToDisplayQString returns the URL as a display-friendly QString.
// Qt: QUrl::toDisplayString(FormattingOptions) → QString.
func (q *QUrl) ToDisplayQString(options FormattingOptions) *qtext.QString {
	return qtext.NewQString(q.ToDisplayString(options))
}

// Adjusted returns a new URL with the given formatting options applied.
// Qt: QUrl::adjusted(FormattingOptions) const.
func (q *QUrl) Adjusted(options FormattingOptions) *QUrl {
	return wrapQUrl(unsafe.Pointer(C.qurl_adjusted(q.handle(), C.int(options))))
}

// ToEncodedBytes returns the URL in its encoded (percent-escaped) form as []byte.
// Qt: QUrl::toEncoded(FormattingOptions) → QByteArray (returned here as []byte).
func (q *QUrl) ToEncodedBytes(options FormattingOptions) []byte {
	var n C.size_t
	p := C.qurl_to_encoded(q.handle(), C.int(options), &n)
	if p == nil || n == 0 {
		return nil
	}
	defer C.qurl_free_string(p)
	return C.GoBytes(unsafe.Pointer(p), C.int(n))
}

// ToEncoded returns the URL in its encoded form as QByteArray.
// Qt: QUrl::toEncoded(FormattingOptions) → QByteArray.
func (q *QUrl) ToEncoded(options FormattingOptions) *qtext.QByteArray {
	return qtext.NewQByteArray(string(q.ToEncodedBytes(options)))
}

// ---- validity ------------------------------------------------------------

// IsValid reports whether the URL is syntactically valid.
// Qt: QUrl::isValid() const.
func (q *QUrl) IsValid() bool { return bool(C.qurl_is_valid(q.handle())) }

// ErrorString returns a human-readable description of the last parsing error.
// Qt: QUrl::errorString() const → QString (returned here as string).
func (q *QUrl) ErrorString() string { return goString(C.qurl_error_string(q.handle())) }

// ErrorQString returns the error string as a QString.
// Qt: QUrl::errorString() const → QString.
func (q *QUrl) ErrorQString() *qtext.QString { return qtext.NewQString(q.ErrorString()) }

// IsEmpty reports whether the URL has no content.
// Qt: QUrl::isEmpty() const.
func (q *QUrl) IsEmpty() bool { return bool(C.qurl_is_empty(q.handle())) }

// Clear resets the URL to empty (invalid).
// Qt: QUrl::clear().
func (q *QUrl) Clear() { C.qurl_clear(q.handle()) }

// ---- scheme --------------------------------------------------------------

// SetScheme sets the scheme from a QString.
// Qt: QUrl::setScheme(const QString&).
func (q *QUrl) SetScheme(s *qtext.QString) {
	if s == nil {
		q.SetSchemeString("")
		return
	}
	q.SetSchemeString(s.GoString())
}

// SetSchemeString sets the scheme from a Go string (UTF-8).
// Qt: QUrl::setScheme(const QString&) — with QString::fromUtf8.
func (q *QUrl) SetSchemeString(s string) {
	var p *C.char
	if len(s) > 0 {
		p = (*C.char)(unsafe.Pointer(unsafe.StringData(s)))
	}
	C.qurl_set_scheme_utf8(q.handle(), p, C.size_t(len(s)))
	runtime.KeepAlive(s)
}

// SchemeString returns the scheme as a Go string.
// Qt: QUrl::scheme() const → QString (returned here as string).
func (q *QUrl) SchemeString() string { return goString(C.qurl_scheme(q.handle())) }

// Scheme returns the scheme as a QString.
// Qt: QUrl::scheme() const → QString.
func (q *QUrl) Scheme() *qtext.QString { return qtext.NewQString(q.SchemeString()) }

// ---- authority -----------------------------------------------------------

// SetAuthorityQString sets the authority component from a QString.
// Qt: QUrl::setAuthority(const QString&, ParsingMode).
func (q *QUrl) SetAuthorityQString(s *qtext.QString, mode ParsingMode) {
	if s == nil {
		q.SetAuthorityString("", mode)
		return
	}
	q.SetAuthorityString(s.GoString(), mode)
}

// SetAuthorityString sets the authority component from a Go string.
// Qt: QUrl::setAuthority(const QString&, ParsingMode) — with QString::fromUtf8.
func (q *QUrl) SetAuthorityString(s string, mode ParsingMode) {
	var p *C.char
	if len(s) > 0 {
		p = (*C.char)(unsafe.Pointer(unsafe.StringData(s)))
	}
	C.qurl_set_authority_utf8(q.handle(), p, C.size_t(len(s)), C.int(mode))
	runtime.KeepAlive(s)
}

// Authority returns the authority component as a Go string.
// Qt: QUrl::authority(ComponentFormattingOptions) const → QString (as string).
func (q *QUrl) Authority(options ComponentFormattingOptions) string {
	return goString(C.qurl_authority(q.handle(), C.int(options)))
}

// AuthorityQString returns the authority component as a QString.
// Qt: QUrl::authority(ComponentFormattingOptions) const → QString.
func (q *QUrl) AuthorityQString(options ComponentFormattingOptions) *qtext.QString {
	return qtext.NewQString(q.Authority(options))
}

// ---- user info / user name / password -----------------------------------

// SetUserInfoString sets the user-info component from a Go string.
// Qt: QUrl::setUserInfo(const QString&, ParsingMode) — with QString::fromUtf8.
//
// Note: Qt permits only TolerantMode or StrictMode here; DecodedMode is
// ignored (with a qWarning).
func (q *QUrl) SetUserInfoString(s string, mode ParsingMode) {
	var p *C.char
	if len(s) > 0 {
		p = (*C.char)(unsafe.Pointer(unsafe.StringData(s)))
	}
	C.qurl_set_user_info_utf8(q.handle(), p, C.size_t(len(s)), C.int(mode))
	runtime.KeepAlive(s)
}

// UserInfo returns the user-info component as a Go string.
// Qt: QUrl::userInfo(ComponentFormattingOptions) const → QString (as string).
func (q *QUrl) UserInfo(options ComponentFormattingOptions) string {
	return goString(C.qurl_user_info(q.handle(), C.int(options)))
}

// UserInfoQString returns the user-info component as a QString.
// Qt: QUrl::userInfo(ComponentFormattingOptions) const → QString.
func (q *QUrl) UserInfoQString(options ComponentFormattingOptions) *qtext.QString {
	return qtext.NewQString(q.UserInfo(options))
}

// SetUserNameString sets the user-name component from a Go string.
// Qt: QUrl::setUserName(const QString&, ParsingMode) — with QString::fromUtf8.
func (q *QUrl) SetUserNameString(s string, mode ParsingMode) {
	var p *C.char
	if len(s) > 0 {
		p = (*C.char)(unsafe.Pointer(unsafe.StringData(s)))
	}
	C.qurl_set_user_name_utf8(q.handle(), p, C.size_t(len(s)), C.int(mode))
	runtime.KeepAlive(s)
}

// UserName returns the user-name component as a Go string.
// Qt: QUrl::userName(ComponentFormattingOptions) const → QString (as string).
func (q *QUrl) UserName(options ComponentFormattingOptions) string {
	return goString(C.qurl_user_name(q.handle(), C.int(options)))
}

// UserNameQString returns the user-name component as a QString.
// Qt: QUrl::userName(ComponentFormattingOptions) const → QString.
func (q *QUrl) UserNameQString(options ComponentFormattingOptions) *qtext.QString {
	return qtext.NewQString(q.UserName(options))
}

// SetPasswordString sets the password component from a Go string.
// Qt: QUrl::setPassword(const QString&, ParsingMode) — with QString::fromUtf8.
func (q *QUrl) SetPasswordString(s string, mode ParsingMode) {
	var p *C.char
	if len(s) > 0 {
		p = (*C.char)(unsafe.Pointer(unsafe.StringData(s)))
	}
	C.qurl_set_password_utf8(q.handle(), p, C.size_t(len(s)), C.int(mode))
	runtime.KeepAlive(s)
}

// Password returns the password component as a Go string.
// Qt: QUrl::password(ComponentFormattingOptions) const → QString (as string).
func (q *QUrl) Password(options ComponentFormattingOptions) string {
	return goString(C.qurl_password(q.handle(), C.int(options)))
}

// PasswordQString returns the password component as a QString.
// Qt: QUrl::password(ComponentFormattingOptions) const → QString.
func (q *QUrl) PasswordQString(options ComponentFormattingOptions) *qtext.QString {
	return qtext.NewQString(q.Password(options))
}

// ---- host / port ---------------------------------------------------------

// SetHostString sets the host component from a Go string.
// Qt: QUrl::setHost(const QString&, ParsingMode) — with QString::fromUtf8.
func (q *QUrl) SetHostString(s string, mode ParsingMode) {
	var p *C.char
	if len(s) > 0 {
		p = (*C.char)(unsafe.Pointer(unsafe.StringData(s)))
	}
	C.qurl_set_host_utf8(q.handle(), p, C.size_t(len(s)), C.int(mode))
	runtime.KeepAlive(s)
}

// Host returns the host component as a Go string.
// Qt: QUrl::host(ComponentFormattingOptions) const → QString (as string).
func (q *QUrl) Host(options ComponentFormattingOptions) string {
	return goString(C.qurl_host(q.handle(), C.int(options)))
}

// HostQString returns the host component as a QString.
// Qt: QUrl::host(ComponentFormattingOptions) const → QString.
func (q *QUrl) HostQString(options ComponentFormattingOptions) *qtext.QString {
	return qtext.NewQString(q.Host(options))
}

// SetPort sets the port (-1 clears it).
// Qt: QUrl::setPort(int).
func (q *QUrl) SetPort(port int) { C.qurl_set_port(q.handle(), C.int(port)) }

// Port returns the port, or defaultPort if not set.
// Qt: QUrl::port(int) const.
func (q *QUrl) Port(defaultPort int) int { return int(C.qurl_port(q.handle(), C.int(defaultPort))) }

// ---- path / filename -----------------------------------------------------

// SetPathString sets the path component from a Go string.
// Qt: QUrl::setPath(const QString&, ParsingMode) — with QString::fromUtf8.
func (q *QUrl) SetPathString(s string, mode ParsingMode) {
	var p *C.char
	if len(s) > 0 {
		p = (*C.char)(unsafe.Pointer(unsafe.StringData(s)))
	}
	C.qurl_set_path_utf8(q.handle(), p, C.size_t(len(s)), C.int(mode))
	runtime.KeepAlive(s)
}

// Path returns the path component as a Go string.
// Qt: QUrl::path(ComponentFormattingOptions) const → QString (as string).
func (q *QUrl) Path(options ComponentFormattingOptions) string {
	return goString(C.qurl_path(q.handle(), C.int(options)))
}

// PathQString returns the path component as a QString.
// Qt: QUrl::path(ComponentFormattingOptions) const → QString.
func (q *QUrl) PathQString(options ComponentFormattingOptions) *qtext.QString {
	return qtext.NewQString(q.Path(options))
}

// FileName returns the last segment of the path as a Go string.
// Qt: QUrl::fileName(ComponentFormattingOptions) const → QString (as string).
func (q *QUrl) FileName(options ComponentFormattingOptions) string {
	return goString(C.qurl_file_name(q.handle(), C.int(options)))
}

// FileNameQString returns the last segment of the path as a QString.
// Qt: QUrl::fileName(ComponentFormattingOptions) const → QString.
func (q *QUrl) FileNameQString(options ComponentFormattingOptions) *qtext.QString {
	return qtext.NewQString(q.FileName(options))
}

// ---- query / fragment ----------------------------------------------------

// HasQuery reports whether a query component is present.
// Qt: QUrl::hasQuery() const.
func (q *QUrl) HasQuery() bool { return bool(C.qurl_has_query(q.handle())) }

// SetQueryString sets the query component from a Go string.
// Qt: QUrl::setQuery(const QString&, ParsingMode) — with QString::fromUtf8.
func (q *QUrl) SetQueryString(s string, mode ParsingMode) {
	var p *C.char
	if len(s) > 0 {
		p = (*C.char)(unsafe.Pointer(unsafe.StringData(s)))
	}
	C.qurl_set_query_utf8(q.handle(), p, C.size_t(len(s)), C.int(mode))
	runtime.KeepAlive(s)
}

// SetQueryQUrlQuery sets the query component from a QUrlQuery.
// The QUrlQuery is read-only here — its lifetime is not transferred; it can
// be reused or closed immediately after the call.
// Qt: QUrl::setQuery(const QUrlQuery&).
//
// Passing nil clears the query (equivalent to setting an empty QUrlQuery).
func (q *QUrl) SetQueryQUrlQuery(query *QUrlQuery) {
	if _, ok := query.tryHandle(); !ok {
		empty := NewQUrlQuery()
		defer empty.Close()
		C.qurl_set_query_qurlquery(q.handle(), empty.handle())
		return
	}
	C.qurl_set_query_qurlquery(q.handle(), query.handle())
	runtime.KeepAlive(query)
}

// Query returns the query component as a Go string.
// Qt: QUrl::query(ComponentFormattingOptions) const → QString (as string).
func (q *QUrl) Query(options ComponentFormattingOptions) string {
	return goString(C.qurl_query(q.handle(), C.int(options)))
}

// QueryQString returns the query component as a QString.
// Qt: QUrl::query(ComponentFormattingOptions) const → QString.
func (q *QUrl) QueryQString(options ComponentFormattingOptions) *qtext.QString {
	return qtext.NewQString(q.Query(options))
}

// HasFragment reports whether a fragment component is present.
// Qt: QUrl::hasFragment() const.
func (q *QUrl) HasFragment() bool { return bool(C.qurl_has_fragment(q.handle())) }

// SetFragmentString sets the fragment component from a Go string.
// Qt: QUrl::setFragment(const QString&, ParsingMode) — with QString::fromUtf8.
func (q *QUrl) SetFragmentString(s string, mode ParsingMode) {
	var p *C.char
	if len(s) > 0 {
		p = (*C.char)(unsafe.Pointer(unsafe.StringData(s)))
	}
	C.qurl_set_fragment_utf8(q.handle(), p, C.size_t(len(s)), C.int(mode))
	runtime.KeepAlive(s)
}

// Fragment returns the fragment component as a Go string.
// Qt: QUrl::fragment(ComponentFormattingOptions) const → QString (as string).
func (q *QUrl) Fragment(options ComponentFormattingOptions) string {
	return goString(C.qurl_fragment(q.handle(), C.int(options)))
}

// FragmentQString returns the fragment component as a QString.
// Qt: QUrl::fragment(ComponentFormattingOptions) const → QString.
func (q *QUrl) FragmentQString(options ComponentFormattingOptions) *qtext.QString {
	return qtext.NewQString(q.Fragment(options))
}

// ---- relations -----------------------------------------------------------

// Resolved returns the absolute URL obtained by resolving relative against q.
// Qt: QUrl::resolved(const QUrl&) const.
func (q *QUrl) Resolved(relative *QUrl) *QUrl {
	if relative == nil {
		return q.Clone()
	}
	return wrapQUrl(unsafe.Pointer(C.qurl_resolved(q.handle(), relative.handle())))
}

// IsRelative reports whether the URL is relative (no scheme).
// Qt: QUrl::isRelative() const.
func (q *QUrl) IsRelative() bool { return bool(C.qurl_is_relative(q.handle())) }

// IsParentOf reports whether q is a parent of url.
// Qt: QUrl::isParentOf(const QUrl&) const.
func (q *QUrl) IsParentOf(url *QUrl) bool {
	if url == nil {
		return false
	}
	return bool(C.qurl_is_parent_of(q.handle(), url.handle()))
}

// Matches reports whether q matches other under the given formatting options.
// Qt: QUrl::matches(const QUrl&, FormattingOptions) const.
func (q *QUrl) Matches(other *QUrl, options FormattingOptions) bool {
	if other == nil {
		return false
	}
	return bool(C.qurl_matches(q.handle(), other.handle(), C.int(options)))
}

// ---- local files ---------------------------------------------------------

// IsLocalFile reports whether the URL refers to a local file.
// Qt: QUrl::isLocalFile() const.
func (q *QUrl) IsLocalFile() bool { return bool(C.qurl_is_local_file(q.handle())) }

// ToLocalFile returns the local file path as a Go string.
// Qt: QUrl::toLocalFile() const → QString (returned here as string).
func (q *QUrl) ToLocalFile() string { return goString(C.qurl_to_local_file(q.handle())) }

// ToLocalFileQString returns the local file path as a QString.
// Qt: QUrl::toLocalFile() const → QString.
func (q *QUrl) ToLocalFileQString() *qtext.QString {
	return qtext.NewQString(q.ToLocalFile())
}

// ---- detach --------------------------------------------------------------

// Detach detaches the URL from any implicitly shared copies.
// Qt: QUrl::detach().
func (q *QUrl) Detach() { C.qurl_detach(q.handle()) }

// IsDetached reports whether the URL is detached from any shared copies.
// Qt: QUrl::isDetached() const.
func (q *QUrl) IsDetached() bool { return bool(C.qurl_is_detached(q.handle())) }

// ---- static parsers -----------------------------------------------------

// FromEncoded parses a URL from its encoded byte representation.
// Qt: QUrl::fromEncoded(QByteArrayView, ParsingMode).
func FromEncoded(b []byte, mode ParsingMode) *QUrl {
	var p *C.uint8_t
	if len(b) > 0 {
		p = (*C.uint8_t)(unsafe.Pointer(&b[0]))
	}
	q := wrapQUrl(unsafe.Pointer(C.qurl_from_encoded(p, C.size_t(len(b)), C.int(mode))))
	runtime.KeepAlive(b)
	return q
}

// FromLocalFile parses a URL from a local file path (Go string).
// Qt: QUrl::fromLocalFile(const QString&) — with QString::fromUtf8.
func FromLocalFile(path string) *QUrl {
	var p *C.char
	if len(path) > 0 {
		p = (*C.char)(unsafe.Pointer(unsafe.StringData(path)))
	}
	q := wrapQUrl(unsafe.Pointer(C.qurl_from_local_file_utf8(p, C.size_t(len(path)))))
	runtime.KeepAlive(path)
	return q
}

// FromLocalFileQString is the QString overload of FromLocalFile.
// Qt: QUrl::fromLocalFile(const QString&).
func FromLocalFileQString(path *qtext.QString) *QUrl {
	if path == nil {
		return FromLocalFile("")
	}
	return FromLocalFile(path.GoString())
}

// FromUserInputString builds a URL from user input (Go string).
// Qt: QUrl::fromUserInput(const QString&, const QString&, UserInputResolutionOptions).
func FromUserInputString(input, workingDir string, options UserInputResolutionOptions) *QUrl {
	var ip, wp *C.char
	if len(input) > 0 {
		ip = (*C.char)(unsafe.Pointer(unsafe.StringData(input)))
	}
	if len(workingDir) > 0 {
		wp = (*C.char)(unsafe.Pointer(unsafe.StringData(workingDir)))
	}
	q := wrapQUrl(unsafe.Pointer(C.qurl_from_user_input_utf8(
		ip, C.size_t(len(input)), wp, C.size_t(len(workingDir)), C.int(options))))
	runtime.KeepAlive(input)
	runtime.KeepAlive(workingDir)
	return q
}

// FromUserInputQString is the QString overload of FromUserInputString.
// Qt: QUrl::fromUserInput(const QString&, const QString&, UserInputResolutionOptions).
func FromUserInputQString(input, workingDir *qtext.QString, options UserInputResolutionOptions) *QUrl {
	is, ws := "", ""
	if input != nil {
		is = input.GoString()
	}
	if workingDir != nil {
		ws = workingDir.GoString()
	}
	return FromUserInputString(is, ws, options)
}

// FromPercentEncodingString decodes percent-encoded bytes into a Go string.
// Qt: QUrl::fromPercentEncoding(const QByteArray&) → QString (as string).
func FromPercentEncodingString(b []byte) string {
	var n C.size_t
	p := C.qurl_from_percent_encoding(uint8Ptr(b), C.size_t(len(b)), &n)
	return goStringN(p, n)
}

// FromPercentEncoding decodes percent-encoded bytes into a QString.
// Qt: QUrl::fromPercentEncoding(const QByteArray&) → QString.
func FromPercentEncoding(b []byte) *qtext.QString {
	return qtext.NewQString(FromPercentEncodingString(b))
}

// ToPercentEncodingBytes percent-encodes a Go string into []byte.
// Qt: QUrl::toPercentEncoding(const QString&, const QByteArray&, const QByteArray&) → QByteArray.
func ToPercentEncodingBytes(s string, exclude, include []byte) []byte {
	var sp *C.char
	if len(s) > 0 {
		sp = (*C.char)(unsafe.Pointer(unsafe.StringData(s)))
	}
	var n C.size_t
	p := C.qurl_to_percent_encoding(sp, C.size_t(len(s)),
		uint8Ptr(exclude), C.size_t(len(exclude)),
		uint8Ptr(include), C.size_t(len(include)), &n)
	runtime.KeepAlive(s)
	runtime.KeepAlive(exclude)
	runtime.KeepAlive(include)
	if p == nil || n == 0 {
		return nil
	}
	defer C.qurl_free_string(p)
	return C.GoBytes(unsafe.Pointer(p), C.int(n))
}

// ToPercentEncoding percent-encodes a QString into QByteArray.
// Qt: QUrl::toPercentEncoding(const QString&, const QByteArray&, const QByteArray&) → QByteArray.
func ToPercentEncoding(s *qtext.QString, exclude, include []byte) *qtext.QByteArray {
	str := ""
	if s != nil {
		str = s.GoString()
	}
	return qtext.NewQByteArray(string(ToPercentEncodingBytes(str, exclude, include)))
}

// FromAceString decodes a Punycode-encoded domain into a Go string.
// Qt: QUrl::fromAce(const QByteArray&, AceProcessingOptions) → QString (as string).
func FromAceString(b []byte) string { return goString(C.qurl_from_ace(uint8Ptr(b), C.size_t(len(b)))) }

// FromAce decodes a Punycode-encoded domain into a QString.
// Qt: QUrl::fromAce(const QByteArray&, AceProcessingOptions) → QString.
func FromAce(b []byte) *qtext.QString { return qtext.NewQString(FromAceString(b)) }

// ToAceBytes encodes a Go string domain into Punycode []byte.
// Qt: QUrl::toAce(const QString&, AceProcessingOptions) → QByteArray.
func ToAceBytes(s string) []byte {
	var sp *C.char
	if len(s) > 0 {
		sp = (*C.char)(unsafe.Pointer(unsafe.StringData(s)))
	}
	var n C.size_t
	p := C.qurl_to_ace(sp, C.size_t(len(s)), &n)
	runtime.KeepAlive(s)
	if p == nil || n == 0 {
		return nil
	}
	defer C.qurl_free_string(p)
	return C.GoBytes(unsafe.Pointer(p), C.int(n))
}

// ToAce encodes a QString domain into Punycode QByteArray.
// Qt: QUrl::toAce(const QString&, AceProcessingOptions) → QByteArray.
func ToAce(s *qtext.QString) *qtext.QByteArray {
	str := ""
	if s != nil {
		str = s.GoString()
	}
	return qtext.NewQByteArray(string(ToAceBytes(str)))
}

// ---- list conversions ---------------------------------------------------

// ToStringListStrings converts a slice of QUrl into a slice of Go strings.
// Qt: QUrl::toStringList(const QList<QUrl>&, FormattingOptions) → QStringList (as []string).
//
// Note: the C ABI joins strings with '\n'; a URL containing a literal newline
// (which QUrl normally percent-encodes as %0A) would be split incorrectly.
func ToStringListStrings(urls []*QUrl, options FormattingOptions) []string {
	if len(urls) == 0 {
		return nil
	}
	cUrls := make([]*C.QUrl, len(urls))
	for i, u := range urls {
		if up, ok := u.tryHandle(); ok {
			cUrls[i] = up
		}
	}
	var n C.size_t
	p := C.qurl_to_string_list((**C.QUrl)(unsafe.Pointer(&cUrls[0])), C.size_t(len(cUrls)),
		C.int(options), &n)
	runtime.KeepAlive(urls)
	if p == nil || n == 0 {
		return nil
	}
	defer C.qurl_free_string(p)
	return strings.Split(C.GoStringN(p, C.int(n)), "\n")
}

// ToStringList converts a slice of QUrl into a QStringList.
// Qt: QUrl::toStringList(const QList<QUrl>&, FormattingOptions) → QStringList.
func ToStringList(urls []*QUrl, options FormattingOptions) *qtext.QStringList {
	strs := ToStringListStrings(urls, options)
	sl := qtext.NewQStringList()
	for _, s := range strs {
		sl.AppendGoString(s)
	}
	return sl
}

// FromStringListStrings parses a slice of Go strings into a slice of *QUrl.
// Qt: QUrl::fromStringList(const QStringList&, ParsingMode) → QList<QUrl>.
//
// The Go strings are concatenated into a single flat byte buffer, and an
// offsets array (len == len(strs)+1) records their boundaries. Both are
// passed to C as plain Go pointers (no Go pointers embedded), satisfying
// cgo's "Go pointer to unpinned Go pointer" rule.
func FromStringListStrings(strs []string, mode ParsingMode) []*QUrl {
	if len(strs) == 0 {
		return nil
	}

	total := 0
	for _, s := range strs {
		total += len(s)
	}
	var data []byte
	if total > 0 {
		data = make([]byte, 0, total)
		for _, s := range strs {
			data = append(data, s...)
		}
	}
	offsets := make([]C.size_t, len(strs)+1)
	offsets[0] = 0
	for i, s := range strs {
		offsets[i+1] = offsets[i] + C.size_t(len(s))
	}

	var dataPtr *C.char
	if len(data) > 0 {
		dataPtr = (*C.char)(unsafe.Pointer(&data[0]))
	}
	var n C.size_t
	arr := C.qurl_from_string_list(
		dataPtr,
		(*C.size_t)(unsafe.Pointer(&offsets[0])),
		C.size_t(len(strs)),
		C.int(mode),
		&n,
	)
	runtime.KeepAlive(data)
	runtime.KeepAlive(offsets)
	if arr == nil {
		return nil
	}
	defer C.qurl_free_string_array(arr)
	count := int(n)
	if count == 0 {
		return nil
	}
	ptrs := unsafe.Slice(arr, count)
	result := make([]*QUrl, count)
	for i := 0; i < count; i++ {
		result[i] = wrapQUrl(unsafe.Pointer(ptrs[i]))
	}
	return result
}

// FromStringList parses a QStringList into a slice of *QUrl.
// Qt: QUrl::fromStringList(const QStringList&, ParsingMode) → QList<QUrl>.
func FromStringList(sl *qtext.QStringList, mode ParsingMode) []*QUrl {
	if sl == nil || sl.IsEmpty() {
		return nil
	}
	strs := make([]string, sl.Size())
	for i := 0; i < sl.Size(); i++ {
		e := sl.At(i)
		strs[i] = e.String()
		e.Close()
	}
	return FromStringListStrings(strs, mode)
}

// ---- helpers -------------------------------------------------------------

func goString(p *C.char) string {
	if p == nil {
		return ""
	}
	defer C.qurl_free_string(p)
	return C.GoString(p)
}

func goStringN(p *C.char, n C.size_t) string {
	if p == nil || n == 0 {
		return ""
	}
	defer C.qurl_free_string(p)
	return C.GoStringN(p, C.int(n))
}

func uint8Ptr(b []byte) *C.uint8_t {
	if len(b) == 0 {
		return nil
	}
	return (*C.uint8_t)(unsafe.Pointer(&b[0]))
}
