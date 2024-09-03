package main

import (
	"fmt"
	"net/http"
)


type apiConfig struct{
	fileServerHits int
}

func (a *apiConfig) incrementMetrics(next http.Handler) http.Handler{
	return http.HandlerFunc(
		func (w http.ResponseWriter, req *http.Request) {
			a.fileServerHits++
			next.ServeHTTP(w,req)
		})
}

func (a *apiConfig) getMetrics(w http.ResponseWriter, req *http.Request){
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(200)
	logText:=fmt.Sprintf(`
<html>

<body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
</body>

</html>`, a.fileServerHits)
	w.Write([]byte(logText))
}

func (a* apiConfig) resetMetrics(w http.ResponseWriter, req *http.Request){
	a.fileServerHits = 0
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	logText := fmt.Sprintf("Hits Reset at: %v", a.fileServerHits)
	w.Write([]byte(logText))
}