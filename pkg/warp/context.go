package warp

import (
	"context"
)

type Context[C context.Context] interface {
	G() C

	BindRequest(req any) error

	Fail(err error)
	Results(result any)
}
