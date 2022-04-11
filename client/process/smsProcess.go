package process

import (
	"customerManager/client/utils"
	"customerManager/common/message"
	"encoding/json"
	"fmt"
)

type smsProcess struct {
}

func (this *smsProcess) SendGroupMes(context string) (err error) {
	//2.准备通过conn发送消息给服务器
	var mes message.SmsMes
	mes.Type = message.SmsMesType
	//3.创建一个LoginMes结构体
	var smsMes message.SmsMes
	smsMes.Context = context

	smsMes.User.UserId = CurUser.UserId
	smsMes.User.UserStatus = CurUser.UserStatus

	//4.将loginMes 序列化
	data, err := json.Marshal(smsMes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}
	//5.把data赋给mes.Data字段
	mes.Data = string(data)

	//6.将mes进行序列化
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}

	tf := &utils.Transfer{
		Conn: CurUser.Conn,
	}
	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("conn.Write(data) fail", err)
		return
	}
	return
}

func outputGroupMes(mes *message.Message) {
	var smsMes message.SmsMes
	err := json.Unmarshal([]byte(mes.Data), &smsMes)
	if err != nil {
		fmt.Println("json.Unmarshal fail err=", err)
		return
	}

	//显示信息
	info := fmt.Sprintf("用户id:\t%d 对大家说:\t%s", smsMes.UserId, smsMes.Context)
	fmt.Println(info)
	fmt.Println()
}
