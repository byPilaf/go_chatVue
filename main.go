package main

import (
	"net/http"
	"webchat/controller"
)

func main() {
	http.Handle("/style/", http.StripPrefix("/style/", http.FileServer(http.Dir("views/style"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("views/js"))))

	http.HandleFunc("/", controller.IndexHandler)
	http.HandleFunc("/login", controller.LoginHandler)
	http.HandleFunc("/getOnlineUserList", controller.GetOnlineUserListHandler)
	http.HandleFunc("/sendMes", controller.GetUserMesHandler)

	http.HandleFunc("/ws", controller.WebSocketHandler) //websocket
	http.ListenAndServe(":10241", nil)
}
