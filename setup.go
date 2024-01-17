package virtualhost

import (
	"fmt"
	"net"

	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
)

func init() {
	plugin.Register("virtualhost", setup)
}

func setup(c *caddy.Controller) error {
	v, err := parseVirtualHost(c)
	if err != nil {
		return plugin.Error("virtualhost", err)
	}

	if err := v.loadHosts(); err != nil {
		return plugin.Error("virtualhost", err)
	}

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		v.Next = next
		return v
	})

	return nil
}

func parseVirtualHost(c *caddy.Controller) (*VirtualHost, error) {
	v := NewVirtualHost()

	for c.Next() {
		args := c.RemainingArgs()

		if len(args) == 0 || len(args) > 2 {
			return v, c.ArgErr()
		}

		err := parseIP(v, args[0])
		if err != nil {
			return v, err
		}

		if len(args) > 1 {
			err := parseIP(v, args[1])
			if err != nil {
				return v, err
			}
		}
	}

	return v, nil
}

func parseIP(v *VirtualHost, s string) error {
	ip := net.ParseIP(s)
	if ip == nil {
		return fmt.Errorf("IP address is not valid %s", s)
	}

	if ip.To4() != nil {
		if v.A != nil {
			return fmt.Errorf("Unable to use more than one IPv4: %s, %s", v.A, ip)
		}
		v.A = ip
	} else {
		if v.AAAA != nil {
			return fmt.Errorf("Unable to use more than one IPv6: %s, %s", v.AAAA, ip)
		}
		v.AAAA = ip
	}

	return nil
}
