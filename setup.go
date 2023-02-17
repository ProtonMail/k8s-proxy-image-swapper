package corednsnetboxplugin

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"

	netboxdnsclient "github.com/ProtonMail/go-netbox-dns/netbox_dns/client"
	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/goware/urlx"
)

func init() { plugin.Register("netbox-dns", setup) }

func setup(c *caddy.Controller) error {
	netboxDns, err := parseNetboxDNS(c)
	if err != nil {
		return plugin.Error("netbox-dns", err)
	}

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		netboxDns.Next = next
		return netboxDns
	})
	ctx, cancel := context.WithCancel(context.Background())
	if err := netboxDns.Run(ctx); err != nil {
		cancel()
		return plugin.Error("netbox-dns", err)
	}
	c.OnShutdown(func() error { cancel(); return nil })
	return nil
}

func parseNetboxDNS(c *caddy.Controller) (*NetboxDNS, error) {
	netboxDNS := New()

	for c.Next() {
		// netbox-dns [zones..]
		args := c.RemainingArgs()
		netboxDNS.Zones = plugin.OriginsFromArgsOrServerBlock(args, c.ServerBlockKeys)

		token := os.Getenv("NETBOX_TOKEN")
		var err error
                var parsedURL *url.URL = nil

		for c.NextBlock() {
			switch c.Val() {
			case "fallthrough":
				netboxDNS.Fall.SetZonesFromArgs(c.RemainingArgs())

			case "url":
				if !c.NextArg() {
					return nil, c.ArgErr()
				}
				parsedURL, err = urlx.Parse(c.Val())
				if err != nil {
					return nil, err
				}
			default:
				return nil, c.Errf("unknown property '%s'", c.Val())
			}
		}

		if parsedURL == nil {
			return nil, c.Errf("url not specified")
		}

		desiredRuntimeClientSchemes := []string{parsedURL.Scheme}
		transport := httptransport.NewWithClient(parsedURL.Host, parsedURL.Path+netboxdnsclient.DefaultBasePath, desiredRuntimeClientSchemes, &http.Client{})
		if token != "" {
			transport.DefaultAuthentication = httptransport.APIKeyAuth("Authorization", "header", fmt.Sprintf("Token %v", token))
		}
		netboxDNS.Client = netboxdnsclient.New(transport, nil)

	}

	return netboxDNS, nil
}
