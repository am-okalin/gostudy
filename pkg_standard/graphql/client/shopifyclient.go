package client

import (
	"fmt"
	"io"
	"net/http"
)

type GraphQLReq struct {
	Query     string      `json:"query"`
	Variables interface{} `json:"variables"`
}

type ShopifyClient struct {
	StoreName string
	Token     string
	client    *http.Client
}

func (c ShopifyClient) getUrl() string {
	return fmt.Sprintf("https://%s.myshopify.com/admin/api/2022-07/graphql.json", c.StoreName)
}

func (c ShopifyClient) Do(body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest("POST", c.getUrl(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("X-Shopify-Access-Token", c.Token)
	req.Header.Add("Content-Type", "application/json")

	return http.DefaultClient.Do(req)
}
