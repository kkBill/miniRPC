package main

import (
	"encoding/gob"
	"log"
	"main/service"
	"network"
)

func main() {
	gob.Register(service.ResponseQueryUser{}) // 应该要取消，不能暴露出来

	addr := "0.0.0.0:2333"
	server := network.CreateServer(addr)

	// 注册服务
	server.Register("queryUser", service.QueryUser)

	go server.Run()

	log.Println("Server is running...")

	// 这么写是很粗暴的
	for {}
}