package mine

type GetParam struct {
	Session_id string `json:"session_id"` //唯一sessionId
}

type MineInfo struct {
	Open_id    string `json:"open_id"`
	Nick_name  string `json:"nick_name"`
	Avatar_url string `json:"avatar_url"`
	Sign       string `json:"sign"`
	Blog_sum   int    `json:"blog_sum"`
}

type UpdateSignParam struct {
	Session_id string `json:"session_id"` //唯一sessionId
	New_sign   string `json:"new_sign"`   //新签名
}
