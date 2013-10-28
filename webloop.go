package webloop

import (
	"errors"
	"github.com/sourcegraph/go-webkit2/webkit2"
	"github.com/sqs/gojs"
	"github.com/sqs/gotk3/glib"
)

// ErrLoadFailed indicates that the View failed to load the requested resource.
var ErrLoadFailed = errors.New("load failed")

// Context stores common settings for a group of Views.
type Context struct{}

// New creates a new Context.
func New() *Context {
	return &Context{}
}

// NewView creates a new View in the context.
func (c *Context) NewView() *View {
	view := make(chan *View, 1)
	glib.IdleAdd(func() bool {
		webView := webkit2.NewWebView()
		v := &View{WebView: webView}
		loadChangedHandler, _ := webView.Connect("load-changed", func(ctx *glib.CallbackContext) {
			loadEvent := webkit2.LoadEvent(ctx.Arg(0).Int())
			switch loadEvent {
			case webkit2.LoadFinished:
				// If we're here, then the load must not have failed, because
				// otherwise we would've disconnected this handler in the
				// load-failed signal handler.
				v.load <- struct{}{}
			}
		})
		webView.Connect("load-failed", func() {
			v.lastLoadErr = ErrLoadFailed
			webView.HandlerDisconnect(loadChangedHandler)
		})
		view <- v
		return false
	})
	return <-view
}

// View represents a WebKit view that can load resources at a given URL and
// query information about them.
type View struct {
	*webkit2.WebView

	load        chan struct{}
	lastLoadErr error

	destroyed bool
}

// Open starts loading the resource at the specified URL.
func (v *View) Open(url string) {
	v.load = make(chan struct{}, 1)
	v.lastLoadErr = nil
	glib.IdleAdd(func() bool {
		if !v.destroyed {
			v.WebView.LoadURI(url)
		}
		return false
	})
}

// Wait waits for the current page to finish loading.
func (v *View) Wait() error {
	<-v.load
	return v.lastLoadErr
}

// URI returns the URI of the current resource in the view.
func (v *View) URI() string {
	uri := make(chan string, 1)
	glib.IdleAdd(func() bool {
		uri <- v.WebView.URI()
		return false
	})
	return <-uri
}

// Title returns the title of the current resource in the view.
func (v *View) Title() string {
	title := make(chan string, 1)
	glib.IdleAdd(func() bool {
		title <- v.WebView.Title()
		return false
	})
	return <-title
}

// EvaluateJavaScript runs the JavaScript in script in the view's context and
// returns the script's result as a Go value.
func (v *View) EvaluateJavaScript(script string) (result *gojs.Value, err error) {
	resultChan := make(chan *gojs.Value, 1)
	errChan := make(chan error, 1)

	glib.IdleAdd(func() bool {
		v.WebView.RunJavaScript(script, func(result *gojs.Value, err error) {
			if err == nil {
				resultChan <- result
			} else {
				errChan <- err
			}
		})
		return false
	})

	select {
	case result = <-resultChan:
		return result, nil
	case err = <-errChan:
		return nil, err
	}
}

// Close closes the view and releases associated resources. Ensure that Close is
// called after all other pending operations on View have returned, or they may
// hang indefinitely.
func (v *View) Close() {
	// TODO(sqs): remove all of the source funcs we added via IdleAdd, etc.,
	// using g_source_remove, to fix "assertion
	// 'WEBKIT_IS_WEB_VIEW(webView) failed" messages.
	v.destroyed = true
	v.Destroy()
}
