package api

import (
	"github.com/go-playground/validator/v10"
	"github.com/ruhancs/bank-go/util"
)

var validCurrency validator.Func = func(fl validator.FieldLevel) bool {
	//ok verifica se é uma string ou nao
	if currency,ok := fl.Field().Interface().(string); ok {
		return util.IsSupportedCurrency(currency)
	}
	return false
}