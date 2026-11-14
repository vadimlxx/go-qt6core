package qdata

// TextElideMode mirrors Qt::TextElideMode from qnamespace.h (Qt 6.11).
type TextElideMode int

const (
	ElideLeft   TextElideMode = 0
	ElideRight  TextElideMode = 1
	ElideMiddle TextElideMode = 2
	ElideNone   TextElideMode = 3
)