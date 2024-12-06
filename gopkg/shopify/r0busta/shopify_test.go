package r0busta

import (
	"context"
	"fmt"
	"github.com/r0busta/go-shopify-graphql/v8"
	"testing"
)

const (
	maxRetries = 3
	shopName   = "yuan-dev"
	token      = "shpat_d266"
)

var ctx = context.Background()
var client = shopify.NewClientWithToken(token, shopName)

func TestListProduct(t *testing.T) {
	// Get all collections
	collections, err := client.Collection.ListAll(ctx)
	t.Log(collections, err)

	// Print out the result
	for _, c := range collections {
		fmt.Println(c.Handle)
	}

	products, err := client.Product.ListAll(ctx)
	t.Log(products, err)
}

func TestBilling(t *testing.T) {
	// 返回gql客户端, 自己调用
	client.GraphQLClient()
}
