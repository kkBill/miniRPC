package public

type User struct {
	Name string
	Age int
}

type ResponseQueryUser struct {
	User     // 结构体的嵌套？？why not like this: user User
	Msg string
}