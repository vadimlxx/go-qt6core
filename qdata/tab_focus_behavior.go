package qdata

// TabFocusBehavior mirrors Qt::TabFocusBehavior from qnamespace.h
// (Qt 6.11).
type TabFocusBehavior int

const (
	NoTabFocus           TabFocusBehavior = 0x00
	TabFocusTextControls TabFocusBehavior = 0x01
	TabFocusListControls TabFocusBehavior = 0x02
	TabFocusAllControls  TabFocusBehavior = 0xff
)