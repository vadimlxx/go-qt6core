package qdatetime

// System mirrors QCalendar::System.
type System int

const (
	Gregorian System = 0
	Julian    System = 8
	Milankovic System = 9
	Jalali    System = 10
	IslamicCivil System = 11
	Last      System = 11
	User      System = -1
)
