package main

import (
	"time"

	"github.com/google/uuid"
	"github.com/hursty1/chirpy/internal/database"
)


type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	HashedPassword string
	IsChirpyRed bool `json:"is_chirpy_red"`
}
type UserLoginResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	IsChirpyRed bool `json:"is_chirpy_red"`
}
func NewUserResponse(u database.User) UserResponse {
	return UserResponse{
		ID: u.ID,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		Email: u.Email,
		IsChirpyRed: u.IsChirpyRed,
	}
}
type ResponseChirp struct {
		ID uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body string `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}
