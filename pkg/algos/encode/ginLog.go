package encode

import (
	"fmt"
	"net/http"
)

const (
	green   = "\033[97;42m"
	white   = "\033[90;47m"
	yellow  = "\033[37;43m"
	red     = "\033[97;41m"
	blue    = "\033[97;44m"
	magenta = "\033[97;45m"
	cyan    = "\033[97;46m"
	reset   = "\033[0m"
)

type GinMsg struct {
	Status int
	Proto  string
	Host   string
	Method string
	Path   string
	Query  string
}

// StatusCodeColor code状态的颜色
func (msg GinMsg) StatusCodeColor() string {
	code := msg.Status

	var StatusColor string

	switch {
	case code >= http.StatusOK && code < http.StatusMultipleChoices:
		StatusColor = green
	case code >= http.StatusMultipleChoices && code < http.StatusBadRequest:
		StatusColor = white
	case code >= http.StatusBadRequest && code < http.StatusInternalServerError:
		StatusColor = yellow
	default:
		StatusColor = red
	}

	return fmt.Sprintf("%v %d %v", StatusColor, msg.Status, reset)
}

// MethodColor 请求方式的颜色
func (msg GinMsg) MethodColor() string {
	method := msg.Method

	var StatusColor string

	switch method {
	case http.MethodGet:
		StatusColor = blue
	case http.MethodPost:
		StatusColor = cyan
	case http.MethodPut:
		StatusColor = yellow
	case http.MethodDelete:
		StatusColor = red
	case http.MethodPatch:
		StatusColor = green
	case http.MethodHead:
		StatusColor = magenta
	case http.MethodOptions:
		StatusColor = white
	default:
		StatusColor = reset
	}

	return fmt.Sprintf("%v %v     %v", StatusColor, msg.Method, reset)
}

// FormatGinLog GIN日记输出格式化
func FormatGinLog(ginMsg GinMsg) string {
	msg := fmt.Sprintf(
		"|%s|%13s|    %s | %s \"%s\" -> %s",
		ginMsg.StatusCodeColor(),
		ginMsg.Proto,
		ginMsg.Host,
		ginMsg.MethodColor(),
		ginMsg.Path,
		ginMsg.Query,
	)

	return msg
}
