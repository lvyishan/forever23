package blog

import (
	"public/tools"
	"time"
)

type AddParam struct {
	Session_id string `json:"session_id"` //唯一sessionId
	Content    string `json:"content"`
	Is_publish int    `json:"is_publish"` //是否发布
}

type GetRecentParam struct {
	Session_id string     `json:"session_id"` //唯一sessionId
	Date       tools.Time `json:"date"`       //日期 2018-07-01
}

type PublicBolgInfo struct {
	Open_id     string    `json:"open_id"`
	Nick_name   string    `json:"nick_name"`
	Avatar_url  string    `json:"avatar_url"`
	Sign        string    `json:"sign"`
	Content     string    `json:"content"`
	Create_time time.Time `json:"create_time"`
	Blog_id     string    `json:"blog_id"`
	Attach      string    `json:"attach"`
}

type DeleteParam struct {
	Session_id string `json:"session_id"` //唯一sessionId
	Blog_id    string `json:"blog_id"`
}
