package fileserve

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
)

const DefaultPort = 8000

type Options struct {
	Port int
	Bind string
	TLS  bool
}

type Fileserve struct {
	handler http.Handler
	options Options
}

func New(root Root, optFns ...func(o *Options)) (*Fileserve, error) {
	options := Options{
		Port: DefaultPort,
		Bind: "0.0.0.0",
		TLS:  false,
	}

	for _, fn := range optFns {
		fn(&options)
	}

	return &Fileserve{
		options: options,
		handler: http.FileServer(root),
	}, nil
}

func (s *Fileserve) Use(middlewares ...func(http.Handler) http.Handler) {
	for _, middleware := range middlewares {
		s.handler = middleware(s.handler)
	}
}

func (s *Fileserve) ListenAndServe() error {
	return s.listenAndServe()
}

func (s *Fileserve) listenAndServe() error {
	addr := fmt.Sprintf("%s:%d", s.options.Bind, s.options.Port)

	if s.options.TLS {
		cert, err := NewSelfSignedCert(sans())
		if err != nil {
			log.Fatalf("cannot generate certificate: %v\n", err)
		}

		tlsConfig := &tls.Config{
			MinVersion:               tls.VersionTLS12,
			PreferServerCipherSuites: true,
			Certificates:             []tls.Certificate{cert},
		}
		server := &http.Server{
			Addr:      addr,
			TLSConfig: tlsConfig,
			Handler:   s.handler,
		}

		return server.ListenAndServeTLS("", "")
	}

	return http.ListenAndServe(addr, s.handler)
}
