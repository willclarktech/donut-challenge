package donutchallenge

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

func getAPIBase() string {
	if os.Getenv("TEST") == "true" {
		return CoinbaseAPIBaseTest
	}
	return CoinbaseAPIBaseProduction
}

func getJSONDataFromResponse(res *http.Response, data interface{}) error {
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, data)
}

type coinbaseError struct {
	Message string `json:"message"`
}

func (e coinbaseError) Error() string {
	return e.Message
}

func handleNon200(res *http.Response) error {
	var data coinbaseError
	err := getJSONDataFromResponse(res, &data)
	if err != nil {
		return err
	}

	msg := "Request failed with status " + strconv.Itoa(res.StatusCode)
	if data.Message != "" {
		msg += ": " + data.Message
	}

	return errors.New(msg)
}
