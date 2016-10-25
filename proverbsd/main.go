package main

import (
	"encoding/json"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

var proverbs = map[string]string{
	"share-memory":         "Don't communicate by sharing memory, share memory by communicating.",
	"concurrency":          "Concurrency is not parallelism.",
	"channel-mutexes":      "Channels orchestrate; mutexes serialize.",
	"interface-size":       "The bigger the interface, the weaker the abstraction.",
	"useful-zero":          "Make the zero value useful.",
	"empty-interrface":     "interface{} says nothing.",
	"gofmt":                "Gofmt's style is no one's favorite, yet gofmt is everyone's favorite.",
	"dependency":           "A little copying is better than a little dependency.",
	"syscall":              "Syscall must always be guarded with build tags.",
	"cgo-build":            "Cgo must always be guarded with build tags.",
	"cgo":                  "Cgo is not Go.",
	"unsafe":               "With the unsafe package there are no guarantees.",
	"clever":               "Clear is better than clever.",
	"reflection":           "Reflection is never clear.",
	"errors":               "Errors are values.",
	"handle-errors":        "Don't just check errors, handle them gracefully.",
	"design-name-document": "Design the architecture, name the components, document the details.",
	"documentation":        "Documentation is for users.",
	"panic":                "Don't panic.",
}

const htmlSrc = `<!DOCTYPE html>
<html>
  <head>
    <title>Go Proverbs</title>
    <link href="https://fonts.googleapis.com/css?family=Vesper+Libre" rel="stylesheet">
    <style type="text/css">
      html, body, .quote { height: 100%; }
      .quote {
	      display: flex;
	      justify-content: center;
	      align-items: center;
	      font-size: 2.25vw;
	      text-align: center;
	      font-family: 'Vesper Libre', serif;
      }
      h2 {
	      width: 80%;
	      padding-bottom: 10%;
      }
      a, a:visited, a:hover {
	      text-decoration: none;
	      color: black;
      }
    </style>
  </head>
  <body>
    <div class="quote">
        <h2><a href="?quote={{ .ID }}" title="link to this quote">&ldquo;{{ .Quote }}&rdquo;</a></h2>
    </div>
  </body>
</html>`

var keys []string

var tmpl *template.Template

type quote struct {
	ID    string
	Quote string
}

func serve(w http.ResponseWriter, r *http.Request) {
	switch r.Header.Get("Return-Error") {
	case "internal":
		w.WriteHeader(http.StatusInternalServerError)
		return
	case "bad-request":
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	rand.Seed(time.Now().UnixNano())
	key := r.URL.Query().Get("quote")
	if key == "" {
		key = keys[rand.Intn(len(keys))]
	}
	q := quote{
		ID:    key,
		Quote: proverbs[key],
	}
	if q.Quote == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if r.Header.Get("Accept") != "application/json" {
		err := tmpl.Execute(w, q)
		if err != nil {
			log.Printf("Error serving request: %+v\n", err)
		}
		return
	}
	data, err := json.Marshal(q)
	if err != nil {
		log.Printf("Error formatting JSON: %+v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = w.Write(data)
	if err != nil {
		log.Printf("Error writing JSON: %+v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func main() {
	for k, _ := range proverbs {
		keys = append(keys, k)
	}
	tmpl = template.Must(template.New("html").Parse(htmlSrc))
	http.HandleFunc("/", serve)
	err := http.ListenAndServe("0.0.0.0:9005", nil)
	if err != nil {
		log.Printf("Failed to listen and serve: %+v\n", err)
		os.Exit(1)
	}
}
