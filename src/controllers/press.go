package controllers

import (
	"errors"
	"io"
	"net/http"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/achun/template"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"

	. "global"
)

func init() {
	InitPress()
}

// Export for Doc viewing easy
func InitPress() {
	Mux.NotFoundHandler = http.HandlerFunc(StaticFile)
}

// wrapper for mux Handler
type HandlerMux func(http.ResponseWriter, *http.Request)

func (f HandlerMux) ServeHTTP(wr http.ResponseWriter, r *http.Request) {
	wr.Header().Set("Server", "TypePress")
	// init KeyContext for controller
	contextSet := make(map[interface{}]interface{})
	contextSet[KeyWriter] = wr
	context.Set(r, KeyContext, contextSet)

	if !FireMuxBefore(wr, r) {
		return
	}

	if contextSet[KeyStopRoute] == nil {
		f(wr, r)
	}

	FireMuxAfter(wr, r)
}

// wrapper for mux router Handler
type HandlerRouter func(http.ResponseWriter, *http.Request)

// ServeHTTP are wrapper for Handler.
// Preset storage container context. use builtin function Get for fetch context.
// Render to View.
func (f HandlerRouter) ServeHTTP(wr http.ResponseWriter, r *http.Request) {
	defer func() {
		if e := recover(); e != nil {
			_, file, line, _ := runtime.Caller(4)
			http.Error(wr, "500 Internal Server Error: "+file+" "+strconv.Itoa(line), 500)
			FireEvent(500, r, "PanicOnHandlerRouter", file, line)
		} else {
			FireEvent(200, r, "deferHandlerRouter")
		}
	}()
	// auto install
	if BlogId == 0 {
		if r.URL.Path != "/install/" {
			http.Redirect(wr, r, "/install/", 302)
			return
		}
	} else if r.URL.Path == "/install/" {
		http.Redirect(wr, r, "/", 302)
		return
	}

	var dir, layout, viewfile string

	// init view dat association to request for template
	dat := map[interface{}]interface{}{}
	context.Set(r, KeyViewDat, dat)
	contextSet := GetContext(r)

	// filter
	if !FireRouteBefore(wr, r) {
		return
	}

	if contextSet[KeyStopRoute] != nil {
		return
	}
	// call controller handler
	f(wr, r)

	if contextSet[KeyStopRoute] != nil {
		return
	}

	// after filter
	if !FireRouteAfter(wr, r) {
		return
	}

	if contextSet[KeyStopRoute] != nil {
		return
	}

	// skip render ify anything setting of KeySkipRender
	if contextSet[KeySkipRender] != nil {
		FireEvent(200, r, "skipRenderOnHandlerRouter")
		return
	}

	// lookup subdirectory of template files
	i, ok := contextSet[KeyViewDir]
	if ok {
		dir = i.(string)
	}
	if dir == "" {
		match := new(mux.RouteMatch)
		if Mux.Match(r, match) {
			dir = match.Route.GetName()
		}
	}

	if dir == "" {
		paths := strings.Split(r.URL.Path, "/")
		dir = filepath.Join(paths...)
	}

	if dir == "" {
		http.Error(wr, "500 Internal Server Error", 500)
		FireEvent(500, r, "directoryOfTemplateIsNotSetOnHandlerRouter")
		return
	}

	// lookup layout file
	i, ok = contextSet[KeyLayoutFile]
	if ok {
		layout, _ = i.(string)
	}

	// lookup template file
	i, ok = contextSet[KeyViewFiles]
	if ok {
		viewfile, _ = i.(string)
	}

	if layout == "" {
		// XMLHttpRequest support
		if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
			layout = "content.html" // {{content}}
		} else {
			layout = TplLayout
		}
	}
	layout = filepath.Join(TplPath, TplName, layout)

	// setting viewfile from request Method, if empty
	if viewfile == "" {
		viewfile = strings.ToLower(r.Method) + TplExt
	}
	viewfile = layout + ";" + viewfile
	viewfiles := strings.Split(viewfile, ";")
	if len(viewfiles) == 1 {
		http.Error(wr, "500 Internal Server Error", 500)
		FireEvent(500, r, "viewFilesIsNotSetOnHandlerRouter")
		return
	}
	for i, v := range viewfiles {
		if !filepath.IsAbs(v) {
			viewfiles[i] = filepath.Join(TplPath, TplName, dir, v)
		}
	}

	// viewfiles[0] is layout, viewfiles[1] is content
	viewfile = filepath.Base(viewfiles[1])
	tpl := template.New("")
	tpl.Builtin()
	tpl.Funcs(map[string]interface{}{
		"request": func() *http.Request {
			return r
		},
		"content": func() string {
			err := tpl.ExecuteTemplate(wr, viewfile, dat)
			if err != nil {
				return err.Error()
			}
			return ""
		},
	})

	_, err := tpl.ParseFiles(viewfiles...)
	if err == nil {
		err = tpl.Execute(wr, dat)
	}
	if err != nil {
		Error(wr, r, 500, err, viewfiles)
	} else {
		FireRenderAfter(r, err)
	}
}

// wrapper for mux router Handler
func Handle(handler http.Handler) *mux.Route {
	return Mux.NewRoute().Handler(HandlerRouter(handler.ServeHTTP))
}

// wrapper for mux router HandlerFunc
func HandleFunc(f func(http.ResponseWriter, *http.Request)) *mux.Route {
	return Mux.NewRoute().Handler(HandlerRouter(f))
}

// wrapper for mux router HandlerFunc and auto ParseForm
func HandlerParseForm(f func(http.ResponseWriter, *http.Request)) HandlerRouter {
	return HandlerRouter(func(w http.ResponseWriter, r *http.Request) {
		if Error(w, r, 400, r.ParseForm()) {
			return
		}
		f(w, r)
	})
}

// wrapper for mux router HandlerFunc and auto ParseForm with CAPTCHA
func HandlerCaptcha(f func(http.ResponseWriter, *http.Request)) HandlerRouter {
	return HandlerRouter(func(w http.ResponseWriter, r *http.Request) {
		if Error(w, r, 400, r.ParseForm()) || !CaptchaOk(w, r) {
			return
		}
		f(w, r)
	})
}

// HandleSignin for user must be logged in to pass
func HandleSignin(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		sess, err := GetSession(r)
		if Error(w, r, 500, err) {
			StopRoute(r)
			return
		}
		if sess.Values["user"] == nil {
			StopRoute(r)
			Error(w, r, 403, errors.New("you must be signed"))
			return
		}
		f(w, r)
	}
}

// wrapper for http.Error
func Error(wr http.ResponseWriter, r *http.Request, code int, err error, i ...interface{}) bool {
	if err != nil {
		SkipRender(r)
		if FireEvent(code, r, "handlerError", err.Error(), i) {
			http.Error(wr, err.Error(), code)
		}
		return true
	}
	return false
}

// wrapper for io.Writer
func Write(wr io.Writer, str string) (int, error) {
	return wr.Write([]byte(str))
}

// wrapper for http.Redirect, for ajax support
func Redirect(wr http.ResponseWriter, r *http.Request, urlStr string, code int) {
	if r.Header.Get("X-Requested-With") != "XMLHttpRequest" {
		http.Redirect(wr, r, urlStr, 302)
		return
	}
	wr.Header().Set("Content-Type", "text/plain; charset=utf-8")
	wr.WriteHeader(code)
	wr.Write([]byte(urlStr))
}

// Subrouter warrper for Mux.Subrouter()
func Subrouter(name string) *mux.Router {
	return Mux.Path("/" + name + "/").Name(name).Subrouter()
}

var (
	// Empty is nothing Handler
	Empty = HandlerRouter(func(wr http.ResponseWriter, r *http.Request) {})
)
