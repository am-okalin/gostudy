package net

import (
	"github.com/jpillora/ipfilter"
	"net"
	"testing"
)

func TestIpFilter(t *testing.T) {
	filter := ipfilter.New(ipfilter.Options{
		BlockedIPs: []string{"93.171.14.0/23"},
	})
	flag := filter.Blocked("93.171.14.1")
	t.Log(flag)
}

func TestIP1(t *testing.T) {
	start := "93.171.14.0"
	end := "93.171.72.255"
	si := net.ParseIP(start)
	ei := net.ParseIP(end)
	t.Log(ei, si)
}

func TestIP2(t *testing.T) {
	ip, net, err := net.ParseCIDR("93.171.14.0/23")
	t.Log(ip, net, err)
}
