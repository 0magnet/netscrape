//go:build js && wasm

// Command browser is netscrape as a standalone Go/wasm module: a thin wrapper
// that resolves the mount element and hands off to netscrape.Open, then blocks
// so its runtime stays alive. This is the binary a host SERVES (dist.BrowserWasm)
// and a shell EXECS (`run browser.wasm`) — one wasm, one Go runtime, on its own.
//
// A host that instead imports netscrape into its OWN wasm calls netscrape.Open
// directly and shares that binary's runtime — no separate module, no duplicated
// runtime. See netscrape.Open.
package main

import (
	"os"
	"syscall/js"

	"github.com/0magnet/netscrape"
)

func main() {
	doc := js.Global().Get("document")
	// The element to draw into: $NETSCRAPE_MOUNT when spawned by shipyard's proc
	// layer, or globalThis.__netscrapeMount (an element id, or the element) when a
	// host instantiates the module directly.
	root := doc.Call("getElementById", os.Getenv("NETSCRAPE_MOUNT"))
	if !root.Truthy() {
		if m := js.Global().Get("__netscrapeMount"); m.Truthy() {
			if m.Type() == js.TypeString {
				root = doc.Call("getElementById", m.String())
			} else {
				root = m
			}
		}
	}
	if !root.Truthy() {
		return
	}
	netscrape.Open(root)
	select {}
}
