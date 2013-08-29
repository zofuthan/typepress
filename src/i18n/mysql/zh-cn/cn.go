package cn

import (
	"i18n/mysql"
)

func init() {
	Lang()
	Dict()
}
func Dict() {
	mysql.Dict["zh-cn"] = map[string]string{
		"User_login": "登录名",
		"User_id":    "用户ID",
		"Term_name":  "分类名",
		"Term_slug":  "缩略名",
		"Site":       "域名",
	}
}

func Lang() {
	mysql.Lang["zh-cn"] = map[int]string{
		0:    "错误编号",
		1062: "'$1' 已被占用请换一个",
	}
}
