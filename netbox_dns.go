package corednsnetboxplugin

import (
	"context"

	"github.com/ProtonMail/go-netbox-dns/netbox_dns/client"
	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/pkg/fall"
	"github.com/miekg/dns"
)

type NetboxDNS struct {
	Next  plugin.Handler
	Zones []string

	Fall   fall.F
	Client *client.NetBoxAPI
}

func New() *NetboxDNS {
	return &NetboxDNS{}
}

func (h *NetboxDNS) Run(ctx context.Context) error {
	// TODO: implement
	return nil
}

func (h *NetboxDNS) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	// TODO
	return dns.RcodeSuccess, nil
}

func (h *NetboxDNS) Name() string { return "netbox-dns" }
