package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/hupe1980/fileserve"
	"github.com/spf13/cobra"
)

const (
	version    = "dev"
	credsParts = 2
)

func main() {
	var (
		bind    string
		port    int
		https   bool
		cors    bool
		nc      bool
		noDir   bool
		noDot   bool
		auths   []string
		headers map[string]string
	)

	rootCmd := &cobra.Command{
		Use:     "fileserve [root]",
		Version: version,
		Short:   "fileserve is a tiny go based file server",
		Args:    cobra.MinimumNArgs(1),
		Example: `- serve the current working dir: fileserve .
- add basic auth: fileserve . -a user1:pass1 -a user2:pass2
- add custom http headers: fileserve . --header Test=ABC --header Foo=Bar
- disable serving of dot files: fileserve . --no-dot`,
		RunE: func(cmd *cobra.Command, args []string) error {
			dir := args[0]

			root, err := fileserve.NewDirRoot(dir, func(o *fileserve.DirRootOptions) {
				o.ShowDirListing = !noDir
				o.ShowDotFiles = !noDot
			})
			if err != nil {
				return err
			}

			fs, err := fileserve.New(root, func(o *fileserve.Options) {
				o.Port = port
				o.Bind = bind
				o.TLS = https
			})
			if err != nil {
				return err
			}

			fs.Use(fileserve.AccessLog())

			if noDot {
				fs.Use(fileserve.NoDot("404 page not found", http.StatusNotFound))
			}

			if cors {
				fs.Use(fileserve.CORS("*"))
			}

			if len(headers) != 0 {
				for k, v := range headers {
					fs.Use(fileserve.HTTPHeader(k, v))
				}
			}

			if nc {
				fs.Use((fileserve.NoCache()))
			}

			if len(auths) != 0 {
				creds, err := parseAuths(auths)
				if err != nil {
					return err
				}

				fs.Use(fileserve.BasicAuth("restricted", creds))
			}

			schema := "http"
			if https {
				schema = "https"
			}
			log.Printf("Serving \"%s\" on: %v://%v:%d\n", dir, schema, bind, port)

			return fs.ListenAndServe()
		},
	}

	rootCmd.Flags().StringVarP(&bind, "bind", "b", "0.0.0.0", "bind to a specific interface")
	rootCmd.Flags().IntVarP(&port, "port", "p", fileserve.DefaultPort, "port to serve on")
	rootCmd.Flags().BoolVarP(&https, "https", "s", false, "serve with a temp self-signed certificate via HTTPS")
	rootCmd.Flags().BoolVarP(&cors, "cors", "", false, "allow cross origin requests to be served")
	rootCmd.Flags().BoolVarP(&nc, "no-cache", "", false, "disable caching for the file server")
	rootCmd.Flags().BoolVarP(&noDir, "no-dir", "", false, "turn off directory listing")
	rootCmd.Flags().BoolVarP(&noDot, "no-dot", "", false, "disable serving of dot files")
	rootCmd.Flags().StringArrayVarP(&auths, "auth", "a", []string{}, "turn on basic auth and set username and password (separate by colon)")
	rootCmd.Flags().StringToStringVarP(&headers, "header", "", map[string]string{}, "add custom http headers")

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func parseAuths(auths []string) (map[string]string, error) {
	creds := make(map[string]string)

	for _, auth := range auths {
		if !strings.Contains(auth, ":") {
			return nil, errors.New("auth must be specified in the format username:password")
		}

		parts := strings.Split(auth, ":")
		if len(parts) > credsParts {
			return nil, errors.New("only one colon is allowed")
		}

		creds[parts[0]] = parts[1]
	}

	return creds, nil
}
