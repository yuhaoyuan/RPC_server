package dal

type UserInfo struct {
	Name     string `json:"name"`
	Pwd      string `json:"pwd"`
	NickName string `json:"nick_name"`
	Picture  string `json:"picture"`
}