package service

import "errors"

type User struct {
	Name string
	Age int
}

type ResponseQueryUser struct {
	User     // 结构体的嵌套？？why not like this: user User
	Msg string
}

func QueryUser(uid int) (ResponseQueryUser, error) {
	// 构造数据（模仿数据库中存储的数据）
	db := make(map[int]User)
	db[0] = User{Name: "Kobe Bryant", Age:39}
	db[1] = User{Name: "LeBron James", Age:32}

	// 执行查询
	if user,ok := db[uid]; ok {
		return ResponseQueryUser{User: user, Msg:"success"}, nil
	}else{
		return ResponseQueryUser{User: User{}, Msg:"fail"}, errors.New("query failed. no such item in database!")
	}
}
