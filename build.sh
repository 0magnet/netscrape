#!/bin/sh
# Build the browser module that netscrape.go embeds. Run after changing
# cmd/browser; the compressed result (browser.wasm.gz) is committed so consumers
# get BrowserWasm() without a wasm build.
set -eu
cd "$(dirname "$0")"
GOOS=js GOARCH=wasm go build -o /tmp/netscrape-browser.wasm ./cmd/browser
gzip -9 -c /tmp/netscrape-browser.wasm > browser.wasm.gz
echo "netscrape: browser.wasm.gz updated ($(du -h browser.wasm.gz | cut -f1))"
