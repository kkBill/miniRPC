package main

import (
	"encoding/gob"
	"errors"
	"log"
	"main/public"
	"network"
)

func queryUser(uid int) (public.ResponseQueryUser, error) {
	// 构造数据（模仿数据库中存储的数据）
	db := make(map[int]public.User)
	db[0] = public.User{Name: "Kobe Bryant", Age:39}
	db[1] = public.User{Name: "LeBron James", Age:32}

	// 执行查询
	if user,ok := db[uid]; ok {
		return public.ResponseQueryUser{User: user, Msg:"success"}, nil
	}else{
		return public.ResponseQueryUser{User: public.User{}, Msg:"fail"}, errors.New("query failed. no such item in database!")
	}
}

func main() {
	gob.Register(public.ResponseQueryUser{})

	addr := "0.0.0.0:2333"
	server := network.CreateServer(addr)

	// 注册服务
	server.Register("queryUser", queryUser)

	log.Println("Server is running...")

	go server.Run()

	// 这里为什么要无限循环，防止主线程退出
	// 这么处理是粗暴的做法，主线程应该等待子线程退出后再退出
	// 具体是哪个api忘了，是关于channel的问题
	for {
	}
}