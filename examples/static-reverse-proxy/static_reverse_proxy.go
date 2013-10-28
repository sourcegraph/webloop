package main

import (
	"flag"
	"github.com/sourcegraph/webloop"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var bind = flag.String("http", ":13000", "HTTP bind address")
var targetURL = flag.String("target", "http://localhost:3000", "base URL of target")
var redirectPrefixesStr = flag.String("redirect-prefixes", "/static,/api,/favicon.ico", "comma-separated list of path prefixes to redirect (not proxy and render)")

func main() {
	flag.Parse()

	log := log.New(os.Stderr, "", 0)

	var redirectPrefixes []string
	if *redirectPrefixesStr != "" {
		redirectPrefixes = strings.Split(*redirectPrefixesStr, ",")
	}

	staticRenderer := &webloop.StaticRenderer{
		TargetBaseURL: *targetURL,
		WaitTimeout:   time.Second * 3,
		Log:           log,
	}
	h := func(w http.ResponseWriter, r *http.Request) {
		for _, rp := range redirectPrefixes {
			if strings.HasPrefix(r.URL.Path, rp) {
				http.Redirect(w, r, *targetURL+r.URL.String(), http.StatusFound)
				return
			}
		}
		staticRenderer.ServeHTTP(w, r)
	}

	http.HandleFunc("/", h)
	log.Printf("Listening on %s and proxying against %s", *bind, *targetURL)
	err := http.ListenAndServe(*bind, nil)
	if err != nil {
		log.Fatalf("ListenAndServe: %s", err)
	}
}
