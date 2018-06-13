package main

import (
	"data"
	"data/config"
	"log"
	"net/http"
	"os"
	"public/server"
)

// 回调函数
func CallBack() {
	var apiroot data.ApiRoot
	apiroot.OnCreat()

	//验证微信公众号
	//	http.HandleFunc("/hotelserver/verify", weixin.M_weixin.VerifyWx)

	//https 支持(单开一个线程)
	//	go func() {
	//		log.Println("https is running at " + config.GetServerHttpsPort() + " port.")
	//		log.Fatal(http.ListenAndServeTLS(":"+config.GetServerHttpsPort(), tools.GetModelPath()+"/pem/cacert.pem", tools.GetModelPath()+"/pem/privatekey.pem", nil))
	//	}()
	//---------------------------end

	//启动http
	log.Println("http is running at " + config.GetServerPort() + " port.")
	log.Fatal(http.ListenAndServe(":"+config.GetServerPort(), nil))
	//---------------------end
}

func main() {
	if config.OnIsDev() || len(os.Args) == 0 {
		CallBack()
	} else {
		server.OnStart(CallBack)
	}
}
