package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
	// "github.com/jinzhu/gorm"
)

//OnlineUsersMap 当前在线用户列表 [用户token]用户模型
var OnlineUsersMap map[string]*User = make(map[string]*User)

//UsersMap 已注册的用户列表
var UsersMap map[string]*User = make(map[string]*User)

//User 用户模型
type User struct {
	// gorm.Model
	Name string `json:"name" gorm:"not null;unique"`
	Pass string `json:"pass" gorm:"not null"`
	//WsConn 注册的ws连接器
	WsConn *websocket.Conn `gorm:"-"`
	Token  string          `gorm:"-"`
	//待发送的私人消息队列
	UserWriteChan chan []byte `gorm:"-"`
	// 接收到的私人消息队列
	UserReadChan chan []byte `gorm:"-"`
	//StatusChan 当前用户状态 1=下线
	StatusChan chan int `gorm:"-"`
}

//CreatChannel 创建频道
func (user *User) CreatChannel() {
	user.UserWriteChan = make(chan []byte, 100)
	user.UserReadChan = make(chan []byte, 100)
	user.StatusChan = make(chan int)
	fmt.Println(user.Name, "上线了")
}

//BeatLine 用户心跳检测
func (user *User) BeatLine() {
	beatMes := Mes{
		Code:    200,
		MesType: HiddenMesType,
	}
	mesJSON, _ := json.Marshal(beatMes)

forE:
	select {
	case <-user.StatusChan: //终止心跳
		return
	default:
		user.UserWriteChan <- mesJSON
		time.Sleep(time.Second * 5)
		goto forE
	}
}

//OffLine 用户下线
func (user *User) OffLine() {
	fmt.Println(user.Name, "下线了")

	//向客户端发送下线消息
	var offLineMes WebSocketMessage
	offLineMes.FromUserName = user.Name
	offLineMes.FromUserToken = user.Token
	offLineMes.MesType = UserStatusMesType
	offLineMes.Code = 200
	offLineMes.Data = "offline"
	offLineMes.SendAllUserMes()
	user.StatusChan <- 1
	//从在线列表排除
	delete(OnlineUsersMap, user.Token)
}

//WaitForSendMes 等待发送数据
func (user *User) WaitForSendMes() {
	for message := range user.UserWriteChan {
		err := user.WsConn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			fmt.Println("user.WsConn.WriteMessage error :", err)
			user.OffLine()
		}
	}
}
