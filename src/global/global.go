package global

import (
	"bufio"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/achun/db"
	"github.com/achun/go-toml"
	"github.com/achun/log"
	"github.com/braintree/manners"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

// global object
var (
	Conf           *toml.TomlTree            // TOML style config support
	Log            log.Logger                // multiLogger support
	GracefulListen *manners.GracefulListener // shutDown support
	Mux            *mux.Router               // http request router
	Db             db.Database               // Database
	DbSql          *sql.DB                   // Raw sql.DB, init by main.go
	SessionStore   sessions.Store            // sessions Store
	DbDriver       string                    // database driver name
	DocRoot        string                    // WEB document root
	BlogId         uint64                    // user_id for top domain blog home
	PWD            string                    // result of os.Getwd()
	Domain, Port   string                    // bind top domain:port
	TplExt         string                    // template file name extension
	TplPath        string                    // template root path
	TplName        string                    // selected name(subdirectory) for template
	TplLayout      string                    // layuout file fullname for template
	ReserveSite    []string                  // reserve name for register site
	Lang           string                    // default language for go code
)

const (
	SessionName = "sessionid"
)

// ContextKey for stores a value in a given request
type ContextKey string

const (
	KeyViewDat    = ContextKey("viewDat") // map[interface{}]interface{}, data for render template
	KeyContext    = ContextKey("context") // context body, do not use for template
	KeySkipRender = "skipRender"          // skip auto render for anything set
	KeyLayoutFile = "layoutFile"          // string, customize name of layout file base on TplPath
	KeyViewDir    = "viewDir"             // string, customize subdirectory for template file base on TplName
	KeyViewFiles  = "viewFiles"           // string, customize name of template file
	KeyWriter     = "writer"              // http.ResponseWriter
	KeyStopRoute  = "stop"                // stop route process
)

// I18n reserved
var I18n func(*http.Request, string, ...interface{}) string

// I18n DbError
var DbError func(code int, r *http.Request, err error) error

// Check Captcha
var CaptchaOk = func(w http.ResponseWriter, r *http.Request) bool {
	return true
}

func init() {
	InitGlobal()
}

// Export for Doc viewing easy
func InitGlobal() {
	if I18n == nil {
		I18n = func(r *http.Request, s string, i ...interface{}) string {
			if len(i) == 0 {
				return s
			}
			if strings.Index(s, "%") == -1 {
				return s + fmt.Sprint(i...)
			}
			return fmt.Sprintf(s, i...)
		}
	}

	if DbError == nil {
		DbError = func(code int, r *http.Request, err error) error {
			if err == nil {
				return nil
			}
			FireEvent(code, r, "dbError", err.Error())
			str, ok := context.GetOk(r, EventKey("dbError"))
			if !ok || str.(string) == "" {
				return err
			}
			return errors.New(str.(string))
		}
	}

	Mux = mux.NewRouter()
	Mux.StrictSlash(true)
	LoadConfig()
}

func OpenDb(dss ...db.DataSource) error {
	var ds db.DataSource
	if DbDriver == "" {
		DbDriver = Conf.GetDefault("db.Driver", "mysql").(string)
	}
	if DbDriver == "" {
		DbDriver = "mysql"
	}
	if len(dss) == 0 {
		ds.Host = Conf.GetDefault("db.Host", "").(string)
		ds.Port = int(Conf.GetDefault("db.Port", 0).(int64))
		ds.Socket = Conf.GetDefault("db.Socket", "").(string)
		ds.Database = Conf.GetDefault("db.Database", "typepress").(string)
		ds.User = Conf.GetDefault("db.User", "").(string)
		ds.Password = Conf.GetDefault("db.Password", "").(string)
	} else {
		ds = dss[0]
	}
	d, err := db.Open(DbDriver, ds)
	if err != nil {
		return err
	}
	Db = d
	return nil
}
func LoadConfig() {
	var err error
	PWD, err = os.Getwd()
	if err != nil {
		panic(err)
	}

	// reserve site list load from conf/reserve_site.txt
	f, err := os.OpenFile("conf/reserve_site.txt", os.O_RDONLY, 0)
	if err == nil {
		defer f.Close()
		br := bufio.NewReader(f)
		for {
			line, err := br.ReadString('\n')
			line = strings.TrimSpace(line)
			if line == "" {
				if err == io.EOF {
					break
				}
				continue
			}
			ReserveSite = append(ReserveSite, line)
		}
	}
	// TOML config
	Conf, err = toml.LoadFile("conf/conf.toml")
	if err != nil {
		panic(err)
	}

	// [blog]
	DocRoot = Conf.GetDefault("blog.root", "").(string)
	if DocRoot == "" {
		DocRoot = "root"
	}
	if !filepath.IsAbs(DocRoot) {
		DocRoot = filepath.Join(PWD, DocRoot)
	}

	DocRoot, err = filepath.Abs(DocRoot)
	if err != nil {
		panic(err)
	}

	stat, err := os.Stat(DocRoot)
	if err != nil {
		panic(err)
	}
	if !stat.IsDir() {
		panic("conf error: must be setting blog.root in conf/conf.toml")
	}

	BlogId = uint64(Conf.GetDefault("blog.userid", int64(0)).(int64))
	Port = strconv.Itoa(int(Conf.GetDefault("blog.port", int64(8080)).(int64)))
	Domain = Conf.GetDefault("blog.domain", "").(string)

	// [template]
	TplExt = Conf.GetDefault("template.ext", ".tmpl").(string)
	if TplExt == "" {
		TplExt = ".tmpl"
	}

	TplPath = Conf.GetDefault("template.path", "views").(string)
	if TplPath == "" {
		TplPath = "views"
	}
	if !filepath.IsAbs(TplPath) {
		TplPath = filepath.Join(PWD, TplPath)
	}
	stat, err = os.Stat(TplPath)
	if err != nil {
		panic(err)
	}
	if !stat.IsDir() {
		panic("conf error: invalid template.path: " + TplPath)
	}

	TplName = Conf.GetDefault("template.name", "typepress.org").(string)
	if TplName == "" {
		TplName = "typepress.org"
	}

	str := filepath.Join(TplPath, TplName)
	stat, err = os.Stat(str)
	if err != nil {
		panic(err)
	}
	if !stat.IsDir() {
		panic("conf error: invalid template.name: " + str)
	}

	TplLayout = Conf.GetDefault("template.layout", "layout.html").(string)
	if TplLayout == "" {
		TplLayout = "layout.html"
	}
	str = filepath.Join(TplPath, TplName, TplLayout)
	stat, err = os.Stat(str)
	if err != nil {
		panic(err)
	}
	if stat.IsDir() {
		panic("conf error: invalid template.layout: " + str)
	}

	// [log]
	logs := Conf.Get("log").(*toml.TomlTree)
	loggers := []log.Logger{}
	if logs != nil {
		keys := logs.Keys()
		for _, key := range keys {
			schema := logs.GetDefault(key+".schema", "").(string)
			if schema == "" {
				continue
			}

			prefix := logs.GetDefault(key+".prefix", "").(string)
			logLevel := logs.GetDefault(key+".level", int64(5)).(int64)
			flags := logs.GetDefault(key+".flags", []interface{}{}).([]interface{})
			equal := logs.GetDefault(key+".equal", false).(bool)

			flag := int(0)
			for _, i := range flags {
				flag = flag | int(i.(int64))
			}
			equalflag := flag
			if equal {
				equalflag = log.LOGLEVEL_EQUAL
			}

			var writer io.WriteCloser
			if schema == ":stderr" {
				writer = os.Stderr
			} else if schema == ":stdout" {
				writer = os.Stdout
			} else {
				if !filepath.IsAbs(schema) {
					schema = filepath.Join(PWD, schema)
				}
				writer, err = NewFileWriter(filepath.Join(schema, key+".log"))
			}
			if err != nil {
				panic(err)
			}
			logger := log.NewLogger(writer, prefix, int(logLevel), flag, equalflag)
			loggers = append(loggers, logger)
		}
	}
	Log = log.MultiLogger(loggers...)
}

// Log easy for Request
func logEasy(code int, r *http.Request, a []interface{}) (ret []interface{}) {
	ret = append(ret, code, " ")
	if r != nil {
		ret = append(ret,
			r.RemoteAddr, " ",
			r.Method, " ",
			fmt.Sprintf("%#v %#v", r.URL.Path, r.UserAgent()),
		)
	}
	for _, i := range a {
		ret = append(ret, fmt.Sprintf(" %#v", i))
	}
	return
}
func LogDebug(code int, r *http.Request, i ...interface{}) {
	Log.Debug(logEasy(code, r, i)...)
}
func LogInfo(code int, r *http.Request, i ...interface{}) {
	Log.Info(logEasy(code, r, i)...)
}
func LogConfig(code int, r *http.Request, i ...interface{}) {
	Log.Config(logEasy(code, r, i)...)
}
func LogWarn(code int, r *http.Request, i ...interface{}) {
	Log.Warn(logEasy(code, r, i)...)
}
func LogError(code int, r *http.Request, i ...interface{}) {
	Log.Error(logEasy(code, r, i)...)
}
func LogAlert(code int, r *http.Request, i ...interface{}) {
	Log.Alert(logEasy(code, r, i)...)
}
func LogFatal(code int, r *http.Request, i ...interface{}) {
	Log.Fatal(logEasy(code, r, i)...)
}

// GetViewDat returns a key-value map stored for a given request with KeyViewDat.
func GetViewDat(r *http.Request) map[interface{}]interface{} {
	return context.Get(r, KeyViewDat).(map[interface{}]interface{})
}

// SetViewDat setting key-value to map stored for a given request with KeyViewDat.
func SetViewDat(r *http.Request, key, value interface{}) {
	mp := context.Get(r, KeyViewDat).(map[interface{}]interface{})
	mp[key] = value
}

// GetContext returns a key-value map stored for a given request with KeyContext.
func GetContext(r *http.Request) map[interface{}]interface{} {
	return context.Get(r, KeyContext).(map[interface{}]interface{})
}

// SetContext setting key-value to map stored for a given request with KeyContext.
func SetContext(r *http.Request, key, value interface{}) {
	mp := context.Get(r, KeyContext).(map[interface{}]interface{})
	mp[key] = value
}

// SkipRender set KeySkipRender of map stored for a given request.
func SkipRender(r *http.Request) {
	mp := context.Get(r, KeyContext).(map[interface{}]interface{})
	mp[KeySkipRender] = true
}

// SkipRender set KeySkipRender of map stored for a given request.
func StopRoute(r *http.Request) {
	mp := context.Get(r, KeyContext).(map[interface{}]interface{})
	mp[KeyStopRoute] = true
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

// GetSession
func GetSession(r *http.Request) (*sessions.Session, error) {
	sess, err := SessionStore.Get(r, SessionName)
	if err != nil {
		r.Header.Del("Cookie")
		sess, err = NewSession(r)
	}
	return sess, err
}

// NewSession
func NewSession(r *http.Request) (*sessions.Session, error) {
	sess, err := SessionStore.New(r, SessionName)
	if err != nil {
		r.Header.Del("Cookie")
		sess, err = SessionStore.New(r, SessionName)
	}
	sess.Options.HttpOnly = true
	return sess, err
}

// SaveSession
func SaveSession(r *http.Request, wr http.ResponseWriter, sess *sessions.Session) error {
	return sess.Save(r, wr)
}
