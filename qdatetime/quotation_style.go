package qdatetime

// QuotationStyle mirrors QLocale::QuotationStyle.
type QuotationStyle int

const (
	StandardQuotation  QuotationStyle = 0
	AlternateQuotation QuotationStyle = 1
)