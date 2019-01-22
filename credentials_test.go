package donutchallenge

import (
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"
)

type envVariables struct {
	Key        string
	Secret     string
	Passphrase string
}

func getDummyRequest() *http.Request {
	body := strings.NewReader("")
	return httptest.NewRequest("POST", "https://dummy.com/fake", body)
}

func resetEnvVariables(initialEnvVariables envVariables) {
	os.Setenv("CB_ACCESS_KEY", initialEnvVariables.Key)
	os.Setenv("CB_ACCESS_SECRET", initialEnvVariables.Secret)
	os.Setenv("CB_ACCESS_PASSPHRASE", initialEnvVariables.Passphrase)
}

func TestGetCredentialsSetsCredentialsFromEnvVariables(t *testing.T) {
	initialEnvVariables := envVariables{
		Key:        os.Getenv("CB_ACCESS_KEY"),
		Secret:     os.Getenv("CB_ACCESS_SECRET"),
		Passphrase: os.Getenv("CB_ACCESS_PASSPHRASE"),
	}
	os.Setenv("CB_ACCESS_KEY", "testkey")
	os.Setenv("CB_ACCESS_SECRET", "testsecret")
	os.Setenv("CB_ACCESS_PASSPHRASE", "testpassphrase")

	credentials := getCredentials()

	table := []struct {
		Actual   string
		Expected string
	}{
		{credentials.CBAccessKey, "testkey"},
		{credentials.CBAccessSecret, "testsecret"},
		{credentials.CBAccessPassphrase, "testpassphrase"},
	}

	for _, test := range table {
		if test.Actual != test.Expected {
			resetEnvVariables(initialEnvVariables)
			t.Errorf("Expected credential to equal '%s', got %s", test.Expected, test.Actual)
		}
	}
	resetEnvVariables(initialEnvVariables)
}

func TestGetCredentialsSetsTimestamp(t *testing.T) {
	now := int(time.Now().Unix())
	credentials := getCredentials()
	timestamp, err := strconv.Atoi(credentials.Timestamp)
	if err != nil || timestamp <= now-1 || timestamp >= now+30 {
		t.Errorf("Expected timestamp to be current time")
	}
}

func TestSetHeaders(t *testing.T) {
	req := getDummyRequest()
	credentials := coinbaseCredentials{
		CBAccessKey:        "testkey",
		CBAccessSecret:     "testsecret",
		CBAccessPassphrase: "testpassphrase",
		Timestamp:          "1548163998",
	}
	signature := "testsignature"
	setHeaders(req, credentials, signature)

	table := []struct {
		Header string
		Value  string
	}{
		{"Content-Type", "application/json"},
		{"CB-ACCESS-KEY", credentials.CBAccessKey},
		{"CB-ACCESS-PASSPHRASE", credentials.CBAccessPassphrase},
		{"CB-ACCESS-TIMESTAMP", credentials.Timestamp},
		{"CB-ACCESS-SIGN", signature},
	}

	for _, test := range table {
		actual := req.Header.Get(test.Header)
		if actual != test.Value {
			t.Errorf("Expected request header %s to equal %s, got %s", test.Header, test.Value, actual)
		}
	}
}

func TestBuildSignaturePropagatesBase64DecodeFailure(t *testing.T) {
	badSecret := "invalid*base64"
	timestamp := "123"
	req := getDummyRequest()
	body := ""
	signature, err := buildSignature(badSecret, timestamp, req, body)

	if signature != "" {
		t.Errorf("Did not expect signature when provided with invalid secret")
	}

	matched, _ := regexp.MatchString("illegal base64 data", err.Error())
	if !matched {
		t.Errorf("Expected invalid base64 secret to throw error")
	}
}

func TestBuildSignatureCreatesValidSignature(t *testing.T) {
	secret := "hC9Ipuo1vx52dA=="
	timestamp := "1548166753"
	req := getDummyRequest()
	body := "Some test body"
	expected := "8Fy8YdojnsoMOchOHwM0SPSq8FFSFZ4+iPV55Ft8990="
	signature, err := buildSignature(secret, timestamp, req, body)

	if err != nil {
		t.Errorf("Expected signature to build")
	}

	if signature != expected {
		t.Errorf("Expected signature to equal %s, got %s", expected, signature)
	}
}
