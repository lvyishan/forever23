package view

//返回消息头
type MapMessageBody struct {
	State bool              `json:"state"`
	Code  int               `json:"code,omitempty"`
	Error string            `json:"error,omitempty"`
	Data  map[string]string `json:"data,omitempty"`
}

//用户缓存结构
type UserCacheBody struct {
	Access_token string
	User_name    string
	Expire_time  int
}
