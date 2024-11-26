package omirpc

import "github.com/go-redis/redis/v8"

type Client struct {
	opts *redis.Options
}

func NewClient(opts *redis.Options) *Client {
	return &Client{opts: opts}
}

func (c *Client) NewRpcClient(serverName string) *RpcClient {
	return newRpcClient(c.opts, serverName)
}

func (c *Client) NewRpcServer(serverName string, address string) *RpcServer {
	return newRpcServer(c.opts, serverName, address)
}
