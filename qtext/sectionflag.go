package qtext

/*
 SectionFlag mirrors QString::SectionFlag from qstring.h (Qt 6.11).

 Флаги комбинируются побитовым ИЛИ и передаются в QString.SectionQChar/SectionQString.
*/

// SectionFlag controls the behavior of QString.Section*.
type SectionFlag int

const (
	SectionDefault             SectionFlag = 0x00
	SectionSkipEmpty           SectionFlag = 0x01
	SectionIncludeLeadingSep   SectionFlag = 0x02
	SectionIncludeTrailingSep  SectionFlag = 0x04
	SectionCaseInsensitiveSeps SectionFlag = 0x08
)