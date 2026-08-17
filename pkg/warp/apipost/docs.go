package apipost

import (
	_ "embed"
	"os"
	"path/filepath"
	"reflect"
	"text/template"
)

// 使用反射创建一个新的 Response[T]
func newResponse(dataType reflect.Type) reflect.Type {
	responseType := reflect.StructOf([]reflect.StructField{
		{
			Name: "Code",
			Type: reflect.TypeOf(0),
			Tag:  `json:"code" comment:"状态码" mock:"2000"`,
		},
		{
			Name: "Data",
			Type: dataType,
			Tag:  `json:"data" comment:"data"`,
		},
	})

	return responseType
}

//go:embed docs.tmpl
var DocsModelTmpl string

type ApiSchema struct {
	UrlPath        string `json:"url_path" comment:"请求地址"`
	FullUrlPath    string `json:"full_url_path" comment:"带环境请求地址"`
	Name           string `json:"name" comment:"名称"`
	RequestSchema  string `json:"request_schema" comment:"请求结构"`
	ResponseSchema string `json:"response_schema" comment:"响应结构"`
}

type ApiDocGenerator struct {
	tml *template.Template
}

func NewApiDocGenerator() (*ApiDocGenerator, error) {
	tml, err := template.New("docs.tmpl").Parse(DocsModelTmpl)
	if err != nil {
		return nil, err
	}

	return &ApiDocGenerator{tml: tml}, nil
}

func (g *ApiDocGenerator) Execute(dirPath string, path string, apiSchema *ApiSchema) error {
	filePath := filepath.Join(dirPath, path) + ".md"
	_ = os.MkdirAll(filepath.Dir(filePath), os.ModePerm)

	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	if err = g.tml.Execute(file, apiSchema); err != nil {
		return err
	}

	return nil
}
