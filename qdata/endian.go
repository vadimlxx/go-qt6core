package qdata

// Endian mirrors QSysInfo::Endian (Qt 6.11).
//
// ByteOrder is exported from qsysinfo.go as a package-level variable
// initialised to whichever of BigEndian / LittleEndian matches the
// current build's byte order.
type Endian int

const (
	// BigEndian is big-endian byte order.
	// Qt: QSysInfo::BigEndian.
	BigEndian Endian = 0

	// LittleEndian is little-endian byte order.
	// Qt: QSysInfo::LittleEndian.
	LittleEndian Endian = 1
)