package httpclient

import (
	"crypto/tls"
	"io"
	"net/http"
)

const (
	MethodGet  = "get"
	MethodPost = "post"
)

type ClientInterface interface {
	//Path 获取请求url
	Path() string
	//RequestBody 构建请求体
	RequestBody() (io.Reader, error)
	//ResponseBody 构建返回体
	ResponseBody(io.Reader) (interface{}, error)
}

type ClientParams struct {
	Domain string
	Api    string
	Method string
	opts   *Options
}

type Api struct {
	request  *http.Request
	response *http.Response
	client   http.Client
	ClientInterface
	ClientParams
}

func NewClient(client ClientInterface, domain, api, method string, options ...Option) *Api {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			// 跳过客户端 HTTPS 验证
			InsecureSkipVerify: true,
		},
	}
	params := ClientParams{
		Domain: domain,
		Api:    api,
		Method: method,
		opts:   loadOptions(options...),
	}
	return &Api{
		request:         nil,
		response:        nil,
		client:          http.Client{Transport: transport},
		ClientInterface: client,
		ClientParams:    params,
	}
}

func (a *Api) Do() error {
	body, err := a.RequestBody()
	if err != nil {
		return err
	}

	request, err := http.NewRequest(MethodPost, a.Path(), body)
	if err != nil {
		return err
	}

	for key, value := range a.opts.headers {
		request.Header.Add(key, value)
	}

	if a.opts.basicAuth != nil {
		request.SetBasicAuth(a.opts.basicAuth.username, a.opts.basicAuth.password)
	}

	reponse, err := a.client.Do(request)
	if err != nil {
		return err
	}

	a.request = request
	a.response = reponse

	return nil
}

func (a *Api) Request() *http.Request {
	if a.request == nil {
		return nil
	}
	return a.request
}

func (a *Api) Response() *http.Response {
	if a.response == nil {
		return nil
	}
	return a.response
}
