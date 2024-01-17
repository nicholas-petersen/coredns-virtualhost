package virtualhost

import (
	"context"
	"net"
	"strings"

	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/request"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/miekg/dns"
)

type VirtualHost struct {
	Next  plugin.Handler
	A     net.IP
	AAAA  net.IP
	Hosts map[string]struct{}
}

func NewVirtualHost() *VirtualHost {
	v := &VirtualHost{
		Hosts: map[string]struct{}{},
	}
	return v
}

func (v VirtualHost) Name() string {
	return "virtualhost"
}

func (v VirtualHost) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	state := request.Request{W: w, Req: r}
	qname := state.Name()

	answers := []dns.RR{}

	if _, ok := v.Hosts[strings.TrimSuffix(qname, ".")]; ok {
		//  Ttl zero; indicating these records should not be cached
		if v.A != nil {
			answers = append(answers, &dns.A{
				Hdr: dns.RR_Header{
					Name:   qname,
					Rrtype: dns.TypeA,
					Ttl:    0,
					Class:  dns.ClassINET,
				},
				A: v.A,
			})
		}

		if v.AAAA != nil {
			answers = append(answers, &dns.AAAA{
				Hdr: dns.RR_Header{
					Name:   qname,
					Rrtype: dns.TypeAAAA,
					Ttl:    0,
					Class:  dns.ClassINET,
				},
				AAAA: v.AAAA,
			})
		}

		hostnameCount.WithLabelValues(qname).Add(1)
		m := new(dns.Msg)
		m.SetReply(r)
		m.Answer = answers

		w.WriteMsg(m)
		return dns.RcodeSuccess, nil
	}

	return plugin.NextOrFailure(v.Name(), v.Next, ctx, w, r)
}

func (v VirtualHost) loadHosts() error {
	ctx := context.Background()

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}

	defer cli.Close()

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		return err
	}

	for _, container := range containers {
		con, err := cli.ContainerInspect(ctx, container.ID)
		if err != nil {
			return err
		}

		for env := range con.Config.Env {
			if strings.Split(con.Config.Env[env], "=")[0] == "VIRTUAL_HOST" {
				host := strings.Split(con.Config.Env[env], "=")[1]
				v.Hosts[strings.TrimSpace(host)] = struct{}{}
			}
		}
	}

	return nil
}
