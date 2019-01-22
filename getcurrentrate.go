package donutchallenge

import (
	"net/http"
	"strconv"
)

const nilPrice = 0.0

type coinbaseRateData struct {
	Price string `json:"price"`
}

func getRateFromResponse(res *http.Response) (float64, error) {
	var data coinbaseRateData
	err := getJSONDataFromResponse(res, &data)
	if err != nil {
		return nilPrice, err
	}

	return strconv.ParseFloat(data.Price, 64)
}

func buildGetCurrentRateURL(currency1 Currency, currency2 Currency) string {
	apiBase := getAPIBase()
	product := string(currency1) + "-" + string(currency2)
	return apiBase + "/products/" + product + "/ticker"
}

func getCurrentRateWithClient(client *http.Client, currency1 Currency, currency2 Currency) (float64, error) {
	productURL := buildGetCurrentRateURL(currency1, currency2)

	res, err := client.Get(productURL)
	if err != nil {
		return nilPrice, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nilPrice, handleNon200(res)
	}

	return getRateFromResponse(res)
}

// GetCurrentRate asks the Coinbase API for the current rate between two
// currencies
func GetCurrentRate(currency1 Currency, currency2 Currency) (float64, error) {
	client := &http.Client{}
	return getCurrentRateWithClient(client, currency1, currency2)
}
