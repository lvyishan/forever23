/*
	路由表
*/
package router

import (
	"data/view/blog"
	"data/view/context"
	"data/view/login"
	"data/view/mine"
	"data/view/verify"
	"data/view/weixin"

	"github.com/ant0ine/go-json-rest/rest"
)

//路由列表
var DefaultRouler = []*rest.Route{
	//管理员登录授权 相关
	Post("/login", login.M_login.OnLogin),              //用户登录
	Post("/check_token", login.M_login.CheckToken),     //验证token
	Post("/refresh_token", login.M_login.RefreshToken), //刷新token管理员
	Post("/go/verify", verify.M_verify.GetVerify),      //获取验证码
	//Post("/go/change_pwd", apiserver.ChangePwd),        //修改密码
	Post("/context/update", context.M_context.OnUpdate), //更新
	Post("/context/get", context.M_context.OnGet),       //获取

	//-----------------微信相关接口---------
	Post("/wx/get_openid", weixin.M_weixin.GetOpenID),        //获取openid
	Post("/wx/update_info", weixin.M_weixin.UpdateUserInfo),  //更新用户微信信息
	Post("/mine/get", mine.M_mine.Get),                       //获取个人信息
	Post("/mine/update_sign", mine.M_mine.UpdateSign),        //更新个性签名
	Post("/blog/add", blog.B_blog.Add),                       //发布记录
	Post("/blog/delete", blog.B_blog.Delete),                 //删除记录
	Post("/blog/get_mine_recent", blog.B_blog.GetMineRecent), //获取给定月份记录
	Post("/blog/get_publish_all", blog.B_blog.GetPublishAll), //获取当年广场记录

	//----------------------------end-----
}

/*
	默认初始化
*/
func OnInitRoter(api *rest.Api) (router rest.App, err error) {
	router, err = rest.MakeRouter(
		DefaultRouler...,
	)

	return
}
