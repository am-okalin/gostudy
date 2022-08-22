package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type CreateSPGRes struct {
	Data struct {
		Products struct {
			Edges []struct {
				Node struct {
					Id string `json:"id"`
				} `json:"node"`
			} `json:"edges"`
		} `json:"products"`
	} `json:"data"`
	Extensions struct {
		Cost struct {
			RequestedQueryCost int `json:"requestedQueryCost"`
			ActualQueryCost    int `json:"actualQueryCost"`
			ThrottleStatus     struct {
				MaximumAvailable   float64 `json:"maximumAvailable"`
				CurrentlyAvailable int     `json:"currentlyAvailable"`
				RestoreRate        float64 `json:"restoreRate"`
			} `json:"throttleStatus"`
		} `json:"cost"`
	} `json:"extensions"`
}

type CreateSPG struct {
	ShopifyClient
	reqObj GraphQLReq
	resObj CreateSPGRes
}

func NewCreateSPG(shopifyClient ShopifyClient) *CreateSPG {
	reqObj := GraphQLReq{
		Query:     fmt.Sprintf("{products(first:%d){edges{node{id}}}}", 5),
		Variables: nil,
	}

	return &CreateSPG{ShopifyClient: shopifyClient, reqObj: reqObj}
}

func (a CreateSPG) request() error {
	reqBytes, err := json.Marshal(a.reqObj)
	if err != nil {
		return err
	}

	res, err := a.Do(bytes.NewReader(reqBytes))
	if err != nil {
		return err
	}

	resBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(resBytes, &a.resObj)
	if err != nil {
		return err
	}

	return nil
}
