package api

import "net/http"

func displayLanding(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("You have successfully connected using this API"))
}
