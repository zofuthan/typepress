package controllers

import (
	. "github.com/achun/typepress/global"
	"mime"
	"net/http"
	"os"
	"path/filepath"
)

func init() {
	Mux.NotFoundHandler = http.HandlerFunc(StaticFile)
}

// Static file, static gzip file support
func StaticFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "404 page not found", 404)
		return
	}
	name, err := filepath.Abs(DocRoot + r.URL.Path)
	if err != nil || len(name) < len(DocRoot) || name[:len(DocRoot)] != DocRoot {
		http.Error(w, "403 Forbidden", 403)
		return
	}
	if tryGzFile(w, r, name) {
		return
	}

	stat, err := os.Stat(name)
	if err != nil || stat.IsDir() {
		http.Error(w, "403 Forbidden", 403)
		return
	}
	http.ServeFile(w, r, name)
}

func tryGzFile(w http.ResponseWriter, r *http.Request, name string) bool {
	ext := filepath.Ext(name)
	if ext == ".gz" {
		ext = filepath.Ext(name[:len(name)-3])
	} else {
		name = name + ".gz"
	}
	stat, err := os.Stat(name)
	if err != nil || stat.IsDir() {
		return false
	}
	w.Header().Set("Content-Encoding", "gzip")
	ctype := mime.TypeByExtension(ext)
	if ctype == "" {
		ctype = "application/octet-stream"
	}
	w.Header().Set("Content-Type", ctype)
	http.ServeFile(w, r, name)
	return true
}
