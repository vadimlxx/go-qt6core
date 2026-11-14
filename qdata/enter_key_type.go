package qdata

// EnterKeyType mirrors Qt::EnterKeyType from qnamespace.h (Qt 6.11).
type EnterKeyType int

const (
	EnterKeyDefault  EnterKeyType = 0
	EnterKeyReturn   EnterKeyType = 1
	EnterKeyDone     EnterKeyType = 2
	EnterKeyGo       EnterKeyType = 3
	EnterKeySend     EnterKeyType = 4
	EnterKeySearch   EnterKeyType = 5
	EnterKeyNext     EnterKeyType = 6
	EnterKeyPrevious EnterKeyType = 7
)