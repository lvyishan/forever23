package blog

import (
	"public/tools"
	"time"
)

func createBlogId(channal string) string {
	tmp := time.Now().Format("20060102150405") + channal
	tmp = tmp + tools.GetRandomString(32-len(tmp))
	return tmp
}
