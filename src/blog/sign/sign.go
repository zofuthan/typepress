package sign

import (
	"net/http"

	. "controllers"
)

func init() {

}
func InitSignIn() {
	sign := Subrouter("signin")
	sign.Methods("get").Handler(Empty)
	sign.Methods("post").Handler(HandlerRouter(func(wr http.ResponseWriter, r *http.Request) {

	}))
}

func InitSignUp() {
	sign := Subrouter("signup")
	sign.Methods("get").Handler(Empty)
	sign.Methods("post").Handler(HandlerRouter(func(wr http.ResponseWriter, r *http.Request) {

	}))
}

func InitSignOut() {
	sign := Subrouter("signout")
	sign.Methods("get").Handler(Empty)
	sign.Methods("post").Handler(HandlerRouter(func(wr http.ResponseWriter, r *http.Request) {

	}))
}
