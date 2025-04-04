package main

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

func setupDatabase(dbPath string) (*sql.DB, error) {
	log.Printf("Conectando a la base de datos en: %s", dbPath)
	// Usar formato de conexión correcto para modernc.org/sqlite
	db, err := sql.Open("sqlite", dbPath) // ¡Cambiado de "sqlite3" a "sqlite"!
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	log.Println("Base de datos conectada exitosamente.")
	// Crear la tabla si no existe
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT NOT NULL UNIQUE,
			password_hash TEXT NOT NULL
		);
	`)
	if err != nil {
		return nil, err
	}

	return db, nil
}
