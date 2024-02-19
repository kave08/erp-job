package sync

// import (
// 	"bytes"
// 	"encoding/json"
// 	"fmt"
// 	"io"
// 	"net/http"
// 	"runtime"
// 	"time"

// 	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
// )

// type client struct {
// 	baseURL    string
// 	httpClient *http.Client
// 	Token      string
// }
// type middleware struct {
// 	UserAgent string
// 	Next      http.RoundTripper
// }

// // NewClient creates a new instance of invoice Client
// func NewClient(app string, version string, baseURL string, token string, timeout time.Duration) Contract {
// 	return &client{
// 		baseURL: baseURL,
// 		httpClient: &http.Client{
// 			Transport: &middleware{
// 				UserAgent: fmt.Sprintf("%s/%s (%s/%s)", app, version, runtime.GOOS, runtime.GOARCH),
// 				Next:      otelhttp.NewTransport(http.DefaultTransport),
// 			},
// 			Timeout: timeout,
// 		},
// 		Token: token,
// 	}
// }

// func (m middleware) RoundTrip(req *http.Request) (res *http.Response, e error) {
// 	req.Header.Set("User-Agent", m.UserAgent)
// 	return m.Next.RoundTrip(req)
// }

// // call will do the http request and decode the response into the v
// func (c *client) call(request *http.Request, v interface{}) error {
// 	request.Header.Set("Content-Type", "application/json")
// 	request.Header.Set("X-Api-Key", c.Token)

// 	resp, err := c.httpClient.Do(request)
// 	if err != nil {
// 		return err
// 	}
// 	defer closeBody(resp)

// 	if resp.StatusCode == http.StatusNoContent {
// 		return nil
// 	} else if resp.StatusCode > 400 {
// 		return fmt.Errorf(http.StatusText(resp.StatusCode))
// 	}

// 	err = json.NewDecoder(resp.Body).Decode(v)

// 	if err != nil {
// 		data, err := io.ReadAll(resp.Body)
// 		if err != nil {
// 			return err
// 		}
// 		return fmt.Errorf("body is not Json: %s", string(data))
// 	}

// 	return nil
// }

// // encode will encode the i into a buffer and return it. default encoding format is json
// func (c *client) encode(i interface{}) (*bytes.Buffer, error) {
// 	buffer := new(bytes.Buffer)
// 	err := json.NewEncoder(buffer).Encode(i)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return buffer, err
// }

// func closeBody(response *http.Response) {
// 	if response != nil {
// 		// to avoid memory leak when reusing http connection
// 		_, _ = io.Copy(io.Discard, response.Body)
// 		_ = response.Body.Close()
// 	}
// }
