package blog

import (
	"data/config"
	"data/view"
	"data/view/weixin"
	"public/message"
	"public/mysqldb"
	"public/tools"
	"time"

	"github.com/ant0ine/go-json-rest/rest"
)

var B_blog Blog

//
type Blog struct {
}

/*
	发布微记
*/
func (o *Blog) Add(w rest.ResponseWriter, r *rest.Request) {
	var req AddParam
	tools.GetRequestJsonObj(r, &req)
	//参数检测
	if !tools.CheckParam(req.Session_id, req.Content) || req.Is_publish == 0 {
		w.WriteJson(message.GetErrorMsg(message.ParameterInvalid))
		return
	}
	//获取session_key
	session := weixin.GetSessionkey(req.Session_id)
	if len(session.Openid) > 0 {
		var db mysqldb.MySqlDB
		defer db.OnDestoryDB()
		orm := db.OnGetDBOrm(config.GetDbUrl())

		var info view.Blog_info_tbl
		blogId := createBlogId(session.Openid)
		info.Blog_id = blogId
		info.Open_id = session.Openid
		info.Content = req.Content
		info.Create_time = time.Now()
		info.Is_publish = req.Is_publish
		err := orm.Create(&info).Error
		if err != nil {
			w.WriteJson(message.GetErrorMsg(message.UnknownError))
			return
		}
		msg := message.GetSuccessMsg(message.NormalMessageId)
		msg.Data = blogId
		w.WriteJson(msg)
	} else {
		w.WriteJson(message.GetErrorMsg(message.NotFindError, "未授权登录"))
	}
	return
}

/*
	获取当月记录
*/
func (o *Blog) GetMineRecent(w rest.ResponseWriter, r *rest.Request) {
	var req GetRecentParam
	tools.GetRequestJsonObj(r, &req)
	//参数检测
	if !tools.CheckParam(req.Session_id) {
		w.WriteJson(message.GetErrorMsg(message.ParameterInvalid))
		return
	}
	//获取session_key
	session := weixin.GetSessionkey(req.Session_id)
	if len(session.Openid) > 0 {
		var db mysqldb.MySqlDB
		defer db.OnDestoryDB()
		orm := db.OnGetDBOrm(config.GetDbUrl())

		mouth0 := time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, time.Local)
		nextMouth0 := mouth0.AddDate(0, 1, 0)
		if !req.Date.Time.IsZero() {
			mouth0 = time.Date(req.Date.Time.Year(), req.Date.Time.Month(), 1, 0, 0, 0, 0, time.Local)
			nextMouth0 = mouth0.AddDate(0, 1, 0)
		}
		var blogs []view.Blog_info_tbl
		orm.Where("open_id = ? and is_vaild = ? and create_time >= ? and create_time < ?",
			session.Openid, 0, mouth0, nextMouth0).Order("create_time desc").Find(&blogs)
		msg := message.GetSuccessMsg(message.NormalMessageId)
		msg.Data = blogs
		w.WriteJson(msg)
	} else {
		w.WriteJson(message.GetErrorMsg(message.NotFindError, "未授权登录"))
	}
	return
}

/*
	删除记录
*/
func (o *Blog) Delete(w rest.ResponseWriter, r *rest.Request) {
	var req DeleteParam
	tools.GetRequestJsonObj(r, &req)
	//参数检测
	if !tools.CheckParam(req.Session_id, req.Blog_id) {
		w.WriteJson(message.GetErrorMsg(message.ParameterInvalid))
		return
	}
	//获取session_key
	session := weixin.GetSessionkey(req.Session_id)
	if len(session.Openid) > 0 {
		var db mysqldb.MySqlDB
		defer db.OnDestoryDB()
		orm := db.OnGetDBOrm(config.GetDbUrl())

		err := orm.Table("blog_info_tbl").Where("blog_id = ?", req.Blog_id).Update("is_vaild", -1).Error
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

/*
	获取广场记录
*/
func (o *Blog) GetPublishAll(w rest.ResponseWriter, r *rest.Request) {
	var req GetRecentParam
	tools.GetRequestJsonObj(r, &req)
	//参数检测
	if !tools.CheckParam(req.Session_id) {
		w.WriteJson(message.GetErrorMsg(message.ParameterInvalid))
		return
	}
	//获取session_key
	session := weixin.GetSessionkey(req.Session_id)
	if len(session.Openid) > 0 {
		var db mysqldb.MySqlDB
		defer db.OnDestoryDB()
		orm := db.OnGetDBOrm(config.GetDbUrl())
		var info []PublicBolgInfo

		year0 := time.Date(time.Now().Year(), 1, 1, 0, 0, 0, 0, time.Local)
		nextYear0 := year0.AddDate(1, 0, 0)

		orm.Table("blog_info_tbl as b").Select("b.open_id,b.content,b.create_time,b.attach,b.blog_id,w.nick_name,w.avatar_url,w.sign").Joins("left join wx_userinfo_tbl as w on b.open_id = w.openid").
			Where("b.is_publish = ? and b.is_vaild = ? and b.create_time >= ? and b.create_time < ?", 1, 0, year0, nextYear0).Order("b.create_time desc").Scan(&info)
		msg := message.GetSuccessMsg(message.NormalMessageId)
		msg.Data = info
		w.WriteJson(msg)
	} else {
		w.WriteJson(message.GetErrorMsg(message.NotFindError, "未授权登录"))
	}
	return
}
