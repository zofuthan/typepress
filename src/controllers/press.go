package controllers

import (
	"net/http"
	"path/filepath"
	"runtime"
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

// ContextKey for stores a value in a given request
type ContextKey int

const (
	KeyViewDat    ContextKey = iota // map[interface{}]interface{}, data for render template
	KeyContext                      // context body, do not use for template
	KeySkipRender                   // skip auto render for anything set
	KeyLayoutFile                   // string, customize name of layout file base on TplPath
	KeyViewDir                      // string, customize subdirectory for template file base on TplName
	KeyViewFiles                    // string, customize name of template file
)

// wrapper for mux Handler
type HandlerMux func(http.ResponseWriter, *http.Request)

func (f HandlerMux) ServeHTTP(wr http.ResponseWriter, r *http.Request) {
	wr.Header().Set("Server", "TypePress")
	if !FireMuxBefore(wr, r) {
		return
	}
	f(wr, r)
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
			http.Error(wr, "500 Internal Server Error", 500)
			_, file, line, _ := runtime.Caller(4)
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

	// setting KeyContext for controller
	contextSet := make(map[interface{}]interface{})
	context.Set(r, KeyContext, contextSet)
	var dir, layout, viewfile string

	// init view dat association to request for template
	dat := map[interface{}]interface{}{}
	context.Set(r, KeyViewDat, dat)

	// filter
	if !FireRouteBefore(wr, r) {
		return
	}

	// call controller handler
	f(wr, r)

	// after filter
	if !FireRouteAfter(wr, r) {
		return
	}

	// skip render ify anything setting of KeySkipRender
	i, ok := contextSet[KeySkipRender]
	if ok {
		FireEvent(200, r, "skipRenderOnHandlerRouter")
		return
	}

	// lookup subdirectory of template files
	i, ok = contextSet[KeyViewDir]
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
func HandlerParseForm(render bool, f func(http.ResponseWriter, *http.Request)) HandlerRouter {
	return HandlerRouter(func(w http.ResponseWriter, r *http.Request) {
		if !render {
			SkipRender(r)
		}
		if Error(w, r, 400, r.ParseForm()) {
			return
		}
		f(w, r)
	})
}

// GetViewDat returns a key-value map stored for a given request.
func GetViewDat(r *http.Request) map[interface{}]interface{} {
	return context.Get(r, KeyViewDat).(map[interface{}]interface{})
}

// SetViewDat setting key-value to map stored for a given request.
func SetViewDat(r *http.Request, key, value interface{}) {
	mp := GetViewDat(r)
	dat := mp[KeyViewDat].(map[interface{}]interface{})
	dat[key] = value
}

// SkipRender set KeySkipRender of map stored for a given request.
func SkipRender(r *http.Request) {
	mp := context.Get(r, KeyContext).(map[interface{}]interface{})
	mp[KeySkipRender] = true
}

// ViewFiles set KeyViewFiles of map stored for a given request.
func ViewFiles(r *http.Request, filenames string) {
	mp := context.Get(r, KeyContext).(map[interface{}]interface{})
	mp[KeyViewFiles] = filenames
}

// LayoutFile set KeyLayoutFile of map stored for a given request.
func LayoutFile(r *http.Request, filename string) {
	mp := context.Get(r, KeyContext).(map[interface{}]interface{})
	mp[KeyLayoutFile] = filename
}

// wrapper for http.Error
func Error(w http.ResponseWriter, r *http.Request, code int, err error, i ...interface{}) bool {
	if err != nil {
		http.Error(w, err.Error(), code)
		FireEvent(code, r, "handlerError", err.Error(), i)
		return true
	}
	return false
}

// wrapper for http.Redirect, for ajax support
func Redirect(w http.ResponseWriter, r *http.Request, urlStr string, code int) {
	if r.Header.Get("X-Requested-With") != "XMLHttpRequest" {
		http.Redirect(w, r, urlStr, code)
		return
	}
	w.Write([]byte(urlStr))
}

// Subrouter warrper for Mux.Subrouter()
func Subrouter(name string) *mux.Route {
	return Mux.Path("/" + name + "/").Name(name).Subrouter()
}

var (
	// Empty is nothing Handler
	Empty = HandlerRouter(func(wr http.ResponseWriter, r *http.Request) {})
)
