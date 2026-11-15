package qdata

// Sizes mirrors QSysInfo::Sizes (Qt 6.11).
//
// The single named member, WordSize, is exported from qsysinfo.go as a
// package-level variable because its value is platform-dependent
// (sizeof(void*) << 3) and therefore cannot be a Go const.
type Sizes int