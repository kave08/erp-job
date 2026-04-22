package syncdata

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func decodeJSONResponse(res *http.Response, target interface{}) error {
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		return fmt.Errorf("http request failed. status: %d, response: %s", res.StatusCode, strings.TrimSpace(string(body)))
	}

	if err := json.NewDecoder(res.Body).Decode(target); err != nil {
		return err
	}

	return nil
}
