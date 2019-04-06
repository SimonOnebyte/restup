package restup

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// RestUp ...
type RestUp struct {
	baseURL   string
	authToken string
	client    *http.Client
}

// NewRestUp returns a new RestUp with the authentication token encoded
func NewRestUp(baseURL string, token string) *RestUp {
	rup := RestUp{}
	rup.baseURL = baseURL
	rup.authToken = fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(token)))
	rup.client = &http.Client{
		Timeout: time.Second * 30,
	}
	return &rup
}

// SetHTTPClient allows you to supply a custom HTTP client, useful for testing
func (rup *RestUp) SetHTTPClient(client *http.Client) {
	rup.client = client
}

// Get performs the requested API Get returning the results as JSON
func (rup *RestUp) Get(url string, out interface{}) error {

	req, reqErr := http.NewRequest(http.MethodGet, url, nil)
	if reqErr != nil {
		return reqErr
	}

	req.Header.Set("Authorization", rup.authToken)
	req.Header.Set("Content-Type", "application/json")

	return rup.doRequestToJSON(req, out)
}

// Post performs the requested API Post returning the results as JSON
func (rup *RestUp) Post(url string, query interface{}, out interface{}) error {

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(query)

	req, reqErr := http.NewRequest(http.MethodPost, url, b)
	if reqErr != nil {
		return reqErr
	}

	req.Header.Set("Authorization", rup.authToken)
	req.Header.Set("Content-Type", "application/json")

	return rup.doRequestToJSON(req, out)
}

func (rup *RestUp) doRequestToJSON(req *http.Request, out interface{}) error {

	res, getErr := rup.client.Do(req)
	if getErr != nil {
		return getErr
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("Status: %s", res.Status)
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		return fmt.Errorf("Error reading body: %s", readErr)
	}

	jsonErr := json.Unmarshal(body, out)
	if jsonErr != nil {
		return fmt.Errorf("Error decodeing JSONL %s", jsonErr)
	}

	return nil
}
