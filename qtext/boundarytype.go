package qtext

// BoundaryType mirrors QTextBoundaryFinder::BoundaryType (Qt 6.11).
//
// Selects the kind of text boundary that a QTextBoundaryFinder detects.
// The numeric values match Qt: Grapheme=0, Word=1, Sentence=2, Line=3.
type BoundaryType int

const (
	// Grapheme is a grapheme cluster boundary (user-perceived character).
	// Qt: QTextBoundaryFinder::Grapheme.
	Grapheme BoundaryType = 0

	// Word is a word boundary (UAX #29 word segmentation).
	// Qt: QTextBoundaryFinder::Word.
	Word BoundaryType = 1

	// Sentence is a sentence boundary (UAX #29 sentence segmentation).
	// Qt: QTextBoundaryFinder::Sentence.
	Sentence BoundaryType = 2

	// Line is a line-break opportunity (UAX #14 line breaking).
	// Qt: QTextBoundaryFinder::Line.
	Line BoundaryType = 3
)