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

func InitConfingEnv() error {
	return godotenv.Load("config.env")
}

func main() {

	if err := InitConfingEnv(); err != nil {
		log.Fatal("Error loading .env file (config.env)")
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

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	httpServer := http.Server{
		Addr:    ":" + port,
		Handler: srv.router,
	}

	ctxStop, stopServer := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer stopServer()

	go func() {
		log.Println("Server started on port: ", port, "")
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
