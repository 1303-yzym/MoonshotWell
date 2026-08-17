package udt

import (
	"reflect"
	"strings"
	"testing"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

var Validator *validator.Validate

func TestMain(t *testing.M) {
	Validator = validator.New(validator.WithRequiredStructEnabled())
	Validator.RegisterTagNameFunc(func(fld reflect.StructField) string {
		if name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]; name != "-" {
			return name
		}
		return ""
	})

	if err := loadCustomValidator(Validator,
		Bool{},
	); err != nil {
		panic(err)
	}

	t.Run()
}

// loadCustomValidator make it possible https://github.com/go-playground/validator/issues/312
// https://github.com/go-playground/validator/issues/470#issuecomment-487990618
func loadCustomValidator(v *validator.Validate, cs ...CustomTypeValidator) error {
	ent := en.New()
	uni := ut.New(ent, ent)
	trans, _ := uni.GetTranslator("en")
	for _, c := range cs {
		ver := c.Validator()

		v.RegisterCustomTypeFunc(ver.CustomTypeFunc, c)

		if err := v.RegisterValidationCtx(ver.TagName, ver.ValidatorFunc); err != nil {
			return err
		}

		if err := v.RegisterTranslation(ver.TagName, trans,
			func(ut ut.Translator) error {
				return ut.Add(ver.TagName, ver.Translation, true)
			},
			func(ut ut.Translator, fe validator.FieldError) string {
				msg, _ := ut.T(ver.TagName, fe.Field())
				return msg
			},
		); err != nil {
			return err
		}
	}
	return nil
}
