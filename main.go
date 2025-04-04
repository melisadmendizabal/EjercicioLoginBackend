package main

import (
	"log"
	"net/http"

	"myapp/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	// Conectar a la base de datos
	db, err := setupDatabase("./users.db")
	if err != nil {
		log.Fatal("CRITICAL: No se pudo conectar a la base de datos:", err)
	}
	defer db.Close()

	// Crear router Chi
	r := chi.NewRouter()

	// Middlewares
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(configureCORS())

	// Rutas Públicas
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("API de Login v1.0"))
	})
	r.Post("/register", handlers.PostRegisterHandler(db))
	r.Post("/login", handlers.PostLoginHandler(db))

	// Ruta para obtener datos públicos del usuario
	r.Get("/users/{userID}", func(w http.ResponseWriter, r *http.Request) {
		// Placeholder implementation for getUserHandler
		userID := chi.URLParam(r, "userID")
		w.Write([]byte("User ID: " + userID))
	})

	port := ":3000"
	log.Printf("Servidor escuchando en puerto %s", port)
	log.Fatal(http.ListenAndServe(port, r))
}
