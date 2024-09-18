package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Yashver1/chirpy/internal/utils"
)

type WebHookPolka struct {
	Event string `json:"event"`
	Data struct {
		UserID int `json:"user_id"`
	} `json:"data"`
}


func (a *apiConfig) WebhookHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		var webhook WebHookPolka

		authHeader := r.Header.Get("Authorization")
		if authHeader == ""{
			utils.RespondWithErr(w,401,"Unauthorized")
			return
		}

		authHeader = strings.TrimPrefix(authHeader,"ApiKey ")
		if authHeader != a.PolkaKey{
			utils.RespondWithErr(w,401,"Unauthorized")
			return
		}

		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&webhook); err != nil {
			utils.RespondWithErr(w,400,"Unable to decode webhook")
			return
		}

		if webhook.Event != "user.upgraded"{
			utils.RespondWithJSON(w,204,nil)
			return
		}

		if err := a.Database.UpgradeUser(webhook.Data.UserID); err!=nil{
			utils.RespondWithErr(w,404,"Unable to upgrade user")
			return
		}

		utils.RespondWithJSON(w,204,nil)
		
	})
}