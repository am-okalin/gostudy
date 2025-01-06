package shurcool

import (
	"context"
	"fmt"
	"github.com/shurcooL/graphql"
	"testing"
	"time"
)

var ctx = context.Background()
var shopName = "yuan-dev"
var token = "shpat_*"
var url = fmt.Sprintf("https://%s.myshopify.com/admin/api/2022-07/graphql.json", shopName)

type Product struct {
	Id              string    `json:"id"`
	Title           string    `json:"title"`
	DescriptionHtml string    `json:"descriptionHtml"`
	Vendor          string    `json:"vendor"`
	ProductType     string    `json:"productType"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
	Variants        struct {
		Edges []struct {
			Node struct {
				Id               string `json:"id"`
				Title            string `json:"title"`
				Price            string `json:"price"`
				Sku              string `json:"sku"`
				AvailableForSale bool   `json:"availableForSale"`
			} `json:"node"`
		} `json:"edges"`
	} `json:"variants" graphql:"variants(first: 5)"`
	Images struct {
		Edges []struct {
			Node struct {
				Id          string `json:"id"`
				OriginalSrc string `json:"originalSrc"`
				AltText     string `json:"altText"`
			} `json:"node"`
		} `json:"edges"`
	} `json:"images" graphql:"images(first: 5)"`
}

func Test1(t *testing.T) {
	//httpClient := oauth2.NewClient(ctx, oauth2.StaticTokenSource(
	//	&oauth2.Token{AccessToken: token},
	//))

	httpClient := NewHTTPClientWithHeaders(nil, map[string]string{
		"x-shopify-access-token": token,
	})

	client := graphql.NewClient(url, httpClient)

	var query struct {
		Products struct {
			Edges []struct {
				Node Product `json:"node"`
			} `json:"edges"`
		} `json:"products" graphql:"products(first: 5)"`
	}

	err := client.Query(context.Background(), &query, nil)
	if err != nil {
		t.Error(err)
	}
	t.Log(query)
}
