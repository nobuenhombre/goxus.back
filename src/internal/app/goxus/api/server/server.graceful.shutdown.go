package server

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// waitOSInterruptSignal blocks until SIGINT or SIGTERM is received.
func (srv *HTTPServer) waitOSInterruptSignal() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")
}

// terminateServer gracefully shuts down the HTTP server with a 5-second timeout.
func (srv *HTTPServer) terminateServer() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := srv.Server.Shutdown(ctx)
	if err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	select {
	case <-ctx.Done():
		log.Println("timeout of 5 seconds.")
	}

	log.Println("Server exiting")
}

// gracefulShutDown waits for an OS signal and then terminates the server.
func (srv *HTTPServer) gracefulShutDown() {
	srv.waitOSInterruptSignal()
	srv.terminateServer()
}
