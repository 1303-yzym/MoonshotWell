package udt

import (
	"database/sql/driver"
	// "reflect"

	"github.com/go-playground/validator/v10"
	// "etopiacity.com/common/pkgs/oas/apipost"
)

// CustomType 自定义类型的接口
type CustomType interface {
	CustomTypeValidator

	String() string

	// gorm
	Scan(val any) error
	Value() (driver.Value, error)

	// json
	MarshalJSON() ([]byte, error)
	UnmarshalJSON(data []byte) (err error)
}

type CustomTypeValidator interface {
	Validator() CustomValidator
	// JsonSchema(_ reflect.Type) *apipost.Properties
}

type CustomValidator struct {
	TagName        string `comment:"标签名称"`
	Translation    string `comment:"错误信息"`
	CustomTypeFunc validator.CustomTypeFunc
	ValidatorFunc  validator.FuncCtx
}
