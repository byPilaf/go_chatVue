package controller

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"webchat/models"
)

//IndexHandler 首页控制器
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		//解析模板
		temp := template.Must(template.ParseFiles("views/index.html"))
		//执行
		temp.Execute(w, "")
		return
	}

	decoder := json.NewDecoder(r.Body)
	var token map[string]string
	err := decoder.Decode(&token)
	if err != nil {
		return
	}

	//响应消息
	var resMes models.Mes

	//当前在线列表
	_, ok := models.OnlineUsersMap[token["user_token"]]
	if !ok {
		resMes.Code = 403
		resMes.Data = "请登陆"
		reJSON, _ := json.Marshal(resMes)
		w.Write(reJSON)
		// fmt.Fprintf(w, `{"code":"403","msg":"请重新登陆"}`)
		return
	}

	resMes.Code = 200
	resMes.Data = "welcome"
	reJSON, _ := json.Marshal(resMes)
	w.Write(reJSON)
	// fmt.Fprintf(w, `{"code":"200","msg":"welcome"}`)
	return
}

//GetOnlineUserListHandler 获取在线用户列表
func GetOnlineUserListHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		var userList []string
		for _, user := range models.OnlineUsersMap {
			userList = append(userList, user.Name)
		}
		data, err := json.Marshal(userList)
		if err == nil {
			resData := string(data)
			fmt.Fprintf(w, `{"code":200,"data":`+resData+`}`)
			return
		}
	}
	//响应消息
	var resMes models.Mes
	resMes.Code = 403
	resMes.MesType = models.ResponseMesType
	resMes.Data = "获取失败"
	reJSON, _ := json.Marshal(resMes)
	w.Write(reJSON)
	// fmt.Fprintf(w, `{"code":403,"msg":"获取失败"}`)
	return
}

//GetUserMesHandler 获取用户发送来的数据
func GetUserMesHandler(w http.ResponseWriter, r *http.Request) {
	//响应消息
	reMes := models.Mes{
		Code:    403,
		MesType: models.SysMesType,
		Data:    "请求错误",
	}

	if r.Method == "POST" {
		//接收json
		decoder := json.NewDecoder(r.Body)
		var message models.Mes //接收到的消息
		err := decoder.Decode(&message)
		if err == nil {
			if models.OnlineUsersMap[message.FromUserToken] != nil {
				reMes.Code = 200
				reMes.MesType = models.ResponseMesType
				reJSON, _ := json.Marshal(reMes)

				//根据消息类型处理不同方法
				switch message.MesType {
				case models.GroupMesType:
					//创建群发的websocekt消息
					var sendMes models.WebSocketMessage
					sendMes.FromUserName = models.OnlineUsersMap[message.FromUserToken].Name
					sendMes.Data = message.Data
					sendMes.MesType = message.MesType
					sendMes.FromUserToken = message.FromUserToken
					sendMes.SendAllUserMes()

					w.Write(reJSON)
					return
				}
			} else {
				reMes.Code = 401
				reMes.Data = "请重新登陆"
			}
		}
	}

	reJSON, _ := json.Marshal(reMes)
	w.Write(reJSON)
	return
}
