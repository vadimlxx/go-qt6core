package qdata

// TransformationMode mirrors Qt::TransformationMode from qnamespace.h
// (Qt 6.11).
type TransformationMode int

const (
	FastTransformation   TransformationMode = 0
	SmoothTransformation TransformationMode = 1
)