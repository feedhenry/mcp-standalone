package middleware

import "net/http"

type Cors struct{}

func (c Cors) Handle(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	headers := w.Header()
	if headers.Get("Access-Control-Allow-Origin") == "" {
		headers.Add("Access-Control-Allow-Origin", "*")
	}
	if headers.Get("Access-Control-Allow-Methods") == "" {
		headers.Add("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE")
	}
	if headers.Get("Access-Control-Allow-Headers") == "" {
		headers.Add("Access-Control-Allow-Headers", "x-auth, content-type")
	}
	if req.Method != "OPTIONS" {
		next(w, req)
		return
	}
}
