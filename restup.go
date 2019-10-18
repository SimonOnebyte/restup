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

// RequestError returns the error context of a failed POST call
type RequestError struct {
	Method    string
	URL       string
	Body      string
	HTTPError error
}

func (e RequestError) Error() string {
	return fmt.Sprintf("error in restup request.\nMethod : %v\nURL   : %v\n%v\n", e.Method, e.URL, e.Body)
}

// RestUp ...
type RestUp struct {
	baseURL   string
	authToken string
	client    *http.Client
	headers   map[string]string
}

// NewRestUp returns a new RestUp with the authentication token encoded
func NewRestUp(baseURL string, token string) *RestUp {
	rup := RestUp{}
	rup.baseURL = baseURL
	rup.authToken = fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(token)))
	rup.client = &http.Client{
		Timeout: time.Second * 30,
	}
	rup.headers = make(map[string]string)

	return &rup
}

// TransportIntercept allows you to supply a custom HTTP client, useful for testing
func (rup *RestUp) TransportIntercept(rt http.RoundTripper) {
	rup.client.Transport = rt
}

// AddHeader allows additional headers to be added to the API request
func (rup *RestUp) AddHeader(name, value string) {
	rup.headers[name] = value
}

// Get performs the requested API Get returning the results as JSON
func (rup *RestUp) Get(action string, out interface{}) error {

	url := rup.baseURL + action
	req, reqErr := http.NewRequest(http.MethodGet, url, nil)
	if reqErr != nil {
		return reqErr
	}

	rup.setReqHeaders(req)
	req.Header.Set("Authorization", rup.authToken)
	req.Header.Set("Content-Type", "application/json")

	err := rup.doRequestToJSON(req, out)
	if err == nil {
		return nil
	}
	return RequestError{
		Method:    "GET",
		URL:       url,
		Body:      err.Error(),
		HTTPError: err,
	}
}

// Post performs the requested API Post returning the results as JSON
func (rup *RestUp) Post(action string, query interface{}, out interface{}) error {

	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(query)

	url := rup.baseURL + action

	req, reqErr := http.NewRequest(http.MethodPost, url, body)
	if reqErr != nil {
		return reqErr
	}

	rup.setReqHeaders(req)
	req.Header.Set("Authorization", rup.authToken)
	req.Header.Set("Content-Type", "application/json")

	err := rup.doRequestToJSON(req, out)
	if err == nil {
		return nil
	}
	return RequestError{
		Method:    "POST",
		URL:       url,
		Body:      err.Error(),
		HTTPError: err,
	}
}

func (rup *RestUp) doRequestToJSON(req *http.Request, out interface{}) error {

	res, getErr := rup.client.Do(req)
	if getErr != nil {
		return getErr
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		return fmt.Errorf("Error reading body: %s", readErr)
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("Status: %s\nBody  : %s", res.Status, body)
	}

	jsonErr := json.Unmarshal(body, out)
	if jsonErr != nil {
		return fmt.Errorf("Error decodeing JSONL %s", jsonErr)
	}

	return nil
}

func (rup *RestUp) setReqHeaders(req *http.Request) {
	if len(rup.headers) > 0 {
		for k, v := range rup.headers {
			req.Header.Set(k, v)
		}
	}
}
