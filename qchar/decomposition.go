package qchar

// Decomposition mirrors QChar::Decomposition. Ends at Fraction (Qt 6.11).
type Decomposition int

const (
    NoDecomposition Decomposition = iota
    Canonical
    Font
    NoBreak
    Initial
    Medial
    Final
    Isolated
    Circle
    Super
    Sub
    Vertical
    Wide
    Narrow
    Small
    Square
    Compat
    Fraction
)

func (d Decomposition) String() string {
    names := [...]string{
        "NoDecomposition", "Canonical", "Font", "NoBreak", "Initial",
        "Medial", "Final", "Isolated", "Circle", "Super", "Sub",
        "Vertical", "Wide", "Narrow", "Small", "Square", "Compat",
        "Fraction",
    }
    if d >= 0 && int(d) < len(names) {
        return names[d]
    }
    return "Decomposition(?)"
}
