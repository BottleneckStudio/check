package check_test

import (
	"testing"

	"github.com/check"
)

func TestCheck(t *testing.T) {

	var testData = []struct {
		name string
		in   string
		out  bool
	}{
		{
			name: "test google.com",
			in:   "google.com",
			out:  true,
		},
		{
			name: "test github.com",
			in:   "github.com",
			out:  true,
		},
		{
			name: "test foobar122zzzxxx333.com",
			in:   "foobar122zzzxxx333.com",
			out:  false,
		},
		{
			name: "test reddit.com",
			in:   "reddit.com",
			out:  true,
		},
		{
			name: "test reddit.com/r/golang",
			in:   "reddit.com/r/golang",
			out:  false,
		},
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			isUp := check.IsUp(tt.in)

			if isUp != tt.out {
				t.Errorf("got %t, want %t", isUp, tt.out)
			}
		})
	}

}
