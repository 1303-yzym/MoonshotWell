package validation

import (
	"github.com/go-playground/validator/v10"
	"github.com/shopspring/decimal"
)

// 自定义错误信息提示.
var validationsMapping = map[string]struct {
	Fn          validator.Func // 验证函数
	translation string         // 默认 {0} field占位符号
	validationsErr
}{
	// 判断是否金额类型 2位小数部位
	"amount": {
		Fn:          Amount,
		translation: "金额必须为['132.22']格式，保留2位小数",
	},
	// 判断是否经纬度类型 6位小数部位
	"geo": {
		Fn:          Geo,
		translation: "经纬度必须为['179.812121']格式，小数不能超过6位",
	},
	// 示例
	"hallo": {
		Fn:          Hallo,
		translation: "{0}值不为hallo",
	},
}

// Amount 判断是否是金额2位小数 decimal.Decimal.
func Amount(fl validator.FieldLevel) bool {
	fieldValue, ok := fl.Field().Interface().(decimal.Decimal)
	if ok {
		// 只能有两位小数部位
		if fieldValue.Equal(fieldValue.Round(2)) {
			return true
		}
	}

	return false
}

// Geo 判断是否是经纬度.
func Geo(fl validator.FieldLevel) bool {
	val, ok := fl.Field().Interface().(decimal.Decimal)
	if !ok {
		return false
	}
	// 只能有6位小数部位
	if val.Equal(val.Round(6)) {
		return true
	}

	return false
}

// Hallo 自定义验证器.
func Hallo(fl validator.FieldLevel) bool {
	if fl.Field().String() == "hallo" {
		// 成功
		return true
	}
	// 失败
	return false
}
