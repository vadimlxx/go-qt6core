package qchar

// Category mirrors QChar::Category.
type Category int

const (
    Mark_NonSpacing Category = iota
    Mark_SpacingCombining
    Mark_Enclosing

    Number_DecimalDigit
    Number_Letter
    Number_Other

    Separator_Space
    Separator_Line
    Separator_Paragraph

    Other_Control
    Other_Format
    Other_Surrogate
    Other_PrivateUse
    Other_NotAssigned

    Letter_Uppercase
    Letter_Lowercase
    Letter_Titlecase
    Letter_Modifier
    Letter_Other

    Punctuation_Connector
    Punctuation_Dash
    Punctuation_Open
    Punctuation_Close
    Punctuation_InitialQuote
    Punctuation_FinalQuote
    Punctuation_Other

    Symbol_Math
    Symbol_Currency
    Symbol_Modifier
    Symbol_Other
)

func (c Category) String() string {
    switch c {
    case Mark_NonSpacing: return "Mark_NonSpacing"
    case Mark_SpacingCombining: return "Mark_SpacingCombining"
    case Mark_Enclosing: return "Mark_Enclosing"
    case Number_DecimalDigit: return "Number_DecimalDigit"
    case Number_Letter: return "Number_Letter"
    case Number_Other: return "Number_Other"
    case Separator_Space: return "Separator_Space"
    case Separator_Line: return "Separator_Line"
    case Separator_Paragraph: return "Separator_Paragraph"
    case Other_Control: return "Other_Control"
    case Other_Format: return "Other_Format"
    case Other_Surrogate: return "Other_Surrogate"
    case Other_PrivateUse: return "Other_PrivateUse"
    case Other_NotAssigned: return "Other_NotAssigned"
    case Letter_Uppercase: return "Letter_Uppercase"
    case Letter_Lowercase: return "Letter_Lowercase"
    case Letter_Titlecase: return "Letter_Titlecase"
    case Letter_Modifier: return "Letter_Modifier"
    case Letter_Other: return "Letter_Other"
    case Punctuation_Connector: return "Punctuation_Connector"
    case Punctuation_Dash: return "Punctuation_Dash"
    case Punctuation_Open: return "Punctuation_Open"
    case Punctuation_Close: return "Punctuation_Close"
    case Punctuation_InitialQuote: return "Punctuation_InitialQuote"
    case Punctuation_FinalQuote: return "Punctuation_FinalQuote"
    case Punctuation_Other: return "Punctuation_Other"
    case Symbol_Math: return "Symbol_Math"
    case Symbol_Currency: return "Symbol_Currency"
    case Symbol_Modifier: return "Symbol_Modifier"
    case Symbol_Other: return "Symbol_Other"
    }
    return "Category(?)"
}
