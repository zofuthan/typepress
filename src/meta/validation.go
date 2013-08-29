package meta

import (
	"fmt"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"global"
)

// all fields set
var Valid Validators
var rv reflect.Value

func init() {
	InitValid()
}

func InitValid() {
	Valid = Validators{
		Comment_agent: func(s string) bool {
			return false
		},
		Comment_ip: func(s string) bool {
			return false
		},
		Comment_content: func(s string) bool {
			return false
		},
		Comment_count: func(s string) bool {
			return false
		},
		Comment_date: func(s string) bool {
			return false
		},
		Comment_id: func(s string) bool {
			return false
		},
		Comment_parent: func(s string) bool {
			return false
		},
		Comment_type: func(s string) bool {
			return false
		},
		Comment_vetoed: func(s string) bool {
			return false
		},
		Commentmeta_id: func(s string) bool {
			return false
		},
		Count: func(s string) bool {
			return false
		},
		Description: func(s string) bool {
			return false
		},
		Guid: func(s string) bool {
			return false
		},
		Link_description: func(s string) bool {
			return false
		},
		Link_id: func(s string) bool {
			return false
		},
		Link_image: func(s string) bool {
			return false
		},
		Link_notes: func(s string) bool {
			return false
		},
		Link_rating: func(s string) bool {
			return false
		},
		Link_rel: func(s string) bool {
			return false
		},
		Link_rss: func(s string) bool {
			return false
		},
		Link_target: func(s string) bool {
			return false
		},
		Link_title: func(s string) bool {
			return false
		},
		Link_updated: func(s string) bool {
			return false
		},
		Link_url: func(s string) bool {
			return false
		},
		Link_vetoed: func(s string) bool {
			return false
		},
		Member_id: func(s string) bool {
			return false
		},
		Member_status: func(s string) bool {
			return false
		},
		Menu_order: func(s string) bool {
			return false
		},
		Meta_key: func(s string) bool {
			return false
		},
		Meta_value: func(s string) bool {
			return false
		},
		Object_id: func(s string) bool {
			return false
		},
		Parent: func(s string) bool {
			return false
		},
		Ping_vetoed: func(s string) bool {
			return false
		},
		Pinged: func(s string) bool {
			return false
		},
		Post_content: func(s string) bool {
			return false
		},
		Post_date: func(s string) bool {
			return false
		},
		Post_excerpt: func(s string) bool {
			return false
		},
		Post_id: func(s string) bool {
			return false
		},
		Post_mime_type: func(s string) bool {
			return false
		},
		Post_modified: func(s string) bool {
			return false
		},
		Post_parent: func(s string) bool {
			return false
		},
		Post_password: func(s string) bool {
			return false
		},
		Post_slug: func(s string) bool {
			return false
		},
		Post_status: func(s string) bool {
			return false
		},
		Post_title: func(s string) bool {
			return false
		},
		Postmeta_id: func(s string) bool {
			return false
		},
		Site: func(s string) bool {
			return len(s) >= 3 && len(s) <= 20 &&
				IsSite(s) &&
				!IsEnum(s, global.ReserveSite)
		},
		Sitemeta_id: func(s string) bool {
			return false
		},
		Taxonomy: func(s string) bool {
			return false
		},
		Term_id: func(s string) bool {
			return false
		},
		Term_name: func(s string) bool {
			return false
		},
		Term_order: func(s string) bool {
			return false
		},
		Term_slug: func(s string) bool {
			return false
		},
		Termtaxonomy_id: func(s string) bool {
			return false
		},
		To_ping: func(s string) bool {
			return false
		},
		User_date: func(s string) bool {
			return false
		},
		User_id: func(s string) bool {
			return IsId(s)
		},
		User_login: func(s string) bool {
			return IsMd5(s)
		},
		User_nicename: func(s string) bool {
			return IsLenRang(s, 1, 50) && !IsBadName(s)
		},
		User_pass: func(s string) bool {
			return IsMd5(s)
		},
		User_status: func(s string) bool {
			return false
		},
		Usermeta_id: func(s string) bool {
			return false
		},
	}
	rv = reflect.ValueOf(Valid)
}

// 所有的字段验证
type Validators struct {
	Comment_agent, Comment_ip,
	Comment_content, Comment_count, Comment_date, Comment_id,
	Comment_parent, Comment_type, Comment_vetoed, Commentmeta_id,
	Count, Description, Guid, Link_description, Link_id, Link_image, Link_notes,
	Link_rating, Link_rel, Link_rss, Link_target, Link_title, Link_updated, Link_url,
	Link_vetoed,
	Member_id, Member_status, Menu_order, Meta_key, Meta_value, Object_id, Parent,
	Ping_vetoed, Pinged,
	Post_content, Post_date, Post_excerpt, Post_id, Post_mime_type, Post_modified, Post_parent,
	Post_password, Post_slug, Post_status, Post_title,
	Postmeta_id,
	Site, Sitemeta_id,
	Taxonomy, Term_id, Term_name, Term_order, Term_slug,
	Termtaxonomy_id, To_ping, User_date, User_email, User_id, User_login,
	User_nicename, User_pass, User_status, Usermeta_id func(s string) bool
}

// 根据 name 匹配并验证 v, 如果匹配不到或者验证失败, 返回 false
func (p *Validators) IsValid(name string, v interface{}) bool {
	name = strings.Title(strings.ToLower(name))
	fn := rv.FieldByName(name)
	if !fn.IsValid() {
		return false
	}
	vv, ok := v.(string)
	if !ok {
		vv = fmt.Sprint(v)
	}
	in := []reflect.Value{reflect.ValueOf(vv)}
	out := fn.Call(in)
	return out[0].Bool()
}

// 验证一个 map ,如果失败, 返回失败的字段名
func (p *Validators) IsValidMap(names map[string]interface{}) string {
	for name, v := range names {
		if !p.IsValid(name, v) {
			return name
		}
	}
	return ""
}

func IsSite(s string) bool {
	if len(s) == 0 {
		return false
	}
	onlyNumber := true
	for _, b := range s {
		if b >= '0' && b <= '9' {
			continue
		}
		onlyNumber = false
		if (b >= 'a' && b <= 'z') || b == '-' || b == '_' {
			continue
		}
		return false
	}
	if onlyNumber || s[0] == '-' || s[0] == '_' {
		return false
	}
	return true
}

func Is09az(s string) bool {
	for _, b := range s {
		if !(b >= '0' && b <= '9') && !(b >= 'a' && b <= 'z') {
			return false
		}
	}
	return true
}

// 拒绝空值的MD5
func IsMd5(s string) bool {
	if len(s) != 32 || s == "d41d8cd98f00b204e9800998ecf8427e" {
		return false
	}
	for _, b := range s {
		if !(b >= '0' && b <= '9') && !(b >= 'a' && b <= 'f') {
			return false
		}
	}
	return true
}

func IsId(s string) bool {
	_, err := strconv.ParseUint(s, 10, 64)
	return err == nil
}

var badName = regexp.MustCompile(`[\x00-\x2F\x7F　\s\v!"#$%&'()*+,\-./:;<=>?@[\\\]^` + "`{|}~]")

func IsBadName(s string) bool {
	return badName.MatchString(s)
}

var emailPattern = regexp.MustCompile("[\\w!#$%&'*+/=?^_`{|}~-]+(?:\\.[\\w!#$%&'*+/=?^_`{|}~-]+)*@(?:[\\w](?:[\\w-]*[\\w])?\\.)+[a-zA-Z0-9](?:[\\w-]*[\\w])?")

func IsEmail(s string) bool {
	return emailPattern.MatchString(s)
}

var ipPattern = regexp.MustCompile("^((2[0-4]\\d|25[0-5]|[01]?\\d\\d?)\\.){3}(2[0-4]\\d|25[0-5]|[01]?\\d\\d?)$")

func IsIp(s string) bool {
	return ipPattern.MatchString(s)
}

func IsUrl(s string) bool {
	if strings.Index(s, "//") == -1 {
		s = "//" + s
	}
	v, err := url.Parse(s)
	return err == nil && v.Host != ""
}

func IsLenMin(s string, l int) bool {
	l--
	for i, _ := range s {
		if i == l {
			return true
		}
	}
	return false
}

func IsLenMax(s string, l int) bool {
	for i, _ := range s {
		if i == l {
			return false
		}
	}
	return true
}

func IsEnum(s string, slice []string) bool {
	for _, ss := range slice {
		if ss == s {
			return true
		}
	}
	return false
}
func IsLenRang(s string, min, max int) bool {
	if min > max {
		min, max = max, min
	}
	i := 0
	for i, _ = range s {
		if i > max {
			return false
		}
	}
	i++
	return i >= min && i <= max
}
