package gss

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	customSignal       = syscall.SIGHUP
	testTimeoutSeconds = 60
)

func getFreePort() int {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		panic(err)
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		panic(err)
	}

	defer l.Close()

	port := l.Addr().(*net.TCPAddr).Port

	return port
}

func createRouter() http.Handler {
	router := gin.New()
	router.GET("/:reqDuration", func(c *gin.Context) {
		requestDuration, _ := strconv.Atoi(c.Param("reqDuration"))

		time.Sleep(time.Duration(requestDuration) * time.Second)

		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	return router
}

func startTestServer(shutdownTimeout int64) (chan error, string) {
	addr := fmt.Sprintf(":%d", getFreePort())
	ch := startTestServerWithAddr(shutdownTimeout, addr)

	return ch, addr
}

func startTestServerWithAddr(shutdownTimeout int64, addr string) chan error {
	settings := &Settings{
		SignalsToWatch:         []os.Signal{customSignal},
		Addr:                   addr,
		ShutdownTimeoutSeconds: &shutdownTimeout,
	}

	ch := startTestServerWithSettings(settings)

	return ch
}

func startTestServerWithSettings(settings *Settings) chan error {
	ch := make(chan error)

	go func() {
		err := StartServerWithSettings(createRouter(), settings)
		ch <- err
	}()

	shortWait()

	return ch
}

func sendRequest(addr string, reqDuration int) error {
	_, err := http.Get(fmt.Sprintf("http://%s/%d", addr, reqDuration))

	return err
}

func sendTestStopSignal() {
	sendStopSignal(customSignal)
}

func sendStopSignal(s os.Signal) {
	p, err := os.FindProcess(os.Getpid())
	if err != nil {
		panic(err)
	}

	p.Signal(s)
}

func shortWait() {
	time.Sleep(500 * time.Millisecond)
}
