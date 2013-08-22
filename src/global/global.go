package global

import (
	"bufio"
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
	"github.com/gorilla/mux"
)

// global object
var (
	Conf           *toml.TomlTree            // TOML style config support
	Log            log.Logger                // multiLogger support
	GracefulListen *manners.GracefulListener // shutDown support
	Mux            *mux.Router               // http request router
	Db             db.Database               // Database
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
)

func init() {
	InitGlobal()
}

// Export for Doc viewing easy
func InitGlobal() {
	log.LogLevelToName[log.LOGLEVEL_DEBUG] = "[D]"
	log.LogLevelToName[log.LOGLEVEL_INFO] = "[I]"
	log.LogLevelToName[log.LOGLEVEL_CONFIG] = "[C]"
	log.LogLevelToName[log.LOGLEVEL_WARN] = "[W]"
	log.LogLevelToName[log.LOGLEVEL_ERROR] = "[E]"
	log.LogLevelToName[log.LOGLEVEL_ALERT] = "[A]"
	log.LogLevelToName[log.LOGLEVEL_FATAL] = "[F]"

	Mux = mux.NewRouter()
	Mux.StrictSlash(true)
	LoadConfig()
}
func OpenDb(dss ...db.DataSource) error {
	var ds db.DataSource
	if len(dss) == 0 {
		DbDriver = Conf.GetDefault("db.Driver", "mysql").(string)
		if DbDriver == "" {
			DbDriver = "mysql"
		}
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
			if err == io.EOF {
				break
			}
			line = strings.TrimSpace(line)
			if line == "" {
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

	// [db]
	OpenDb()

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
	ret = append(ret,
		code, " ",
		r.RemoteAddr, " ",
		r.Method, " ",
		fmt.Sprintf("%#v %#v", r.URL.Path, r.UserAgent()),
	)
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
