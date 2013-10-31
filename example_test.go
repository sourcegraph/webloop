package webloop_test

import (
	"fmt"
	"github.com/sourcegraph/webloop"
	"os"
)

func Example() {
	ctx := webloop.New()
	view := ctx.NewView()
	defer view.Close()
	view.Open("http://google.com")
	err := view.Wait()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load URL: %s", err)
		os.Exit(1)
	}
	res, err := view.EvaluateJavaScript("document.title")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to run JavaScript: %s", err)
		os.Exit(1)
	}
	fmt.Printf("JavaScript returned: %q\n", res)
	// output:
	// JavaScript returned: "Google"
}
