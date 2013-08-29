package mysql

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/context"

	g "global"
)

var fetch map[int]func(string) []string

var Lang = map[string]map[int]string{}
var Dict = map[string]map[string]string{}

func init() {
	InitFetch()
	g.OnEvent(Fire)
}

func Fire(code int, r *http.Request, key string, i ...interface{}) bool {
	if key != "dbError" || len(i) == 0 {
		return true
	}
	s, ok := i[0].(string)
	if !ok {
		return true
	}
	pos := strings.Index(s, ":")
	if pos < 6 {
		return true
	}
	code, err := strconv.Atoi(s[6:pos])
	if err != nil {
		return true
	}
	f, ok := fetch[code]
	if !ok {
		return true
	}
	a := f(s[pos:])
	if a == nil {
		return true
	}
	s = g.GetContext(r)["lang"].(string)
	lang, ok := Lang[s]
	if !ok {
		return true
	}
	dict, ok := Dict[s]
	if !ok {
		dict = map[string]string{}
	}
	s = lang[code]
	if s == "" {
		return true
	}
	for i, v := range a {
		if key, ok = dict[strings.Title(v)]; ok {
			v = key
		}
		s = strings.Replace(s, "$"+strconv.Itoa(i), v, -1)
	}
	s = lang[0] + " " + strconv.Itoa(code) + ": " + s
	context.Set(r, g.EventKey("dbError"), s)
	return true
}

func InitFetch() {
	fetch = map[int]func(string) []string{
		1062: func(s string) []string {
			// Error 1062: Duplicate entry '%s' for key '%s'
			a := strings.Split(s, "'")
			if len(a) < 4 {
				return nil
			}
			return []string{a[1], a[3]}
		},
	}
}
