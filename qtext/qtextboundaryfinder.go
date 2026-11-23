package qtext

/*
 Краткое описание реализации:

   QTextBoundaryFinder — обёртка над одноимённым C++-классом через
   непрозрачный указатель (C ABI). Go-структура содержит только этот
   указатель; владение обеспечивается финализатором и явным Close().

   В Qt 6.11 QTextBoundaryFinder — настоящий класс (не alias), поэтому
   forward-декларация class QTextBoundaryFinder; в qtextboundaryfinder_c.h
   корректна. Это отличает его от QStringList (там — using-alias).

   enum-ы BoundaryType и BoundaryReason вынесены в отдельные файлы
   (qtextboundaryfinder_boundarytype.go и qtextboundaryfinder_boundaryreason.go),
   BoundaryReasons — alias на BoundaryReason (эргономичный эквивалент
   Q_DECLARE_FLAGS).

   Конструкторы принимают либо *QString (сохраняет семантическое различие
   null/empty), либо QStringView (zero-copy, вызывает QStringView-перегрузку
   QTextBoundaryFinder(BoundaryType, QStringView, ...)). Буфер/буфер-размер
   не пробрасываем — Qt сам аллоцирует внутренний буфер.

   У QTextBoundaryFinder нет protected-секций, нет статических функций,
   нет free-функций. swap() — member, поэтому реализован как метод Swap().

 Что сознательно не вошло (пока):

   - Конструктор с явным buffer / bufferSize:
     QTextBoundaryFinder(BoundaryType, const QChar*, qsizetype,
                         unsigned char*, qsizetype) — это оптимизация
     для многократного использования без реаллокации; пользователь Go
     всё равно не сможет разделить буфер с Qt без unsafe-хрупкости.
   - Move-конструктор и move-assignment — в Go перенос через указатель.
   - Q_DECLARE_SHARED-семантика (copy-on-write) — в Go копия через Clone().
   - Java/STL-style итераторы — не применимы в Go; используйте
     ToNextBoundary/ToPreviousBoundary.
*/

/*
#cgo pkg-config: Qt6Core
#cgo CXXFLAGS: -std=c++17 -fPIC -Wno-sfinae-incomplete
#cgo LDFLAGS: -lstdc++
#include <stdint.h>
#include <stdbool.h>
#include <stddef.h>
#include "qtextboundaryfinder_c.h"
*/
import "C"

import (
	"runtime"
	"sync/atomic"
	"unsafe"
)

// QTextBoundaryFinder wraps a Qt6 QTextBoundaryFinder object.

type qtextbfState struct {
	p atomic.Pointer[byte]
}

func (s *qtextbfState) cpp() *C.QTextBoundaryFinder {
	if bp := s.p.Load(); bp != nil {
		return (*C.QTextBoundaryFinder)(unsafe.Pointer(bp))
	}
	return nil
}

func (s *qtextbfState) claim() *C.QTextBoundaryFinder {
	if bp := s.p.Swap(nil); bp != nil {
		return (*C.QTextBoundaryFinder)(unsafe.Pointer(bp))
	}
	return nil
}

// Each *QTextBoundaryFinder owns a heap-allocated C++ QTextBoundaryFinder;
// the finalizer (and explicit Close) frees it.
type QTextBoundaryFinder struct {
	// p *C.QTextBoundaryFinder
	state   *qtextbfState
	cleanup runtime.Cleanup
}

func wrapQTextBoundaryFinder(p unsafe.Pointer) *QTextBoundaryFinder {
	if p == nil {
		return nil
	}

	st := &qtextbfState{}
	st.p.Store((*byte)(p))
	q := &QTextBoundaryFinder{state: st}
	q.cleanup = runtime.AddCleanup(q, func(st *qtextbfState) {
		if cpp := st.claim(); cpp != nil {
			C.qtextboundaryfinder_delete(cpp)
		}
	}, st)

	return q
}

func (q *QTextBoundaryFinder) handle() *C.QTextBoundaryFinder {
	if cpp, ok := q.tryHandle(); ok {
		return cpp
	}
	panic("qtextboundaryfinder: use after Close (or nil receiver)")
}

func (q *QTextBoundaryFinder) tryHandle() (*C.QTextBoundaryFinder, bool) {
	if q == nil || q.state == nil {
		return nil, false
	}
	if cpp := q.state.cpp(); cpp != nil {
		return cpp, true
	}
	return nil, false
}

// Close releases the underlying C++ object immediately and disables the finalizer.
// Qt: QTextBoundaryFinder::~QTextBoundaryFinder().
func (q *QTextBoundaryFinder) Close() {
	if q == nil || q.state == nil {
		return
	}

	q.cleanup.Stop()
	if cpp := q.state.claim(); cpp != nil {
		C.qtextboundaryfinder_delete(cpp)
	}
}

// ---- construction --------------------------------------------------------

// NewQTextBoundaryFinder returns a default-constructed (invalid) finder.
// Qt: QTextBoundaryFinder() — default constructor; isValid() returns false.
func NewQTextBoundaryFinder(arg ...any) *QTextBoundaryFinder {
	if len(arg) == 0 || arg[0] == nil {
		return wrapQTextBoundaryFinder(unsafe.Pointer(C.qtextboundaryfinder_new_empty()))
	}

	if len(arg) == 1 {
		b, ok := arg[0].(*QTextBoundaryFinder)
		if ok {
			return b.Clone()
		}
	}

	if len(arg) == 2 {
		t, ok := arg[0].(BoundaryType)

		if ok {
			switch s := arg[1].(type) {
			case string:
				return NewQTextBoundaryFinderFromString(t, s)
			case *QString:
				return NewQTextBoundaryFinderFromQString(t, s)
			case QStringView:
				return NewQTextBoundaryFinderFromView(t, s)
			}
		}
	}

	panic("NewQTextBoundaryFinder invalid argument")
}

// NewQTextBoundaryFinderFromQString creates a finder for the given QString.
// Qt: QTextBoundaryFinder(BoundaryType, const QString&).
func NewQTextBoundaryFinderFromQString(t BoundaryType, s *QString) *QTextBoundaryFinder {
	if s == nil || s.handle() == nil {
		return NewQTextBoundaryFinder()
	}
	q := wrapQTextBoundaryFinder(unsafe.Pointer(C.qtextboundaryfinder_new_from_qstring(C.int(t), s.handle())))
	runtime.KeepAlive(s)
	return q
}

// NewQTextBoundaryFinderFromString creates a finder for a Go string (UTF-8).
// Qt: QTextBoundaryFinder(BoundaryType, const QString&) — with QString::fromUtf8.
func NewQTextBoundaryFinderFromString(t BoundaryType, s string) *QTextBoundaryFinder {
	qs := NewQStringFromString(s)
	defer qs.Close()
	return NewQTextBoundaryFinderFromQString(t, qs)
}

// NewQTextBoundaryFinderFromView creates a finder over the given QStringView.
// The view's underlying buffer is copied by Qt (see note below).
// Qt: QTextBoundaryFinder(BoundaryType, QStringView, unsigned char*, qsizetype).
//
// Note: Qt 6.11's QStringView constructor still stores a QString internally
// (QTextBoundaryFinder::s) plus a QStringView into it, so the data is copied
// regardless of the view's origin. The Go-side keepAlive is therefore not
// needed beyond the call, but we still call KeepAlive defensively.
func NewQTextBoundaryFinderFromView(t BoundaryType, v QStringView) *QTextBoundaryFinder {
	q := wrapQTextBoundaryFinder(unsafe.Pointer(C.qtextboundaryfinder_new_from_view(
		C.int(t), (*C.uint16_t)(v.Utf16Unsafe()), C.ptrdiff_t(v.Size()))))
	runtime.KeepAlive(v)
	return q
}

// Clone returns a deep copy of the finder (independent internal state).
// Qt: QTextBoundaryFinder(const QTextBoundaryFinder&) — copy constructor.
func (q *QTextBoundaryFinder) Clone() *QTextBoundaryFinder {
	return wrapQTextBoundaryFinder(unsafe.Pointer(C.qtextboundaryfinder_new_copy(q.handle())))
}

// Swap exchanges the contents of q and other.
// Qt: QTextBoundaryFinder::swap(QTextBoundaryFinder&).
func (q *QTextBoundaryFinder) Swap(other *QTextBoundaryFinder) {
	C.qtextboundaryfinder_swap(q.handle(), other.handle())
	runtime.KeepAlive(other)
}

// ---- accessors -----------------------------------------------------------

// IsValid reports whether the finder has been properly initialized with a string.
// Qt: QTextBoundaryFinder::isValid() const.
func (q *QTextBoundaryFinder) IsValid() bool {
	return bool(C.qtextboundaryfinder_is_valid(q.handle()))
}

// Type returns the boundary type this finder was constructed for.
// Qt: QTextBoundaryFinder::type() const.
func (q *QTextBoundaryFinder) Type() BoundaryType {
	return BoundaryType(C.qtextboundaryfinder_type(q.handle()))
}

// String returns a copy of the string this finder operates on.
// Qt: QTextBoundaryFinder::string() const.
func (q *QTextBoundaryFinder) String() *QString {
	return wrapQString(unsafe.Pointer(C.qtextboundaryfinder_string(q.handle())))
}

// ---- navigation ---------------------------------------------------------

// ToStart moves the finder to position 0.
// Qt: QTextBoundaryFinder::toStart().
func (q *QTextBoundaryFinder) ToStart() {
	C.qtextboundaryfinder_to_start(q.handle())
}

// ToEnd moves the finder to the end of the string.
// Qt: QTextBoundaryFinder::toEnd().
func (q *QTextBoundaryFinder) ToEnd() {
	C.qtextboundaryfinder_to_end(q.handle())
}

// Position returns the current position in the string.
// Qt: QTextBoundaryFinder::position() const.
func (q *QTextBoundaryFinder) Position() int {
	return int(C.qtextboundaryfinder_position(q.handle()))
}

// SetPosition sets the current position. Panics if pos is negative.
// Qt: QTextBoundaryFinder::setPosition(qsizetype).
func (q *QTextBoundaryFinder) SetPosition(pos int) {
	if pos < 0 {
		panic("qtextboundaryfinder: SetPosition negative position")
	}
	C.qtextboundaryfinder_set_position(q.handle(), C.ptrdiff_t(pos))
}

// ToNextBoundary advances to the next boundary and returns the new position,
// or -1 if there is no next boundary.
// Qt: QTextBoundaryFinder::toNextBoundary().
func (q *QTextBoundaryFinder) ToNextBoundary() int {
	return int(C.qtextboundaryfinder_to_next_boundary(q.handle()))
}

// ToPreviousBoundary retreats to the previous boundary and returns the new
// position, or -1 if there is no previous boundary.
// Qt: QTextBoundaryFinder::toPreviousBoundary().
func (q *QTextBoundaryFinder) ToPreviousBoundary() int {
	return int(C.qtextboundaryfinder_to_previous_boundary(q.handle()))
}

// ---- predicates ----------------------------------------------------------

// IsAtBoundary reports whether the finder is currently on a boundary.
// Qt: QTextBoundaryFinder::isAtBoundary() const.
func (q *QTextBoundaryFinder) IsAtBoundary() bool {
	return bool(C.qtextboundaryfinder_is_at_boundary(q.handle()))
}

// BoundaryReasons returns the bitmask of reasons why the current position
// is a boundary. NotAtBoundary (0) means "not on a boundary".
// Qt: QTextBoundaryFinder::boundaryReasons() const.
func (q *QTextBoundaryFinder) BoundaryReasons() BoundaryReasons {
	return BoundaryReasons(C.qtextboundaryfinder_boundary_reasons(q.handle()))
}
