package main

import (
	"fmt"

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
	client := omirpc.NewClient(&redis.Options{Addr: redisAddr, Password: password}).NewRpcClient("rpc_server")
	client.SkipVerify()
	user := User{Name: "lili"}
	resp, err := client.Post("/hello", &user)
	if err == nil {
		resp.Read(&user)
		fmt.Println(user)
	} else {
		fmt.Println(err)
	}
}
