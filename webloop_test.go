package webloop

import (
	"github.com/sqs/gotk3/gtk"
	"net/http"
	"runtime"
	"testing"
)

func init() {
	gtk.Init(nil)
	go func() {
		runtime.LockOSThread()
		gtk.Main()
	}()
}

var ctx Context

func TestNew(t *testing.T) {
	New()
}

func TestContext_NewView(t *testing.T) {
	view := ctx.NewView()
	defer view.Close()
}

func TestView_Open(t *testing.T) {
	runtime.LockOSThread()

	setup()
	defer teardown()
	mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {})

	view := ctx.NewView()
	defer view.Close()

	view.Open(server.URL)
}

func TestView_Wait(t *testing.T) {
	runtime.LockOSThread()

	setup()
	defer teardown()

	view := ctx.NewView()
	defer view.Close()

	loaded := false
	mux.HandleFunc("/abc", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("abc"))
		loaded = true
	})

	url := server.URL + "/abc"
	view.Open(server.URL + "/abc")
	view.Wait()

	if gotURI := view.URI(); url != gotURI {
		t.Errorf("want URI %q, got %q", url, gotURI)
	}

	if !loaded {
		t.Error("!loaded")
	}
}

func TestView_Wait_multi(t *testing.T) {
	runtime.LockOSThread()

	setup()
	defer teardown()

	view := ctx.NewView()
	defer view.Close()

	loaded := 0
	mux.HandleFunc("/abc", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("abc"))
		loaded++
	})
	mux.HandleFunc("/xyz", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("xyz"))
		loaded++
	})

	url1 := server.URL + "/abc"
	view.Open(url1)
	view.Wait()
	if gotURI := view.URI(); url1 != gotURI {
		t.Errorf("want URI %q, got %q", url1, gotURI)
	}

	url2 := server.URL + "/xyz"
	view.Open(url2)
	view.Wait()

	if gotURI := view.URI(); url2 != gotURI {
		t.Errorf("want URI %q, got %q", url2, gotURI)
	}

	if wantLoaded := 2; wantLoaded != loaded {
		t.Errorf("want loaded == %d, got %d", wantLoaded, loaded)
	}
}
