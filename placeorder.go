package donutchallenge

import (
	"encoding/json"
	"net/http"
	"strings"
)

// CoinbaseOrderParams is a struct for holding parameters when placing an order
// via the Coinbase API
type CoinbaseOrderParams struct {
	ClientOID   string           `json:"client_oid,omitempty"`
	Type        OrderType        `json:"type,omitempty"`
	Side        OrderSide        `json:"side"`
	ProductID   string           `json:"product_id"`
	STP         OrderSTP         `json:"stp,omitempty"`
	Stop        OrderStop        `json:"stop,omitempty"`
	StopPrice   string           `json:"stop_price,omitempty"`
	Price       string           `json:"price,omitempty"`
	Size        string           `json:"size,omitempty"`
	TimeInForce OrderTimeInForce `json:"time_in_force,omitempty"`
	CancelAfter OrderCancelAfter `json:"cancel_after,omitempty"`
	PostOnly    bool             `json:"post_only,omitempty"`
	Funds       string           `json:"funds,omitempty"`
}

func buildBody(params CoinbaseOrderParams) (string, error) {
	body, err := json.Marshal(params)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func buildPlaceOrderURL() string {
	apiBase := getAPIBase()
	return apiBase + CoinbaseOrdersEndpoint
}

func buildRequest(params CoinbaseOrderParams, credentials coinbaseCredentials) (*http.Request, error) {
	const method = "POST"
	orderURL := buildPlaceOrderURL()

	body, err := buildBody(params)
	if err != nil {
		return nil, err
	}

	bodyReader := strings.NewReader(body)
	req, err := http.NewRequest(method, orderURL, bodyReader)
	if err != nil {
		return nil, err
	}

	signature, err := buildSignature(credentials.CBAccessSecret, credentials.Timestamp, req, body)
	if err != nil {
		return nil, err
	}

	setHeaders(req, credentials, signature)

	return req, nil
}

func placeOrderWithClient(client *http.Client, params CoinbaseOrderParams, credentials coinbaseCredentials) error {
	req, err := buildRequest(params, credentials)
	if err != nil {
		return err
	}

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return handleNon200(res)
	}

	return nil
}

// PlaceOrder places an order via the Coinbase API
func PlaceOrder(params CoinbaseOrderParams) error {
	client := &http.Client{}
	credentials := getCredentials()
	return placeOrderWithClient(client, params, credentials)
}
