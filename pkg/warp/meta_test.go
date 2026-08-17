package warp

import (
	"testing"
)

type (
	Req struct {
		Meta `path:"hallo" method:"post" comment:"request"`
	}
)

func TestMeta(t *testing.T) {
	var req Req
	_, err := parseMeta(req)
	if err != nil {
		return
	}
}

func BenchmarkMeta(b *testing.B) {
	var req *Req
	for b.Loop() {
		_, _ = parseMeta(req)
	}
}
