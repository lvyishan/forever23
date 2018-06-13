/*
	构造restful基础框架类
*/
package data

import (
	"data/config"
	"data/view/file"
	"log"
	"net/http"

	"data/router"
	"public/tools"

	"github.com/ant0ine/go-json-rest/rest"
)

var api *rest.Api

//导出路由
type ApiRoot struct {
}

//构建路由
func (*ApiRoot) OnCreat() {
	api = rest.NewApi()
	if config.OnIsDev() {
		api.Use(rest.DefaultDevStack...) //DefaultProdStack
	} else {
		api.Use(rest.DefaultProdStack...)
	}
	router, err := router.OnInitRoter(api)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	api.SetApp(router)

	//文件路径
	http.Handle(config.Url_host+config.File_host+"/",
		http.StripPrefix(config.Url_host+config.File_host,
			http.FileServer(http.Dir(tools.GetCurrentDirectory()+config.File_host))))
	//统一上传文件
	http.HandleFunc(config.Url_host+"/api/v1/file/upload", file.O_file.Upload)
	http.HandleFunc(config.Url_host+"/api/v1/file/uploadone", file.O_file.UploadOne)
	//指定api默认路由
	http.Handle(config.Url_host+"/api/", http.StripPrefix(config.Url_host+"/api", api.MakeHandler()))
}
