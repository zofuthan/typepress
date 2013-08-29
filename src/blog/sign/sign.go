package sign

import (
	"encoding/gob"
	"net/http"

	. "controllers"
	g "global"
	"meta"
	"models"
)

func init() {
	gob.Register(meta.Users{})
	InitSignIn()
	InitSignUp()
	InitSignOut()
}
func InitSignIn() {
	sign := Subrouter("signin")
	sign.Methods("get").Handler(Empty)
	sign.Methods("post").Handler(HandlerCaptcha(func(wr http.ResponseWriter, r *http.Request) {
		g.SkipRender(r)
		if !g.ValuesExists(r.Form, "User_login", "User_pass") {
			http.Error(wr, "", 403)
			return
		}

		item, err := models.Users.Find(r, 0)
		if Error(wr, r, 409, err) {
			return
		}
		if len(item) == 0 {
			Write(wr, g.I18n(r, "Incorrect User_login or User_pass"))
			return
		}
		sess, err := g.NewSession(r)
		if Error(wr, r, 500, err) {
			return
		}
		user := meta.Users{}
		g.ToStruct(item, &user)
		sess.Values["user"] = user
		if Error(wr, r, 409, g.SaveSession(r, wr, sess)) {
			return
		}
		Redirect(wr, r, "/", 200)
	}))
}

func InitSignUp() {
	sign := Subrouter("signup")
	sign.Methods("get").Handler(Empty)
	sign.Methods("post").Handler(HandlerCaptcha(func(wr http.ResponseWriter, r *http.Request) {
		g.SkipRender(r)
		_, err := models.Users.Append(r)
		if Error(wr, r, 409, err) {
			return
		}
		g.FilterValues(r.Form, "User_login", "User_pass")

		item, err := models.Users.Find(r, 0)
		if Error(wr, r, 409, err) {
			return
		}
		if len(item) == 0 {
			Write(wr, g.I18n(r, "Incorrect User_login or User_pass"))
			return
		}
		sess, err := g.NewSession(r)
		if Error(wr, r, 500, err) {
			return
		}
		user := meta.Users{}
		g.ToStruct(item, &user)
		sess.Values["user"] = user
		if Error(wr, r, 409, g.SaveSession(r, wr, sess)) {
			return
		}
		Redirect(wr, r, "/", 200)
	}))
}

func InitSignOut() {
	HandleFunc(func(wr http.ResponseWriter, r *http.Request) {
		g.SkipRender(r)
		co, err := r.Cookie(g.SessionName)
		if err != nil || co.Value == "" {
			Redirect(wr, r, "/", 200)
			return
		}

		sess, err := g.GetSession(r)
		if Error(wr, r, 409, err) {
			Redirect(wr, r, "/", 200)
			return
		}

		sess.Options.MaxAge = -1
		g.SaveSession(r, wr, sess)
		Redirect(wr, r, "/", 200)
	}).Path("/signout/")
}
