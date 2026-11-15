package qtext

// BoundaryReason mirrors QTextBoundaryFinder::BoundaryReason (Qt 6.11).
//
// The values form a bitmask. When QTextBoundaryFinder::boundaryReasons()
// returns NotAtBoundary (0), the finder is not sitting on a boundary.
// Otherwise the low 5 bits indicate a break opportunity (BreakOpportunity
// is a mask covering those bits), and the remaining bits carry the specific
// reason.
//
// Naming note: Qt names the soft-hyphen flag QTextBoundaryFinder::SoftHyphen,
// which would collide at Go package scope with QChar::SpecialCharacter::SoftHyphen
// (U+00AD) — also exported from package qtext. To disambiguate, the
// boundary-reason constant is exported as BoundarySoftHyphen; QChar's
// U+00AD keeps the short name SoftHyphen because it is by far the more
// commonly used Unicode constant. All other BoundaryReason constants match
// the Qt spelling.
type BoundaryReason int

const (
	// NotAtBoundary means the current position is not on a boundary.
	// Qt: QTextBoundaryFinder::NotAtBoundary (== 0).
	NotAtBoundary BoundaryReason = 0

	// BreakOpportunity is a mask of the low 5 bits. Any non-zero result
	// in that mask means the position is a break opportunity.
	// Qt: QTextBoundaryFinder::BreakOpportunity (== 0x1f).
	BreakOpportunity BoundaryReason = 0x1f

	// StartOfItem marks the start of a text item.
	// Qt: QTextBoundaryFinder::StartOfItem (== 0x20).
	StartOfItem BoundaryReason = 0x20

	// EndOfItem marks the end of a text item.
	// Qt: QTextBoundaryFinder::EndOfItem (== 0x40).
	EndOfItem BoundaryReason = 0x40

	// MandatoryBreak marks a mandatory break (e.g. newline).
	// Qt: QTextBoundaryFinder::MandatoryBreak (== 0x80).
	MandatoryBreak BoundaryReason = 0x80

	// BoundarySoftHyphen marks a soft-hyphen break opportunity.
	// Qt: QTextBoundaryFinder::SoftHyphen (== 0x100) — renamed to
	// BoundarySoftHyphen to avoid a name clash with QChar::SoftHyphen
	// (U+00AD) in the same qtext package.
	BoundarySoftHyphen BoundaryReason = 0x100
)

// BoundaryReasons is a bitmask of BoundaryReason flags.
//
// In C++, Q_DECLARE_FLAGS(BoundaryReasons, BoundaryReason) creates a distinct
// QFlags<BoundaryReason> type with an implicit constructor taking
// BoundaryReason, so individual flags can be passed anywhere a BoundaryReasons
// is expected. Go has no implicit conversions, so BoundaryReasons is declared
// as a type alias for BoundaryReason — the ergonomic equivalent of Qt's QFlags
// behavior.
type BoundaryReasons = BoundaryReason