# fileserve
![Build Status](https://github.com/hupe1980/fileserve/workflows/build/badge.svg) 
[![Go Reference](https://pkg.go.dev/badge/github.com/hupe1980/fileserve.svg)](https://pkg.go.dev/github.com/hupe1980/fileserve)
> Golang-based simple file server to serve static files of the current working directory
- File sharing in LAN or home network
- Web application testing
- Personal web site hosting or demonstrating

## Features
- Directory listing
- Cors or NoCache headers
- HTTPS with self-signed certificate
- BasicAuth

## How to use
```bash
fileserve is a tiny go based file server

Usage:
  fileserve [root] [flags]

Examples:
serve the current working dir: fileserve .

Flags:
  -a, --auth string   turn on basic auth and set username and password (separate by colon)
  -b, --bind string   bind to a specific interface (default "0.0.0.0")
      --cors          allow cross origin requests to be served
  -h, --help          help for fileserve
  -s, --https         serve with a temp self-signed certificate via HTTPS
      --no-cache      disable caching for the file server
  -p, --port int      port to serve on (default 8000)
  -v, --version       version for fileserve
```

## License
[MIT](LICENCE)
