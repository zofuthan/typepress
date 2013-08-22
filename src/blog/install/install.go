package install

import (
	"database/sql"
	"errors"
	"github.com/achun/db"
	"io/ioutil"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"

	. "controllers"
	. "global"
	"models"
)

func init() {
	if BlogId == 0 {
		InitInstall()
	}
}

// Export for Doc viewing easy
func InitInstall() {
	inst := Mux.Path("/install/").Name("install").Subrouter()
	inst.Methods("get").Handler(HandlerRouter(func(w http.ResponseWriter, r *http.Request) {
		dat := GetViewDat(r)
		dat["Port"] = Port
		dat["Domain"] = Domain
	}))
	inst.Methods("post").Handler(HandlerParseForm(true, func(w http.ResponseWriter, r *http.Request) {
		var err error
		// setting domain
		mp, ok := FetchValueExists(r.Form, "Domain", "Port")
		if !ok {
			http.Error(w, "Invalid arguments for Domain", 400)
			return
		}

		ul, err := url.Parse(mp["Domain"].(string))
		if Error(w, r, 400, err) {
			return
		}
		Domain = ul.Host
		Conf.Set("blog.domain", Domain)
		port, err := strconv.ParseUint(mp["Port"].(string), 10, 16)
		if port != 0 {
			Port = strconv.Itoa(int(port))
			Conf.Set("blog.port", Port)
		}

		// setting database
		mp, ok = FetchValueExists(r.Form, "Host", "Database", "User", "Password")
		if !ok {
			http.Error(w, "Invalid arguments for DB", 400)
			return
		}
		if mp["Database"].(string) == "" {
			mp["Database"] = "typepress"
		}
		if mp["Host"].(string) == "" {
			mp["Host"] = "127.0.0.1:3306"
		}
		if mp["User"].(string) == "" {
			mp["User"] = "root"
		}
		mp["Port"] = ""
		mp["Socket"] = ""
		if len(mp["Host"].(string)) > 0 && mp["Host"].(string)[0] == '/' {
			mp["Socket"], mp["Host"] = mp["Host"], ""
		} else {
			a := strings.Split(mp["Host"].(string), ":")
			mp["Host"] = a[0]
			if len(a) > 1 {
				mp["Port"] = a[1]
			}
		}
		mp["Port"], _ = strconv.Atoi(mp["Port"].(string))

		ds := db.DataSource{
			Host:     mp["Host"].(string),
			Port:     mp["Port"].(int),
			Socket:   mp["Socket"].(string),
			User:     mp["User"].(string),
			Password: mp["Password"].(string),
		}
		switch DbDriver {
		case "mysql":
			ds.Database = "mysql"
		case "postgresql":
			ds.Database = "pg_database"
		case "sqlite":
			ds.Database = mp["Database"].(string)
		case "mongo":
			ds.Database = ""
		}

		err = OpenDb(ds)
		if Error(w, r, 400, err) {
			return
		}

		// auto create database
		err = InitDb(mp["Database"].(string))
		if err == nil {
			err = Db.Use(mp["Database"].(string))
		}
		if err == nil {
			Conf.Set("db.Host", ds.Host)
			Conf.Set("db.Port", ds.Port)
			Conf.Set("db.Socket", ds.Socket)
			Conf.Set("db.User", ds.User)
			Conf.Set("db.Password", ds.Password)
			Conf.Set("db.Database", mp["Database"].(string))
			err = Conf.SaveToFile()
		}
		if Error(w, r, 400, err) {
			return
		}

		// setting first user
		mp, ok = FetchValueExists(r.Form, "User_login", "User_pass")
		if !ok {
			Error(w, r, 400, errors.New("Invalid arguments for admin user"))
			return
		}
		mp["Site"] = ""
		mp["User_nicename"] = "TypePress"
		ids, err := models.Users.Append(mp)
		if Error(w, r, 400, err) {
			return
		}
		if len(ids) != 1 {
			Error(w, r, 500, errors.New("Users.Append() returns length of []db.Id != 1"), ids)
		}
		BlogId, err = strconv.ParseUint(string(ids[0]), 10, 64)
		if Error(w, r, 500, err) {
			return
		}
		Conf.Set("blog.userid", int64(BlogId))
		if Error(w, r, 500, Conf.SaveToFile()) {
			return
		}
		Redirect(w, r, "/", 200)
	}))
}

func InitDb(name string) (err error) {
	if DbDriver == "mongo" {
		return errors.New("does not support MongoDB")
	}
	str, err := loadSql()
	if err == nil {
		sqldb, ok := Db.Driver().(*sql.DB)
		if !ok {
			return errors.New("can not get Db.Driver()")
		}
		if DbDriver != "sqlite" {
			_, err = sqldb.Exec("CREATE DATABASE IF NOT EXISTS `" + name + "` DEFAULT CHARACTER SET utf8 ")
			if err == nil {
				_, err = sqldb.Exec("USE `" + name + "`")
			}
			if err != nil {
				return
			}
		}
		sqls := strings.Split(str, "CREATE TABLE ")
		for _, sql := range sqls {
			sql = strings.TrimSpace(sql)
			if sql == "" {
				continue
			}
			_, err = sqldb.Exec("CREATE TABLE " + sql)
			if err != nil {
				return
			}
		}
	}
	return
}

func loadSql() (string, error) {
	filename := filepath.Join(PWD, "conf", DbDriver+".sql")
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(buf), nil
}
