package qdata

// PaintDeviceFlags mirrors QInternal::PaintDeviceFlags.
type PaintDeviceFlags int

const (
	UnknownDevice    PaintDeviceFlags = 0x00
	Widget           PaintDeviceFlags = 0x01
	Pixmap           PaintDeviceFlags = 0x02
	Image            PaintDeviceFlags = 0x03
	Printer          PaintDeviceFlags = 0x04
	Picture          PaintDeviceFlags = 0x05
	Pbuffer          PaintDeviceFlags = 0x06
	FramebufferObject PaintDeviceFlags = 0x07
	CustomRaster     PaintDeviceFlags = 0x08
	PaintBuffer      PaintDeviceFlags = 0x0a
	OpenGL           PaintDeviceFlags = 0x0b
)
