package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"myapp/models"

	"golang.org/x/crypto/bcrypt"
)

// Define la constante para UNIQUE constraint.
const ErrConstraintUnique = 2067

// PostUserRegisterHandler maneja el registro de nuevos usuarios.
func PostUserRegisterHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.RegisterRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Printf("Error decoding register request: %v", err)
			response := models.NewErrorResponse("Invalid request body")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
			return
		}

		if req.Username == "" || req.Password == "" {
			response := models.NewErrorResponse("Username and password cannot be empty")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Error hashing password for user %s: %v", req.Username, err)
			response := models.NewErrorResponse("Internal server error during registration setup")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
			return
		}

		result, err := db.ExecContext(r.Context(),
			"INSERT INTO users(username, password_hash) VALUES(?, ?)",
			req.Username, string(hashedPassword),
		)

		userID, err := result.LastInsertId()
		if err != nil {
			log.Printf("Error getting last insert ID after registering user %s: %v", req.Username, err)
			response := models.NewErrorResponse("Registration partially successful, but failed to retrieve user ID")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
			return
		}

		log.Printf("User '%s' (ID: %d) registered successfully via PostUserRegisterHandler.", req.Username, userID)
		registerData := models.RegisterSuccessData{
			UserID:   userID,
			Username: req.Username,
		}

		response := models.NewSuccessResponse(registerData)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	}
}
