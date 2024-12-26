package net

import (
	"github.com/jpillora/ipfilter"
	"testing"
)

func TestIpFilter(t *testing.T) {
	filter := ipfilter.New(ipfilter.Options{
		BlockedIPs: []string{"93.171.14.0/23"},
	})
	flag := filter.Blocked("93.171.14.1")
	t.Log(flag)
}
