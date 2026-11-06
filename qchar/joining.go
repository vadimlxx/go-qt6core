package qchar

// JoiningType mirrors QChar::JoiningType.
type JoiningType int

const (
    Joining_None JoiningType = iota
    Joining_Causing
    Joining_Dual
    Joining_Right
    Joining_Left
    Joining_Transparent
)

func (j JoiningType) String() string {
    names := [...]string{
        "Joining_None", "Joining_Causing", "Joining_Dual",
        "Joining_Right", "Joining_Left", "Joining_Transparent",
    }
    if j >= 0 && int(j) < len(names) {
        return names[j]
    }
    return "JoiningType(?)"
}
