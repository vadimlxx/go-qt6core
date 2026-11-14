package qtext

/*
 NormalizationForm mirrors QString::NormalizationForm from qstring.h (Qt 6.11).

 Используется в QString.Normalized(mode) и QString.NormalizedWithVersion(mode, version).
*/

// NormalizationForm selects the Unicode normalization form used by QString.Normalized.
type NormalizationForm int

const (
	NormalizationForm_D  NormalizationForm = 0 // Canonical decomposition (NFD)
	NormalizationForm_C  NormalizationForm = 1 // Canonical composition (NFC)
	NormalizationForm_KD NormalizationForm = 2 // Compatibility decomposition (NFKD)
	NormalizationForm_KC NormalizationForm = 3 // Compatibility composition (NFKC)
)