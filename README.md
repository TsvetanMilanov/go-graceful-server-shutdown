# go-graceful-server-shutdown
Go package which gracefully stops servers

[![GoDoc](https://godoc.org/github.com/TsvetanMilanov/go-graceful-server-shutdown/gss?status.svg)](https://godoc.org/github.com/TsvetanMilanov/go-graceful-server-shutdown/gss)
![Go](https://github.com/TsvetanMilanov/go-graceful-server-shutdown/workflows/Go/badge.svg?branch=master)
![Create Release](https://github.com/TsvetanMilanov/go-graceful-server-shutdown/workflows/Create%20Release/badge.svg)

## Quick Start
```Go
// Create server based on the provided handler.
// The server will listen on port 8080 and there will be
// no timeout for the connection draining.
// StartServer will block the program execution until the server is closed.
err := gss.StartServer(http.DefaultServeMux)
```
