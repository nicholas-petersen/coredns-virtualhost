package virtualhost

import (
	"strings"
	"testing"

	"github.com/coredns/caddy"
)


func TestSetup(t *testing.T) {
	tests := []struct{
			input	string
			shouldErr bool
			expectedErr string
	}{
		// positive
		{"virtualhost 192.168.0.100", false, ""},
		{"virtualhost fe80::8770:87a8:3d30:84da", false, ""},
		{"virtualhost 192.168.0.100 fe80::8770:87a8:3d30:84da", false, ""},
		{"virtualhost fe80::8770:87a8:3d30:84da 192.168.0.100", false, ""},
		// negative
		{"virtualhost 192.168.0.", true, "IP address is not valid 192.168.0"},
		{"virtualhost 192.168.0.100 192.168.0.101", true, "Unable to use more than one IPv4: 192.168.0.100, 192.168.0.101"},
		{"virtualhost fe80::8770:87a8:3d30:", true, "IP address is not valid fe80::8770:87a8:3d30:"},
		{"virtualhost fe80::8770:87a8:3d30:84da fe80::8770:87a8:3d30:84da", true, "Unable to use more than one IPv6: fe80::8770:87a8:3d30:84da, fe80::8770:87a8:3d30:84da"},
		{"virtualhost 192.168.0.100 fe80::8770:87a8:3d30:84da 192.168.0.101", true, "Wrong argument count or unexpected line ending after '192.168.0.101'"},
	}

	for i, test := range tests {
		c := caddy.NewTestController("dns", test.input)
		_, err := parseVirtualHost(c)

		if test.shouldErr && err == nil {
			t.Errorf("Test %d: expected error but found %s for input %s", i, err, test.input)
		}

		if err != nil {
			if !test.shouldErr {
				t.Fatalf("Test %d: expected no error but found one for input %s, got: %v", i, test.input, err)
			}

			if !strings.Contains(err.Error(), test.expectedErr) {
				t.Error(err)
				t.Errorf("Test %d: expected error to contain: %v, found error: %v, input: %s", i, test.expectedErr, err, test.input)
			}
		}
	}
}
