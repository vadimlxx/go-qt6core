// Package qdata provides Go bindings for Qt 6.11 QtCore types that are
// shared across multiple packages of go-qt6core (Qt:: namespace enums and
// small value types) without depending on qtext.
//
// This file documents the QSysInfo binding.
//
// QSysInfo is a purely static class that reports information about the
// platform Qt was built for and is currently running on. It has two nested
// enums (Sizes, Endian) and eleven static getters.
//
// # API surface
//
// Enum types:
//   - Sizes  (qsysinfo_sizes.go)  — Sizes enum; the single named value
//                                   (WordSize) is exposed in qsysinfo.go.
//   - Endian (qsysinfo_endian.go) — BigEndian, LittleEndian; the
//                                   platform-dependent ByteOrder is exposed
//                                   in qsysinfo.go.
//
// Package-level variables:
//   - WordSize  Sizes  — QSysInfo::WordSize (sizeof(void*) << 3)
//   - ByteOrder Endian — QSysInfo::ByteOrder (BigEndian or LittleEndian)
//
// Static getters (Qt: QSysInfo::*):
//   - QSysInfoBuildCpuArchitecture()   ← buildCpuArchitecture()   → string
//   - QSysInfoCurrentCpuArchitecture() ← currentCpuArchitecture() → string
//   - QSysInfoBuildAbi()               ← buildAbi()               → string
//   - QSysInfoKernelType()             ← kernelType()             → string
//   - QSysInfoKernelVersion()          ← kernelVersion()          → string
//   - QSysInfoProductType()            ← productType()            → string
//   - QSysInfoProductVersion()         ← productVersion()         → string
//   - QSysInfoPrettyProductName()      ← prettyProductName()      → string
//   - QSysInfoMachineHostName()        ← machineHostName()        → string
//   - QSysInfoMachineUniqueId()        ← machineUniqueId()        → []byte
//   - QSysInfoBootUniqueId()           ← bootUniqueId()           → []byte
//
// # Return types: no qtext dependency
//
// Qt's QSysInfo returns QString and QByteArray. To avoid a qdata → qtext
// import cycle, this binding converts both to plain Go types inside the C ABI
// helper:
//
//   - QString-returning getters produce a NUL-terminated UTF-8 buffer on the
//     C++ side (malloc'd); the Go side copies it into a Go string via
//     C.GoString and then frees the C buffer via qsysinfo_free_string.
//
//   - QByteArray-returning getters produce a raw byte buffer plus an explicit
//     length on the C++ side (malloc'd); the Go side copies exactly n bytes
//     into a []byte via C.GoBytes — preserving any embedded NUL bytes — and
//     then frees the C buffer. A null or empty QByteArray maps to a nil
//     []byte, mirroring Qt's own null-vs-empty distinction as closely as
//     []byte allows (both cases have len(x) == 0).
//
// Not implemented
//
//   - Nothing from the QSysInfo public API in qsysinfo.h is missing: the
//     class has no protected members, no non-static members, no free
//     functions, and no operators.
//   - Some platform-specific data (machineUniqueId / bootUniqueId) may
//     legitimately return nil on certain platforms — that mirrors Qt's own
//     behaviour, not a bug in this binding.
//
// This file is part of package qdata of github.com/vadimlxx/go-qt6core.
// Requires Qt >= 6.11 and a C++17 compiler; links via pkg-config Qt6Core.
package qdata
