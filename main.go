package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	err := godotenv.Load("database.env")
	if err != nil {
		log.Fatal("Error loading .env file (database.env)")
	}
	err = godotenv.Load("token.env")
	if err != nil {
		log.Fatal("Error loading .env file (token.env)")
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	duration := 5 * time.Second
	pool, err := EstablishDb(ctx, duration)

	if err != nil {
		log.Fatal("Failed to connect: ", err)
	}
	defer func() {
		pool.Close()
		log.Println("Connection with DB closed")
	}()
	fmt.Println("Connected")
	select {
	case <-ctx.Done():
		log.Println("Shutting down...")
		return
	default:
	}

	// TODO : SERVER STRUCT W/ CHI AND SERVICES

	userService := NewUserService(NewUserPgRepository(pool))
	taskService := NewTaskService(NewTaskPgRepository(pool))

	srv := NewServer(userService, taskService)

	httpServer := http.Server{
		Addr:    ":8080",
		Handler: srv.router,
	}

	ctxStop, stopServer := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer stopServer()

	go func() {
		log.Println("Server started")
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	<-ctxStop.Done()

	log.Println("Shutting down...")
	ctxShutDown, shutDown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutDown()
	if err := httpServer.Shutdown(ctxShutDown); err != nil {
		log.Fatal("Server shutdown: ", err)
	} else {
		log.Println("Server shut down gracefully")
	}

}
