package main

import (
	"customerManager/client/process"
	"fmt"
	"os"
)

//定义两个变量，一个表示id,一个表示用户密码
var userId int
var userPwd string
var userName string

func main() {

	var key int
	//var loop = true

	for {
		fmt.Println("欢迎登录多人聊天系统")
		fmt.Println("\t\t\t 1.登录聊天室")
		fmt.Println("\t\t\t 2.注册用户")
		fmt.Println("\t\t\t 3.退出系统")
		fmt.Println("\t\t\t 请选择(1-3):")

		fmt.Scanf("%d\n", &key)
		switch key {
		case 1:
			fmt.Println("登录聊天室")
			fmt.Println("请输入用户的id")
			fmt.Scanf("%d\n", &userId)
			fmt.Println("请输入用户的密码")
			fmt.Scanf("%s\n", &userPwd)
			//完成登录
			//1.创建一个UserProcess的实例
			up := &process.UserProcess{}
			up.Login(userId, userPwd)

			//loop = false
		case 2:
			fmt.Println("注册用户")
			fmt.Println("请输入用户的id")
			fmt.Scanf("%d\n", &userId)
			fmt.Println("请输入用户的密码")
			fmt.Scanf("%s\n", &userPwd)
			fmt.Println("请输入用户的姓名")
			fmt.Scanf("%s\n", &userName)
			up := &process.UserProcess{}
			up.Register(userId, userPwd, userName)
			//loop = false
		case 3:
			fmt.Println("退出系统")
			os.Exit(0)
		default:
			fmt.Println("你的输入有误，请重新输入")
		}
	}

	//根据用户的输入显示新的聊天信息
	//if key == 1 {
	//说明用户要登陆

	//先把登录的函数写到另外一个文件，比如login.go

	//因为使用了新的程序结构

	//这里我们会需要重新调用
	//	err := login(userId, userPwd)
	//	if err != nil {
	//		fmt.Println("登陆失败了wuwu")
	//	} else {
	//		fmt.Println("登陆成功hahah")
	//	}
	//
	//} else if key == 2 {
	//fmt.Println("进行用户注册")
	//}
	//
	//}
}
