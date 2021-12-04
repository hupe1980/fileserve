package main

import (
	"errors"
	"fmt"
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
		bind  string
		port  int
		https bool
		cors  bool
		nc    bool
		auth  string
	)

	rootCmd := &cobra.Command{
		Use:     "fileserve [root]",
		Version: version,
		Short:   "fileserve is a tiny go based file server",
		Args:    cobra.MinimumNArgs(1),
		Example: "serve the current working dir: fileserve .",
		RunE: func(cmd *cobra.Command, args []string) error {
			root := args[0]

			fs, err := fileserve.New(root, func(o *fileserve.Options) {
				o.Port = port
				o.Bind = bind
				o.TLS = https
			})
			if err != nil {
				return err
			}

			fs.Use(fileserve.AccessLog())

			if cors {
				fs.Use(fileserve.CORS("*"))
			}

			if nc {
				fs.Use((fileserve.NoCache()))
			}

			if auth != "" {
				creds, err := parseAuth(auth)
				if err != nil {
					return err
				}

				fs.Use(fileserve.BasicAuth("restricted", creds))
			}

			schema := "http"
			if https {
				schema = "https"
			}
			fmt.Printf("Serving \"%s\" on: %v://%v:%d\n", root, schema, bind, port)

			return fs.ListenAndServe()
		},
	}

	rootCmd.Flags().StringVarP(&bind, "bind", "b", "0.0.0.0", "bind to a specific interface")
	rootCmd.Flags().IntVarP(&port, "port", "p", fileserve.DefaultPort, "port to serve on")
	rootCmd.Flags().BoolVarP(&https, "https", "s", false, "serve with a temp self-signed certificate via HTTPS")
	rootCmd.Flags().BoolVarP(&cors, "cors", "", false, "allow cross origin requests to be served")
	rootCmd.Flags().BoolVarP(&nc, "no-cache", "", false, "disable caching for the file server")
	rootCmd.Flags().StringVarP(&auth, "auth", "a", "", "turn on basic auth and set username and password (separate by colon)")

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func parseAuth(auth string) (map[string]string, error) {
	creds := make(map[string]string)

	if !strings.Contains(auth, ":") {
		return nil, errors.New("auth must be specified in the format username:password")
	}

	parts := strings.Split(auth, ":")
	if len(parts) > credsParts {
		return nil, errors.New("only one colon is allowed")
	}

	creds[parts[0]] = parts[1]

	return creds, nil
}
