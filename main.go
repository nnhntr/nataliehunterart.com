package main

import (
	"bytes"
	"embed"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	//go:embed images
	images embed.FS

	//go:embed styles
	styles embed.FS

	//go:embed index.html
	indexHTML []byte
)

func main() {
	start := time.Now()
	mux := http.NewServeMux()
	mux.Handle("/images/", http.FileServer(http.FS(images)))
	mux.Handle("/styles/", http.FileServer(http.FS(styles)))
	mux.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet {
			http.Error(res, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}
		if req.URL.Path != "/" {
			http.Error(res, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		http.ServeContent(res, req, "index.html", start, bytes.NewReader(indexHTML))
	})
	log.Fatalln(http.ListenAndServe(":"+os.Getenv("PORT"), mux))
}
