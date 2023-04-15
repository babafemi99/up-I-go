package main

import (
	"context"
	"fmt"
	"github.com/babafemi99/up-I-go/cmd"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	r := mux.NewRouter()
	port := "8000"

	r.HandleFunc("/upload", cmd.HandleUpload).Methods("POST")
	r.HandleFunc("/download", cmd.HandleDownload).Methods("GET")

	server := http.Server{
		Addr:        fmt.Sprintf(":%s", port),
		Handler:     r,
		IdleTimeout: 120 * time.Second,
	}
	go func() {
		log.Printf("------------ SERVER STARTING ON PORT: %s ------------\n", port)
		err := server.ListenAndServe()
		if err != nil {
			log.Printf("ERROR STARTING SERVER: %v", err)
			os.Exit(1)
		}

	}()
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, os.Kill)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigChan
	log.Printf("Closing now, We've gotten signal: %v", sig)

	ctx := context.Background()
	server.Shutdown(ctx)

}
