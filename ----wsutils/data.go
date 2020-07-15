package wsutils

//Data 传输的数据
type Data struct {
	IP       string   `json:"ip"`
	Type     string   `json:"type"`
	From     string   `json:"from"`
	Content  string   `json:"content"`
	User     string   `json:"user"`
	UserList []string `json:"user_list"`
}
