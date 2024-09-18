package main

import (
	"encoding/json"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/golang-jwt/jwt/v5"

	"github.com/Yashver1/chirpy/internal/utils"
)



func (a *apiConfig) CreateChirpHandler() http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		defer req.Body.Close()
		var chirp utils.Chirp
		database := a.Database

		authHeader := req.Header.Get("Authorization")
		if authHeader == ""{
			utils.RespondWithErr(w,401,"Unauthorized")
			return
		}

		authHeader = strings.TrimPrefix(authHeader,"Bearer ")
		token, err := jwt.ParseWithClaims(authHeader,&jwt.RegisteredClaims{},func(token *jwt.Token) (interface{},error){
			return []byte(a.JwtSecret),nil
		})

		if err!=nil{
			utils.RespondWithErr(w,500,"Unable to parse JWT")
			return
		}

		userId,err := token.Claims.GetSubject()
		if err!=nil{
			utils.RespondWithErr(w,500,"Server Error")
			return
		}

		parsedId,err := strconv.Atoi(userId)
		if err!=nil{
			utils.RespondWithErr(w,500,"Server Error")
			return
		}

		decoder := json.NewDecoder(req.Body)
		if err := decoder.Decode(&chirp); err!=nil{
			utils.RespondWithErr(w,500,"Unable to parse request")
			return
		}


		

		///validation///

		if utf8.RuneCountInString(chirp.Body) > 140 {
			utils.RespondWithErr(w, 400, "Chirp is too long")
			return
		}

		jsonArray := strings.Split(chirp.Body, " ")
		for i := 0; i < len(jsonArray); i++ {
			lowerString := strings.ToLower(jsonArray[i])
			if lowerString == "kerfuffle" || lowerString == "sharbert" || lowerString == "fornax" {
				jsonArray[i] = "****"
			}
		}

		resp := strings.Join(jsonArray, " ")
		///validationEnd///

		createdChirp, err := database.CreateChirp(resp,parsedId)
		if err != nil {
			utils.RespondWithErr(w, 500, "unable to create Chirp")
			return
		}

		utils.RespondWithJSON(w, 201, createdChirp)
	})
}

func (a *apiConfig) GetAllChirpsHandler() http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		database := a.Database
		chirps, err := database.GetChirps()
		if err != nil {
			utils.RespondWithErr(w, 500, "unable to get chirps")
			return
		}

		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].Id < chirps[j].Id
		})

		utils.RespondWithJSON(w, 200, chirps)

	})
}

func (a *apiConfig) GetChirpHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		database := a.Database
		idString := r.PathValue("chirpID")
		chirp, err := database.GetChirp(idString)
		if err != nil {
			utils.RespondWithErr(w, 404, "Chirp does not exist")
			return
		}

		utils.RespondWithJSON(w, 200, chirp)
	})
}


func (a *apiConfig) DeleteChirpHandler() http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		defer r.Body.Close()

		idstring := r.PathValue("chirpID")

		authHeader := r.Header.Get("Authorization")
		if authHeader == ""{
			utils.RespondWithErr(w,401,"Unauthorized")
			return
		}

		authHeader = strings.TrimPrefix(authHeader,"Bearer ")
		token, err := jwt.ParseWithClaims(authHeader,&jwt.RegisteredClaims{},func(token *jwt.Token) (interface{},error){
			return []byte(a.JwtSecret),nil
		})

		if err!=nil{
			utils.RespondWithErr(w,500,"Unable to parse JWT")
			return
		}

		userId,err := token.Claims.GetSubject()
		if err!=nil{
			utils.RespondWithErr(w,500,"Server Error")
			return
		}

		parsedId,err := strconv.Atoi(userId)
		
		if err!=nil{
			utils.RespondWithErr(w,500,"Server Error")
			return
		}

		if err:= a.Database.DeleteChirp(idstring,parsedId); err!=nil{
			utils.RespondWithErr(w,403,"Unauthorized")
			return
		}

		utils.RespondWithJSON(w,204,"Chirp deleted successfully")

	})
}
