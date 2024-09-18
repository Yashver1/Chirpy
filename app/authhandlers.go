package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Yashver1/chirpy/internal/utils"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Password string `json:"password"`
	Email    string `json:"email"`
	Deadline int   `json:"expires_in_seconds,omitempty"`
}

//simple Password hashed based auth

func (a *apiConfig) PasswordLoginHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		database := a.Database

		var login LoginRequest
		type Response struct {
			Id    int    `json:"id"`
			Email string `json:"email"`
		}

		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&login); err != nil {
			utils.RespondWithErr(w, 400, "Unable to parse request")
			return
		}

		users, err := database.GetUsers()
		if err != nil {
			utils.RespondWithErr(w, 400, fmt.Sprintf("%s", err))
			return
		}

		var foundUser utils.User

		for _, element := range users {
			if element.Email == login.Email {
				foundUser = element
				break
			}
		}

		if err := bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(login.Password)); err != nil {
			utils.RespondWithErr(w, 401, "Incorrect email or password")
			return
		}

		response := Response{
			Id:    foundUser.Id,
			Email: foundUser.Email,
		}

		utils.RespondWithJSON(w, 200, response)
	})
}

func (a* apiConfig) DeleteRefreshTokenHandler() http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		defer r.Body.Close()
		authHeader := r.Header.Get("Authorization")
		if authHeader == ""{
			utils.RespondWithErr(w,400,"No authorization header sent")
			return
		}

		refreshtoken := strings.TrimPrefix(authHeader,"Bearer ")
		err := a.Database.DeleteToken(refreshtoken)
		if err!=nil{
			utils.RespondWithErr(w,500,"Error deleting token")
			return
		}

		utils.RespondWithJSON(w,204,nil)
		})
}



//refresh token based auth

func (a *apiConfig) RefreshTokenHandler() http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		defer r.Body.Close()
		database := a.Database

		type Response struct {
			Token string `json:"token"`
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader == ""{
			utils.RespondWithErr(w,400,"No authorization header sent")
			return
		}

		refreshtoken := strings.TrimPrefix(authHeader,"Bearer ")
		userId, err := database.GetToken(refreshtoken)
		if err!=nil{
			utils.RespondWithErr(w,401,"Invalid refresh token")
			return
		}

		currentTime := time.Now().UTC()
		tokenExpiry := currentTime.Add(1 * time.Hour)

		claims := jwt.RegisteredClaims{
			Issuer: "chirpy",
			IssuedAt: jwt.NewNumericDate(currentTime),
			ExpiresAt: jwt.NewNumericDate(tokenExpiry),
			Subject: strconv.Itoa(userId),
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256,claims)
		signedJWT, err := token.SignedString([]byte(a.JwtSecret))
		if err!=nil{
			utils.RespondWithErr(w,500,"Unable to generate JWT")
			return
		}

		response := Response{
			Token: signedJWT,
		}

		utils.RespondWithJSON(w,200,response)
		
	

	})

}



//jwt based auth

func (a *apiConfig) JWTLoginHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		defer r.Body.Close()
		var login LoginRequest

		type Response struct {
			Id    int    `json:"id"`
			Email string `json:"email"`
			Token string `json:"token"`
			RefreshToken string `json:"refresh_token"`
		}

		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&login); err != nil {
			utils.RespondWithErr(w, 400, "Unable to Parse Response")
			return
		}


		users, err := a.Database.GetUsers()
		var foundUser utils.User

		for _,value := range users{
			if value.Email == login.Email{
				foundUser = value
				break
			}
		}

		if err:= bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(login.Password)); err!=nil{
			utils.RespondWithErr(w,401, "Incorrect email or password")
			return
		}


		currentTime := time.Now().UTC()
		expirationTime := currentTime.Add(1 * time.Hour)
		

		claims := jwt.RegisteredClaims{
			Issuer: "chirpy",
			IssuedAt: jwt.NewNumericDate(currentTime),
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			Subject: strconv.Itoa(foundUser.Id),
		
		}



		token := jwt.NewWithClaims(jwt.SigningMethodHS256,claims)
		signedJWT, err := token.SignedString([]byte(a.JwtSecret))
		if err!=nil{
			utils.RespondWithErr(w,500,"Error processing request")
			return
		}

		refreshToken := make([]byte, 32)
		_, err = rand.Read(refreshToken)
		if err!=nil{
			utils.RespondWithErr(w,500,"Error generating refresh token")
			return
		}

		refreshTokenString := hex.EncodeToString(refreshToken)
		if err := a.Database.CreateToken(foundUser.Id, refreshTokenString); err != nil {
			utils.RespondWithErr(w, 500, "Error generating refresh token")
			return
		}


		response := Response{
			Id: foundUser.Id,
			Email: foundUser.Email,
			Token: signedJWT,
			RefreshToken: refreshTokenString,
		}

		utils.RespondWithJSON(w,200,response)


	})
}