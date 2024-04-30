package httpclient

type Options struct {
	headers map[string]string
}

func NewOptions() *Options {
	return &Options{}
}

func (o *Options) SetHeader(key string, value string) *Options {
	if o.headers == nil {
		o.headers = make(map[string]string)
	}
	o.headers[key] = value
	return o
}
