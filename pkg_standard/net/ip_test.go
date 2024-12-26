package net

import (
	"net"
	"testing"
)

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
