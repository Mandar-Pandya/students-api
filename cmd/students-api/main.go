package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Mandar-Pandya/students-api/internal/config"
)

func main() {
	// setup config
	cfg := config.MustLoad()

	// setup router
	router := http.NewServeMux()

	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to students api"))
	})

	// setup server

	server := http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}

	slog.Info("Server started", slog.String("address", cfg.Addr))

	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal("Unable to start the server")
		}
	}()

	<-done

	slog.Info("shutting down the server")

	ctx,cancel := context.WithTimeout(context.Background(), 5 * time.Second)

	defer cancel()

	if err := server.Shutdown(ctx);
	err != nil{
		slog.Error("Failed to shut down server",slog.String("error",err.Error()))
	}

	slog.Info("server shutdown successfully")

}
