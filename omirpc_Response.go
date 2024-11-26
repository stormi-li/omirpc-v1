package omirpc

import (
	"fmt"
	"io"
	"net/http"

	"github.com/vmihailenco/msgpack/v5"
)

type Response struct {
	HttpResponse *http.Response
}

// Read 读取响应的 Body 并解码到 v
func (response *Response) Read(v any) error {
	if response.HttpResponse.Body == nil {
		return fmt.Errorf("response body is nil")
	}

	defer response.HttpResponse.Body.Close()

	// 读取 Body 内容
	body, err := io.ReadAll(response.HttpResponse.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if err := msgpack.Unmarshal(body, v); err != nil {
		return fmt.Errorf("failed to decode response body using msgpack: %w", err)
	}

	return nil
}
