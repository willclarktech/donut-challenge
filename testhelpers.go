package donutchallenge

import (
	"io"
	"io/ioutil"
	"strings"
)

func createBodyFromString(stringBody string) io.ReadCloser {
	return ioutil.NopCloser(strings.NewReader(stringBody))
}
