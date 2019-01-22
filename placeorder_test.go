package donutchallenge

import (
	"encoding/json"
	"io/ioutil"
	"testing"
)

func TestBuildBodyMarshalsJSON(t *testing.T) {
	table := []struct {
		Params CoinbaseOrderParams
		Body   string
	}{
		{CoinbaseOrderParams{
			Side:      BUY,
			ProductID: "BTC-USD",
		}, "{\"side\":\"buy\",\"product_id\":\"BTC-USD\"}"},
		{CoinbaseOrderParams{
			ClientOID: "my_tag",
			Side:      SELL,
			ProductID: "BTC-ETH",
			STP:       CO,
			PostOnly:  true,
		}, "{\"client_oid\":\"my_tag\",\"side\":\"sell\",\"product_id\":\"BTC-ETH\",\"stp\":\"co\",\"post_only\":true}"},
	}

	for _, test := range table {
		body, err := buildBody(test.Params)

		if err != nil {
			t.Errorf("Expected body to build, got error: %s", err)
		}

		if body != test.Body {
			t.Errorf("Expected body '%s', got '%s'", test.Body, body)
		}
	}
}

func TestBuildRequestIncludesAllData(t *testing.T) {
	params := CoinbaseOrderParams{
		ClientOID: "my_tag",
		Side:      SELL,
		ProductID: "BTC-ETH",
		STP:       CO,
		PostOnly:  true,
	}
	credentials := coinbaseCredentials{
		CBAccessKey:        "testkey",
		CBAccessSecret:     "dGVzdHNlY3JldA==",
		CBAccessPassphrase: "testpassphrase",
		Timestamp:          "1548166753",
	}

	req, err := buildRequest(params, credentials)

	if err != nil {
		t.Errorf("Expected request to build successfully, got error: %s", err.Error())
	}

	if req.Method != "POST" {
		t.Errorf("Expected request method to be set to POST")
	}

	expectedURL := getAPIBase() + "/orders"
	if req.URL.String() != expectedURL {
		t.Errorf("Expected request URL to be %s, got %s", expectedURL, req.URL)
	}

	if req.Header.Get("CB-ACCESS-KEY") != credentials.CBAccessKey {
		t.Errorf("Expected request credential headers to be set")
	}

	if req.Header.Get("CB-ACCESS-SIGN") == "" {
		t.Errorf("Expected request signature header to be set")
	}

	body, _ := ioutil.ReadAll(req.Body)
	var data CoinbaseOrderParams
	json.Unmarshal(body, &data)
	if data != params {
		t.Errorf("Expected request body to match params, got %v", data)
	}
}

func TestPlaceOrderSuccess(t *testing.T) {
	err := PlaceOrder(CoinbaseOrderParams{
		ProductID: "BTC-USD",
		Side:      BUY,
		Size:      "0.03",
		Price:     "3200.11",
	})

	if err != nil {
		t.Errorf("Expected order to be placed successfully, got error: %s", err.Error())
	}
}

func TestPlaceOrderHandlesFailure(t *testing.T) {
	err := PlaceOrder(CoinbaseOrderParams{
		ProductID: "UKN-XXX",
		Side:      BUY,
		Size:      "4.23",
		Price:     "3200.11",
	})

	expected := "Request failed with status 404: Product not found"

	if err == nil || err.Error() != expected {
		t.Errorf("Expected order to fail with '%s', got '%s'", expected, err.Error())
	}
}
