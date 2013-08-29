package controllers

import (
	"errors"
	"mime"
	"net/http"
	"os"
	"path/filepath"

	. "global"
)

// Static file, static Gzip file support
func StaticFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		Error(w, r, 404, errors.New("404 Not Found"))
		return
	}
	name, err := filepath.Abs(DocRoot + r.URL.Path)
	if err != nil || len(name) < len(DocRoot) || name[:len(DocRoot)] != DocRoot {
		Error(w, r, 403, errors.New("403 Forbidden"))
		return
	}
	if TryGzFile(w, r, name) {
		return
	}

	stat, err := os.Stat(name)
	if err != nil || stat.IsDir() {
		Error(w, r, 403, errors.New("403 Forbidden"))
		return
	}
	http.ServeFile(w, r, name)
}

func TryGzFile(w http.ResponseWriter, r *http.Request, name string) bool {
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
