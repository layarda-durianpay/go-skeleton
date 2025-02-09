package kafka

type configReader struct {
	beforeFunc BeforeFunc
	afterFunc  AfterFunc
}

type Option interface {
	apply(c *configReader)
}

type optionFunc func(*configReader)

func (o optionFunc) apply(c *configReader) {
	o(c)
}

func newConfig(opts ...Option) *configReader {
	conf := &configReader{}

	for _, opt := range opts {
		opt.apply(conf)
	}

	return conf
}

func WithAfterFunc(af AfterFunc) Option {
	return optionFunc(func(cr *configReader) {
		cr.afterFunc = af
	})
}

func WithBeforeFunc(bf BeforeFunc) Option {
	return optionFunc(func(cr *configReader) {
		cr.beforeFunc = bf
	})
}
