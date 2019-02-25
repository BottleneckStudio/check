package check

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

// Response represents the result of the site to check if its okay
type Response struct {
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

// IsUp is our little helper to check if the site is up.
func IsUp(site string) bool {
	reqURL := fmt.Sprintf("https://isitup.org/%s.json", site)

	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return false
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

	resp, err := client.Do(req)
	if err != nil {
		return false
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false
	}

	respo := Response{}

	err = json.Unmarshal(body, &respo)
	if err != nil {
		return false
	}

	if respo.StatusCode == statusInvalidDomain {
		return false
	}

	if respo.StatusCode == statusDown {
		return false
	}

	return respo.StatusCode == statusSuccess
}
