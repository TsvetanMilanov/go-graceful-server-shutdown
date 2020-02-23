package gss_test

import (
	"fmt"
	"net/http"

	"github.com/TsvetanMilanov/go-graceful-server-shutdown/gss"
)

func Example() {
	// Create server based on the provided handler.
	// The server will listen on port 8080 and there will be
	// no timeout for the connection draining.
	// StartServer will block the program execution until the server is closed.
	err := gss.StartServer(http.DefaultServeMux)

	fmt.Println(err)
}
