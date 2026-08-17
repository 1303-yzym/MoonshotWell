package udt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBool_Constant(t *testing.T) {
	c := True()
	c.val = 1
	assert.Equal(t, int8(2), True().val)

	f := False()
	f.val = 2
	assert.Equal(t, int8(1), False().val)
}

func TestBool_Conversion(t *testing.T) {
	bo := new(Bool)

	err := bo.Scan(int64(2))
	assert.NoError(t, err)

	err = bo.Scan([]byte("2"))
	assert.NoError(t, err)

	err = bo.Scan("2")
	assert.NoError(t, err)

	err = bo.Scan("a")
	assert.Errorf(t, err, "types.bool value overflow")

	err = bo.Scan(int16(2))
	assert.Errorf(t, err, "types.bool type err")

	tr := True()
	value, err := tr.Value()
	assert.NoError(t, err)
	assert.Equal(t, int64(2), value)

	jsonText, err := tr.MarshalJSON()
	assert.NoError(t, err)
	assert.Equal(t, "2", string(jsonText))

	jsonBo := new(Bool)
	err = jsonBo.UnmarshalJSON(jsonText)
	assert.NoError(t, err)
	assert.Equal(t, True(), *jsonBo)

	err = jsonBo.UnmarshalJSON([]byte("3"))
	assert.Error(t, err)

	err = jsonBo.UnmarshalJSON([]byte("0"))
	assert.NoError(t, err)
	assert.Equal(t, False(), *jsonBo)
}

func TestBool_Validator(t *testing.T) {
	err := Validator.Var(True(), "bool")
	assert.NoError(t, err)

	err = Validator.Var(False(), "bool")
	assert.NoError(t, err)

	err = Validator.Var(Bool{val: 8}, "bool")
	assert.Error(t, err)
}
