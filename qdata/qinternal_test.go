package qdata

import (
	"testing"
	"unsafe"
)

func TestQInternalTypesExist(t *testing.T) {
	// Compile-time check: QInternal and its enums exist and have expected shapes.
	var _ QInternal
	var _ Callback = EventNotifyCallback
	var _ PaintDeviceFlags = Widget
	var _ RelayoutType = RelayoutNormal
	var _ DockPosition = LeftDock
	var _ CallMode = Testing
	var _ SaveStateRule = KeepSavedState
}

func TestQInternalCallbackConstants(t *testing.T) {
	if EventNotifyCallback != 0 {
		t.Fatalf("EventNotifyCallback = %d, want 0", EventNotifyCallback)
	}
	if LastCallback != 1 {
		t.Fatalf("LastCallback = %d, want 1", LastCallback)
	}
}

func TestQInternalPaintDeviceFlags(t *testing.T) {
	cases := []struct {
		name string
		got  PaintDeviceFlags
		want PaintDeviceFlags
	}{
		{"UnknownDevice", UnknownDevice, 0x00},
		{"Widget", Widget, 0x01},
		{"Pixmap", Pixmap, 0x02},
		{"Image", Image, 0x03},
		{"Printer", Printer, 0x04},
		{"Picture", Picture, 0x05},
		{"Pbuffer", Pbuffer, 0x06},
		{"FramebufferObject", FramebufferObject, 0x07},
		{"CustomRaster", CustomRaster, 0x08},
		{"PaintBuffer", PaintBuffer, 0x0a},
		{"OpenGL", OpenGL, 0x0b},
	}
	for _, c := range cases {
		if c.got != c.want {
			t.Errorf("%s = %d, want %d", c.name, c.got, c.want)
		}
	}
}

func TestQInternalRelayoutType(t *testing.T) {
	if RelayoutNormal != 0 || RelayoutDragging != 1 || RelayoutDropped != 2 {
		t.Fatalf("RelayoutType values: %d, %d, %d", RelayoutNormal, RelayoutDragging, RelayoutDropped)
	}
}

func TestQInternalDockPosition(t *testing.T) {
	if LeftDock != 0 || RightDock != 1 || TopDock != 2 || BottomDock != 3 || DockCount != 4 {
		t.Fatalf("DockPosition values: %d, %d, %d, %d, %d",
			LeftDock, RightDock, TopDock, BottomDock, DockCount)
	}
}

func TestQInternalCallMode(t *testing.T) {
	if Testing != 0 || Live != 1 {
		t.Fatalf("CallMode: Testing=%d Live=%d", Testing, Live)
	}
}

func TestQInternalSaveStateRule(t *testing.T) {
	if KeepSavedState != 0 || ClearSavedState != 1 {
		t.Fatalf("SaveStateRule: Keep=%d Clear=%d", KeepSavedState, ClearSavedState)
	}
}

func TestQInternalRegisterUnregisterNil(t *testing.T) {
	// Passing nil callback must not crash. Qt's registerCallback checks the
	// callback slot and stores nullptr; unregisterCallback on a never-registered
	// callback returns false.
	if RegisterCallback(EventNotifyCallback, nil) {
		t.Log("RegisterCallback(nil) returned true — Qt version may allow it; not an error")
	}
	if UnregisterCallback(EventNotifyCallback, nil) {
		t.Log("UnregisterCallback(nil) returned true — not an error")
	}
}

func TestQInternalActivateNoCallbacks(t *testing.T) {
	// No callbacks registered — activateCallbacks returns false, must not crash.
	var params unsafe.Pointer
	if ActivateCallbacks(EventNotifyCallback, params) {
		t.Fatal("ActivateCallbacks with no registered callbacks must return false")
	}
}

func TestQInternalInvalidCallbackSlot(t *testing.T) {
	// Out-of-range Callback value must be rejected safely.
	if RegisterCallback(Callback(-1), nil) {
		t.Fatal("RegisterCallback(-1) must return false")
	}
	if RegisterCallback(Callback(999), nil) {
		t.Fatal("RegisterCallback(999) must return false")
	}
}
