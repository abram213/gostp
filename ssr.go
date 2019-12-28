package gostp

import (
	"context"
	"net/http"
	"path/filepath"
	"time"

	"github.com/mafredri/cdp"
	"github.com/mafredri/cdp/devtool"
	"github.com/mafredri/cdp/protocol/dom"
	"github.com/mafredri/cdp/protocol/page"
	"github.com/mafredri/cdp/rpcc"

	"github.com/langaner/crawlerdetector"
)

var detector = crawlerdetector.New()

// CachedPage contains info about rendered page
type CachedPage struct {
	HTML           string
	ExpirationTime time.Time
}

// CachedPages contains rendered page
var CachedPages = make(map[string]CachedPage)

// SSR checks if rendered page exist in memory and not expired
func SSR(w http.ResponseWriter, r *http.Request) {
	if detector.IsCrawler(r.Header.Get("User-Agent")) {
		if val, ok := CachedPages[r.RequestURI]; !ok || time.Now().After(val.ExpirationTime) {
			w.Write([]byte(generatePage(r.RequestURI)))
		} else {
			w.Write([]byte(CachedPages[r.RequestURI].HTML))

		}
	} else {
		http.ServeFile(w, r, filepath.Join(Settings.WorkDir, "dist/index.html"))
	}
}

func generatePage(pageURL string) string {
	deleteExpiredPages() //before render new page - delete all expired
	var renderedHTML string
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Use the DevTools HTTP/JSON API to manage targets (e.g. pages, webworkers).
	devt := devtool.New(Settings.SSRdevtools)
	pt, err := devt.Get(ctx, devtool.Page)
	if err != nil {
		pt, _ = devt.Create(ctx)
	}

	// Initiate a new RPC connection to the Chrome DevTools Protocol target.
	conn, _ := rpcc.DialContext(ctx, pt.WebSocketDebuggerURL)
	defer conn.Close() // Leaving connections open will leak memory.

	c := cdp.NewClient(conn)

	// Open a DOMContentEventFired client to buffer this event.
	domContent, _ := c.Page.DOMContentEventFired(ctx)

	defer domContent.Close()

	// Enable events on the Page domain, it's often preferrable to create
	// event clients before enabling events so that we don't miss any.
	c.Page.Enable(ctx)

	// Create the Navigate arguments with the optional Referrer field set.
	navArgs := page.NewNavigateArgs(Settings.SSRhost + pageURL)
	c.Page.Navigate(ctx, navArgs)

	time.Sleep(time.Millisecond * time.Duration(Settings.SSRMillisecondWait))

	// Wait until we have a DOMContentEventFired event.
	domContent.Recv()

	// Fetch the document root node. We can pass nil here
	// since this method only takes optional arguments.
	doc, _ := c.DOM.GetDocument(ctx, nil)

	// Get the outer HTML for the page.
	result, _ := c.DOM.GetOuterHTML(ctx, &dom.GetOuterHTMLArgs{
		NodeID: &doc.Root.NodeID,
	})

	renderedHTML = result.OuterHTML

	CachedPages[pageURL] = CachedPage{HTML: renderedHTML, ExpirationTime: time.Now().Local().Add(time.Second * time.Duration(Settings.SSRexpiration))}
	return renderedHTML
}

func deleteExpiredPages() {
	for index, loopPage := range CachedPages {
		if time.Now().After(loopPage.ExpirationTime) {
			delete(CachedPages, index)
		}
	}
}
