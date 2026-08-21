package snowflake_test

import (
	"testing"

	"github.com/1303-yzym/MoonshotWell/pkg/algos/id/snowflake"
	"github.com/stretchr/testify/assert"
)

func TestSnowflake(t *testing.T) {
	ss := snowflake.New(1023)
	t.Log(ss.GenerateID())

	assert.Panics(t, func() {
		snowflake.New(1024)
	})
}

func BenchmarkSnowflake(b *testing.B) {
	ss := snowflake.New(21)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ss.GenerateID()
	}
}
