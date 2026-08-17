package code

import (
	"errors"
	"fmt"

	"github.com/1303-yzym/MoonshotWell/pkg/i18n"
)

type Code struct {
	HCode int    `json:"h_code" comment:"http的code状态"`
	Code  int    `json:"code" comment:"code"`
	Msg   string `json:"msg" comment:"msg"` // 默认消息,无国际化展示默认
	Err   any    `json:"err" comment:"err"`
}

func New(code int, msg string) Code {
	return Code{
		Code: code,
		Msg:  msg,
	}
}

// MsgError 自定义错误消息函数.
func MsgError(msg string) Code {
	return Code{
		Code: 999,
		Msg:  msg,
	}
}

// 实现 Stringer.
func (c Code) String() string {
	return fmt.Sprintf("HCode=%d,Code=%d,Msg=%s,Err=%v", c.HCode, c.Code, c.Msg, c.Err)
}

// 实现 error.
func (c Code) Error() string {
	return c.String()
}

func (c Code) GetCode() int {
	return c.Code
}

func (c Code) Decode(local i18n.I18n, lang ...string) (int, int, string, any) {
	if msg, err := local.T(c.Code, lang...); err == nil {
		return c.HCode, c.Code, msg, c.Err
	}

	return c.HCode, c.Code, c.Msg, c.Err
}

func (c Code) SetCode(code int) Code {
	c.Code = code

	return c
}

func (c Code) SetHttpCode(httpCode int) Code {
	c.HCode = httpCode

	return c
}

func (c Code) SetMsg(message string) Code {
	c.Msg = message

	return c
}

func (c Code) SetError(err any) Code {
	c.Err = err

	return c
}

// JoinError 外层错误的code msg为准.
func (c Code) joinError(err error) Code {
	var codeErr Code
	if errors.As(err, &codeErr) {
		if c.Code == codeErr.Code {
			return codeErr
		} else {
			c.Err = errors.Join(c, codeErr)

			return c
		}
	} else {
		c.Err = errors.Join(c, err)
	}

	return c
}

func (c Code) JoinErrors(errs ...error) Code {
	result := c

	for _, err := range errs {
		result = result.joinError(err)
	}

	return result
}

func (c Code) JoinErrorMsg(format string, a ...any) Code {
	var newErr error
	if len(a) == 0 {
		newErr = errors.New(format)
	} else {
		newErr = fmt.Errorf(format, a...)
	}

	c.Err = errors.Join(c, newErr)

	return c
}
