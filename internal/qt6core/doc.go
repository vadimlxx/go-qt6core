// Package qt6core provides internal, cross-package C++ bridge helpers
// for the go-qt6core library.
//
// It is NOT part of the public API. It exists solely so that sibling
// packages (qstring, qstringlist, qbytearray, qjson, qsettings, ...)
// can share one implementation of QString ↔ UTF-8/UTF-16/Latin-1 and
// QByteArray ↔ QString conversions, without a circular dependency
// between qstring and qstringlist.
//
// All exported functions exchange objects via unsafe.Pointer:
//   - Functions named QStringFrom*/QByteArrayFrom* return a newly
//     allocated C++ object; the caller owns it and must free it via
//     the corresponding *Delete function (typically from a finalizer
//     in the owning public package).
//   - Functions named QStringTo* return freshly allocated Go copies;
//     no ownership transfer of C++ objects takes place.
//
// Requirements: Qt >= 6.11 (checked at compile time via
// check_qt_version.h), C++17 compiler, pkg-config Qt6Core.
package qt6core