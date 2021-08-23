package rules

import (
	"context"
	"net"

	"socks-proxy/socks5"
)

type All []socks5.RuleSet

func (s *All) Allow(ctx context.Context, req *socks5.Request) (context.Context, bool) {
	for _, ruleSet := range *s {
		if ctx, allow := ruleSet.Allow(ctx, req); !allow {
			return ctx, false
		}
	}
	return ctx, true
}

//

type BlockDestNets struct {
	Nets []net.IPNet
}

func (b *BlockDestNets) Allow(ctx context.Context, req *socks5.Request) (context.Context, bool) {
	for _, net := range b.Nets {
		if net.Contains(req.DestAddr.IP) {
			return ctx, false
		}
	}
	return ctx, true
}
