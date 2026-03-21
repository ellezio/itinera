package handler

import "net/http"

func redirect(w http.ResponseWriter, r *http.Request, url string, status int) {
	if r.Header.Get("Hx-Request") == "true" {
		w.Header().Set("Hx-Redirect", url)
		w.WriteHeader(status)
	} else {
		http.Redirect(w, r, url, status)
	}
}
