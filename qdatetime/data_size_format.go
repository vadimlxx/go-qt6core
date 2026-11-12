package qdatetime

// DataSizeFormat mirrors QLocale::DataSizeFormat.
type DataSizeFormat int

const (
	DataSizeBase1000      DataSizeFormat = 1
	DataSizeSIQuantifiers DataSizeFormat = 2

	DataSizeIecFormat         DataSizeFormat = 0
	DataSizeTraditionalFormat DataSizeFormat = DataSizeSIQuantifiers
	DataSizeSIFormat          DataSizeFormat = DataSizeBase1000 | DataSizeSIQuantifiers
)

// DataSizeFormats is a bitmask of DataSizeFormat flags.
type DataSizeFormats int