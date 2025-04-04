package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"myapp/models"

	"golang.org/x/crypto/bcrypt"
)

// PostLoginHandler maneja los intentos de login de usuarios.
func PostLoginHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Printf("Error decoding login request: %v", err)
			response := models.NewErrorResponse("Invalid request body")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
			return
		}

		if req.Username == "" || req.Password == "" {
			response := models.NewErrorResponse("Username and password are required")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
			return
		}

		var storedHash string
		var userID int64
		err := db.QueryRowContext(r.Context(),
			"SELECT id, password_hash FROM users WHERE username = ?",
			req.Username,
		).Scan(&userID, &storedHash)

		if err != nil {
			response := models.NewErrorResponse("Invalid username or password")
			statusCode := http.StatusUnauthorized

			if err != sql.ErrNoRows {
				log.Printf("Error querying user '%s': %v", req.Username, err)
				response = models.NewErrorResponse("Internal server error")
				statusCode = http.StatusInternalServerError
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(statusCode)
			json.NewEncoder(w).Encode(response)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(req.Password))
		if err != nil {
			response := models.NewErrorResponse("Invalid username or password")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(response)
			return
		}

		log.Printf("Login successful for user ID: %d (%s)", userID, req.Username)
		loginData := models.LoginSuccessData{
			UserID:   userID,
			Username: req.Username,
		}

		response := models.NewSuccessResponse(loginData)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}
