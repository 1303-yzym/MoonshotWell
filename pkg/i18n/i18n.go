package i18n

import (
	_ "embed"
	"errors"

	"github.com/1303-yzym/MoonshotWell/pkg/sugar"
)

// I18n 国际化接口
type I18n interface {
	T(code int, selectLang ...string) (string, error)
}

const (
	Default = ZH
	ZH      = "zh_cn"
	EN      = "en_us"
)

var (
	NotLocal = errors.New("no local")
	NotMsg   = errors.New("no msg")
)

var M = sugar.NewSingleton[I18n](func() I18n {
	return newI18n(Default)
})

type i18n struct {
	DefaultLang string
}

func newI18n(defaultLang string) *i18n {
	return &i18n{DefaultLang: defaultLang}
}

func (l *i18n) T(code int, selectLang ...string) (string, error) {
	var lang string
	if len(selectLang) < 1 {
		lang = Default
	} else {
		lang = selectLang[0]
	}

	switch lang {
	case ZH:
		if s, ok := ZhCn[code]; ok {
			return s, nil
		}
	case EN:
		if s, ok := EnUs[code]; ok {
			return s, nil
		}
	default:
		return "", NotLocal
	}

	return "", NotMsg
}
