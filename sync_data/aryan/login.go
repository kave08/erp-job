package asyncdata

import (
	"erp-job/common"
	"erp-job/config"
	"erp-job/utility/logger"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"go.uber.org/zap"
)

// Login represents the Aryan ERP system's data synchronization service.
type Login struct {
	log        *zap.SugaredLogger
	httpClient *http.Client
	baseURL    string
}

// NewLogin initializes and returns a new Aryan service instance.
func NewLogin() *Login {
	return &Login{
		log:     logger.Logger(),
		baseURL: config.Cfg.AryanApp.BaseURL,
		httpClient: &http.Client{
			Timeout: config.Cfg.AryanApp.Timeout,
		},
	}
}

// Login attempts to log in to the Fararavand ERP system with the provided credentials and place.
func (l *Login) Login(user, pass, place string) error {
	data := url.Values{}
	data.Set("user", user)
	data.Set("pass", pass)
	data.Set("place", place)

	req, err := http.NewRequest(http.MethodPost, l.baseURL+
		fmt.Sprintf(common.Login), strings.NewReader(data.Encode()))
	if err != nil {
		l.log.Errorw("failed to create login request",
			"error", err)
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := l.httpClient.Do(req)
	if err != nil {
		l.log.Errorw("failed to send login request",
			"error", err)
		return err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		l.log.Errorw("failed to read login response body",
			"error", err)
		return err
	}

	if res.StatusCode != http.StatusOK {
		l.log.Errorw("login failed",
			"status", res.StatusCode,
			"response", string(body))
		return fmt.Errorf("login failed with status: %d, response: %s", res.StatusCode, string(body))
	}

	l.log.Info("Login successful", "response", string(body))

	return nil
}
