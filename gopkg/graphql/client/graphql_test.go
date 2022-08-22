package client

import (
	"net/http"
	"testing"
)

func Test1(t *testing.T) {
	sc := ShopifyClient{
		StoreName: "redmagicdev",
		Token:     "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",
		client:    &http.Client{},
	}
	products := NewPageProducts(sc)
	err := products.request()
	t.Log(err)
}
