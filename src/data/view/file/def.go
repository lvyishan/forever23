package file

//删除文件参数
type DeleteParam struct {
	File_name string `json:"file_name"` //文件地址
	Business  string `json:"business"`  //文件所在业务
}
