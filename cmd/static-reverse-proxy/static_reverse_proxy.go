package main

import (
	"flag"
	"fmt"
	"github.com/sourcegraph/webloop"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var bind = flag.String("http", ":13000", "HTTP bind address")
var targetURL = flag.String("target", "http://localhost:3000", "base URL of target")
var waitTimeout = flag.Duration("wait", time.Second*3, "timeout for pages to set window.$renderStaticReady")
var returnUnfinishedPages = flag.Bool("unfinished", false, "return unfinished pages at wait timeout (instead of erroring)")
var removeScripts = flag.Bool("remove-scripts", false, "remove <script> tags")
var redirectPrefixesStr = flag.String("redirect-prefixes", "/static,/api,/favicon.ico", "comma-separated list of path prefixes to redirect to the target (not proxy and render)")

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr)
		fmt.Fprintf(os.Stderr, "static-reverse-proxy proxies a dynamic JavaScript application and serves\n")
		fmt.Fprintf(os.Stderr, "an equivalent statically rendered HTML website to clients. It uses a headless\n")
		fmt.Fprintf(os.Stderr, "WebKit browser instance to render the static HTML.\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n\n")
		fmt.Fprintf(os.Stderr, "\tstatic-reverse-proxy [options]\n\n")
		fmt.Fprintf(os.Stderr, "The options are:\n\n")
		flag.PrintDefaults()
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr)
		fmt.Fprintf(os.Stderr, "Example usage:\n\n")
		fmt.Fprintf(os.Stderr, "\tTo proxy a dynamic application at http://example.com and serve an equivalent\n")
		fmt.Fprintf(os.Stderr, "\tstatically rendered HTML website on http://localhost:13000\n")
		fmt.Fprintf(os.Stderr, "\t    $ static-reverse-proxy -target=http://example.com -bind=:13000\n\n")
		fmt.Fprintf(os.Stderr, "Notes:\n\n")
		fmt.Fprintf(os.Stderr, "\tBecause a headless WebKit instance is used, your $DISPLAY must be set. Use\n")
		fmt.Fprintf(os.Stderr, "\tXvfb if you are running on a machine without an existing X server. See\n")
		fmt.Fprintf(os.Stderr, "\thttps://sourcegraph.com/github.com/sourcegraph/webloop/readme for more info.\n")
		fmt.Fprintln(os.Stderr)
		os.Exit(1)
	}
	flag.Parse()

	log := log.New(os.Stderr, "", 0)

	var redirectPrefixes []string
	if *redirectPrefixesStr != "" {
		redirectPrefixes = strings.Split(*redirectPrefixesStr, ",")
	}

	staticRenderer := &webloop.StaticRenderer{
		TargetBaseURL:         *targetURL,
		WaitTimeout:           *waitTimeout,
		ReturnUnfinishedPages: *returnUnfinishedPages,
		RemoveScripts:         *removeScripts,
		Log:                   log,
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
