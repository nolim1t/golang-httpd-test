# Test GoLang HTTPD
![Golang Gopher](https://gitlab.com/nolim1t/golang-httpd-test/-/raw/master/golang.png)

[![pipeline status](https://gitlab.com/nolim1t/golang-httpd-test/badges/master/pipeline.svg)](https://gitlab.com/nolim1t/golang-httpd-test/-/commits/master) 
[![Build on tag](https://github.com/nolim1t/golang-httpd-test/workflows/Docker%20build%20on%20tag/badge.svg)](https://github.com/nolim1t/golang-httpd-test/actions?query=workflow%3A%22Docker+build+on+tag%22)
[![Build on push](https://github.com/nolim1t/golang-httpd-test/workflows/Docker%20build%20on%20push/badge.svg)](https://github.com/nolim1t/golang-httpd-test/actions?query=workflow%3A%22Docker+build+on+push%22)
![Version](https://img.shields.io/github/v/release/nolim1t/golang-httpd-test?sort=semver) 
[![Docker Pulls Count](https://img.shields.io/docker/pulls/nolim1t/golang-httpd-test.svg?style=flat)](https://hub.docker.com/r/nolim1t/golang-httpd-test)
## About

The purpose of this project is purely for learning how go build interfaces and a HTTPD API server in GOLANG.

I do like go, because the packaging system works with git seamlessly (because decentralization!)

Eventually, I'd like to use this for querying external based libraries (such as lnd / bitcoind), and utilizing a config file. 

## Directory structure

- `common`:  contains some useful utilities.
- `pineclient` : contains the pineclient package for reading stuff from the PINEphone.
- `bitcoind` : contains bitcoind client package for reading bitcoind
- `btcprice` : contains btc price utilities
- `go.mod` : contains a list of all the go modules and defines the base package name.
- `main.go` : Defines the entry point which binds all the modules together.

## Building - Docker

### Build Arguments

* `VERSION` - defines the version to be timestamped inside the --version identifier in the binary
* `VER_GO` - defines the version for GO to be used (default: 1.15)
* `VER_ALPINE` - defines the version of alpine to be used. Must also support the version of go or it will fail  (default: 3.12)
* `TAGS` - defines any build tags to be used.

### Example

Here is an example of building this project in docker (on your current environment). You may substitute this with `buildx build` which uses very similar parameters.

```bash
docker build --build-arg VERSION=0.0.1 \
             --build-arg VER_GO=1.15 \
             --build-arg VER_ALPINE=3.12 \
             --build-arg TAGS="static_build" \
             -t nolim1t/httpd:0.0.1 \
             .
```

## TODO

- [x] Configuration File support 
- [x] Static Directory serving
- [ ] Toggle between dev and production mode
- [x] `Dockerfile` support
- [x] Docker buildx support to push to both github and gitlab
- [ ] Tidy up code and make this a template fully configurable server in a MVC format
