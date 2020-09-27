package utils

import (
	"avito/internal/pkg/models"
	"github.com/mailru/easyjson"
	"github.com/valyala/fasthttp"
)

func Exchange(baseCurrency string, amount float64, otherCurrency string) (float64, error) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	req.SetRequestURI("https://api.exchangeratesapi.io/latest?base="+ baseCurrency + "&symbols=" + otherCurrency)

	err := fasthttp.Do(req, resp)
	if err != nil {
		return 0.0, err
	}

	rate := &models.Rates{}
	err = easyjson.Unmarshal(resp.Body(), rate)
	if err != nil {
		return 0.0, err
	}

	return amount * rate.Rates[otherCurrency], nil
}
