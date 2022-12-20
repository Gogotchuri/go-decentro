package decentro

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

var _ CollectionsAPI = &client{}

type client struct {
	httpClient http.Client

	baseAPIURL         string
	clientID           string
	clientSecret       string
	moduleSecret       string
	providerSecret     string
	defaultPayeeNumber string
}

func CreateCollectionsClient(clientID, clientSecret, moduleSecret, providerSecret, payeeNumber string, isStaging bool) CollectionsAPI {
	url := ProdAPIURL
	if isStaging {
		url = StagingAPIURL
	}
	return &client{
		httpClient:         http.Client{},
		clientID:           clientID,
		clientSecret:       clientSecret,
		moduleSecret:       moduleSecret,
		providerSecret:     providerSecret,
		defaultPayeeNumber: payeeNumber,
		baseAPIURL:         url,
	}
}

/** NewRequest constructs a new http.Request, Marshal payload to json bytes */
func (c client) newRequest(method, url string, payload interface{}) (*http.Request, error) {
	headers := http.Header{}
	headers.Set("Content-Type", "application/json")
	headers.Set("Accept", "application/json")

	headers.Set("client_id", c.clientID)
	headers.Set("client_secret", c.clientSecret)
	headers.Set("module_secret", c.moduleSecret)
	headers.Set("provider_secret", c.providerSecret)
	var buf io.Reader
	if payload != nil {
		b, err := json.Marshal(&payload)
		if err != nil {
			return nil, err
		}
		buf = bytes.NewBuffer(b)
	}
	req, err := http.NewRequest(method, url, buf)
	if err != nil {
		return nil, err
	}
	req.Header = headers
	return req, nil
}

func (c client) send(method string, path string, body interface{}, dst interface{}) error {
	req, err := c.newRequest(method, c.baseAPIURL+path, body)
	if err != nil {
		return err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if dst == nil {
		return nil
	}
	if w, ok := dst.(io.Writer); ok {
		if _, err := io.Copy(w, resp.Body); err != nil {
			return err
		}
		return nil
	}

	return json.NewDecoder(resp.Body).Decode(dst)
}
