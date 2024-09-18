package main

import (
	"fmt"
	"net/http"

	"github.com/Yashver1/chirpy/internal/utils"
)

type apiConfig struct {
	FileServerHits int
	Database  *utils.DB
	JwtSecret string
	PolkaKey string
}

func (a *apiConfig) incrementMetrics(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, req *http.Request) {
			a.FileServerHits++
			next.ServeHTTP(w, req)
		})
}

func (a *apiConfig) getMetrics(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(200)
	logText := fmt.Sprintf(`
<html>

<body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
</body>

</html>`, a.FileServerHits)
	w.Write([]byte(logText))
}

func (a *apiConfig) resetMetrics(w http.ResponseWriter, req *http.Request) {
	a.FileServerHits = 0
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	logText := fmt.Sprintf("Hits Reset at: %v", a.FileServerHits)
	w.Write([]byte(logText))
}
