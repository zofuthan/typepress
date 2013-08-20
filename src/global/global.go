package global

import (
	"github.com/achun/go-toml"
	"github.com/achun/log"
	"github.com/braintree/manners"
	"github.com/gorilla/mux"
	"os"
	"path/filepath"
)

// global object
var (
	Conf           *toml.TomlTree
	Log            log.Logger
	GracefulListen *manners.GracefulListener
	Mux            *mux.Router
	DocRoot        string
)

// Regester func() on shutDown
func OnShutDown(f func()) {
	onShutDown = append(onShutDown, f)
}

// Fire func() on shutDown
func FireShutDown() {
	for _, f := range onShutDown {
		f()
	}
	onShutDown = []func(){}
}

var (
	err        error
	onShutDown []func()
)

func init() {
	Mux = mux.NewRouter()
	LoadConfig()
}

func LoadConfig() {
	Conf, err = toml.LoadFile("conf/conf.toml")
	if err != nil {
		panic(err)
	}
	DocRoot = Conf.GetDefault("blog.root", "").(string)
	if DocRoot == "" {
		DocRoot, err = os.Getwd()
		if err != nil {
			panic(err)
		}
		DocRoot = DocRoot + "/root"
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
	// setting Log
	logs := Conf.Get("log").(*toml.TomlTree)
	loggers := []log.Logger{}
	if logs != nil {
		keys := logs.Keys()
		for _, k := range keys {
			key := "log." + k + "."
			schema := Conf.GetDefault(key+"schema", "").(string)
			if schema == "" {
				continue
			}
			prefix := Conf.GetDefault(key+"prefix", "").(string)
			logLevel := Conf.GetDefault(key+"level", 5).(int64)
			flags := Conf.GetDefault(key+"flags", []interface{}{}).([]interface{})
			flag := int(0)
			for _, i := range flags {
				flag = flag | int(i.(int64))
			}
			writer, err := NewFileWriter(schema, k+".log")
			if err != nil {
				panic(err)
			}
			logger := log.NewLogger(writer, prefix, int(logLevel), flag)
			loggers = append(loggers, logger)
		}
	}
	Log = log.MultiLogger(loggers...)
}
