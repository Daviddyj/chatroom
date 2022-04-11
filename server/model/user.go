package model

//定义一个用户的结构体

type User struct {
	//确定字段信息
	//用户信息的json字符串key和结构体字段对应的tag名字一致！！！
	UserId   int    `json:"userId"`
	UserPwd  string `json:"userPwd"`
	UserName string `json:"userName"`
}
