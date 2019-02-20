package models

// Cost and other cloud constants
const (
	// Cost constants
	DefaultCPUCostPerCPUPerHour    = "0.024"
	DefaultMemCostPerGBPerHour     = "0.01"
	DefaultStorageCostPerGBPerHour = "0.00013888888"
	DefaultStorageCostInFloat64    = 0.00013888888

	// Cloud provider constants
	AWS = "aws"

	// Time constants
	HoursInMonth = 720

	// Other constants
	PriceError = -1.0
)
