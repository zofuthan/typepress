package cn

import (
	"i18n"
)

func init() {
	Lang()
	Dict()
}
func Dict() {
	i18n.Dict["zh-cn"] = map[string]string{
		"User_login": "登录名",
		"User_id":    "用户ID",
		"Term_name":  "分类名",
		"Term_slug":  "缩略名",
		"Site":       "域名",
		"Users":      "用户",
		"Append()":   "添加",
		"Find()":     "查询",
		"FindAll()":  "查询",
		"Query()":    "查询",
		"Update()":   "更新",
	}
}

func Lang() {
	i18n.Lang["zh-cn"] = map[string]string{
		"%s must be exists on %s %s":    "'$1' '$2' 操作必须提供 '$0' 数据",
		"not enough arguments on %s %s": "'$0' '$1' 操作缺少足够的参数",
		"%s is invalid on %s %s":        "'$1' '$2' 参数 '$0' 无效",
	}
}
