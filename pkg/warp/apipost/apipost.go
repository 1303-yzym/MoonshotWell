package apipost

import (
	"fmt"
	"io"
	"log"
	"reflect"
	"runtime"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"

	apipostGen "etopiacity.com/common/pkgs/oas/apipost/gen"
	"etopiacity.com/common/pkgs/warp"
)

type ApiPost struct {
	generator *ApiDocGenerator
	routers   []*warp.RouteInfo
}

func NewApiPost(routes []*warp.RouteInfo) (*ApiPost, error) {
	generator, err := NewApiDocGenerator()
	if err != nil {
		return nil, err
	}

	return &ApiPost{
		generator: generator,
		routers:   routes,
	}, nil
}

func (h *ApiPost) WriteApiTable(w io.Writer) error {
	t := table.NewWriter()
	t.AppendHeader(table.Row{"method", "path", "handlers", "middlewares", "comment"})
	t.SetColumnConfigs([]table.ColumnConfig{
		{
			Name:     "path",
			Align:    text.AlignLeft,
			Colors:   text.Colors{text.FgHiBlue},
			WidthMin: 10,
		},
		{
			Name:  "comment",
			Align: text.AlignLeft,
		},
	})

	for _, info := range h.routers {
		var handlerName string

		if info.HandlerValue != nil {
			name := runtime.FuncForPC(info.HandlerValue.Pointer()).Name()

			handlerNames := strings.Split(name, "/")
			if len(handlerNames) > 0 {
				handlerName = handlerNames[len(handlerNames)-1]
			}
		}

		t.AppendRow(table.Row{
			info.Method,
			info.Path,
			handlerName,
			strings.Join(info.Middlewares, ","),
			info.Comment,
		})
	}

	_, err := io.WriteString(w, t.Render()+"\n")
	if err != nil {
		return err
	}

	_, err = io.WriteString(w, fmt.Sprintf("total http api: %d \n", len(h.routers)))
	if err != nil {
		return err
	}

	return nil
}

func (h *ApiPost) GenApiJsonSchema(dirPath string, urlPrefix string) {
	for _, route := range h.routers {
		if strings.Contains(route.Path, ":") {
			continue
		}

		if route.HandlerValue != nil {
			funcType := route.HandlerValue.Type()
			if funcType.Kind() == reflect.Func {
				// 参数类型
				QType := funcType.In(1) // 第二个参数是 Q
				// 返回值类型
				RType := funcType.Out(0) // 第一个返回值是 R

				requestSchema, err := apipostGen.GenFromType(QType)
				if err != nil {
					log.Fatal(err)
				}

				responseSchema, err := apipostGen.GenFromType(newResponse(RType))
				if err != nil {
					log.Fatal(err)
				}

				if err = h.generator.Execute(dirPath, route.Path, &ApiSchema{
					UrlPath:        route.Path,
					FullUrlPath:    urlPrefix + route.Path,
					Name:           route.Comment,
					RequestSchema:  string(requestSchema),
					ResponseSchema: string(responseSchema),
				}); err != nil {
					log.Fatal(err)
				}
			}
		}
	}
}
