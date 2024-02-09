package helpers

import (
	"net/url"
	"strings"
)

func EncodeURLParams(inputURL string, URLParams map[string]string) string {
	params := url.Values{}

	for key, value := range URLParams {
		params.Add(key, value)
	}

	return strings.TrimRight(inputURL, "?") + "?" + params.Encode()
}
