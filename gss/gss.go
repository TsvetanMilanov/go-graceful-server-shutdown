package gss

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"
)

// StartServer creates server from the provided handler, starts it
// and handles its shutdown based on the default settings.
// This function is blocking.
func StartServer(router http.Handler) error {
	settings := new(Settings)
	err := StartServerWithSettings(router, settings)

	return err
}

// StartServerWithSettings creates server from the provided handler, starts it
// and handles its shutdown based on the provided settings.
// This function is blocking.
func StartServerWithSettings(router http.Handler, settings *Settings) error {
	srv := createServer(router, settings)

	err := StartCustomServerWithSettings(srv, settings)

	return err
}

// StartCustomServerWithSettings starts the provided server and handles its shutdown
// based on the provided settings.
// This function is blocking.
func StartCustomServerWithSettings(srv *http.Server, settings *Settings) error {
	startServerChan := startServer(srv)

	quit := make(chan os.Signal)
	signal.Notify(quit, settings.getSignalsToWatch()...)

	select {
	case <-quit:
		signal.Stop(quit)
		break
	case err := <-startServerChan:
		if err != nil {
			signal.Stop(quit)
			close(quit)
			return err
		}
	}

	ctx, cancel := getContext(settings)
	if cancel != nil {
		defer cancel()
	}

	err := srv.Shutdown(ctx)

	return err
}

func getContext(settings *Settings) (ctx context.Context, cancel context.CancelFunc) {
	if settings.getShutdownTimeoutSeconds() == -1 {
		ctx = context.Background()
	} else {
		ctx, cancel = context.WithTimeout(context.Background(), time.Duration(settings.getShutdownTimeoutSeconds())*time.Second)
	}

	return ctx, cancel
}

func startServer(srv *http.Server) chan error {
	startServerChan := make(chan error)
	go func() {
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			startServerChan <- err
		} else {
			startServerChan <- nil
		}
	}()

	return startServerChan
}

func createServer(router http.Handler, settings *Settings) *http.Server {
	srv := &http.Server{
		Addr:    settings.getAddr(),
		Handler: router,
	}

	return srv
}
