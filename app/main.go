package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/Yashver1/chirpy/internal/utils"
	"github.com/joho/godotenv"
)

func healthCheckHandler(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))

}

func main() {

	godotenv.Load()

	const port = "8080"
	const healthPath = "/healthz"
	const apiPath = "/api"
	const metricsPath = "/metrics"
	const resetPath = "/reset"
	const adminPath = "/admin"
	const chirpPath = "/chirps"
	const loginPath = "/login"
	const dbPath = "../storage/database.json"

	dbg := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()

	if *dbg {

		os.Remove(dbPath)

	}

	DB, err := utils.NewDB(dbPath)
	if err != nil {
		log.Fatal(err)
	}


	apiConfig := apiConfig{
		FileServerHits: 0,
		Database: DB,
		JwtSecret: os.Getenv("JWT_SECRET"),
	}

	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("."))
	mux.Handle("/app/", apiConfig.incrementMetrics(http.StripPrefix("/app/", fs)))

	mux.HandleFunc("GET "+apiPath+healthPath, healthCheckHandler)
	mux.HandleFunc("GET "+adminPath+metricsPath, apiConfig.getMetrics)
	mux.HandleFunc(apiPath+resetPath, apiConfig.resetMetrics)
	mux.Handle("POST "+apiPath+chirpPath, apiConfig.CreateChirpHandler())
	mux.Handle("GET "+apiPath+chirpPath, apiConfig.GetAllChirpsHandler())
	mux.Handle("GET "+apiPath+chirpPath+"/{chirpID}", apiConfig.GetChirpHandler())
	mux.Handle("POST "+apiPath+"/users", apiConfig.CreateUserHandler())
	mux.Handle("POST "+apiPath+loginPath, apiConfig.JWTLoginHandler())
	mux.HandleFunc("PUT "+apiPath + "/users", apiConfig.UpdateUserHandler)
	mux.Handle("POST " + apiPath + "/revoke", apiConfig.DeleteRefreshTokenHandler())
	mux.Handle("POST " + apiPath + "/refresh", apiConfig.RefreshTokenHandler())
	mux.Handle("DELETE "+apiPath+chirpPath+"/{chirpID}", apiConfig.DeleteChirpHandler())

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Listening at port: %s", port)
	log.Fatal(server.ListenAndServe())
}
