package commerce

import (
	"context"
	goshopify "github.com/bold-commerce/go-shopify/v4"
	"testing"
)

const (
	testApiVersion = "9999-99"
	maxRetries     = 3

	shopName = "yuan-dev"
	token    = "shpat_d266"
)

var (
	client *goshopify.Client
	app    goshopify.App
	ctx    = context.TODO()
)

func Test1(t *testing.T) {

	app = goshopify.App{
		ApiKey:    "apikey",
		ApiSecret: "hush",
	}
	client = goshopify.MustNewClient(app, shopName, token)

	// Fetch the number of products.
	numProducts, err := client.Product.Count(ctx, nil)
	t.Log(numProducts, err)
}
