package config

import (
	"fmt"
	"os"
	"strings"
)

const (
	test_file = `
{
"以下为通用配置":"--------------------------------------------",
"serial_number":"1.0",
"service_name":"hotelserver",
"service_displayname":"hotelserver",
"sercice_desc":"hotelserver",
"http_port":"8011",
"https_port":"8012",
"db_url":"lvyishan:lvyishan^qwer!@tcp(rm-uf6m89g4jeob19s2mo.mysql.rds.aliyuncs.com:3306)/hotel_db?charset=utf8&parseTime=True&loc=Local",
"leveldb_dir":"./database",
"token_type":"nomal",
"app_id":"wwwthings",
"app_secret":"4EE0A9A43B9B911C067BEE5CC50A9972",
"is_dev":true,
"oauth2_url":"http://192.168.0.152:7001/oauth2/api/v1",
"apiserver_url":"http://192.168.0.152:7031/apiserver/api/v1",
"register_url":"http://192.168.0.152:7021/register/api/v1",
"es_addr_url":"http://192.168.0.152:9200",
"以下为逻辑配置":"--------------------------------------------",

"wx_appid":"wxe12fb98bdebdff01",
"wx_secret":"2888ca8d524d160620759d436e46de1f",
"wx_key":"3",
"wx_mch_id":"4",
"sign_api_key":"apiserver"
}
`
)

//判断是否在测试环境下使用
func IsRunTesting() bool {
	if len(os.Args) > 1 {
		fmt.Println(os.Args[1])
		return strings.HasPrefix(os.Args[1], "-test")
	}
	return false
}
