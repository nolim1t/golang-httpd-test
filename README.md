# Test GoLang HTTPD
![Golang Gopher](https://gitlab.com/nolim1t/golang-httpd-test/-/raw/master/golang.png)

## About

The purpose of this project is purely for learning how go build interfaces and a HTTPD API server in GOLANG.

I do like go, because the packaging system works with git seamlessly (because decentralization!)

Eventually, I'd like to use this for querying external based libraries (such as lnd / bitcoind), and utilizing a config file. 

## Directory structure

- `common`:  contains some useful utilities.
- `pineclient` : contains the pineclient package for reading stuff from the PINEphone.
- `go.mod` : contains a list of all the go modules and defines the base package name.
- `main.go` : Defines the entry point which binds all the modules together.

## TODO

- [x] Configuration File support 
- [ ] Static Directory serving
- [ ] Toggle between dev and production mode
- [ ] `Dockerfile` support
- [ ] Docker buildx support to push to both github and gitlab
- [ ] Tidy up code and make this a template fully configurable server in a MVC format
