package main

import "net/http"

func handleRequest(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, 200, struct{}{})
}
func handleErr(w http.ResponseWriter, r *http.Request) {
	respondWithError(w, 399, "Something went wrong")
}
