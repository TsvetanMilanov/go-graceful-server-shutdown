package gss

import (
	"fmt"
	"net/http"
	"os"
	"syscall"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStartServer(t *testing.T) {
	for _, s := range []os.Signal{syscall.SIGINT, syscall.SIGTERM} {
		ch := make(chan error)

		go func() {
			err := StartServer(createRouter())
			ch <- err
		}()

		shortWait()

		reqCh := make(chan error)
		go func() {
			reqErr := sendRequest(fmt.Sprintf(":8080"), 2)
			reqCh <- reqErr
		}()

		shortWait()
		sendStopSignal(s)

		assert.NoError(t, <-ch)
		assert.NoError(t, <-reqCh)
	}
}

func TestStartServerWithSettings(t *testing.T) {
	ch, addr := startTestServer(0)

	reqErr := sendRequest(addr, 0)
	assert.NoError(t, reqErr)

	sendTestStopSignal()

	err := <-ch

	assert.NoError(t, err)
}

func TestStartCustomServerWithSettings(t *testing.T) {
	for _, s := range []os.Signal{syscall.SIGINT, syscall.SIGTERM} {
		ch := make(chan error)
		addr := fmt.Sprintf(":%d", getFreePort())

		go func() {
			err := StartCustomServerWithSettings(&http.Server{Addr: addr, Handler: http.DefaultServeMux}, &Settings{})
			ch <- err
		}()

		shortWait()

		reqCh := make(chan error)
		go func() {
			reqErr := sendRequest(addr, 2)
			reqCh <- reqErr
		}()

		shortWait()
		sendStopSignal(s)

		assert.NoError(t, <-ch)
		assert.NoError(t, <-reqCh)
	}
}

func TestGracefulShutdown(t *testing.T) {
	// Set high timeout to ensure that the method will exit after all requests are drained
	// and there is still timeout.
	ch, addr := startTestServer(testTimeoutSeconds * 10)

	reqCh := make(chan error)
	go func() {
		reqErr := sendRequest(addr, 2)
		reqCh <- reqErr
	}()

	shortWait()
	sendTestStopSignal()
	shortWait()

	reqErr := sendRequest(addr, 3)
	assert.Error(t, reqErr)

	assert.NoError(t, <-reqCh)
	assert.NoError(t, <-ch)
}

func TestGracefulShutdownTimeout(t *testing.T) {
	ch, addr := startTestServer(1)

	// Run the request longer than the test timeout to ensure the shutdown
	// timeout will return an error.
	go sendRequest(addr, testTimeoutSeconds*10)

	shortWait()
	sendTestStopSignal()

	assert.Error(t, <-ch)
}

func TestServerListenError(t *testing.T) {
	ch := startTestServerWithAddr(0, "invalid")

	assert.Error(t, <-ch)
}

func TestInfiniteTimeout(t *testing.T) {
	ch, addr := startTestServer(-1)

	reqCh := make(chan error)
	go func() {
		reqErr := sendRequest(addr, 1)
		reqCh <- reqErr
	}()

	shortWait()
	sendTestStopSignal()

	assert.NoError(t, <-ch)
	assert.NoError(t, <-reqCh)
}

func TestStopChannel(t *testing.T) {
	stopChan := make(chan bool)
	addr := fmt.Sprintf(":%d", getFreePort())
	settings := &Settings{
		Addr:            addr,
		ShutdownChannel: stopChan,
	}

	ch := startTestServerWithSettings(settings)

	reqCh := make(chan error)
	go func() {
		reqErr := sendRequest(addr, 1)
		reqCh <- reqErr
	}()

	shortWait()
	stopChan <- true

	assert.NoError(t, <-ch)
	assert.NoError(t, <-reqCh)
}
