package tracing

type options struct {
	opNameFunc OpNameFunc
}

type Option func(*options)

type OpNameFunc func(method string) string

// WithOpName customizes the trace Operation name
func WithOpName(f OpNameFunc) Option {
	return func(o *options) {
		o.opNameFunc = f
	}
}
