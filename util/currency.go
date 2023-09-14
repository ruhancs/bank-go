package util

const (
	USD = "USD"
	EUR = "EUR"
	BRL = "BRL"
)

func IsSupportedCurrency(curency string) bool {
	switch curency {
	case USD,EUR,BRL:
		return true
	}
	return false
}