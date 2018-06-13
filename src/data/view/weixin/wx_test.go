package weixin

import (
	"data/config"
	"fmt"
	"public/mysqldb"
	"testing"
)

func Test_order(t *testing.T) {
	var info Wx_userinfo
	info.Openid = "3"
	info.Nick_name = "test12"
	info.Avatar_url = "aaa123"
	info.Gender = "1"
	info.City = "2"
	info.Province = ""
	info.Country = "ccb"
	info.Language = "cn"
	var db mysqldb.MySqlDB
	defer db.OnDestoryDB()
	db.OnGetDBOrm(config.GetDbUrl())
	fmt.Println(UpdateInfo(info))
}
