package util

const (
	USD = "USD"
	EUR = "EUR"
	CAD = "CAD"
	GBP = "GBP"
	JPY = "JPY"
)

func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, EUR, CAD, GBP, JPY:
		return true
	}
	return false
}