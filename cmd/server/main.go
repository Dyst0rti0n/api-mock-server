package main

import (
	"advanced-mock-server/config"
	"advanced-mock-server/internal/handlers"
	"advanced-mock-server/internal/middleware"
	"advanced-mock-server/mockdb"
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"golang.org/x/time/rate"

	_ "advanced-mock-server/docs"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "modernc.org/sqlite"
)

// @title Advanced Mock Server API
// @version 1.0
// @description This is an advanced mock server built with Go.
// @host localhost:8080
// @BasePath /api/v1

func initDB(db *sql.DB) {
	mockdb.ResetDatabase(db)
}

func main() {
	cfg := config.LoadConfig()

	db, err := sql.Open("sqlite3", "./mockserver.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	initDB(db)

	r := mux.NewRouter()

	r.Use(middleware.DetailedLoggingMiddleware)
	r.Use(middleware.CORSMiddleware)
	r.Use(middleware.JWTAuthMiddleware)

	limiter := rate.NewLimiter(1, 5) // 1 request per second, with a burst of 5
	r.Use(middleware.RateLimitMiddleware(limiter))

	r.HandleFunc("/api/v1/resource", handlers.GetResource(db)).Methods("GET")
	r.HandleFunc("/api/v1/resource", handlers.CreateResource(db)).Methods("POST")
	r.HandleFunc("/api/v1/resource/{id}", handlers.UpdateResource(db)).Methods("PUT")
	r.HandleFunc("/api/v1/resource/{id}", handlers.DeleteResource(db)).Methods("DELETE")
	r.HandleFunc("/api/v1/reset", handlers.ResetDatabaseHandler(db)).Methods("POST")

	r.HandleFunc("/api/v1/auth", handlers.GenerateToken).Methods("POST")

	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	srv := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: r,
	}

	go func() {
		log.Printf("Server is running on port %s", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
