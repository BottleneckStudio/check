// Package check identifies if a given site is up or not.
// Data is from https://isitup.org.
//
// Example:
// check := check.New("google.com")
//
// check.IsUp()
// check.IP()
// check.Verbose()
package check

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"sync"
	"time"
)

// IsUpResponse represents the result of the site to check if its okay
type IsUpResponse struct {
	Domain       string  `json:"domain"`
	Port         int     `json:"port"`
	StatusCode   int     `json:"status_code"`
	ResponseIP   string  `json:"response_ip"`
	ResponseCode int     `json:"response_code"`
	ResponseTime float64 `json:"response_time"`
}

const (
	statusSuccess       = 1
	statusDown          = 2
	statusInvalidDomain = 3
)

// Check custom type for our package to check if a site is up or not.
type Check struct {
	mu     *sync.RWMutex
	Res    *IsUpResponse
	Req    *http.Request
	Client *http.Client
}

// New returns a pointer to check that defines the different methods to check a site.
func New(site string) *Check {
	reqURL := fmt.Sprintf("https://isitup.org/%s.json", site)
	req, err := http.NewRequest("GET", reqURL, nil)

	if err != nil {
		return nil
	}
	req.Header.Set("User-Agent", "https://github.com/BottleneckStudio/check")

	var netTransport = &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 5 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 5 * time.Second,
	}

	var client = &http.Client{
		Timeout:   time.Second * 10,
		Transport: netTransport,
	}

	return &Check{
		mu:     &sync.RWMutex{},
		Res:    new(IsUpResponse),
		Req:    req,
		Client: client,
	}
}

// IsUp checks whether the given site is up or not.
func (c *Check) IsUp() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	body, err := getResponseBody(c.Client, c.Req)
	if err != nil {
		return false
	}

	err = json.Unmarshal(body, &c.Res)
	if err != nil {
		return false
	}

	if c.Res.StatusCode == statusInvalidDomain {
		return false
	}

	if c.Res.StatusCode == statusDown {
		return false
	}

	return c.Res.StatusCode == statusSuccess
}

// IP returns the IP of the given site.
func (c *Check) IP() string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	body, err := getResponseBody(c.Client, c.Req)
	if err != nil {
		return ""
	}
	err = json.Unmarshal(body, &c.Res)
	if err != nil {
		return ""
	}

	return c.Res.ResponseIP
}

// Verbose returns a prettified JSON format as a string.
// Example output:
// {
//     "domain": "google.com",
//     "port": 80,
//     "status_code": 1,
//     "response_ip": "216.58.201.46",
//     "response_code": 301,
//     "response_time": 0.007
// }
func (c *Check) Verbose() string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	iur := IsUpResponse{}

	body, err := getResponseBody(c.Client, c.Req)
	if err != nil {
		return ""
	}
	err = json.Unmarshal(body, &iur)
	if err != nil {
		return ""
	}

	return iur.String()
}

func (iur IsUpResponse) String() string {
	return fmt.Sprintf(`
		{
			"domain": %s,
			"port": %d,
			"status_code": %d,
			"response_ip": %s,
			"response_code": %d,
			"response_time": %.2f
		}
	`, iur.Domain, iur.Port, iur.StatusCode, iur.ResponseIP, iur.ResponseCode, iur.ResponseTime)
}

func getResponseBody(c *http.Client, req *http.Request) ([]byte, error) {
	resp, err := c.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}
