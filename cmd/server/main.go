package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/swaggo/http-swagger/example/go-chi/docs"
	"github.com/yourusername/exchange-rate-service/internal/api"
	"github.com/yourusername/exchange-rate-service/internal/cache"
	"github.com/yourusername/exchange-rate-service/internal/repository"
	"github.com/yourusername/exchange-rate-service/internal/service"
)

func main() {
	fmt.Println("Exchange Rate Service starting up...")

	// Setup dependencies
	repo := repository.NewAPIRepository()
	cache := cache.NewRateCache()
	svc := service.NewService(repo, cache)
	endpoints := api.MakeEndpoints(svc)
	handler := api.MakeHTTPHandler(endpoints)

	mux := http.NewServeMux()
	mux.Handle("/", handler)
	// Serve Swagger UI and static docs at /docs/ (note trailing slash)
	mux.Handle("/docs/", http.StripPrefix("/docs/", http.FileServer(http.Dir("docs"))))
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := ":" + port
	fmt.Printf("Listening on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
