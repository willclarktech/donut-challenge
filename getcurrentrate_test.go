package donutchallenge

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func TestGetRateFromResponseReturnsPrice(t *testing.T) {
	table := []struct {
		Body  string
		Price float64
	}{
		{"{\"price\":\"3200.44\"}", 3200.44},
		{"{\"price\":\"3200.0\"}", 3200.0},
		{"{\"price\":\"0.00001\"}", 0.00001},
	}

	for _, test := range table {
		res := &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(strings.NewReader(test.Body)),
		}
		price, err := getRateFromResponse(res)

		if err != nil {
			t.Errorf("Expected price to be extracted successfully, got error: %s", err.Error())
		}

		if price != test.Price {
			t.Errorf("Expected price to equal %f, got %f", test.Price, price)
		}
	}
}

func TestBuildGetCurrentRateURLCreatesCorrectURL(t *testing.T) {
	table := []struct {
		Currency1 Currency
		Currency2 Currency
		Result    string
	}{
		{BTC, USD, getAPIBase() + "/products/BTC-USD/ticker"},
		{USD, ETH, getAPIBase() + "/products/USD-ETH/ticker"},
		{LOOM, MANA, getAPIBase() + "/products/LOOM-MANA/ticker"},
		{BTC, MANA, getAPIBase() + "/products/BTC-MANA/ticker"},
	}

	for _, test := range table {
		url := buildGetCurrentRateURL(test.Currency1, test.Currency2)
		if url != test.Result {
			t.Errorf("Expected URL to be %s, got %s", test.Result, url)
		}
	}
}

func TestGetCurrentRateReturnsPositiveRate(t *testing.T) {
	rate, err := GetCurrentRate(BTC, USD)

	if err != nil {
		t.Errorf("Expected current rate to be successfully retrieved")
	}

	if !(rate > 0.0) {
		t.Errorf("Expected rate to be positive")
	}
}

func TestGetCurrentRateHandlesUnknownProduct(t *testing.T) {
	rate, err := GetCurrentRate(USD, BTC)

	if rate != 0.0 {
		t.Errorf("Expected rate to be 0.0")
	}

	if err == nil || err.Error() != "Request failed with status 404: NotFound" {
		t.Errorf("Expected 404 error for unknown product")
	}
}
