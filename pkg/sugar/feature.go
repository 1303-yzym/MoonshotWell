package sugar

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func PrintEsRequest(reqFun func(ctx context.Context) (*http.Request, error)) {
	req, _ := reqFun(context.Background())
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(req.Body)

	body, err := io.ReadAll(req.Body)
	if err != nil {
		fmt.Println("Error reading body:", err)

		return
	}
	// 检查是否是有效的 JSON
	var prettyJSON bytes.Buffer

	err = json.Indent(&prettyJSON, body, "", "  ")
	if err != nil {
		fmt.Println("Error parsing JSON:", err)

		return
	}

	// 格式化输出 JSON
	fmt.Println(prettyJSON.String())
}
