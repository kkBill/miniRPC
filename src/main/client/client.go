package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net"

	"main/service"
	"network"
)

func main() {
	gob.Register(service.ResponseQueryUser{})

	addr := "0.0.0.0:2333"
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Printf("dial error: %v\n", err)
	}
	log.Println("connect successfully...")

	client := network.CreateClient(conn)

	// 本地的函数声明（只有声明，没有实现）
	var queryUser func(int) (service.ResponseQueryUser, error)

	// 远程调用
	client.Call("queryUser", &queryUser)
	user, err := queryUser(1)
	fmt.Println(user)

	conn.Close()
}