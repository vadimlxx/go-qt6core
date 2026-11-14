package qdata

// TimerId mirrors Qt::TimerId (enum class) from qnamespace.h (Qt 6.11).
//
// The only named value is Invalid (0); any other value is a live timer id
// returned by QObject::startTimer.
type TimerId int

const (
	// TimerIdInvalid is the null timer id (Qt::TimerId::Invalid).
	TimerIdInvalid TimerId = 0
)