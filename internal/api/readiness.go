package api

import "net/http"

func handlerReadiness(w http.ResponseWriter, r *http.Request) { //handler function for readiness endpoint
	w.Header().Add("Content-type", "text/plain; charset=utf-8") //sets content-type header
	w.WriteHeader(http.StatusOK)                                // sets status code to 200
	w.Write([]byte(http.StatusText(http.StatusOK)))             // write the body text
}
