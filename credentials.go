package donutchallenge

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"net/http"
	"os"
	"strconv"
	"time"
)

type coinbaseCredentials struct {
	CBAccessKey        string
	CBAccessPassphrase string
	CBAccessSecret     string
	Timestamp          string
}

func getCredentials() coinbaseCredentials {
	return coinbaseCredentials{
		CBAccessKey:        os.Getenv("CB_ACCESS_KEY"),
		CBAccessSecret:     os.Getenv("CB_ACCESS_SECRET"),
		CBAccessPassphrase: os.Getenv("CB_ACCESS_PASSPHRASE"),
		Timestamp:          strconv.FormatInt(time.Now().Unix(), 10),
	}
}

func setHeaders(req *http.Request, credentials coinbaseCredentials, signature string) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("CB-ACCESS-KEY", credentials.CBAccessKey)
	req.Header.Set("CB-ACCESS-PASSPHRASE", credentials.CBAccessPassphrase)
	req.Header.Set("CB-ACCESS-TIMESTAMP", credentials.Timestamp)
	req.Header.Set("CB-ACCESS-SIGN", signature)
}

func buildSignature(secret string, timestamp string, req *http.Request, body string) (string, error) {
	payload := timestamp + req.Method + CoinbaseOrdersEndpoint + body
	key, err := base64.StdEncoding.DecodeString(secret)
	if err != nil {
		return "", err
	}

	mac := hmac.New(sha256.New, key)
	_, err = mac.Write([]byte(payload))
	if err != nil {
		return "", err
	}

	signature := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	return signature, nil
}
