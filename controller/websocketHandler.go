package controller

import (
	"fmt"
	"net/http"
	"webchat/models"

	"github.com/gorilla/websocket"
)

//upgrader 将http请求提升为websocket
var upgrader = &websocket.Upgrader{
	ReadBufferSize:  1024, //读缓冲
	WriteBufferSize: 1024, //写缓冲
	CheckOrigin: func(r *http.Request) bool {
		return true //允许跨域
	},
}

// WebSocketHandler websocket连接控制器
func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	//获取userToken
	if r.Method == "GET" {
		userTokenList := r.URL.Query()["user_token"]
		userToken := userTokenList[0]

		wsConn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println("upgrader.Upgrade(w, r, nil) error :", err)
			return
		}

		//绑定已在线用户
		user := models.OnlineUsersMap[userToken]
		user.WsConn = wsConn
		oldConn := user.WsConn
		if oldConn != nil {
			user.StatusChan <- 1
			close(user.StatusChan)
		} else {
			user.CreatChannel()
			go user.WaitForSendMes()
			go user.BeatLine() //心跳检测
		}

		//发送上线通知
		var mes models.WebSocketMessage
		mes.Data = "online"
		mes.FromUserToken = user.Token
		mes.FromUserName = user.Name
		mes.MesType = models.UserStatusMesType
		mes.SendAllUserMes()
	}
}
