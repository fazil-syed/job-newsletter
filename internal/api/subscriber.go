package api

import (
	"encoding/json"
	"net/http"

	"jobs-newsletter/internal/db"
	"jobs-newsletter/internal/utils"
)

type SubscribeRequest struct {
	Email string `json:"email"`
}

func SubscribeHandler(w http.ResponseWriter, r *http.Request) {
	var email string

	// form (browser)
	if r.Method == http.MethodPost {
		email = r.FormValue("email")
	}

	// json (api)
	if email == "" {
		var req SubscribeRequest
		json.NewDecoder(r.Body).Decode(&req)
		email = req.Email
	}

	if email == "" {
		http.Error(w, "invalid email", http.StatusBadRequest)
		return
	}

	sub := db.Subscriber{Email: email}

	err := db.DB.Create(&sub).Error
	if err != nil {
		w.Write([]byte("already subscribed"))
		return
	}
	utils.RenderHTML(w, "subscribed.html")
}
func UnsubscribeHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		utils.RenderHTML(w, "unsubscribe.html")
		return
	case http.MethodPost:

		email := r.URL.Query().Get("email")

		if email == "" {
			http.Error(w, "missing email", http.StatusBadRequest)
			return
		}

		db.DB.Where("email = ?", email).Delete(&db.Subscriber{})

		utils.RenderHTML(w, "unsubscribed.html")

		// w.Write([]byte("You have been unsubscribed"))

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)

	}
}
