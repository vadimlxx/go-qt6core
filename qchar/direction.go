package qchar

// Direction mirrors QChar::Direction.
type Direction int

const (
    DirL Direction = iota
    DirR
    DirEN
    DirES
    DirET
    DirAN
    DirCS
    DirB
    DirS
    DirWS
    DirON
    DirLRE
    DirLRO
    DirAL
    DirRLE
    DirRLO
    DirPDF
    DirNSM
    DirBN
    DirLRI
    DirRLI
    DirFSI
    DirPDI
)

func (d Direction) String() string {
    names := [...]string{
        "DirL", "DirR", "DirEN", "DirES", "DirET", "DirAN", "DirCS",
        "DirB", "DirS", "DirWS", "DirON", "DirLRE", "DirLRO", "DirAL",
        "DirRLE", "DirRLO", "DirPDF", "DirNSM", "DirBN",
        "DirLRI", "DirRLI", "DirFSI", "DirPDI",
    }
    if d >= 0 && int(d) < len(names) {
        return names[d]
    }
    return "Direction(?)"
}
