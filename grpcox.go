package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/gusaul/grpcox/core"
	"github.com/gusaul/grpcox/handler"
	"go.uber.org/zap"
)

func main() {
	// logging conf
	f, err := os.OpenFile("log/grpcox.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// start app
	addr := "0.0.0.0:6969"
	if value, ok := os.LookupEnv("BIND_ADDR"); ok {
		addr = value
	}
	muxRouter := mux.NewRouter()
	handler.Init(muxRouter)
	var wait time.Duration = time.Second * 15

	srv := &http.Server{
		Addr:         addr,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      muxRouter,
	}

	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	sugar := logger.Sugar()
	sugar.Infow("failed to fetch URL",
		// Structured context as loosely typed key-value pairs.
		"url", "dummy_url",
		"attempt", 3,
		"backoff", time.Second,
	)
	sugar.Infof("Failed to fetch URL: %s", "dummy_url")
	sugar.Desugar().Fire("Fire is burning")
	logger.Fire("there is a fire")

	fmt.Println("Service started on", addr)
	go func() {
		log.Fatal(srv.ListenAndServe())
	}()

	c := make(chan os.Signal, 1)

	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	err = removeProtos()
	if err != nil {
		log.Printf("error while removing protos: %s", err.Error())
	}

	srv.Shutdown(ctx)
	log.Println("shutting down")
	os.Exit(0)
}

// removeProtos will remove all uploaded proto file
// this function will be called as the server shutdown gracefully
func removeProtos() error {
	log.Println("removing proto dir from /tmp")
	err := os.RemoveAll(core.BasePath)
	return err
}
