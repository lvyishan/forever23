package weixin

import (
	"crypto/md5"
	"data/config"
	"data/view"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"public/message"
	"public/mycache"
	"public/mylog"
	"public/mysqldb"
	"public/tools"
	"strconv"
	"strings"
	"time"

	wxpay "gopkg.in/go-with/wxpay.v1"
)

const (
	// 微信支付商户平台证书路径
	certFile   = "/cert/apiclient_cert.pem"
	keyFile    = "/cert/apiclient_key.pem"
	rootcaFile = "/cert/rootca.pem"
)

/*
	通过jscode获取用户openID
*/
func GetOpenID(jscode string) string {
	info := config.GetWxInfo()
	var url = "https://api.weixin.qq.com/sns/jscode2session?appid=" + info.AppID + "&secret=" +
		info.AppSecret + "&js_code=" + jscode + "&grant_type=authorization_code&trade_type=JSAPI"

	resp, e := http.Get(url)
	if e != nil {
		mylog.Error(e)
		return ""
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		mylog.Error(e)
		return ""
	}
	return SaveSessionKey(body)
}

/*
	保存用户sessionkey
*/
func SaveSessionKey(body []byte) (sessionId string) {
	v := WxSessionKey{}
	json.Unmarshal(body, &v)
	if len(v.Openid) > 0 {
		var db mysqldb.MySqlDB
		defer db.OnDestoryDB()
		orm := db.OnGetDBOrm(config.GetDbUrl())

		var tmp view.Wx_userinfo_tbl
		orm.Where("openid = ?", v.Openid).Find(&tmp)

		if tmp.Id == 0 {
			var t view.Wx_userinfo_tbl
			t.Openid = v.Openid
			orm.Create(&t)
		}
		sessionId = tools.UniqueId()
		//保存缓存
		cache := mycache.OnGetCache("session_key")
		cache.Add(sessionId, v, 30*time.Minute)
	}
	return
}

/*
	获取session_key信息
*/
func GetSessionkey(sessionId string) (info WxSessionKey) {
	cache := mycache.OnGetCache("session_key")
	tp, b := cache.Value(sessionId)
	if b {
		info = tp.(WxSessionKey)
	}
	return
}

/*
	更新用户信息
*/
func UpdateInfo(openId string, info Wx_userinfo) bool {
	var db mysqldb.MySqlDB
	defer db.OnDestoryDB()
	orm := db.OnGetDBOrm(config.GetDbUrl())

	err := orm.Where("openid = ?", openId).Updates(&info).Error
	if err != nil {
		return false
	} else {
		return true
	}

}

/*
	统一下单接口
	open_id:用户唯一标识
	price : 预支付价钱
	price_body ： 支付描述
	order_id ： 商户订单号
*/
func Unifiedorder(open_id string, price int, price_body, order_id, client_ip string) message.MessageBody {
	if !tools.CheckParam(open_id, order_id) || price <= 0 { //参数检测
		return message.GetErrorMsg(message.ParameterInvalid)
	}

	//微信支付统一下单
	wx_info := config.GetWxInfo()

	c := wxpay.NewClient(wx_info.AppID, wx_info.AppSecret, wx_info.Key)
	// 附着商户证书
	err := c.WithCert(tools.GetModelPath()+certFile, tools.GetModelPath()+keyFile, tools.GetModelPath()+rootcaFile)
	if err != nil {
		log.Fatal(err)
	}

	params := make(wxpay.Params)
	// 查询企业付款接口请求参数
	params.SetString("appid", c.AppId)
	params.SetString("mch_id", c.MchId)
	params.SetString("body", price_body)
	params.SetInt64("total_fee", int64(price*10))
	params.SetString("spbill_create_ip", client_ip)
	params.SetString("notify_url", wx_info.Notify_url)
	params.SetString("trade_type", "JSAPI")
	params.SetString("openid", open_id)
	params.SetString("nonce_str", tools.GetRandomString(32)) // 随机字符串
	params.SetString("out_trade_no", order_id)               // 商户订单号
	params.SetString("sign", c.Sign(params))                 // 签名 c.Sign(params)

	// 查询企业付款接口请求URL
	url := "https://api.mch.weixin.qq.com/pay/unifiedorder"

	// 发送查询企业付款请求
	ret, err := c.Post(url, params, true)
	if err != nil {
		log.Fatal(err)
	}
	//-----------------------end

	//ret["order_id"] = order_tbl.Order_id
	fmt.Println(ret)

	if ret["result_code"] == "SUCCESS" {
		dd := make(map[string]string)
		dd["timeStamp"] = strconv.FormatInt(tools.GetUtcTime(time.Now()), 10)
		dd["nonceStr"] = tools.GetRandomString(32)
		dd["package"] = "prepay_id=" + ret["prepay_id"]
		dd["signType"] = "MD5"
		dd["paySign"] = "MD5"
		//appId=wxd678efh567hg6787&nonceStr=5K8264ILTKCH16CQ2502SI8ZNMTM67VS&package=prepay_id=&signType=MD5&timeStamp=1490840662&key=qazwsxedcrfvtgbyhnujmikolp111111
		str := "appId=" + wx_info.AppID + "&nonceStr=" + dd["nonceStr"] + "&package=" + dd["package"] + "&signType=MD5&timeStamp=" + dd["timeStamp"] + "&key=" + wx_info.Key
		by := md5.Sum([]byte(str))
		dd["paySign"] = strings.ToUpper(fmt.Sprintf("%x", by))
		dd["order_id"] = order_id

		msg := message.GetSuccessMsg(message.NormalMessageId)
		msg.Data = dd
		return msg
	}

	msg := message.GetErrorMsg(message.InValidOp)
	msg.Data = ret
	return msg
}

//统一查询接口
func OnSelectData(open_id, order_id string) (bool, message.MessageBody) {
	if !tools.CheckParam(open_id, order_id) { //参数检测
		return false, message.GetErrorMsg(message.ParameterInvalid)
	}

	b := false

	//查询支付信息-
	wx_info := config.GetWxInfo()

	c := wxpay.NewClient(wx_info.AppID, wx_info.Mch_id, wx_info.Key)
	// 附着商户证书
	err := c.WithCert(tools.GetModelPath()+certFile, tools.GetModelPath()+keyFile, tools.GetModelPath()+rootcaFile)
	if err != nil {
		log.Fatal(err)
	}

	params := make(wxpay.Params)
	// 查询企业付款接口请求参数
	params.SetString("appid", c.AppId)
	params.SetString("mch_id", c.MchId)
	params.SetString("out_trade_no", order_id)               //商户订单号
	params.SetString("nonce_str", tools.GetRandomString(32)) // 随机字符串
	params.SetString("sign", c.Sign(params))                 // 签名 c.Sign(params)

	// 查询企业付款接口请求URL
	url := "https://api.mch.weixin.qq.com/pay/orderquery"

	// 发送查询企业付款请求
	ret, err := c.Post(url, params, true)
	if err != nil {
		log.Fatal(err)
	}
	//-----------------------end

	msg := message.GetSuccessMsg(message.NormalMessageId)

	/*
		SUCCESS—支付成功
		REFUND—转入退款
		NOTPAY—未支付
		CLOSED—已关闭
		REVOKED—已撤销（刷卡支付）
		USERPAYING--用户支付中
		PAYERROR--支付失败(其他原因，如银行返回失败)
	*/
	if ret["trade_state"] == "SUCCESS" {
		b = false
	} else if ret["trade_state"] == "CLOSED" || ret["trade_state"] == "PAYERROR" {
		b = false
	} else if ret["trade_state"] == "NOTPAY" {
		b = true
	}

	if b { //成功，更新，续费会员
		msg.State = true
	} else {
		msg.State = false
	}

	msg.Data = ret
	return b, msg
}
