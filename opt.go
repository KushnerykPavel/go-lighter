package lighter

type Opt[T any] func(*T)

func (o Opt[T]) Apply(opt *T) {
	o(opt)
}

type (
	WsOpt = Opt[WebsocketClient]
)

func WithOptDebugMode(l logger) WsOpt {
	return func(w *WebsocketClient) {
		w.debug = true
		w.logger = l
	}
}
