package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	log.Println("Starting server...")

	// Load dev env from .env file
	asOrg := flag.String("as", "", "")
	flag.Parse()

	var err error
	if *asOrg != "" {
		err = godotenv.Load(fmt.Sprintf(".env.%s", *asOrg))
	} else {
		err = godotenv.Load(".env")
	}
	if err != nil {
		log.Fatalf("Error loading .env file: %v\n", err)
	}

	// initialize data sources
	ds, err := initDS()
	if err != nil {
		log.Fatalf("Unable to initialize data sources: %v\n", err)
	}

	router, err := inject(ds)
	if err != nil {
		log.Fatalf("Failure to inject data sources: %v\n", err)
	}

	srv := &http.Server{
		Addr:    ":" + os.Getenv("PORT"),
		Handler: router,
	}

	// Graceful server shutdown - https://github.com/gin-gonic/examples/blob/master/graceful-shutdown/graceful-shutdown/server.go
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to initialize server: %v\n", err)
		}
	}()

	log.Printf("Listening on port %v\n", srv.Addr)

	// Wait for kill signal of channel
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// This blocks until a signal is passed into the quit channel
	<-quit

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// shutdown data sources
	ds.Gateway.Close()
	if err := ds.Couch.Close(context.TODO()); err != nil {
		log.Fatalf("A problem occurred gracefully shutting down data sources: %v\n", err)
	}

	// Shutdown server
	log.Println("Shutting down server...")
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v\n", err)
	}
}
