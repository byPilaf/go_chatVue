package models

import (
	"encoding/json"
	"fmt"
)

var (
	//WebSocketChann 向客户端发送websocket数据管道
	//群发消息
	WebSocketChann chan []byte
)

func init() {
	WebSocketChann = make(chan []byte)
}

//WebSocketMessage websocket消息
type WebSocketMessage struct {
	Mes
	SendToUserToken string `json:"send_to_user_token"`
}

//SendAllUserMes 发送全体消息
func (wsMes *WebSocketMessage) SendAllUserMes() {
	//结构化
	sendMes, err := json.Marshal(wsMes)
	if err != nil {
		fmt.Println(err)
		return
	}

	//放入管道
	WebSocketChann <- sendMes
}
