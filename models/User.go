package models

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/jinzhu/gorm"
)

//OnlineUsersMap 当前在线用户列表 [用户token]用户模型
var OnlineUsersMap map[string]*User = make(map[string]*User)

//User 用户模型
type User struct {
	gorm.Model
	Name string `json:"name" gorm:"not null;unique"`
	Pass string `json:"pass" gorm:"not null"`
	//WsConn 注册的ws连接器
	WsConn *websocket.Conn
	Mutex  sync.Mutex
	Token  string
	//todo 待发送的私人消息队列
	//todo 接收到的私人消息队列
}

//BeatLine 用户心跳检测
func (user *User) BeatLine() {
	beatMes := Mes{
		Code:    200,
		MesType: HiddenMesType,
	}
	mesJSON, _ := json.Marshal(beatMes)

	for {
		//这里并发写入了???
		err := user.WsConn.WriteMessage(websocket.TextMessage, mesJSON)
		if err != nil {
			fmt.Println("err:", err)
			go user.OffLine()
			break
		}
		time.Sleep(time.Second * 5)
	}

	return
}

//OffLine 用户下线
func (user *User) OffLine() {
	//从在线列表排除
	delete(OnlineUsersMap, user.Token)
	//向客户端发送下线消息
	offLineMes := Mes{
		FromUserName:  user.Name,
		FromUserToken: user.Token,
		MesType:       UserStatusMesType,
		Code:          200,
		Data:          "offline",
	}
	mesJSON, _ := json.Marshal(offLineMes)
	WebSocketChann <- mesJSON
}
