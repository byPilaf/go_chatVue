package controller

import (
	"encoding/json"
	"html/template"
	"net/http"
	"webchat/models"
	"webchat/utils"
)

//LoginHandler 用户登录/注册
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	//定义回复信息
	if r.Method == "GET" {
		//解析模板
		temp := template.Must(template.ParseFiles("views/login.html"))
		//执行
		temp.Execute(w, "")
	} else {
		//回复消息
		var reMes models.Mes
		reMes.MesType = models.ResponseMesType

		//接收json
		decoder := json.NewDecoder(r.Body)
		var loginUser models.User
		err := decoder.Decode(&loginUser)
		if err != nil || len(loginUser.Name) < 4 || len(loginUser.Pass) < 4 {
			reMes.Code = 401
			reMes.Data = "请输入有效值"
			reJSON, _ := json.Marshal(reMes)
			w.Write(reJSON)
			// fmt.Fprintf(w, `{"code":401, "msg":"请输入有效值"}`)
			return
		}

		//查询数据库
		// utils.DbMysql.Where("name = ?", &loginUser.Name).First(&user)
		// if user.ID != 0 {
		// 	if loginUser.Pass != user.Pass {
		// 		reMes.Code = 403
		// 		reMes.Data = "密码错误"
		// 		reJSON, _ := json.Marshal(reMes)
		// 		w.Write(reJSON)
		// 		// fmt.Fprintf(w, `{"code":403, "msg":"密码错误"}`)
		// 		return
		// 	}
		// } else {
		// 	//注册
		// 	utils.DbMysql.NewRecord(loginUser)
		// 	utils.DbMysql.Create(&loginUser)
		// 	user = loginUser
		// }

		//查询列表
		if u, ok := models.UsersMap[loginUser.Name]; ok {
			//登陆
			if u.Pass != loginUser.Pass {
				reMes.Code = 403
				reMes.Data = "密码错误"
				reJSON, _ := json.Marshal(reMes)
				w.Write(reJSON)
				return
			}
		} else {
			//注册
			models.UsersMap[loginUser.Name] = &loginUser
		}

		//在用户列表中判断是否已经登陆
		for _, onlineUser := range models.OnlineUsersMap {
			if onlineUser.Name == loginUser.Name {
				//已经在线
				onlineUser.OffLine() //顶掉
			}
		}

		//生成user token
		userToken := utils.GetToken()
		loginUser.Token = userToken                   //添加token
		models.OnlineUsersMap[userToken] = &loginUser //添加到在线用户列表

		reMes.Code = 200
		reMes.Data = "登陆成功"
		reMes.FromUserToken = userToken
		reJSON, _ := json.Marshal(reMes)
		w.Write(reJSON)
		// mes := `{"code":200,"msg":"登陆成功","user_token":"` + userToken + `"}`
		// fmt.Fprintf(w, mes)
		return
	}
}
