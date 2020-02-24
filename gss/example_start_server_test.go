package gss_test

import (
	"fmt"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/TsvetanMilanov/go-graceful-server-shutdown/gss"
)

func ExampleStartServer() {
	// StartServer will block the program execution until the server is closed.
	err := gss.StartServer(http.DefaultServeMux)

	fmt.Println(err)
}

func ExampleStartServerWithSettings() {
	var shutdownTimeoutSeconds int64 = 30
	settings := &gss.Settings{
		// Set the server address.
		Addr: ":5678",
		// Set the signals which will trigger the graceful server shutdown.
		SignalsToWatch: []os.Signal{syscall.SIGTERM, syscall.SIGKILL, syscall.SIGINT},
		// Set shutdown timeout.
		ShutdownTimeoutSeconds: &shutdownTimeoutSeconds,
	}

	// StartServerWithSettings will block the program execution until the server is closed.
	err := gss.StartServerWithSettings(http.DefaultServeMux, settings)

	fmt.Println(err)
}

func ExampleStartCustomServerWithSettings() {
	settings := &gss.Settings{
		// Set the signals which will trigger the graceful server shutdown.
		SignalsToWatch: []os.Signal{syscall.SIGTERM, syscall.SIGKILL, syscall.SIGINT},
	}

	srv := &http.Server{Handler: http.DefaultServeMux, Addr: ":5678"}

	// StartServerWithSettings will block the program execution until the server is closed.
	err := gss.StartCustomServerWithSettings(srv, settings)

	fmt.Println(err)
}

func ExampleStartServerWithSettings_shutdownchannel() {
	shutdownChannel := make(chan bool)
	settings := &gss.Settings{
		ShutdownChannel: shutdownChannel,
	}

	serverChannel := make(chan error)
	go func() {
		err := gss.StartServerWithSettings(http.DefaultServeMux, settings)

		serverChannel <- err
	}()

	time.Sleep(500 * time.Millisecond)

	shutdownChannel <- true

	err := <-serverChannel

	fmt.Println(err)
	// Output: <nil>
}
