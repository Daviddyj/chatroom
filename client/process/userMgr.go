package process

import (
	"customerManager/client/model"
	"customerManager/common/message"
	"fmt"
)

var onlineUsers map[int]*message.User = make(map[int]*message.User, 10)
var CurUser model.CurUser

func outputOnlineUser() {
	fmt.Println("当前在线用户列表")
	for id, _ := range onlineUsers {
		fmt.Println("用户id:\t", id)
	}
}

//编写i一个方法，处理返回的NofifyUserStatusMes
func updateUserStatus(notifyUserstatusMes *message.NotifyUserStatusMes) {
	user, ok := onlineUsers[notifyUserstatusMes.UserId]
	if !ok {
		user = &message.User{
			UserId: notifyUserstatusMes.UserId,
		}
	}
	user.UserStatus = notifyUserstatusMes.Status
	onlineUsers[notifyUserstatusMes.UserId] = user
	outputOnlineUser()

}
