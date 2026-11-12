package qdata

// DataSizeFormat mirrors QLocale::DataSizeFormat.
type DataSizeFormat int

const (
	// Single-bit values, for internal use.
	DataSizeBase1000      DataSizeFormat = 1
	DataSizeSIQuantifiers DataSizeFormat = 2

	// Named flag combinations, matching the API-level enum values in Qt.
	DataSizeIecFormat         DataSizeFormat = 0
	DataSizeTraditionalFormat DataSizeFormat = DataSizeSIQuantifiers
	DataSizeSIFormat          DataSizeFormat = DataSizeBase1000 | DataSizeSIQuantifiers
)

// DataSizeFormats is a bitmask of DataSizeFormat flags.
//
// In C++, Q_DECLARE_FLAGS(DataSizeFormats, DataSizeFormat) creates a distinct
// QFlags<DataSizeFormat> type with an implicit constructor taking
// DataSizeFormat. In Go, DataSizeFormats is a type alias for DataSizeFormat —
// so single-flag constants (DataSizeIecFormat, DataSizeSIFormat, …) can be
// passed directly wherever DataSizeFormats is expected, without explicit casts.
type DataSizeFormats = DataSizeFormat
