package qtext

// UnicodeVersion mirrors QChar::UnicodeVersion from Qt 6.11.
// Note: the first value is Unicode_Unassigned (offset 0), and the
// version "2.1.2" — not "2.1.9".
type UnicodeVersion int

const (
    Unicode_Unassigned UnicodeVersion = iota
    Unicode_1_1
    Unicode_2_0
    Unicode_2_1_2
    Unicode_3_0
    Unicode_3_1
    Unicode_3_2
    Unicode_4_0
    Unicode_4_1
    Unicode_5_0
    Unicode_5_1
    Unicode_5_2
    Unicode_6_0
    Unicode_6_1
    Unicode_6_2
    Unicode_6_3
    Unicode_7_0
    Unicode_8_0
    Unicode_9_0
    Unicode_10_0
    Unicode_11_0
    Unicode_12_0
    Unicode_12_1
    Unicode_13_0
    Unicode_14_0
    Unicode_15_0
    Unicode_15_1
    Unicode_16_0
    Unicode_17_0
)

func (u UnicodeVersion) String() string {
    names := [...]string{
        "Unicode_Unassigned",
        "Unicode_1_1", "Unicode_2_0", "Unicode_2_1_2", "Unicode_3_0",
        "Unicode_3_1", "Unicode_3_2", "Unicode_4_0", "Unicode_4_1",
        "Unicode_5_0", "Unicode_5_1", "Unicode_5_2", "Unicode_6_0",
        "Unicode_6_1", "Unicode_6_2", "Unicode_6_3", "Unicode_7_0",
        "Unicode_8_0", "Unicode_9_0", "Unicode_10_0", "Unicode_11_0",
        "Unicode_12_0", "Unicode_12_1", "Unicode_13_0", "Unicode_14_0",
        "Unicode_15_0", "Unicode_15_1", "Unicode_16_0", "Unicode_17_0",
    }
    if u >= 0 && int(u) < len(names) {
        return names[u]
    }
    return "UnicodeVersion(?)"
}
