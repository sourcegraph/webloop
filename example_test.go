package webloop_test

import (
	"fmt"
	"os"
	"runtime"

	"github.com/gotk3/gotk3/gtk"
	"github.com/sourcegraph/webloop"
)

func Example() {
	gtk.Init(nil)
	go func() {
		runtime.LockOSThread()
		gtk.Main()
	}()

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
