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
		http.Error(w, "Invalid username length", http.StatusBadRequest)
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
	user, err := rt.db.GetUserByName(req.Name)
	if errors.Is(err, database.ErrUserDoesNotExist) {
		newID, genErr := generateNewID()
		if genErr != nil {
			ctx.Logger.WithError(genErr).Error("Failed to generate user ID")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		newUser := database.User{
			Id:    newID,
			Name:  req.Name,
			Photo: photoBytes,
		}
		user, err = rt.db.CreateUser(newUser)
		if err != nil {
			ctx.Logger.WithError(err).Error("Failed to create user")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	} else if err != nil {
		ctx.Logger.WithError(err).Error("Failed to query user by name")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Always respond with JSON containing the identifier
	resp := struct {
		Identifier string `json:"identifier"`
	}{
		Identifier: user.Id,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		ctx.Logger.WithError(err).Error("Failed to encode login response")
	}
}
