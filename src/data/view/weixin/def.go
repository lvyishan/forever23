package weixin

//用户信息
type Wx_userinfo struct {
	Session_id string `gorm:"-" json:"session_id"` //sessionId
	Nick_name  string `json:"nick_name"`           //用户昵称
	Avatar_url string `json:"avatar_url"`          //用户头像，最后一个数值代表正方形头像大小（有0、46、64、96、132数值可选，0代表640*640正方形头像），用户没有头像时该项为空。若用户更换头像，原有头像URL将失效。
	Gender     string `json:"gender"`              //用户的性别，值为1时是男性，值为2时是女性，值为0时是未知
	City       string `json:"city"`                //用户所在城市
	Province   string `json:"province"`            //用户所在省份
	Country    string `json:"country"`             //用户所在国家
	Language   string `json:"language"`            //用户的语言，简体中文为zh_CN
}

//jscode
type GetOpenIDParam struct {
	Jscode string `json:"jscode"`
}

//
type WxSessionKey struct {
	Openid      string `json:"openid"`
	Session_key string `json:"session_key"`
	Unionid     string `json:"unionid"`
}
