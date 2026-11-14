package qdata

// MouseEventSource mirrors Qt::MouseEventSource from qnamespace.h (Qt 6.11).
type MouseEventSource int

const (
	MouseEventNotSynthesized           MouseEventSource = 0
	MouseEventSynthesizedBySystem      MouseEventSource = 1
	MouseEventSynthesizedByQt          MouseEventSource = 2
	MouseEventSynthesizedByApplication MouseEventSource = 3
)