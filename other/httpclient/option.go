package httpclient

//Option 选项方法
type Option func(*Options)

//loadOptions 加载选项
func loadOptions(options ...Option) *Options {
	opts := NewOptions()
	for _, option := range options {
		option(opts)
	}
	return opts
}

type basicAuth struct {
	username string
	password string
}

//Options 选项内容
type Options struct {
	headers   map[string]string
	basicAuth *basicAuth
}

//NewOptions 新建选项
func NewOptions() *Options {
	return &Options{}
}
func WithBasicAuth(username, password string) Option {
	return func(o *Options) {
		o.basicAuth = &basicAuth{
			username: username,
			password: password,
		}
	}
}
func WithContentType(contentType string) Option {
	return func(o *Options) {
		o.headers["Content-Type"] = contentType
	}
}
func WithHeader(h map[string]string) Option {
	return func(o *Options) {
		o.headers = h
	}
}
