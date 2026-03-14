package main

import (
	"embed"
	"flag"
	"html/template"
	"log"
	"net/http"
)

//go:embed index.tmpl
var tmplFS embed.FS

func main() {
	addr := flag.String("addr", ":9000", "Server address")
	flag.Parse()

	// parse template from embedded filesystem
	tmpl := template.Must(template.ParseFS(tmplFS, "index.tmpl"))

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err := tmpl.Execute(w, nil)
		if err != nil {
			http.Error(w, "Template error", http.StatusInternalServerError)
			log.Println(err)
		}
	})

	log.Printf("Starting Simple CORS example server on %s", *addr)
	log.Fatal(http.ListenAndServe(*addr, mux))
}
