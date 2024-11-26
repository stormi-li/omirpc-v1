package omirpc

import (
	"net/http"

	"github.com/go-redis/redis/v8"
	"github.com/stormi-li/omicert-v1"
	"github.com/stormi-li/omiserd-v1"
)

type RpcServer struct {
	ServerName     string
	Address        string
	ServerRegister *omiserd.Register
	Weight         int
}

func newRpcServer(opts *redis.Options, serverName, address string) *RpcServer {
	return &RpcServer{
		ServerName:     serverName,
		Address:        address,
		ServerRegister: omiserd.NewClient(opts, omiserd.Server).NewRegister(serverName, address),
	}
}

func (RpcServer) AddHandleFunc(pattern string, handleFunc func(w ResponseWriter, r *Request)) {
	http.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		handleFunc(ResponseWriter{HttpResponseWriter: w}, &Request{HttpRequest: r})
	})
}

func (rpcServer *RpcServer) Start(weight int, credential *omicert.Credential) {
	rpcServer.ServerRegister.RegisterAndServe(weight, func(port string) {
		omicert.ListenAndServeTLS(port, credential)
	})
}
