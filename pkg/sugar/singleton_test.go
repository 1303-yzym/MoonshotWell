package sugar

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type DC struct {
	Name string
}

func (d *DC) SetName(name string) {
	d.Name = name
}

var myObjectSingleton = NewSingleton[*DC](func() *DC {
	return &DC{Name: "SingletonInstance"}
})

func TestSingleton(t *testing.T) {
	e := "singleton"
	myObjectSingleton.Load().SetName(e)
	assert.Equal(t, myObjectSingleton.Load().Name, e)
}
