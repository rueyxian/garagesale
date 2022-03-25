package web

// Middleware
type Middleware func(HandlerFunc) HandlerFunc

// ChainMiddleware
func WrapMiddleware(mws []Middleware, fn HandlerFunc) HandlerFunc {
	for i := len(mws) - 1; i >= 0; i-- {
		mw := mws[i]
		fn = mw(fn)
	}
	return fn
}
