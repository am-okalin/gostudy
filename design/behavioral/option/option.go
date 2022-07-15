package option

const (
	DefaultTimeout = 2
	DefaultRetries = 3
)

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

//Options 选项内容
type Options struct {
	timeout int //超时时间
	retries int //重试次数
}

//NewOptions 新建选项
func NewOptions() *Options {
	return &Options{
		timeout: DefaultTimeout,
		retries: DefaultRetries,
	}
}

func WithRetries(retries int) Option {
	return func(o *Options) {
		o.retries = retries
	}
}

func WithTimeout(timeout int) Option {
	return func(o *Options) {
		o.timeout = timeout
	}
}
