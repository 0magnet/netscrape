// Package netscrape is a web browser written in Go/wasm. Its chrome — a tab
// strip, address bar, back/forward/reload — is DOM built with syscall/js; each
// tab is a sandboxed <iframe>. A page is fetched over a host-supplied transport
// (clearnet, or skywire's dmsg mesh), rendered into a sandboxed srcdoc with its
// stylesheets and images inlined, and its navigation relayed back to the chrome.
// The browser is Go; only the rendering (the iframe) and the network (the
// transport) are delegated.
//
// A host serves BrowserWasm() at a URL and includes LoaderJS(), then calls
// globalThis.Netscrape.open(element, opts) to mount it — see loader.js.
//
// The previous JavaScript engine (browse.js, the SkywireBrowse panel) lives on
// the `js` branch.
package netscrape

import (
	"bytes"
	"compress/gzip"
	_ "embed"
	"io"
	"sync"
)

//go:embed browser.wasm.gz
var browserWasmGz []byte

//go:embed loader.js
var loaderJS []byte

// BrowserWasmGz is the compressed browser module, for a consumer that inlines
// it into a page rather than serving it as a separate fetch.
func BrowserWasmGz() []byte { return browserWasmGz }

// LoaderJS returns loader.js, which defines globalThis.Netscrape.open.
func LoaderJS() []byte { return loaderJS }

var (
	browserOnce sync.Once
	browserWasm []byte
)

// BrowserWasm is the browser module as served (e.g. at /netscrape.wasm),
// inflated once on first use and kept.
func BrowserWasm() []byte {
	browserOnce.Do(func() {
		zr, err := gzip.NewReader(bytes.NewReader(browserWasmGz))
		if err != nil {
			return
		}
		defer zr.Close() //nolint:errcheck
		if b, err := io.ReadAll(zr); err == nil {
			browserWasm = b
		}
	})
	return browserWasm
}
