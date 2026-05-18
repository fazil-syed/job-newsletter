package utils

import (
	"jobs-newsletter/static"
	"net/http"
)

func RenderHTML(w http.ResponseWriter, filename string) {
	data, err := static.StaticFiles.ReadFile(filename)
	if err != nil {
		http.Error(w, "file not found", http.StatusInternalServerError)
		return
	}

	w.Write(data)
}
