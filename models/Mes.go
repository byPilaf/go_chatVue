package models

const (
	//SysMesType 系统消息类型
	SysMesType = iota
	//UserStatusMesType 用户状态消息类型
	UserStatusMesType
	//UserMesType 用户聊天消息类型
	UserMesType
	//GroupMesType 群聊消息类型
	GroupMesType
	//HiddenMesType 隐藏消息
	HiddenMesType
	//ResponseMesType 响应消息
	ResponseMesType
)

//Mes 消息
type Mes struct {
	FromUserToken string `json:"from_user_token"`
	FromUserName  string `json:"from_user_name"`
	Code          int    `json:"code"`
	MesType       int    `json:"mes_type"`
	Data          string `json:"data"`
}
