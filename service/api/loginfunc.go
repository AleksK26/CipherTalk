package api

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/AleksK26/WASA_AleksK_2024-25/service/api/reqcontext"
	"github.com/AleksK26/WASA_AleksK_2024-25/service/database"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if len(req.Name) < 3 || len(req.Name) > 16 {
		http.Error(w, "Username must be between 3 and 16 characters", http.StatusBadRequest)
		return
	}

	if req.Password == "" {
		http.Error(w, "Password is required", http.StatusBadRequest)
		return
	}

	var photoBytes []byte
	if req.Photo != "" {
		var err error
		photoBytes, err = base64.StdEncoding.DecodeString(req.Photo)
		if err != nil {
			ctx.Logger.WithError(err).Error("Invalid base64 photo data")
			http.Error(w, "Invalid photo data", http.StatusBadRequest)
			return
		}
	}

	existingUser, err := rt.db.GetUserByName(req.Name)

	if req.Mode == "signup" {
		// Sign up: user must NOT exist
		if err == nil {
			http.Error(w, "Username already taken", http.StatusConflict)
			return
		}
		if !errors.Is(err, database.ErrUserDoesNotExist) {
			ctx.Logger.WithError(err).Error("Failed to query user by name")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		newID, genErr := generateNewID()
		if genErr != nil {
			ctx.Logger.WithError(genErr).Error("Failed to generate user ID")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		newUser := database.User{
			Id:       newID,
			Name:     req.Name,
			Photo:    photoBytes,
			Password: req.Password,
		}
		createdUser, createErr := rt.db.CreateUser(newUser)
		if createErr != nil {
			ctx.Logger.WithError(createErr).Error("Failed to create user")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(struct {
			Identifier string `json:"identifier"`
		}{Identifier: createdUser.Id})
		return
	}

	// Default mode: "signin"
	if errors.Is(err, database.ErrUserDoesNotExist) {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to query user by name")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// If user has no password yet (legacy account), accept any password and set it
	if existingUser.Password != "" && existingUser.Password != req.Password {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(struct {
		Identifier string `json:"identifier"`
	}{Identifier: existingUser.Id})
}
