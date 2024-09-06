package main

import (
	"log"
	"net/http"

	"github.com/Yashver1/chirpy/internal/utils"
)

func healthCheckHandler(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))

}

func main() {

	const port = "8080"
	const healthPath = "/healthz"
	const apiPath = "/api"
	const metricsPath = "/metrics"
	const resetPath = "/reset"
	const adminPath = "/admin"
	const chirpPath = "/chirps"
	const dbPath = "../storage/database.json"

	apiConfig := apiConfig{}
	DB, err := utils.NewDB(dbPath)
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("."))
	mux.Handle("/app/", apiConfig.incrementMetrics(http.StripPrefix("/app/", fs)))

	mux.HandleFunc("GET "+apiPath+healthPath, healthCheckHandler)
	mux.HandleFunc("GET "+adminPath+metricsPath, apiConfig.getMetrics)
	mux.HandleFunc(apiPath+resetPath, apiConfig.resetMetrics)
	mux.Handle("POST "+apiPath+chirpPath, utils.CreateChirpHandler(DB))
	mux.Handle("GET "+apiPath+chirpPath, utils.GetAllChirpsHandler(DB))

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Listening at port: %s", port)
	log.Fatal(server.ListenAndServe())
}
