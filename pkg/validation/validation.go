// Package validation 验证器
package validation

import (
	"log"

	"github.com/1303-yzym/MoonshotWell/pkg/udt"
	"go.uber.org/zap"

	"github.com/go-playground/validator/v10"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
)

// https://github.com/go-playground/validator/blob/master/_examples/translations/main.go
// 自定义字典和错误信息.
type validationsErr struct {
	customRegisFunc validator.RegisterTranslationsFunc
	customTransFunc validator.TranslationFunc
}

var uni = ut.New(en.New())

// RegisterErr 收集注册错误.
func RegisterErr(err error) {
	if err != nil {
		zap.L().Named("runtime").Panic("validator " + err.Error())
	}
}

// RegisterValidations 注册验证器.
func RegisterValidations() {
	trans, _ := uni.GetTranslator("zh_cn")

	v, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		zap.L().Named("runtime").Panic("register validation  err")
	}

	// 验证器tag默认为binding
	v.SetTagName("v")

	// 注册tag验证器
	for tag, val := range validationsMapping {
		if val.Fn != nil {
			RegisterErr(v.RegisterValidation(tag, val.Fn))
		}

		RegisterErr(registerTranslations(v, tag, val.translation, trans, val.validationsErr))
	}
	// 自定义类型错误
	if err := registerCustomValidator(v, udt.Bool{}); err != nil {
		RegisterErr(err)
	}
}

func registerCustomValidator(v *validator.Validate, cs ...udt.CustomTypeValidator) error {
	trans, _ := uni.GetTranslator("zh_cn")
	for _, c := range cs {
		ver := c.Validator()
		v.RegisterCustomTypeFunc(ver.CustomTypeFunc, c)

		if err := v.RegisterValidationCtx(ver.TagName, ver.ValidatorFunc); err != nil {
			return err
		}

		if err := registerTranslations(v, ver.TagName, ver.Translation, trans, validationsErr{}); err != nil {
			return err
		}
	}

	return nil
}

func registerTranslations(v *validator.Validate, tag, translation string, trans ut.Translator, t validationsErr) (err error) {
	if t.customTransFunc != nil && t.customRegisFunc != nil {
		err = v.RegisterTranslation(tag, trans, t.customRegisFunc, t.customTransFunc)
	} else if t.customTransFunc != nil && t.customRegisFunc == nil {
		err = v.RegisterTranslation(tag, trans, registrationFunc(tag, translation, true), t.customTransFunc)
	} else if t.customTransFunc == nil && t.customRegisFunc != nil {
		err = v.RegisterTranslation(tag, trans, t.customRegisFunc, translateFunc)
	} else {
		err = v.RegisterTranslation(tag, trans, registrationFunc(tag, translation, true), translateFunc)
	}

	return
}

func registrationFunc(tag string, translation string, override bool) validator.RegisterTranslationsFunc {
	return func(ut ut.Translator) (err error) {
		if err = ut.Add(tag, translation, override); err != nil {
			return
		}

		return
	}
}

func translateFunc(ut ut.Translator, fe validator.FieldError) string {
	t, err := ut.T(fe.Tag(), fe.StructField())
	if err != nil {
		log.Printf("警告: 翻译字段错误: %#v", fe)

		return fe.(error).Error()
	}

	return t
}
