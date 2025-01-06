package client

import (
	"net/http"
	"testing"
)

func Test1(t *testing.T) {
	sc := ShopifyClient{
		StoreName: "yuan-dev",
		Token:     "shpat_*",
		client:    &http.Client{},
	}
	products := NewPageProducts(sc)
	err := products.request()
	t.Log(err)
}
