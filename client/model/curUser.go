package model

import (
	"customerManager/common/message"
	"net"
)

type CurUser struct {
	Conn net.Conn
	message.User
}
