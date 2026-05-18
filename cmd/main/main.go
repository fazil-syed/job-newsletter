package main

import (
	"fmt"
	"jobs-newsletter/internal/api"
	"jobs-newsletter/internal/config"
	"jobs-newsletter/internal/db"
	"jobs-newsletter/static"
	"log"
	"net/http"
)

func AuthMiddleware(next http.HandlerFunc, apiKey string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := r.Header.Get("x-api-key")

		if key != apiKey {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}

func main() {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	db.Init(cfg.Postgres)
	db.DB.AutoMigrate(&db.Subscriber{}, &db.Post{}, &db.Event{})
	port := cfg.Server.Port

	fs := http.FileServer(http.FS(static.StaticFiles))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/index.html")
	})
	http.HandleFunc("/subscribe", api.SubscribeHandler)
	http.HandleFunc("/unsubscribe", api.UnsubscribeHandler)
	http.HandleFunc("/open", api.OpenHandler)
	http.HandleFunc("/click", api.ClickHandler)
	http.HandleFunc("/post", AuthMiddleware(api.CreatePostHandler(cfg), cfg.App.APIKey))
	log.Printf("server running on :%d", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
