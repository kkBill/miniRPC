package main

import (
	"encoding/gob"
	"log"
	"net"

	"main/public"
	"network"
)

func main() {
	// ??
	gob.Register(public.ResponseQueryUser{})

	addr := "0.0.0.0:2333"
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Printf("dial error: %v\n", err)
	}
	log.Println("connect successfully...")

	client := network.CreateClient(conn)

	// 本地的函数声明（只有声明，没有实现）
	var correctQuery func(int) (public.ResponseQueryUser, error)
	//var wrongQuery func(int) (public.ResponseQueryUser, error)

	// 远程调用
	client.Call("queryUser", &correctQuery)

	user, err := correctQuery(1) // 卡在这里了

	log.Println("xxxxxx")
	if err != nil {
		log.Printf("query error: %v\n", err)
	} else {
		log.Printf("query result: %v %v %v\n", user.Name, user.Age, user.Msg)
	}

	//user, err = correctQuery(2)
	//if err!=nil {
	//	log.Printf("query error: %v\n", err)
	//}else {
	//	log.Printf("query result: %v %v %v\n", user.Name, user.Age, user.Msg)
	//}

	//client.Call("queryUser", &wrongQuery)
	//user, err = wrongQuery(1)
	//if err!=nil {
	//	log.Printf("query error: %v\n", err)
	//}else {
	//	log.Printf("query result: %v %v %v\n", user.Name, user.Age, user.Msg)
	//}

	conn.Close()
}
