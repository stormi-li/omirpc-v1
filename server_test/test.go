package main

import (
	"github.com/go-redis/redis/v8"
	"github.com/stormi-li/omirpc-v1"
)

var redisAddr = "118.25.196.166:3934"
var password = "12982397StrongPassw0rd"

type User struct {
	ID    int
	Name  string
	Email string
}

func main() {
	server := omirpc.NewClient(&redis.Options{Addr: redisAddr, Password: password}).NewRpcServer("rpc_server", "118.25.196.166:9988")
	server.AddHandleFunc("/hello", func(w omirpc.ResponseWriter, r *omirpc.Request) {
		var user User
		r.Read(&user)
		user.Name = "rpc server: hello " + user.Name + "!"
		w.Write(&user)
	})
	server.Start(1, nil)
}
