package main

import (
	"bytes"
	"embed"
	"html/template"
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

	//go:embed index.html.template
	indexHTML string
)

type PageData struct {
	Year int
}

func main() {
	templates := template.Must(template.New("").Parse(indexHTML))
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
		var buf bytes.Buffer
		if err := templates.Execute(&buf, PageData{
			Year: time.Now().Year(),
		}); err != nil {
			http.Error(res, "failed to render page", http.StatusInternalServerError)
			return
		}
		http.ServeContent(res, req, "index.html", start, bytes.NewReader(buf.Bytes()))
	})
	log.Println("starting server...")
	log.Fatalln(http.ListenAndServe(":"+os.Getenv("PORT"), mux))
}
