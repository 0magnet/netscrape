// netscrape loader — mount the Go/wasm browser into a page element and wire its
// network transport. A host calls:
//
//   Netscrape.open(element, {
//     wasmURL:      "/netscrape.wasm",     // where the module is served
//     fetchDmsg:    fn(pkHost, m, path, body) -> {status, body, headers},
//     fetchClearnet:fn(url, m, body)        -> {status, body, headers},
//     fetch:        fn(url) -> Response,     // override the whole transport
//   })
//
// The browser (cmd/browser) reads globalThis.__netscrapeMount for its element
// and calls globalThis.__netscrapeFetch(url) for every page and subresource, so
// the same Go browser serves a plain clearnet page, a skywire mesh visor, or
// anything else — only this transport differs.
(function () {
	function meshHost(h) {
		return /\.(dmsg|skysocks|skynet)$/i.test(h) || /^[0-9a-f]{66}$/i.test(h);
	}
	function respond(r) {
		var h = new Headers();
		if (r && r.headers) {
			try { for (var k in r.headers) h.set(k, r.headers[k]); } catch (e) { /* ignore */ }
		}
		return new Response((r && r.body) || new Uint8Array(0), { status: (r && r.status) || 200, headers: h });
	}
	globalThis.Netscrape = {
		open: function (el, opts) {
			opts = opts || {};
			if (typeof globalThis.Go !== "function") {
				console.error("netscrape: the Go runtime (wasm_exec.js) is not on the page");
				return null;
			}
			globalThis.__netscrapeFetch = opts.fetch || function (url) {
				var u;
				try { u = new URL(url); } catch (e) { return fetch(url); }
				var path = (u.pathname || "/") + (u.search || "");
				if (opts.fetchDmsg && meshHost(u.hostname)) {
					return Promise.resolve(opts.fetchDmsg(u.hostname, "GET", path, null)).then(respond);
				}
				if (opts.fetchClearnet) {
					return Promise.resolve(opts.fetchClearnet(url, "GET", null)).then(respond);
				}
				return fetch(url);
			};
			globalThis.__netscrapeMount = el;
			fetch(opts.wasmURL || "/netscrape.wasm")
				.then(function (r) { return r.arrayBuffer(); })
				.then(function (buf) {
					var go = new globalThis.Go();
					return WebAssembly.instantiate(buf, go.importObject).then(function (res) { go.run(res.instance); });
				})
				.catch(function (e) { el.textContent = "netscrape failed to load: " + e; });
			return el;
		},
	};
})();
