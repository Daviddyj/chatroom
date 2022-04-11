package process

import (
	"customerManager/client/utils"
	"customerManager/common/message"
	"encoding/json"
	"fmt"
	"net"
	"os"
)

//显示登录成功后的界面..

func ShowMenu() {
	fmt.Println("---------恭喜xxx登录成功----------")
	fmt.Println("---------1.显示在线用户列表")
	fmt.Println("---------2.发送消息----------")
	fmt.Println("---------3.信息列表----------")
	fmt.Println("---------4.退出系统----------")
	fmt.Println("---------请选择（1-4）----------")
	var key int
	var context string
	smsProcess := smsProcess{}
	fmt.Scanf("%d\n", &key)
	switch key {
	case 1:
		fmt.Println("显示在线用户列表")
		outputOnlineUser()
	case 2:
		fmt.Println("你想对大家说点什么,请输入:")
		fmt.Scanf("%s\n", &context)
		smsProcess.SendGroupMes(context)

	case 3:
		fmt.Println("信息列表")
	case 4:
		fmt.Println("你选择退出了系统")
		os.Exit(0)
	default:
		fmt.Println("你输入的选项不对，请重新输入")
	}
}

//和服务器端保持通讯
func serverProcessMes(Conn net.Conn) {
	//创建一个transfer实例，不停的读取服务端发送的消息
	tf := &utils.Transfer{
		Conn: Conn,
	}
	for {
		fmt.Println("客户端正在等待服务器发送的消息")
		fmt.Println("真的跑到这一步了吗")
		mes, err := tf.ReadPkg()
		if err != nil {
			fmt.Println("tf.ReadPkg err=", err)
			return
		}
		fmt.Printf("mes=%v", mes)
		switch mes.Type {
		case message.NotifyUserStatusMesType:

			var notifyUserStatusMes message.NotifyUserStatusMes
			json.Unmarshal([]byte(mes.Data), &notifyUserStatusMes)
			updateUserStatus(&notifyUserStatusMes)
		case message.SmsMesType:
			outputGroupMes(&mes)
		default:
			fmt.Println("服务器返回了未知的消息类型")
		}
	}
}
