package process

import (
	"customerManager/client/utils"
	"customerManager/common/message"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"net"
)

type UserProcess struct {
	//暂时不需要字段
}

func (this *UserProcess) Login(userId int, userPwd string) (err error) {
	//下一部开始开始定制协议
	//fmt.Printf("userId=%d userPwd=%s", userId, userPwd)
	//return nil

	//1.链接到服务器
	conn, err := net.Dial("tcp", "localhost:8889")
	if err != nil {
		fmt.Println("net.Dial err=", err)
		return
	}
	//延时关闭
	defer conn.Close()

	//2.准备通过conn发送消息给服务器
	var mes message.Message
	mes.Type = message.LoginMesType
	//3.创建一个LoginMes结构体
	var loginMes message.LoginMes
	loginMes.UserId = userId
	loginMes.UserPwd = userPwd

	//4.将loginMes 序列化
	data, err := json.Marshal(loginMes)
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
	//7.到这个时候 data就是我们要发送的消息
	//7.1先把data的长度发送给服务器
	//先获取到data的长度-> 转换成一个表示长度的byte切片
	var pkgLen uint32
	pkgLen = uint32(len(data))
	var buf [4]byte
	binary.BigEndian.PutUint32(buf[0:4], pkgLen)
	//现在发送长度
	n, err := conn.Write(buf[:4])
	if n != 4 || err != nil {
		fmt.Println("conn.Write(bytes) fail", err)
		return
	}
	fmt.Println("客户端发送的消息长度成功=", len(data), string((data)))

	//发送消息本身
	_, err = conn.Write(data)
	if err != nil {
		fmt.Println("conn.Write(data) fail", err)
		return
	}

	//这里还需要处理服务器端返回的消息
	//休眠20
	//time.Sleep(20 * time.Second)
	//fmt.Println("休眠了20..")
	//创建一个Transfer 实例

	tf := &utils.Transfer{
		Conn: conn,
	}
	mes, err = tf.ReadPkg() //mes就是

	if err != nil {
		fmt.Println("readPkg(conn) err=", err)
	}
	//将mes的Data部分反序列化成LoginResMes
	var loginResMes message.LoginResMes
	err = json.Unmarshal([]byte(mes.Data), &loginResMes)
	if loginResMes.Code == 200 {
		CurUser.Conn = conn
		CurUser.UserId = userId
		CurUser.UserStatus = message.UserOnline

		fmt.Println("用户登录成功,状态码200")
		fmt.Println("当前在线的用户列表如下:")
		for _, v := range loginResMes.UserIds {
			if v == userId {
				continue
			}
			fmt.Println("用户id:\t", v)

			//完成客户端的onlineUser 完成初始化
			user := &message.User{
				UserId:     v,
				UserStatus: message.UserOnline,
			}
			onlineUsers[v] = user
		}
		fmt.Print("\n\n")

		//这里我们还需要在客户端启动一个协程
		//该协程保持和服务器端的通讯，如果服务器有数据推送给客户端
		//则接收并显示在客户端的终端
		go serverProcessMes(conn)

		//1.循环显示我们登录成功的菜单
		for {
			ShowMenu()
		}

	} else if loginResMes.Code == 500 {
		fmt.Println(loginResMes.Error)
		err = errors.New("状态码500")
	} else if loginResMes.Code == 403 {
		fmt.Println(loginResMes.Error)
		err = errors.New("状态码403")
	}

	return
}
func (this *UserProcess) Register(userId int, userPwd string, userName string) (err error) {
	//1.链接到服务器
	conn, err := net.Dial("tcp", "localhost:8889")
	if err != nil {
		fmt.Println("net.Dial err=", err)
		return
	}
	//延时关闭
	defer conn.Close()

	//2.准备通过conn发送消息给服务器
	var mes message.Message
	mes.Type = message.RegisterMesType
	//3.创建一个RegisterMes结构体
	var registerMes message.Register
	registerMes.User.UserId = userId
	registerMes.User.UserPwd = userPwd
	registerMes.User.UserName = userName

	//4.将registerMes 序列化
	data, err := json.Marshal(registerMes)
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
	//7.到这个时候 data就是我们要发送的消息
	//7.1先把data的长度发送给服务器
	//先获取到data的长度-> 转换成一个表示长度的byte切片
	var pkgLen uint32
	pkgLen = uint32(len(data))
	var buf [4]byte
	binary.BigEndian.PutUint32(buf[0:4], pkgLen)
	//现在发送长度
	n, err := conn.Write(buf[:4])
	if n != 4 || err != nil {
		fmt.Println("conn.Write(bytes) fail", err)
		return
	}
	fmt.Println("客户端发送的消息长度成功=", len(data), string((data)))

	//发送消息本身
	_, err = conn.Write(data)
	if err != nil {
		fmt.Println("conn.Write(data) fail", err)
		return
	}
	//这里还需要处理服务器端返回的消息
	//休眠20
	//time.Sleep(20 * time.Second)
	//fmt.Println("休眠了20..")
	//创建一个Transfer 实例

	tf := &utils.Transfer{
		Conn: conn,
	}

	mes, err = tf.ReadPkg() //mes就是
	if err != nil {
		fmt.Println("readPkg(conn) err=", err)
	}
	//将mes的Data部分反序列化成RegisterResMes
	var registerResMes message.RegisterResMes
	err = json.Unmarshal([]byte(mes.Data), &registerResMes)
	if registerResMes.Code == 200 {
		fmt.Println("用户注册成功,状态码200")

		//这里我们还需要在客户端启动一个协程
		//该协程保持和服务器端的通讯，如果服务器有数据推送给客户端
		//则接收并显示在客户端的终端
		//go serverProcessMes(conn)
	} else if registerResMes.Code == 505 {
		err = errors.New("状态码505")
		fmt.Println(registerResMes.Error)

	}
	return
}
