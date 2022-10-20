package ipadress

import (
	"github.com/jpillora/ipfilter"
	"testing"
)

func Test1(t *testing.T) {
	f := ipfilter.New(ipfilter.Options{
		BlockedIPs: []string{"152.206.0.0/15"},
	})

	a := f.Blocked("152.206.0.1")
	t.Log(a)
}
