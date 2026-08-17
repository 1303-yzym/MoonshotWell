package udt

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"reflect"
	"strconv"

	"github.com/go-playground/validator/v10"
)

const (
	BoolInvalid int8 = 0
	BoolFalse   int8 = 1
	BoolTrue    int8 = 2
)

func True() Bool  { return InBool(true) }
func False() Bool { return InBool(false) }

func InBool(b bool) Bool {
	v := BoolFalse
	if b {
		v = BoolTrue
	}

	return Bool{val: v}
}

var _ CustomType = (*Bool)(nil)

// Bool 自定义布尔类型 1.false 2.true
type Bool struct {
	val int8
}

func (b Bool) String() string {
	if b.val == BoolTrue {
		return "true"
	} else {
		return "false"
	}
}

func (b Bool) toBool() (bool, error) {
	switch b.val {
	case BoolInvalid:
		return false, nil
	case BoolFalse:
		return false, nil
	case BoolTrue:
		return true, nil
	default:
		return false, errors.New("invalid input for boolean conversion")
	}
}

func (b Bool) ToVal() bool {
	bo, _ := b.toBool()

	return bo
}

func (b Bool) ToNum() int8 {
	return b.val
}

// Scan sql.Scanner
func (b *Bool) Scan(val any) error {
	var vs int8

	switch v := val.(type) {
	case int64:
		vs = int8(v)
	case string:
		if atom, err := strconv.Atoi(v); err != nil {
			return errors.New("types.bool value overflow")
		} else {
			vs = int8(atom)
		}
	case []byte:
		if atom, err := strconv.Atoi(string(v)); err != nil {
			return errors.New("types.bool value overflow")
		} else {
			vs = int8(atom)
		}
	default:
		return errors.New("types.bool type err")
	}

	tmp := Bool{val: vs}

	bo, err := tmp.toBool()
	if err != nil {
		return err
	}

	*b = InBool(bo)

	return nil
}

// Value driver.Valuer
func (b Bool) Value() (driver.Value, error) {
	if b.val == BoolInvalid {
		return int64(BoolFalse), nil
	}

	return int64(b.val), nil
}

// MarshalJSON json.Marshaler
func (b Bool) MarshalJSON() ([]byte, error) {
	if b.val == BoolInvalid {
		return json.Marshal(BoolFalse)
	}

	return json.Marshal(b.val)
}

// UnmarshalJSON json.Unmarshaler
func (b *Bool) UnmarshalJSON(data []byte) (err error) {
	var val int8
	if err = json.Unmarshal(data, &val); err != nil {
		return err
	}

	tmp := Bool{val: val}

	bo, err := tmp.toBool()
	if err != nil {
		return
	}

	*b = InBool(bo)

	return
}

func (b Bool) Validator() CustomValidator {
	return CustomValidator{
		TagName:     "bool",
		Translation: "{0} 值必须为[1 = false || 2 = true]",
		CustomTypeFunc: func(field reflect.Value) interface{} {
			if val, ok := field.Interface().(Bool); ok {
				return val.val
			}

			return int8(0)
		},
		ValidatorFunc: func(ctx context.Context, fl validator.FieldLevel) bool {
			bo, ok := fl.Field().Interface().(int8)
			if !ok {
				return false
			}

			switch bo {
			case 1, 2, 3:
				return true
			default:
				return false
			}
		},
	}
}
