package qdata

// Initialization mirrors Qt::Initialization (and the inline constant
// Qt::Uninitialized) from qnamespace.h (Qt 6.11).
//
// Used as a tag in constructors like QString(qsizetype, Qt::Initialization)
// and QList(Qt::Initialization), which allocate without initializing.
type Initialization int

const (
	// InitializationUninitialized matches Qt::Initialization::Uninitialized
	// (exposed in C++ as Qt::Uninitialized).
	InitializationUninitialized Initialization = 0
)