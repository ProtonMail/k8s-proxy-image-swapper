package corednsnetboxplugin

import (
	"testing"

	"github.com/coredns/caddy"
)

func TestSetup(t *testing.T) {
	tests := []struct {
		input     string
		shouldErr bool
	}{
		{`netbox-dns {
				url test
			}`, false},
		{`netbox-dns {
				url test
                                fallthrough
			}`, false},
		{`netbox-dns example1.org {
				url test
			}`, false},
		{`netbox-dns example1.org example2.org {
				url test
			}`, false},

		// fails
		{`netbox-dns`, true}, // no url
	}
	for _, test := range tests {
		c := caddy.NewTestController("dns", test.input)
		if err := setup(c); (err != nil) != test.shouldErr {
			t.Errorf("Unexpected errors: %v", err)
		}
	}
}
