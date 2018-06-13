package config

//配置文件信息
type Config struct {
	Cfg_Base
	Wx_AppID      string `json:"wx_appid"`
	Wx_Secret     string `json:"wx_secret"`
	Wx_key        string `json:"wx_key"` //Api_key
	WX_Mch_id     string `json:"wx_mch_id"`
	WX_Notify_url string `json:"wx_notify_url"`
	Sign_api_key  string `json:"sign_api_key,omitempty"` //验签key
}

//微信配置信息
type WxInfo struct {
	AppID      string
	AppSecret  string
	Key        string
	Mch_id     string
	Notify_url string
}

//获取信息
func GetWxInfo() WxInfo {
	return WxInfo{_map.Wx_AppID,
		_map.Wx_Secret,
		_map.Wx_key,
		_map.WX_Mch_id,
		_map.WX_Notify_url}
}

//
type ApiServerUrl struct {
	GetUserInfo string
}

//
func GetFromApiServer() ApiServerUrl {
	return ApiServerUrl{
		_map.Apiserver_url + "/auto_menu/get_one_user"}
}

//
func GetSignKey() string {
	return _map.Sign_api_key
}
