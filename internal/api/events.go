package api

import (
	"fmt"
	"jobs-newsletter/internal/db"
	"net/http"
)

func parseUint(s string) uint {
	var id uint
	fmt.Sscanf(s, "%d", &id)
	return id
}
func OpenHandler(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	postID := r.URL.Query().Get("post_id")

	db.DB.Create(&db.Event{
		Subscriber: email,
		Type:       "open",
		PostID:     parseUint(postID),
	})

	w.Header().Set("Content-Type", "image/png")
	pixel := []byte{
		71, 73, 70, 56, 57, 97, 1, 0, 1, 0, 128, 0, 0,
		0, 0, 0, 255, 255, 255, 33, 249, 4, 1, 0, 0,
		1, 0, 44, 0, 0, 0, 0, 1, 0, 1, 0, 0, 2, 2,
		68, 1, 0, 59,
	}

	w.Write(pixel)
}
func ClickHandler(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	postID := r.URL.Query().Get("post_id")
	url := r.URL.Query().Get("url")

	db.DB.Create(&db.Event{
		Subscriber: email,
		Type:       "click",
		PostID:     parseUint(postID),
	})

	http.Redirect(w, r, url, http.StatusFound)
}
