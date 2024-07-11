package util

const (
	USD = "USD"
	EUR = "EUR"
	CAD = "CAD"
	BDT = "BDT"
)

func SupportedCurrency(currency string) bool {
	switch currency {
	case USD, EUR, CAD, BDT:
		return true
	}
	return false
}
