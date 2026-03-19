package api

import (
	"encoding/json"
	"net/http"

	"jobs-newsletter/internal/config"
	"jobs-newsletter/internal/db"
	"jobs-newsletter/internal/email"
)

type CreatePostRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func CreatePostHandler(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreatePostRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil || req.Title == "" || req.Content == "" {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		post := db.Post{
			Title:   req.Title,
			Content: req.Content,
		}

		if err := db.DB.Create(&post).Error; err != nil {
			http.Error(w, "failed to create post", http.StatusInternalServerError)
			return
		}

		go email.SendPost(cfg, post)

		w.Write([]byte("post created and sending started"))
	}
}
