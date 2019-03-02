package check_test

import (
	"testing"

	"github.com/check/check"
)

func TestCheck(t *testing.T) {

	var testData = []struct {
		name        string
		in          string
		expectedIPs []string
		out         bool
	}{
		{
			name:        "test google.com",
			in:          "google.com",
			expectedIPs: []string{"216.58.201.46", "172.217.17.142"},
			out:         true,
		},
		{
			name:        "test github.com",
			in:          "github.com",
			expectedIPs: []string{"140.82.118.3", "140.82.118.4"},
			out:         true,
		},
		{
			name:        "test foobar122zzzxxx333.com",
			in:          "foobar122zzzxxx333.com",
			expectedIPs: []string{"foobar122zzzxxx333.com", ""},
			out:         false,
		},
		{
			name:        "test reddit.com",
			in:          "reddit.com",
			expectedIPs: []string{"151.101.193.140", "151.101.129.140", "151.101.65.140", "151.101.1.140"},
			out:         true,
		},
		{
			name:        "test reddit.com/r/golang",
			in:          "reddit.com/r/golang",
			expectedIPs: []string{"reddit.com/r/golang", ""},
			out:         false,
		},
	}

	for _, tt := range testData {
		check := check.New(tt.in)
		t.Run(tt.name, func(t *testing.T) {
			// isUp := check.IsUp(tt.in)

			if up := check.IsUp(); up != tt.out {
				t.Errorf("got %t, want %t", up, tt.out)
			}

			ip := check.IP()

			if !contains(tt.expectedIPs, ip) {
				t.Errorf("IP: %s is not in the list", ip)
			}

			v := check.Verbose()

			if v == "" {
				t.Errorf("%s is Empty", v)
			}

			t.Log(ip)
			t.Log(v)
		})
	}

}

func contains(arr []string, item string) bool {
	for _, a := range arr {
		if a == item {
			return true
		}
	}
	return false
}
