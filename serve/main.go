// Command serve is a static file server for the demo page, so the wasm can
// be loaded over http rather than file:// (which blocks WebAssembly.instantiate).
package main

import (
	"log"
	"net/http"

	"github.com/0magnet/calvin/clihelp"
	"github.com/spf13/cobra"
)

func main() {
	var addr, dir string
	cmd := &cobra.Command{
		Use:          "serve",
		Short:        "serve the demo page over HTTP",
		SilenceUsage: true,
		Run: func(*cobra.Command, []string) {
			log.Printf("serving %s on http://localhost%s", dir, addr)
			log.Fatal(http.ListenAndServe(addr, http.FileServer(http.Dir(dir)))) //nolint:gosec
		},
	}
	cmd.Flags().StringVar(&addr, "addr", ":8812", "listen address")
	cmd.Flags().StringVar(&dir, "dir", "docs", "directory to serve")
	clihelp.Init(cmd, "serve")
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
