package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/alipourhabibi/go-examples/gzip-rest/handlers"
)

func main() {
	mux := http.NewServeMux()
	gzh := handlers.NewGzipHandler()
	mux.HandleFunc("/images/", gzh.MiddleWare(http.HandlerFunc(download)))
	err := http.ListenAndServe(":8080", mux)
	log.Fatal(err)
}

func download(w http.ResponseWriter, r *http.Request) {
	varString := strings.TrimPrefix(r.URL.Path, "/images/")
	vars := strings.Split(varString, "/")
	name := vars[0]
	image, err := ioutil.ReadFile("images/" + name)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("image not found"))
		return
	}
	w.Write(image)
}
