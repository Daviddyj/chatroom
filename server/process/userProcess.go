package process

import (
	"customerManager/common/message"
	"customerManager/server/model"
	"customerManager/server/utils"
	"encoding/json"
	"fmt"
	"net"
)

type UserProcess struct {
	//字段？
	Conn   net.Conn
	UserId int
}

func (this *UserProcess) ServerProcessLogin(mes *message.Message) (err error) {
	//1. 先从mes中取出 mes.Data，并直接反序列化成LoginMes
	var loginMes message.LoginMes
	err = json.Unmarshal([]byte(mes.Data), &loginMes)
	if err != nil {
		fmt.Println("json.Unmarshal fail err=", err)
		return
	}
	//先声明一个resMes
	var resMes message.Message
	resMes.Type = message.LoginResMesType
	//在声明一个 LoginResMes, 并完成赋值
	var LoginResMes message.LoginResMes
	//如果用户id=100，密码=123456，认为合法，否则不合法
	//if loginMes.UserId == 100 && loginMes.UserPwd == "daiyijie" {
	//	LoginResMes.Code = 200
	//	fmt.Println("到核验用户密码这步了")
	//} else {
	//	LoginResMes.Code = 500
	//	LoginResMes.Error = "该用户不存在，请注册后再使用"
	//}

	//1.使用model.MyUserDao到redis去验证
	user, err := model.MyUserDao.Login(loginMes.UserId, loginMes.UserPwd)

	if err != nil {
		//LoginResMes.Code = 500
		//LoginResMes.Error = "该用户不存在，请注册后再使用"
		//这里我们先测试成功，然后我们再返回具体的错误信息
		if err == model.ERROR_USER_NOTEXISTS {
			LoginResMes.Code = 500
			LoginResMes.Error = err.Error()

		} else if err == model.ERROR_USER_PWD {
			LoginResMes.Code = 403
			LoginResMes.Error = err.Error()

		} else {
			LoginResMes.Code = 505
			LoginResMes.Error = "服务器内部错误"
		}
	} else {
		LoginResMes.Code = 200
		this.UserId = loginMes.UserId
		userMgr.AddOlineUser(this)
		//通知其他用户我上线了
		this.NotifyOtherOnlineUser(loginMes.UserId)

		fmt.Println(user, "登陆成功")
		for id, _ := range userMgr.onlineUsers {
			LoginResMes.UserIds = append(LoginResMes.UserIds, id)
		}
	}

	//将 loginResMes序列化
	data, err := json.Marshal(LoginResMes)
	if err != nil {
		fmt.Println("json.Marshal fail", err)
		return
	}
	//将data赋值给resMes
	resMes.Data = string(data)
	//对resMes序列化进行发送
	data, err = json.Marshal(resMes)
	if err != nil {
		fmt.Println("json.Marshal fail", err)
	}
	//发送
	//我们因为使用了分层模式(MVC),我们先创建一个Transfer实例，然后读取
	tf := &utils.Transfer{
		Conn: this.Conn,
	}
	err = tf.WritePkg(data)

	return
}

func (this *UserProcess) ServerProcessRegister(mes *message.Message) (err error) {
	var registerMes message.Register
	err = json.Unmarshal([]byte(mes.Data), &registerMes)
	if err != nil {
		fmt.Println("json.Unmarshal fail err=", err)
		return
	}
	//先声明一个resMes
	var resMes message.Message
	resMes.Type = message.RegisterMesType
	//在声明一个 RegisterResMes, 并完成赋值
	var RegisterResMes message.RegisterResMes
	//1.使用model.MyUserDao到redis去完成注册
	err = model.MyUserDao.Register(&registerMes.User)
	if err != nil {
		if err == model.ERROR_USER_EXISTS {
			RegisterResMes.Code = 505
			RegisterResMes.Error = err.Error()
		} else {
			RegisterResMes.Code = 506
			RegisterResMes.Error = "注册发生未知错误"
		}

	} else {
		RegisterResMes.Code = 200
	}
	data, err := json.Marshal(RegisterResMes)
	if err != nil {
		fmt.Println("json.Marshal fail", err)
		return
	}
	//将data赋值给resMes
	resMes.Data = string(data)
	//对resMes序列化进行发送
	data, err = json.Marshal(resMes)
	if err != nil {
		fmt.Println("json.Marshal fail", err)
	}
	//发送
	//我们因为使用了分层模式(MVC),我们先创建一个Transfer实例，然后读取
	tf := &utils.Transfer{
		Conn: this.Conn,
	}
	err = tf.WritePkg(data)

	return

}

func (this *UserProcess) NotifyOtherOnlineUser(userId int) {
	for id, up := range userMgr.onlineUsers {
		if id == userId {
			continue
		}
		up.NotifyMeOnline(userId)
	}
}

func (this *UserProcess) NotifyMeOnline(userId int) {
	var mes message.Message
	mes.Type = message.NotifyUserStatusMesType
	var notifyUserStatusMes message.NotifyUserStatusMes
	notifyUserStatusMes.UserId = userId
	notifyUserStatusMes.Status = message.UserOnline

	data, err := json.Marshal(notifyUserStatusMes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}
	tf := &utils.Transfer{
		Conn: this.Conn,
	}
	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("NodifyMeOnline err=", err)
		return
	}
}
