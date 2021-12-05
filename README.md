# fileserve
![Build Status](https://github.com/hupe1980/fileserve/workflows/build/badge.svg) 
[![Go Reference](https://pkg.go.dev/badge/github.com/hupe1980/fileserve.svg)](https://pkg.go.dev/github.com/hupe1980/fileserve)
> Golang-based simple file server to serve static files of the current working directory
- File sharing in LAN or home network
- Web application testing
- Personal web site hosting or demonstrating

## Features
- Directory listing
- Filter dot files
- Cors or NoCache headers
- HTTPS with self-signed certificate
- BasicAuth
- Custom HTTP headers

## Installing
You can install the pre-compiled binary in several different ways

### homebrew tap:
```bash
brew tap hupe1980/fileserve
brew install fileserve
```

### snapcraft:
[![Get it from the Snap Store](https://snapcraft.io/static/images/badges/en/snap-store-black.svg)](https://snapcraft.io/fileserve)
```bash
sudo snap install fileserve
```

### scoop:
```bash
scoop bucket add fileserve https://github.com/hupe1980/fileserve-bucket.git
scoop install fileserve
```

### deb/rpm/apk:

Download the .deb, .rpm or .apk from the [releases page](https://github.com/hupe1980/fileserve/releases) and install them with the appropriate tools.

### manually:
Download the pre-compiled binaries from the [releases page](https://github.com/hupe1980/fileserve/releases) and copy to the desired location.

## How to use
```bash
fileserve is a tiny go based file server

Usage:
  fileserve [root] [flags]

Examples:
- serve the current working dir: fileserve .
- add basic auth: fileserve . -a user1:pass1 -a user2:pass2
- add custom http headers: fileserve . --header Test=ABC --header Foo=Bar
- disable serving of dot files: fileserve . --no-dot

Flags:
  -a, --auth stringArray        turn on basic auth and set username and password (separate by colon)
  -b, --bind string             bind to a specific interface (default "0.0.0.0")
      --cors                    allow cross origin requests to be served
      --header stringToString   add custom http headers (default [])
  -h, --help                    help for fileserve
  -s, --https                   serve with a temp self-signed certificate via HTTPS
      --no-cache                disable caching for the file server
      --no-dir                  turn off directory listing
      --no-dot                  disable serving of dot files
  -p, --port int                port to serve on (default 8000)
  -v, --version                 version for fileserve
```

Or use the module to create your own protable:
* [nextjs-portable example](_examples/nextjs-portable) 

## License
[MIT](LICENCE)
