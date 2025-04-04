package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"myapp/models"

	"golang.org/x/crypto/bcrypt"
)

func PostRegisterHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.RegisterRequest
		w.Header().Set("Content-Type", "application/json")

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Printf("Error decoding register request: %v", err)
			response := models.NewErrorResponse("Invalid request body")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
			return
		}

		if req.Username == "" || req.Password == "" {
			response := models.NewErrorResponse("Username and password cannot be empty")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Error hashing password for user %s: %v", req.Username, err)
			response := models.NewErrorResponse("Internal server error during registration setup")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
			return
		}

		result, err := db.ExecContext(r.Context(),
			"INSERT INTO users(username, password_hash) VALUES(?, ?)",
			req.Username, string(hashedPassword),
		)

		if err != nil {
			var statusCode int
			response := models.NewErrorResponse("Internal server error")
			statusCode = http.StatusInternalServerError

			// Check for unique constraint error
			if strings.Contains(err.Error(), "UNIQUE constraint failed: users.username") {
				response = models.NewErrorResponse("Username already in use")
				statusCode = http.StatusConflict
			} else {
				log.Printf("Error inserting user %s: %v", req.Username, err)
			}

			w.WriteHeader(statusCode)
			json.NewEncoder(w).Encode(response)
			return
		}

		userID, err := result.LastInsertId()
		if err != nil {
			log.Printf("Error getting last insert ID after registering user %s: %v", req.Username, err)
			response := models.NewErrorResponse("Registration partially successful, but failed to retrieve user ID")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
			return
		}

		log.Printf("User '%s' (ID: %d) registered successfully.", req.Username, userID)
		registerData := models.RegisterSuccessData{
			UserID:   userID,
			Username: req.Username,
		}

		response := models.NewSuccessResponse(registerData)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	}
}
