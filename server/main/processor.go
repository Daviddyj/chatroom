package main

import (
	"customerManager/common/message"
	process2 "customerManager/server/process"
	"customerManager/server/utils"
	"errors"
	"fmt"
	"io"
	"net"
)

//先创建一个Processor的结构体
type Processor struct {
	Conn net.Conn
}

func (this *Processor) ServerProcessMes(mes *message.Message) (err error) {

	switch mes.Type {
	case message.LoginMesType:
		//处理登录
		//创建一个UserProcess实例
		up := &process2.UserProcess{
			Conn: this.Conn,
		}
		//处理登录的逻辑
		err = up.ServerProcessLogin(mes)
	case message.RegisterMesType:
		//处理注册
		up := &process2.UserProcess{
			Conn: this.Conn,
		}
		//处理登录的逻辑
		err = up.ServerProcessRegister(mes)
	case message.SmsMesType:
		//处理群发
		//创建一个UserProcess实例
		smsProcess := &process2.SmsProcess{}
		//处理登录的逻辑
		//err = up.ServerProcessLogin(mes)
		smsProcess.SendGroupMes(mes)

	default:
		fmt.Println("消息类型不存在，无法处理")
	}
	return err
}

func (this *Processor) process3() (error error) {
	for {
		//创建一个Transfer实例完成读包任务
		tf := &utils.Transfer{
			Conn: this.Conn,
		}
		mes, err := tf.ReadPkg()
		if err != nil {
			if err == io.EOF {
				fmt.Println("客户端退出，服务器也退出..")
				return
			} else {
				err = errors.New("read pkg header error")
				return
			}
		}
		//这里我们将读取数据包，直接封装成一个函数readPkg(),返回Message,Err
		fmt.Println("mes=", mes)
		err = this.ServerProcessMes(&mes)
		if err != nil {
			return
		}
	}
}
