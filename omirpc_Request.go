package omirpc

import (
	"fmt"
	"io"
	"net/http"

	"github.com/vmihailenco/msgpack/v5"
)

type Request struct {
	HttpRequest *http.Request
}

func (request *Request) Read(v any) error {
	// 确保读取 Body 的内容
	body, err := io.ReadAll(request.HttpRequest.Body)
	if err != nil {
		return fmt.Errorf("failed to read body: %w", err)
	}

	// 解码到目标对象
	if err := msgpack.Unmarshal(body, v); err != nil {
		return fmt.Errorf("failed to unmarshal body: %w", err)
	}

	return nil
}
