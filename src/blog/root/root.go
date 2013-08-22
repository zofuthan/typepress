package root

import (
	"net/http"

	. "controllers"
)

func init() {
	InitRoot()
}

// Export for Doc viewing easy
func InitRoot() {
	HandleFunc(func(w http.ResponseWriter, r *http.Request) {
	}).Path("/").Name("root").Methods("get")
}