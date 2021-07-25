package palgate

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const (
	authTokenHeaderName = "x-bt-user-token"

	getLogURL        = "/v1/bt/user/log"
	getUserURL       = "/v1/bt/user"
	getUsersURL      = "/v1/bt/device/%s/users-v2"
	updateUserURL    = "/v1/bt/device/%s/user"
	openGateURL      = "/v1/bt/device/%s/open-gate"
	deviceDetailsURL = "/v1/bt/device/%s/"
	listDevicesURL   = "/v1/bt/devices/"

	getLogParamID       = "id"
	getUserParamID      = "id"
	getUsersParamSkip   = "skip"
	getUsersParamLimit  = "limit"
	getUsersParamFilter = "filter"
	openGateParamOutput = "outputNum"
)

type Client interface {
	Log() (*GetLogResponse, error)
}

type client struct {
	httpClient *http.Client

	config *Config
}

func NewClient(config *Config) Client {
	return &client{
		httpClient: &http.Client{},
		config:     config,
	}
}

func (c *client) Log() (*GetLogResponse, error) {
	url, _ := url.Parse(c.config.ServerAddress + getLogURL)
	query := url.Query()
	query.Set(getLogParamID, c.config.GateID)
	url.RawQuery = query.Encode()
	url.Scheme = "HTTPS"

	req, _ := http.NewRequest(http.MethodGet, url.String(), nil)
	req.Header.Add(authTokenHeaderName, c.config.AuthToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error executing HTTP request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error response from server: %s", resp.Status)
	}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	var logResponse GetLogResponse
	err = json.Unmarshal(bytes, &logResponse)
	if err != nil {
		return nil, fmt.Errorf("error parsing response body: %v", err)
	}

	return &logResponse, nil
}
