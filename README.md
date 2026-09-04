# netscrape

A web browser written in Go/wasm.

The chrome — a tab strip, an address bar, back/forward/reload, history — is DOM
built from Go with `syscall/js`. Each tab is a sandboxed `<iframe>`. A page is
fetched over a **host-supplied transport**, transcoded (rendered into a
sandboxed `srcdoc` with its stylesheets inlined as `<style>` and its images as
`data:` URIs, via a parent/sandbox relay), and its links and form submits are
relayed back to the chrome. The browser is Go; only the rendering (the iframe)
and the network (the transport) are delegated.

That transport is the seam. The same browser renders a plain clearnet page, or
a [skywire](https://github.com/skycoin/skywire) site addressed by public key
over the dmsg mesh — the host decides, by what it plugs into the fetch.

## Use

A host serves the module and includes the loader, then mounts it into an
element and supplies the network:

```go
import "github.com/0magnet/netscrape"

serveBytes("/netscrape.wasm", "application/wasm", netscrape.BrowserWasm())
serveJS(netscrape.LoaderJS()) // defines globalThis.Netscrape.open
```

```js
Netscrape.open(document.getElementById("browser"), {
  wasmURL:       "/netscrape.wasm",
  fetchDmsg:     (pk, method, path, body) => ...,   // → {status, body, headers}
  fetchClearnet: (url, method, body)      => ...,   // → same shape
  // or fetch: (url) => Response   to replace the whole transport
});
```

Mesh hosts (`*.dmsg`, `*.skysocks`, a 66-hex public key) route through
`fetchDmsg`; everything else through `fetchClearnet`; absent either, a plain
same-origin `fetch`. The browser reads `globalThis.__netscrapeMount` for its
element and `globalThis.__netscrapeFetch(url)` for every request.

`./build.sh` rebuilds `browser.wasm.gz` (embedded by `netscrape.go`) from
`cmd/browser`.

## The JavaScript engine

netscrape was a 3.6k-line JavaScript engine (`browse.js`) — a mesh browser with
its own transcoder, a `SkywireBrowse` panel, address-bar channel dispatch, and
window management, grown inside skywire. That version lives on the
[`js` branch](https://github.com/0magnet/netscrape/tree/js). `main` is the Go
rewrite; consumers move over as it reaches parity.
