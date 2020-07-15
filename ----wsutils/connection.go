package wsutils

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

//connection 抽象需要的数据结构
type connection struct {
	//ws 连接器
	ws *websocket.Conn
	//数据
	data *Data
	//管道(发送)
	send chan []byte
}

//抽象ws连接器
//hub 处理ws中的各种逻辑
type hub struct {
	//注册的连接器
	connections map[*connection]bool
	//从连接器发送的信息
	broadcast chan []byte
	//从连接器注册请求
	resistder chan *connection
	//销毁请求
	unregister chan *connection
}

//实现ws的读和写
//ws写数据
func (conn *connection) writer() {
	//从管道遍历数据
	for message := range conn.send {
		//数据写出
		conn.ws.WriteMessage(websocket.TextMessage, message)
	}

	defer conn.ws.Close()
}

//ws读数据
func (conn *connection) reader() {
	//循环读取数据
	for {
		_, message, err := conn.ws.ReadMessage()
		if err != nil {
			//将此连接移除
			h.unregister <- conn
			break
		}

		//读取数据
		json.Unmarshal(message, &conn.data)
		//根据data的type判断该做什么
		switch conn.data.Type {
		case "login":
			//弹出窗口
			conn.data.User = conn.data.Content
			conn.data.From = conn.data.Content
		case "user":
			conn.data.Type = "user"
			dataB, _ := json.Marshal(conn.data)
			h.broadcast <- dataB
		case "logout":
			conn.data.Type = "logout"
			//删除用户列表
			//推送下线消息
			dataB, _ := json.Marshal(conn.data)
			h.broadcast <- dataB
			h.unregister <- conn
		default:
			fmt.Println("other")
		}
	}

}

//定义一个升级器
//将http请求升级位websocket请求
var upgreader = &websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true //允许跨域访问
	},
}

//ws 的回调函数
func wsHandle(w http.ResponseWriter, r *http.Request) {
	//1. 获取ws对象
	wsConn, err := upgreader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	//2. 创建连接对象
	//初始化连接对象
	c := &connection{
		send: make(chan []byte, 128),
		ws:   wsConn,
		data: &Data{},
	}
	//在ws中注册一下
	h.resistder <- c
	//ws将数据读写运行
	go c.writer()
	c.reader()

	//注销
	defer func() {
		c.data.Type = "logout"
		dataB, _ := json.Marshal(c.data)
		h.broadcast <- dataB
		h.unregister <- c
	}()
}
