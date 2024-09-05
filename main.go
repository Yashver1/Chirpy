package main

import (
	"log"
	"net/http"
)

func healthCheckHandler(w http.ResponseWriter, req *http.Request){

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))


}

func main(){
	
	const port = "8080"
	const healthPath = "/healthz"
	const apiPath = "/api"
	const metricsPath = "/metrics"
	const resetPath = "/reset"
	const adminPath = "/admin"
	const chirpPath = "/chirps"

	apiConfig := apiConfig{}

	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("./app"))
	mux.Handle("/app/",apiConfig.incrementMetrics(http.StripPrefix("/app/",fs)))

	mux.HandleFunc("GET " + apiPath + healthPath, healthCheckHandler)
	mux.HandleFunc("GET " + adminPath + metricsPath, apiConfig.getMetrics)
	mux.HandleFunc(apiPath + resetPath, apiConfig.resetMetrics)
	mux.HandleFunc("POST " + apiPath + chirpPath, chirpValidateHandler)

	server := &http.Server{
		Addr: ":" + port,
		Handler: mux,
	}



	log.Printf("Listening at port: %s",port)
	log.Fatal(server.ListenAndServe())
}
