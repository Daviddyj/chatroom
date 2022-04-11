package model

import (
	"customerManager/common/message"
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
)

//定义一个UserDao 结构体
//完成对User结构体的各种操作
var MyUserDao *UserDao

type UserDao struct {
	pool *redis.Pool
}

func NewUserDao(pool *redis.Pool) (userDao *UserDao) {
	userDao = &UserDao{
		pool: pool,
	}
	return
}

//在UserDao应该提供哪些方法给我们
//根据一个用户id 返回一个User实例+err
func (this *UserDao) getUserById(conn redis.Conn, id int) (user *User, err error) {
	//通过给定 id 去 redis 查询这个用户
	res, err := redis.String(conn.Do("HGet", "users", id))
	if err != nil {
		//错误!
		if err == redis.ErrNil { //表示在 users 哈希中，没有找到对应 id
			err = ERROR_USER_NOTEXISTS
		}
		return
	}
	user = &User{}
	//这里我们需要把 res 反序列化成 User 实例
	err = json.Unmarshal([]byte(res), user)
	if err != nil {
		fmt.Println("json.Unmarshal err=", err)
		return
	}
	return
}
func (this *UserDao) Login(userId int, userPwd string) (user *User, err error) {
	//先从 UserDao 的连接池中取出一根连接
	conn := this.pool.Get()
	defer conn.Close()
	user, err = this.getUserById(conn, userId)
	if err != nil {
		return
	}
	//这时证明这个用户是获取到.
	if user.UserPwd != userPwd {
		err = ERROR_USER_PWD
		return
	}
	return
}
func (this *UserDao) Register(user *message.User) (err error) {
	//先从 UserDao 的连接池中取出一根连接
	conn := this.pool.Get()
	defer conn.Close()
	_, err = this.getUserById(conn, user.UserId)
	if err == nil {
		err = ERROR_USER_EXISTS
		return
	}
	//这时证明这个用户是没有被注册的
	data, err := json.Marshal(user)
	if err != nil {
		return
	}
	//入库
	_, err = conn.Do("Hset", "users", user.UserId, string(data))
	if err != nil {
		fmt.Println("保存注册用户错误")
		return
	}
	return
}
