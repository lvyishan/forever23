package mine

import (
	"data/config"
	"data/view/weixin"
	"public/message"
	"public/mysqldb"
	"public/tools"

	"github.com/ant0ine/go-json-rest/rest"
)

//
var M_mine Mine

//
type Mine struct {
}

/*
	获取个人信息
*/
func (o *Mine) Get(w rest.ResponseWriter, r *rest.Request) {
	var req GetParam
	tools.GetRequestJsonObj(r, &req)
	//参数检测
	if !tools.CheckParam(req.Session_id) {
		w.WriteJson(message.GetErrorMsg(message.ParameterInvalid))
		return
	} //获取session_key
	session := weixin.GetSessionkey(req.Session_id)
	if len(session.Openid) > 0 {
		var db mysqldb.MySqlDB
		defer db.OnDestoryDB()
		orm := db.OnGetDBOrm(config.GetDbUrl())

		var mine MineInfo
		orm.Where("openid = ?", session.Openid).Find(&mine)
		//总记录数
		orm.Raw("select count(*) as blog_sum FROM Blog_info_tbl WHERE open_id = ? and is_vaild = 0",
			session.Openid).Scan(&mine)

		msg := message.GetSuccessMsg(message.NormalMessageId)
		msg.Data = mine
		w.WriteJson(msg)
	} else {
		w.WriteJson(message.GetErrorMsg(message.NotFindError, "未授权登录"))
	}
	return
}

/*
	更新个性签名
*/
func (o *Mine) UpdateSign(w rest.ResponseWriter, r *rest.Request) {
	var req UpdateSignParam
	tools.GetRequestJsonObj(r, &req)
	//参数检测
	if !tools.CheckParam(req.Session_id, req.New_sign) {
		w.WriteJson(message.GetErrorMsg(message.ParameterInvalid))
		return
	} //获取session_key
	session := weixin.GetSessionkey(req.Session_id)
	if len(session.Openid) > 0 {
		var db mysqldb.MySqlDB
		defer db.OnDestoryDB()
		orm := db.OnGetDBOrm(config.GetDbUrl())

		err := orm.Table("wx_userinfo_tbl").Where("openid = ?", session.Openid).Update("sign", req.New_sign).Error
		if err != nil {
			w.WriteJson(message.GetErrorMsg(message.UnknownError))
			return
		}
		w.WriteJson(message.GetSuccessMsg())
	} else {
		w.WriteJson(message.GetErrorMsg(message.NotFindError, "未授权登录"))
	}
	return
}
