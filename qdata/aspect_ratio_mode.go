package qdata

// AspectRatioMode mirrors Qt::AspectRatioMode from qnamespace.h (Qt 6.11).
type AspectRatioMode int

const (
	IgnoreAspectRatio          AspectRatioMode = 0
	KeepAspectRatio            AspectRatioMode = 1
	KeepAspectRatioByExpanding AspectRatioMode = 2
)