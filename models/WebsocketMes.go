package models

import (
	"encoding/json"
	"fmt"
)

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
	for _, user := range OnlineUsersMap {
		user.UserWriteChan <- sendMes
	}
}
