package wsutils

import "encoding/json"

//将连接器对象初始化
var h = hub{
	//注册的连接器
	connections: make(map[*connection]bool),
	//从连接器发送的信息
	broadcast: make(chan []byte),
	//从连接器注册请求
	resistder: make(chan *connection),
	//销毁请求
	unregister: make(chan *connection),
}

//ws控制器

//处理ws的逻辑实现
func (h *hub) run() {
	//监听数据管道, 在后端处理管道数据
	for {
		//根据不同的数据管道,处理不同的逻辑
		select {
		//注册
		case c := <-h.resistder:
			//标识注册
			h.connections[c] = true
			//组装data数据
			c.data.IP = c.ws.RemoteAddr().String()
			//更新类型
			c.data.Type = "handshake"
			//用户列表
			// c.data.UserList =
			//序列化
			dataB, _ := json.Marshal(c.data)
			//将数据放入数据管道
			c.send <- dataB
		//注销
		case c := <-h.unregister:
			//判断map中存在要删除的数据
			if _, ok := h.connections[c]; ok {
				delete(h.connections, c)
				close(c.send)
			}
		case data := <-h.broadcast:
			//处理数据流转, 将数据同步到所有用户
			//c为具体的每一个连接
			for c := range h.connections {
				select {
				case c.send <- data:
				default:
					//防止死循环
					delete(h.connections, c)
					close(c.send)
				}
			}
		}
	}
}
