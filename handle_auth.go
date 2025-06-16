package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/hursty1/chirpy/internal/auth"
	"github.com/hursty1/chirpy/internal/database"
)


func (cfg *apiConfig)handleAddUser(rw http.ResponseWriter, req *http.Request) {
	type UserRequest struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}
	decoder := json.NewDecoder(req.Body)
	// fmt.Println(decoder)
	user := UserRequest{}
	err := decoder.Decode(&user)
	if err != nil {
		responseWithError(rw, 400, fmt.Sprintf("Error decoding body %s", err))
		return
	}
	hashed_password, err := auth.HashPassword(user.Password)
	if err != nil {
		responseWithError(rw, 400, fmt.Sprintf("%s", err))
	}
	// log.Printf("User email is: %s\n", user.Email)
	userParam := database.CreateUserParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Email: user.Email,
		HashedPassword: hashed_password,

	}
	created_usr, err := cfg.db.CreateUser(context.Background(),userParam)
	if err != nil {
		responseWithError(rw, 400, fmt.Sprintf("Error creating user: %s", err))
		return
	}
	response_user := UserLoginResponse{
		ID: created_usr.ID,
		CreatedAt: created_usr.CreatedAt,
		UpdatedAt: created_usr.UpdatedAt,
		Email: created_usr.Email,
		IsChirpyRed: created_usr.IsChirpyRed,
	}
	responseWithJSON(rw,201,response_user)
	return
}


func (cfg *apiConfig)HandleLogin(rw http.ResponseWriter, req *http.Request){
	type UserRequest struct {
		Email string `json:"email"`
		Password string `json:"password"`
		// ExpiresIn int `json:"expires_in_seconds"`
	}
	type UserLoginResponseLogin struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
		Token 	  string    `json:"token"`
		RefreshToken string `json:"refresh_token"`
		IsChirpyRed bool `json:"is_chirpy_red"`
	}
	decoder := json.NewDecoder(req.Body)
	defer req.Body.Close()
	user := UserRequest{}
	err := decoder.Decode(&user)
	if err != nil {
		responseWithError(rw, 400, fmt.Sprintf("Error decoding body %s", err))
		return
	}
	db_user, err := cfg.db.GetUserFromEmail(req.Context(), user.Email)
	if err != nil {
		responseWithError(rw, 400, fmt.Sprintf("Wrong Email address, or user does not exist"))
		return
	}
	pw_err := auth.CheckPasswordHash(db_user.HashedPassword, user.Password)
	if pw_err != nil {
		responseWithError(rw, 401, fmt.Sprintf("Wrong Password!!"))
		return
	}
	
	expires := 60*60
	// fmt.Println("Expires in: %d", expires)
	jwt, err := auth.MakeJWT(db_user.ID,cfg.secret,time.Duration(expires * int(time.Second)))
	if err != nil {
		responseWithError(rw, 400, fmt.Sprintf("Error Creating JWT: %s", err))
		return 
	}

	//refresh token
	refresh_token, err := auth.MakeRefreshToken()
	if err != nil {
		responseWithError(rw, 400, fmt.Sprintf("Unable to create fresh token %s", err))
		return
	}

	tokenExpires := time.Hour * 24 * 60 // 60 days
	dateExpires := time.Now().Add(tokenExpires)
	token_params := database.AddRefreshTokenParams{
		Token: refresh_token,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID: db_user.ID,
		ExpiresAt: dateExpires,
	}

	db_token, err := cfg.db.AddRefreshToken(req.Context(),token_params)
	if err != nil {
		responseWithError(rw, 400, fmt.Sprintf("Unable to create fresh token %s", err))
		return
	}

	user_response := UserLoginResponseLogin{
		ID: db_user.ID,
		CreatedAt: db_user.CreatedAt,
		UpdatedAt: db_user.UpdatedAt,
		Email: db_user.Email,
		Token: jwt,
		RefreshToken: db_token.Token,
		IsChirpyRed: db_user.IsChirpyRed,
	}


	responseWithJSON(rw,200,user_response)
}


func (cfg *apiConfig)HandleRefresh(rw http.ResponseWriter, req *http.Request){
	type RefreshResponse struct {
		Token string `json:"token"`
	}
	refresh_token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		responseWithError(rw, 401, fmt.Sprintf("Unable to create refresh token: %s", err))
		return
	}

	db_token, err := cfg.db.FetchFreshToken(req.Context(), refresh_token)
	if err != nil {
		responseWithError(rw, 401, fmt.Sprintf("Unable to create refresh token. Does not exist: "))
		return
	}

	ok, err := auth.CheckRefreshToken(db_token)
	if !ok || err != nil {
		responseWithError(rw, 401, fmt.Sprintf("Refresh Token is expired or revoked"))
		return
	}

	//refresh token good create new token and return
	token, err := auth.MakeJWT(db_token.UserID,cfg.secret,time.Duration(60*int(time.Second)))
	if err != nil {
		responseWithError(rw, 401, fmt.Sprintf("Unable to create JWT."))
		return
	}
	rr := RefreshResponse{
		Token: token,
	}
	responseWithJSON(rw,200,rr)
}

func (cfg *apiConfig)HandleRevoke(rw http.ResponseWriter, req *http.Request){
	refresh_token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		responseWithError(rw, 401, fmt.Sprintf("Unable to create refresh token: %s", err))
		return
	}

	_, err = cfg.db.FetchFreshToken(req.Context(), refresh_token)
	if err != nil {
		responseWithError(rw, 401, fmt.Sprintf("Unable to create refresh token. Does not exist: "))
		return
	}

	_, err = cfg.db.RevokeRefreshToken(req.Context())
	if err != nil {
		responseWithError(rw, 401, fmt.Sprintf("Unable to revoke Token."))
	}

	responseWithOutBody(rw, 204)
}


func (cfg *apiConfig)HandleUserUpdate(rw http.ResponseWriter, req *http.Request){
	type UserUpdateRequestBody struct {
		Password string `json:"password"`
		Email string `json:"email"`
	}
	
	access_token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		responseWithError(rw, 401, fmt.Sprintf("Unable to create refresh token: %s", err))
		return
	}

	decoder := json.NewDecoder(req.Body)
	defer req.Body.Close()
	user := UserUpdateRequestBody{}
	err = decoder.Decode(&user)
	if err != nil {
		responseWithError(rw, 400, fmt.Sprintf("Error decoding body %s", err))
		return
	}

	UserID, err := auth.ValidateJWT(access_token, cfg.secret)
	if err != nil {
		responseWithError(rw, 401, fmt.Sprintf("Invalid Access token: %s", err))
		return
	}
	hashed_password, err := auth.HashPassword(user.Password)
	if err != nil {
		responseWithError(rw, 401, fmt.Sprintf("Invalid password unable to hash: %s", err))
		return
	}
	updateParams := database.UpdateEmailAndPasswordParams{
		Email: user.Email,
		HashedPassword: hashed_password,
		ID: UserID,
	}

	new_user, err := cfg.db.UpdateEmailAndPassword(req.Context(),updateParams)
	if err != nil {
		responseWithError(rw, 401, fmt.Sprintf("Invalid password unable to hash: %s", err))
		return
	}
	user_response := UserLoginResponse{
		ID: new_user.ID,
		CreatedAt: new_user.CreatedAt,
		UpdatedAt: new_user.UpdatedAt,
		Email: new_user.Email,
		IsChirpyRed: new_user.IsChirpyRed,
	}
	responseWithJSON(rw,200,user_response)
}