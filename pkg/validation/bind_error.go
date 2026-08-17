package validation

import (
	"errors"
	"fmt"

	"github.com/duke-git/lancet/v2/strutil"
	"github.com/go-playground/validator/v10"
)

func ParseValidationErrors(err error) any {
	var errs validator.ValidationErrors
	if errors.As(err, &errs) {
		var errMses []string
		// trans, _ := uni.GetTranslator("zh_cn")
		for i := 0; i < len(errs); i++ {
			errI := errs[i]

			var errMsg string
			// "Field" => Param=,Msg=
			errMsg = fmtI("Field", strutil.SnakeCase(errI.Field())) +
				fmtI("Param", fmt.Sprintf("%v", errI.Value())) +
				// fmtI("Msg", errI.Translate(trans)) +
				fmtI("Struct", errI.StructField()) +
				fmtI("VTag", errI.ActualTag(), true)
			errMses = append(errMses, errMsg)
		}

		return errMses
	}

	return err.Error()
}

// Fmt tag=str.
func fmtI(tag, str string, end ...bool) string {
	if len(end) > 0 && end[0] {
		return tag + "=" + str
	} else {
		return tag + "=" + str + ", "
	}
}
