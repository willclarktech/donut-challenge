package donutchallenge

import (
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"
)

func resetEnvVariable(initialEnvVariable string) {
	os.Setenv("TEST", initialEnvVariable)
}

func TestGetAPIBaseRespondsToEnvVariable(t *testing.T) {
	initialEnvVariable := os.Getenv("TEST")

	os.Setenv("TEST", "true")
	apiBase := getAPIBase()

	if apiBase != CoinbaseAPIBaseTest {
		resetEnvVariable(initialEnvVariable)
		t.Errorf("Expected API base to be %s, got %s", CoinbaseAPIBaseTest, apiBase)
	}

	os.Setenv("TEST", "false")
	apiBase = getAPIBase()

	if apiBase != CoinbaseAPIBaseProduction {
		resetEnvVariable(initialEnvVariable)
		t.Errorf("Expected API base to be %s, got %s", CoinbaseAPIBaseProduction, apiBase)
	}

	resetEnvVariable(initialEnvVariable)
}

type JSONTestType struct {
	Key1 string `json:"key1"`
	Key2 string `json:"key2,omitempty"`
	Key3 int    `json:"key_3,omitempty"`
}

func TestGetJSONDataFromResponseReturnsUnmarshalledJSON(t *testing.T) {
	table := []struct {
		Body   string
		Result JSONTestType
	}{
		{"{\"key1\":\"value1\",\"key2\":\"value2\",\"key_3\":3}", JSONTestType{
			Key1: "value1",
			Key2: "value2",
			Key3: 3,
		}},
		{"{}", JSONTestType{
			Key1: "",
			Key2: "",
			Key3: 0,
		}},
	}

	for _, test := range table {
		res := &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(strings.NewReader(test.Body)),
		}
		var data JSONTestType
		err := getJSONDataFromResponse(res, &data)

		if err != nil {
			t.Errorf("Expected JSON data not to throw, got %s", err.Error())
		}

		if data != test.Result {
			t.Errorf("Expected data to equal %v, got %v", test.Result, data)
		}
	}
}

func TestHandleNon200ReturnsHelpfulError(t *testing.T) {
	table := []struct {
		StatusCode   int
		Body         string
		ErrorMessage string
	}{
		{http.StatusNotFound, "{\"message\": \"NotFound\"}", "Request failed with status 404: NotFound"},
		{http.StatusInternalServerError, "{}", "Request failed with status 500"},
	}

	for _, test := range table {
		res := &http.Response{
			StatusCode: test.StatusCode,
			Body:       ioutil.NopCloser(strings.NewReader(test.Body)),
		}
		err := handleNon200(res)
		if err == nil || err.Error() != test.ErrorMessage {
			t.Errorf("Expected error message to be '%s', got '%s'", test.ErrorMessage, err.Error())
		}
	}
}
