package main

import (
	"customerManager/server/model"
	"fmt"
	"net"
	"time"
)

func initUserDao() {
	model.MyUserDao = model.NewUserDao(pool)
}

//处理和客户端的通讯
func process(conn net.Conn) {
	//这里需要延时关闭
	defer conn.Close()

	//这里调用总控
	processor := &Processor{
		Conn: conn,
	}
	err := processor.process3()
	if err != nil {
		fmt.Println("客户端和服务器通讯协程错误=", err)
		return
	}
}

func init() {
	//提示信息
	initPool("localhost:6379", 16, 0, 300*time.Second)
	initUserDao()
}

func main() {

	fmt.Println("服务器在8889端口监听")

	listen, err := net.Listen("tcp", "0.0.0.0:8889")

	if err != nil {
		fmt.Println("net.Listen err=", err)
		return
	}
	defer listen.Close()
	//一旦监听成功，就等待客户端来链接服务器

	for {
		fmt.Println("等待客户端来连接服务器")
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("listen.Accept err=", err)
		}
		//一旦连接成功，则启动一个协程和客户端保持通讯。。
		go process(conn)

	}

}
