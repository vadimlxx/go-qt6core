// Package qt6core is an internal helper package for the go-qt6core library.
//
// It provides exactly two things:
//
//   - check_qt_version.h — compile-time minimum Qt version check
//     (#included by every *_c.cpp in package qtext).
//   - Version() — the runtime Qt version query, for diagnostics.
//
// Before the qtext unification, this package also hosted the cross-package
// QString/QByteArray bridge (UTF-8/UTF-16/Latin-1 conversions). After QString
// and QByteArray became part of package qtext, that bridge is no longer
// needed: qtext has its own direct C bindings (qstring_c.*, qbytearray_c.*),
// and no package cycles exist to work around.
//
// NOT part of the public API. External packages should use qtext directly.
//
// Requirements: Qt >= 6.11 (checked at compile time), C++17 compiler,
// pkg-config Qt6Core.
package qt6core
