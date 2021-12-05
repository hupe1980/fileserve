package main

import (
	"embed"
	"io/fs"
	"log"

	"github.com/hupe1980/fileserve"
)

//go:embed nextjs/out
//go:embed nextjs/out/_next
//go:embed nextjs/out/_next/static/chunks/pages/*.js
//go:embed nextjs/out/_next/static/*/*.js
var nextFS embed.FS

func main() {
	// Root at the `out` folder generated by the Next.js app.
	distFS, err := fs.Sub(nextFS, "nextjs/out")
	if err != nil {
		log.Fatal(err)
	}

	root, err := fileserve.NewFSRoot(distFS, func(o *fileserve.RootOptions) {
		o.HideDirListing = true
		o.HideDotFiles = true
	})
	if err != nil {
		log.Fatal(err)
	}

	// The static Next.js app will be served under `/`.
	fs, err := fileserve.New(root)
	if err != nil {
		log.Fatal(err)
	}

	// Start HTTP server at :8000.
	log.Println("Starting HTTP server at http://localhost:8000 ...")
	log.Fatal(fs.ListenAndServe())
}
