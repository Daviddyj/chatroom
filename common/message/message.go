package message

const (
	LoginMesType            = "LoginMes"
	LoginResMesType         = "LoginResMes"
	RegisterMesType         = "Register"
	NotifyUserStatusMesType = "NotifyUserStatus"
	SmsMesType              = "SmsMesType"
)

const (
	UserOnline = iota
	Useroffline
	UserBusyStatus
)

type Message struct {
	Type string `json:"type"` //消息类型
	Data string `json:"data"` //
}

type LoginMes struct {
	UserId   int    `json:"userId"`
	UserPwd  string `json:"userPwd"`
	UserName string `json:"userName"`
}

type LoginResMes struct {
	Code    int    `json:"code"` //返回状态码  500表示用户未注册  200表示登陆成功
	Error   string `json:"error"`
	UserIds []int
}

type Register struct {
	User User `json:"user"`
}
type RegisterResMes struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
}
type NotifyUserStatusMes struct {
	UserId int `json:"userId"`
	Status int `json:"status"`
}

type SmsMes struct {
	Context string `json:"context"`
	User
	Type string `json:"type"`
	Data string `json:"data"`
}
type CurUser struct {
	UserId     int
	UserStatus int
}
