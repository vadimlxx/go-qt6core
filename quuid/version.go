package quuid

// Version mirrors QUuid::Version (Qt 6.11).
type Version int

const (
	VerUnknown    Version = -1
	Time          Version = 1
	EmbeddedPOSIX Version = 2
	Md5           Version = 3
	Name          Version = Md5
	Random        Version = 4
	Sha1          Version = 5
	UnixEpoch     Version = 7
)