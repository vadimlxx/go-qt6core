package qdatetime

// MeasurementSystem mirrors QLocale::MeasurementSystem.
type MeasurementSystem int

const (
	MetricSystem      MeasurementSystem = 0
	ImperialUSSystem  MeasurementSystem = 1
	ImperialUKSystem  MeasurementSystem = 2
	ImperialSystem    = ImperialUSSystem // Qt 4 compatibility
)