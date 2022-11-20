package testclient

import (
	"io"
	"kit/httpclient"
)

type CampaignGetUrlParam struct {
}

type CampaignGetReq struct {
}
type CampaignGetRes struct {
}
type CampaignGetApi struct {
	reqUrlParam CampaignGetUrlParam
	reqBody     CampaignGetReq
	resBody     CampaignGetRes
	httpclient.ClientParams
}

func NewCampaignGetApi(reqUrlParam CampaignGetUrlParam) *CampaignGetApi {
	return &CampaignGetApi{reqUrlParam: reqUrlParam}
}

func (c CampaignGetApi) Path() string {
	return ""
}

func (c CampaignGetApi) RequestBody() (io.Reader, error) {
	//TODO implement me
	panic("implement me")
}

func (c CampaignGetApi) ResponseBody(reader io.Reader) (interface{}, error) {
	//TODO implement me
	panic("implement me")
}
