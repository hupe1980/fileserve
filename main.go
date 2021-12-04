package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

const (
	defaultPort          = 8080
	defaultCompressLevel = 0 //no compression
)

func main() {
	flag.Usage = func() {
		fmt.Print("Usage: fileserve [options] <dir>\n\nOptions:\n")
		flag.PrintDefaults()
	}

	creds := flag.String("a", "", `Require basic authentication with given credentials (e.g. -a "bob:secret")`)
	bind := flag.String("b", "0.0.0.0", `Bind to a specific interface`)
	port := flag.Int("p", defaultPort, "Port to serve on. 8080 by default for HTTP, 8443 for HTTPS")
	https := flag.Bool("s", false, "Creates a temporary self-signed certificate to serve via HTTPS")
	compress := flag.Int("c", defaultCompressLevel, "Compression level to use. Levels range from 1 (BestSpeed) to 9 (BestCompression)")

	flag.Parse()

	dir := flag.Arg(0)

	if dir == "" {
		fmt.Print("Directory is mssing. Use . for the current working dir\n\n")
		flag.Usage()
		os.Exit(1)
	}

	root, err := Root(dir)
	if err != nil {
		log.Fatal(err)
	}

	handler := http.FileServer(root)

	handler = RequestLog(handler)

	if *compress != 0 {
		handler = Compress(handler, *compress)
	}

	if *creds != "" {
		user, pass, err := ParseCreds(*creds)
		if err != nil {
			log.Fatal(err)
		}

		handler = BasicAuth(handler, "restricted", map[string]string{user: pass})
	}

	http.Handle("/", handler)

	if *https {
		if !isFlagPassed("p") {
			*port = 8443
		}

		cert, err := SelfSignedCert(sans())
		if err != nil {
			log.Fatalf("cannot generate certificate: %v\n", err)
		}

		tlsConfig := &tls.Config{
			MinVersion:               tls.VersionTLS12,
			PreferServerCipherSuites: true,
			Certificates:             []tls.Certificate{cert},
		}

		server := &http.Server{
			Addr:      fmt.Sprintf("%s:%d", *bind, *port),
			TLSConfig: tlsConfig,
		}

		log.Fatal(server.ListenAndServeTLS("", ""))
	} else {
		log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", *bind, *port), nil))
	}
}

func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
			return
		}
	})
	return found
}
