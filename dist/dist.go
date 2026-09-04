// Package dist carries netscrape's pre-built standalone browser module and its
// loader — for a host that SERVES the browser as a separate wasm fetch (or
// inlines it into a page) rather than importing the browser into its own wasm.
//
// A host serves BrowserWasm() at a URL and includes LoaderJS(), then calls
// globalThis.Netscrape.open(element, opts) to mount it — see loader.js.
//
// A host that instead compiles the browser INTO its own wasm imports the parent
// package and calls netscrape.Open — no separate module, no duplicated Go
// runtime. Use this package only when you want the standalone binary.
package dist

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
