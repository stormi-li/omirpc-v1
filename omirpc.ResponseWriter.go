package omirpc

import (
	"fmt"
	"net/http"

	"github.com/vmihailenco/msgpack"
)

type ResponseWriter struct {
	HttpResponseWriter http.ResponseWriter
}

func (response ResponseWriter) Write(v any) error {
	// 序列化为 MsgPack
	data, err := msgpack.Marshal(v)
	if err != nil {
		return fmt.Errorf("failed to marshal response: %w", err)
	}

	// 写入响应
	_, err = response.HttpResponseWriter.Write(data)
	if err != nil {
		return fmt.Errorf("failed to write response: %w", err)
	}

	return nil
}
