package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/hursty1/chirpy/internal/auth"
	"github.com/hursty1/chirpy/internal/database"
)


func validateChirp(rw http.ResponseWriter, req *http.Request) {
	type validateChirp struct {
		Body string `json:"body"`
	}
	type responseMsg struct {
		valid bool `json:"valid"`
		CleanedBody string `json:"cleaned_body"`
	}
	
	decoder := json.NewDecoder(req.Body)
	params := validateChirp{}
	err := decoder.Decode(&params)
	
	if err != nil {
		responseWithError(rw, 400, fmt.Sprintf("Error Decoding Body %s", err))
		return
	}
	if len(params.Body) > 140 {
		responseWithError(rw, 400, "Chirp is too long.")
		return
	}
	cleaned := cleanString(params.Body)
	respBody := responseMsg {
		valid: true,
		CleanedBody: cleaned,
	}
	responseWithJSON(rw,200, respBody)
}


func (cfg *apiConfig)handleAddChirp(rw http.ResponseWriter, req *http.Request) {
	type validChrip struct {
		Body string `json:"body"`
		UserId uuid.UUID `json:"user_id"`
	}
	
	//header
	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		// fmt.Printf("Error getting token from request %s\n", err)
		responseWithError(rw, 401, "Unauthorized")
		return
	}
	//validate
	userID, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		// fmt.Printf("Error validating JWT %s\n", err)
		responseWithError(rw, 401, "Unauthorized")
		return
	}
	
	decoder := json.NewDecoder(req.Body)
	params := validChrip{}
	err = decoder.Decode(&params)
	

	if err != nil {
		responseWithError(rw, 400, fmt.Sprintf("Error Decoding Body %s", err))
		return
	}
	if len(params.Body) > 140 {
		responseWithError(rw, 400, "Chirp is too long.")
		return
	}

	newChrip := database.Chirp{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Body: params.Body,
		UserID: userID,
	}

	dbChirp, err := cfg.db.CreateChirp(context.Background(), database.CreateChirpParams(newChrip))
	if err != nil {
		responseWithError(rw, 400, fmt.Sprintf("Error saving chirp to the database %s", err))
		return
	}
	responseWithJSON(rw, 201, ResponseChirp(dbChirp))
}



func (cfg *apiConfig)handleGetAllChirps(rw http.ResponseWriter, req *http.Request) {
	s := req.URL.Query().Get("author_id")
	sort_url := req.URL.Query().Get("sort")
	var chirps []database.Chirp
	var err error
	if s != "" {
		user_id, err := uuid.Parse(s)
		if err != nil {
			responseWithError(rw, 400, fmt.Sprintf("%s", err))
			return 
		}
		chirps, err = cfg.db.GetAuthorChirps(req.Context(), user_id)
	} else {
		chirps, err = cfg.db.GetAllChirps(req.Context())
		if sort_url == "desc" {
			sort.Slice(chirps, func(i, j int) bool {return chirps[j].CreatedAt.Before(chirps[i].CreatedAt)})
		} else {
			sort.Slice(chirps, func(i, j int) bool {return chirps[i].CreatedAt.Before(chirps[j].CreatedAt)})
		}
	}
	if err != nil {
			responseWithError(rw, 400, fmt.Sprintf("Error Decoding Body %s", err))
			return
		}
	responseChirps := make([]ResponseChirp, len(chirps))
	for index, chirp := range chirps {
		responseChirps[index] = ResponseChirp(chirp)
	}
	responseWithJSON(rw, 200, responseChirps)
}


func (cfg *apiConfig)handleGetChirpById(rw http.ResponseWriter, req *http.Request){
	id := req.PathValue("chirpID")
	

	req_uuid, err := uuid.Parse(id)
	if err != nil {
		// fmt.Println("Error parsing UUID")
		responseWithError(rw, 400, fmt.Sprintf("Error saving chirp to the database %s", err))
		return
	}
	chirp, err := cfg.db.GetChirpById(req.Context(), req_uuid)
	if err != nil {
		responseWithError(rw, 404, fmt.Sprintf("Unable to find chirp with id:%s", id))
		return
	}
	// fmt.Println(chirp)
	responseWithJSON(rw, 200, ResponseChirp(chirp))
}

func (cfg *apiConfig)HandleDeleteChripById(rw http.ResponseWriter, req *http.Request){
	id := req.PathValue("chirpID")
	
	req_uuid, err := uuid.Parse(id)
	if err != nil {
		// fmt.Println("Error parsing UUID")
		responseWithError(rw, 400, fmt.Sprintf("Error finding chirp to the database %s", err))
		return
	}
	access_token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		responseWithError(rw, 401, fmt.Sprintf("Unable to create refresh token: %s", err))
		return
	}
	UserID, err := auth.ValidateJWT(access_token, cfg.secret)
	if err != nil {
		responseWithError(rw, 401, fmt.Sprintf("Invalid Access token: %s", err))
		return
	}
	chirp, err := cfg.db.GetChirpById(req.Context(), req_uuid)
	if err != nil {
		responseWithError(rw, 404, fmt.Sprintf("Unable to find chirp with id:%s", id))
		return
	}

	if chirp.UserID != UserID {
		responseWithError(rw, 403, fmt.Sprintf("Unable to delete chrip for chrip you do not own"))
		return
	}

	err  = cfg.db.DeleteChripById(req.Context(), chirp.ID)
	if err != nil {
		responseWithError(rw, 403, fmt.Sprintf("Unable to delete chrip."))
		return
	}

	responseWithOutBody(rw, 204)
}

func (cfg *apiConfig)HandleChirpyRed(rw http.ResponseWriter, req *http.Request){
		
	type ChirpyRedBody struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}
	//auth check
	polka_key, err := auth.GetPolkaApiKey(req.Header)
	if err != nil || polka_key != cfg.Polka_key{
		// fmt.Println(err)
		// fmt.Println(cfg.Polka_key)
		// fmt.Println(polka_key)
		responseWithError(rw, 401, fmt.Sprintf("Unauthorized."))
		return
	}
	decoder := json.NewDecoder(req.Body)
	res_body := ChirpyRedBody{}
	err = decoder.Decode(&res_body)
	if err != nil {
		responseWithError(rw, 400, fmt.Sprintf("Unable to process body %s", err))
		return
	}
	if res_body.Event != "user.upgraded" {
		responseWithOutBody(rw, 204) //not processing
		return 
	}
	user_id, err := uuid.Parse(res_body.Data.UserID)
	if err != nil {
		responseWithError(rw, 400, "")
		return
	}
	update, err := cfg.db.UpdateIsRed(req.Context(), user_id)
	if err != nil {
		responseWithError(rw, 404, "Unable to upgrade user")
		return
	}
	// fmt.Println(update)
	if update.IsChirpyRed != true {
		responseWithError(rw, 400, "unable to update user to red")
	}

	responseWithOutBody(rw, 204)
}