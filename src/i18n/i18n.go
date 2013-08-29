package i18n

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	g "global"
)

var (
	Dict = map[string]map[string]string{}
	Lang = map[string]map[string]string{}
)

func init() {
	g.OnMuxBefore(SetLang)
	g.I18n = I18n
}

func I18n(r *http.Request, s string, i ...interface{}) string {
	lang := g.GetContext(r)["lang"].(string)
	mp, ok := Lang[lang]
	var format string
	if ok {
		format, ok = mp[s]
	}
	if !ok {
		if len(i) == 0 {
			return s
		}
		if strings.Index(s, "%") == -1 {
			return s + fmt.Sprint(i...)
		}
		return fmt.Sprintf(s, i...)
	}
	dict, ok := Dict[lang]

	for k, v := range i {
		key := "$" + strconv.Itoa(k)
		str := fmt.Sprint(v)
		dic, ok := dict[str]
		if ok {
			str = dic
		}
		format = strings.Replace(format, key, str, -1)
	}
	return format
}

func SetLang(wr http.ResponseWriter, r *http.Request) bool {
	// setting lang
	var lang string
	ck, err := r.Cookie("lang")
	if err == nil {
		lang = ck.Value
	} else {
		lang = r.Header.Get("Accept-Language")
		lang = strings.Split(lang, ",")[0]
	}
	if lang == "" {
		lang = g.Lang
	}
	mp := g.GetContext(r)
	mp["lang"] = lang
	return true
}
