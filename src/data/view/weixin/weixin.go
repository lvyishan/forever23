package weixin

import (
	"public/message"
	"public/tools"

	"github.com/ant0ine/go-json-rest/rest"
)

//
var M_weixin Weixin

//
type Weixin struct {
}

/*
	获取openid
*/
func (o *Weixin) GetOpenID(w rest.ResponseWriter, r *rest.Request) {
	var req GetOpenIDParam
	tools.GetRequestJsonObj(r, &req)

	//参数检测
	if !tools.CheckParam(req.Jscode) {
		w.WriteJson(message.GetErrorMsg(message.ParameterInvalid))
		return
	}
	msg := message.GetSuccessMsg()
	result := GetOpenID(req.Jscode)
	msg.Data = result
	w.WriteJson(msg)
}

/*
	更新用户信息
*/
func (o *Weixin) UpdateUserInfo(w rest.ResponseWriter, r *rest.Request) {
	var req Wx_userinfo
	tools.GetRequestJsonObj(r, &req)

	//参数检测
	if !tools.CheckParam(req.Session_id) {
		w.WriteJson(message.GetErrorMsg(message.ParameterInvalid))
		return
	}
	//获取session_key
	session := GetSessionkey(req.Session_id)
	if len(session.Openid) > 0 {
		result := UpdateInfo(session.Openid, req)
		msg := message.GetSuccessMsg()
		msg.Data = result
		w.WriteJson(msg)
	} else {
		w.WriteJson(message.GetErrorMsg(message.NotFindError, "未授权登录"))
	}
	return
}
