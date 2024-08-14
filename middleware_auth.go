package main

import (
	"fmt"
	"net/http"

	"github.com/Anas-Sayed0/rss-agg/internal/auth"
	"github.com/Anas-Sayed0/rss-agg/internal/database"
)

type authHandler func(http.ResponseWriter, *http.Request, database.User)

func (apiCfg *apiConfig) middlewareAuth(handler authHandler) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
		apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		responseWithError(w, 403, fmt.Sprintf("Error getting API key: %s",err))
		return
	}
	user, err:=apiCfg.DB.GetUserByAPIKey(r.Context(), apiKey)
	if err != nil {
		responseWithError(w, 400, fmt.Sprintf("Error getting user: %s",err))
		return
	}
	handler(w, r, user)
	}
}