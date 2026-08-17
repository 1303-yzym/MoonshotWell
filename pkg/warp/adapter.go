package warp

import (
	"context"
	"net/http"
)

type Adapter[O context.Context, C Context[O]] interface {
	WrapContext(original O) C
	RegisterHandler(route *Route[O, C])
	Handler() http.Handler
}
