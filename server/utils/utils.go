package utils

import (
	"customerManager/common/message"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"net"
)

//这里我们将这些方法关联到结构体中
type Transfer struct {
	//分析它应该有哪些字段
	Conn net.Conn
	Buf  [8064]byte //这是传输时使用的缓冲
}

func (this *Transfer) ReadPkg() (mes message.Message, err error) {
	//buf := make([]byte, 8096)
	fmt.Println("读取客户端发送的数据")
	//conn.Read在conn没有被关闭的情况下，会阻塞
	//如果客户端关闭了conn 则，就不会阻塞
	_, err = this.Conn.Read(this.Buf[0:4])
	if err != nil {
		return

	}
	//根据buf[:4]  转成一个uint32类型
	var pkhLen uint32
	pkhLen = binary.BigEndian.Uint32(this.Buf[0:4])

	//根据pkgLen读取消息内容
	n, err := this.Conn.Read(this.Buf[:pkhLen])
	if uint32(n) != pkhLen || err != nil {
		fmt.Println("conn.Read fail err=", err)
	}
	//把pkgLen反序列化成 -> message.Message
	//技术就是一层窗户纸  &mes
	err = json.Unmarshal(this.Buf[:pkhLen], &mes)
	if err != nil {
		err = errors.New("read pkg body error")
		return
	}
	return
}

func (this *Transfer) WritePkg(data []byte) (err error) {
	//先发送给一个长度给对方
	var pkgLen uint32
	pkgLen = uint32(len(data))
	//var buf [4]byte
	binary.BigEndian.PutUint32(this.Buf[0:4], pkgLen)
	//现在发送长度
	n, err := this.Conn.Write(this.Buf[:4])
	if n != 4 || err != nil {
		fmt.Println("conn.Write(bytes) fail", err)
		return
	}
	n, err = this.Conn.Write(data)
	if n != int(pkgLen) || err != nil {
		fmt.Println("conn.Write(bytes) fail", err)
		return
	}
	return

}
