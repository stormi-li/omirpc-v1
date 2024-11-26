package omirpc

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"math/rand/v2"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/stormi-li/omiserd-v1"
	"github.com/vmihailenco/msgpack/v5"
)

type RpcClient struct {
	ServerName       string
	Discover         *omiserd.Discover
	ListenHandleFunc func(serverName, oldAddress string, discover *omiserd.Discover) string
	monitor          *omiserd.Monitor
	Address          string
	HttpClient       *http.Client
}

func newRpcClient(opts *redis.Options, serverName string) *RpcClient {
	return &RpcClient{
		ServerName: serverName,
		Discover:   omiserd.NewClient(opts, omiserd.Server).NewDiscover(),
		ListenHandleFunc: func(serverName, oldAddress string, discover *omiserd.Discover) string {
			if !discover.IsAlive(serverName, oldAddress) {
				address := discover.GetByWeight(serverName)
				if len(address) != 0 {
					return address[rand.IntN(len(address))]
				}
			}
			return ""
		},
		HttpClient: &http.Client{},
	}
}
func (rpcClient *RpcClient) SkipVerify() {
	rpcClient.HttpClient = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // 忽略证书验证
			},
		},
	}
}

func (rpcClient *RpcClient) Post(pattern string, v any) (*Response, error) {
	if rpcClient.monitor == nil {
		rpcClient.monitor = rpcClient.Discover.NewMonitor(rpcClient.ServerName)
		rpcClient.Address = rpcClient.ListenHandleFunc(rpcClient.ServerName, "", rpcClient.Discover)
		rpcClient.monitor.Address = rpcClient.Address
		go rpcClient.monitor.ListenAndConnect(2*time.Second, rpcClient.ListenHandleFunc, func(address string) {
			rpcClient.Address = address
		})
	}
	// 构造请求 URL
	url := "https://" + rpcClient.Address + pattern

	// 将 v 序列化为 JSON 数据
	jsonData, err := msgpack.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request data: %w", err)
	}

	// 发起 POST 请求
	resp, err := rpcClient.HttpClient.Post(url, "application/json", bytes.NewReader(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to send POST request: %w", err)
	}

	return &Response{HttpResponse: resp}, nil
}
