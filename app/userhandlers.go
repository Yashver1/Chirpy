package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/Yashver1/chirpy/internal/utils"
	"github.com/golang-jwt/jwt/v5"
)

func (a *apiConfig) CreateUserHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		type Response struct {
			Id    int    `json:"id"`
			Email string `json:"email"`
			IsChirpyRed bool `json:"is_chirpy_red"`
		}

		var user utils.User
		defer r.Body.Close()

		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&user); err != nil {
			utils.RespondWithErr(w, 500, "Unable to parse new user json")
			return
		}

		newUser, err := a.Database.CreateUser(user)
		if err != nil {
			utils.RespondWithErr(w, 500, "Unable to create new user")
			return
		}

		response := Response{
			Id:    newUser.Id,
			Email: newUser.Email,
			IsChirpyRed: newUser.IsChirpyRed,
		}

		utils.RespondWithJSON(w, 201, response)

	})
}

//with jwt header update user

func (a *apiConfig) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	type Response struct{
		Id int `json:"id"`
		Email string `json:"email"`
	}

	var login LoginRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&login); err != nil {
		utils.RespondWithErr(w, 400, "Error parsing request")
		return
	}

	authHeader := r.Header.Get("Authorization")
	authHeader = strings.TrimPrefix(authHeader,"Bearer ")
	token, err := jwt.ParseWithClaims(authHeader,&jwt.RegisteredClaims{},func(token *jwt.Token) (interface{},error){
		return []byte(a.JwtSecret), nil
	})

	if err!=nil{
		utils.RespondWithErr(w,401,"Unauthorized")
		return 
	}

	currentId,err := token.Claims.GetSubject()
	if err!=nil{
		utils.RespondWithErr(w,500,"Server Error")
		return
	}

	parsedId , err := strconv.Atoi(currentId)
	if err!=nil{
		utils.RespondWithErr(w,500,"Server Error")
		return
	}

	if err:= a.Database.UpdateUser(parsedId,login.Email,login.Password); err!=nil{
		utils.RespondWithErr(w,500, "Server Error")
		return
	}

	response := Response{
		Id: parsedId,
		Email: login.Email,
	}

	utils.RespondWithJSON(w,200,response)

}
