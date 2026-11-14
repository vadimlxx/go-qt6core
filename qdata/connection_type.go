package qdata

// ConnectionType mirrors Qt::ConnectionType from qnamespace.h (Qt 6.11).
//
// Used by QObject::connect / QObject::disconnect. Unique and SingleShot are
// flags that can be OR-ed with AutoConnection / DirectConnection / etc.
type ConnectionType int

const (
	AutoConnection           ConnectionType = 0
	DirectConnection         ConnectionType = 1
	QueuedConnection         ConnectionType = 2
	BlockingQueuedConnection ConnectionType = 3
	UniqueConnection         ConnectionType = 0x80
	SingleShotConnection     ConnectionType = 0x100
)