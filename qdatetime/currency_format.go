package qdatetime

// CurrencySymbolFormat mirrors QLocale::CurrencySymbolFormat.
type CurrencySymbolFormat int

const (
	CurrencyIsoCode     CurrencySymbolFormat = 0
	CurrencySymbol      CurrencySymbolFormat = 1
	CurrencyDisplayName CurrencySymbolFormat = 2
)