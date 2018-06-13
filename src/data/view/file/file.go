package file

import (
	"data/config"
	"net/http"
	"os"
	"public/message"
	"public/myhttp"
	"public/mylog"
	"public/tools"
	"regexp"
	"strings"
)

//
var O_file File

//
type File struct {
}

//文件信息
type fileInfo struct {
	File_name string `json:"file_name"` //文件地址
}

/*
	上传文件
*/
func (u *File) Upload(w http.ResponseWriter, r *http.Request) {
	mp := make(map[string]interface{})
	mp["status"] = 0
	//form提交
	if r.Method == "POST" {
		//判断是否存在文件上传
		ret, fileNames := myhttp.UploadMoreFile(r, "")
		if ret == false || len(fileNames) == 0 { //没有上传文件
			mp["msg"] = "请选择需要上传的文件。"
			myhttp.WriteJson(w, mp)
			return
		}

		host := "http://"
		if r.TLS != nil {
			host = "https://"
		}
		host += r.Host

		var file_names []string
		for _, v := range fileNames {
			file_names = append(file_names, host+config.Url_host+v)
		}

		mp["status"] = 1
		if len(file_names) == 1 {
			mp["url"] = file_names[0]
		} else {
			mp["url"] = file_names
		}

		// msg := message.GetSuccessMsg(message.NormalMessageId)
		// msg.Data = file_names
		myhttp.WriteJson(w, mp)

	} else {
		mp["msg"] = "method not allowed"
		myhttp.WriteJson(w, mp) //message.MessageBody{State: false, Code: 405, Error: "method not allowed", Data: nil})
	}
	return
}

/*
	删除单个文件
	url 为本程序运行目录的相对目录
*/
func (u *File) DeleteOne(path string) bool {

	path = tools.GetModelPath() + path
	if !tools.CheckFileIsExist(path) {
		return false
		//err := os.MkdirAll(tools.GetModelPath()+config.File_host+"/"+_dir+"/", os.ModePerm) //生成多级目录
	}

	//删除资源
	err := os.Remove(path)
	if err != nil {
		mylog.Error(err)
		return false
	}

	return true
}

//上传一个
func (u *File) UploadOne(w http.ResponseWriter, r *http.Request) {
	//form提交
	if r.Method == "POST" {
		//判断是否存在文件上传
		ret, fileNames := myhttp.UploadMoreFile(r, "")
		if ret == false || len(fileNames) == 0 { //没有上传文件
			myhttp.WriteJson(w, message.GetErrorMsg(message.InValidOp))
			return
		}

		host := "http://"
		if r.TLS != nil {
			host = "https://"
		}
		host += r.Host

		var file_names []fileInfo
		for _, v := range fileNames {
			file_names = append(file_names, fileInfo{File_name: host + config.Url_host + v})
		}

		msg := message.GetSuccessMsg(message.NormalMessageId)
		msg.Data = file_names
		myhttp.WriteJson(w, msg)

	} else {
		myhttp.WriteJson(w, message.MessageBody{State: false, Code: 405, Error: "method not allowed", Data: nil})
	}
	return
}

/*
	删除单个文件
*/
// func (u *File) DeleteOne(w rest.ResponseWriter, r *rest.Request) {
// 	var req DeleteParam
// 	tools.GetRequestJsonObj(r, &req)
// 	//参数检测
// 	if !tools.CheckParam(req.File_name, req.Business) {
// 		w.WriteJson(message.GetErrorMsg(message.ParameterInvalid, req))
// 		return
// 	}
// 	var db mysqldb.MySqlDB
// 	defer db.OnDestoryDB()
// 	orm := db.OnGetDBOrm(config.GetDbUrl())
// 	//判断文件是否保存
// 	if req.Business == "basic" { //门店基本信息
// 		//判断文件是否存在
// 		var basic view.Basic_info_tbl
// 		orm.Where("picture like ?", "%"+req.File_name+"%").Find(&basic)
// 		if basic.Id == 0 {
// 			w.WriteJson(message.GetErrorMsg(message.FileNotExisted))
// 			return
// 		}
// 		oldPicture := basic.Picture
// 		newsStr, ret := removeStrArray(oldPicture, ",", req.File_name)
// 		if ret {
// 			arrayFile := strings.Split(req.File_name, "/")
// 			dir := arrayFile[3]
// 			file := arrayFile[4]
// 			res := "file/" + dir + "/" + file //源文件路径

// 			err := os.Remove(res) //删除资源
// 			if err != nil {
// 				w.WriteJson(message.GetErrorMsg(message.FileNotExisted))
// 				return
// 			}
// 			w.WriteJson(message.GetSuccessMsg(message.NormalMessageId))
// 			return
// 		} else {
// 			w.WriteJson(message.GetErrorMsg(message.UnknownError))
// 			return
// 		}
// 	} else {
// 		w.WriteJson(message.GetErrorMsg(message.ParameterInvalid))
// 		return
// 	}

// }

func removeStrArray(old, split, deletes string) (news string, ret bool) {
	//判断是否包含要删除的
	b1, _ := regexp.MatchString(deletes, old)
	if b1 {
		//判断是否包含分割号
		b2, _ := regexp.MatchString(split, old)
		if b2 {
			array := strings.Split(old, split)
			//移除元素
			for k, v := range array {
				if v == deletes {
					array = append(array[:k], array[k+1:]...)
				}

			}
			//拼接
			for i := 0; i < len(array); i++ {
				if i != len(array)-1 {
					news += array[i] + split
				} else {
					news += array[i]
				}

			}
			return news, true
		}
	}

	return "", false
}
