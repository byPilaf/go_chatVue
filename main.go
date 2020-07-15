package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"webchat/controller"
	"webchat/models"

	"github.com/gorilla/websocket"
)

func main() {
	go sendMes()
	go WS()
	http.Handle("/style/", http.StripPrefix("/style/", http.FileServer(http.Dir("views/style"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("views/js"))))

	http.HandleFunc("/", controller.IndexHandler)
	http.HandleFunc("/login", controller.LoginHandler)
	http.HandleFunc("/getOnlineUserList", controller.GetOnlineUserListHandler)
	http.HandleFunc("/sendMes", controller.GetUserMesHandler)

	http.HandleFunc("/ws", controller.WebSocketHandler) //websocket
	http.ListenAndServe(":7999", nil)
}

//WS 监听ws
func WS() {
	//要发送的数据
	for message := range models.WebSocketChann {
		for _, user := range models.OnlineUsersMap {
			err := user.WsConn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				fmt.Println("user.WsConn.WriteMessage err=", err)
			}
		}
	}
}

func sendMes() {
	var mes models.Mes
	mes.MesType = 0
	for {
		fmt.Printf("输入消息:")
		fmt.Scanln(&mes.Data)

		mesJSON, _ := json.Marshal(mes)
		models.WebSocketChann <- mesJSON
	}
}
