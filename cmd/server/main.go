package main

import (
	"database/sql"
	"log"
	"net/http"
	"projetpostgre/internal/handler"
	"projetpostgre/internal/repository"
	"projetpostgre/router"

	_ "github.com/lib/pq"
)

func main() {
	dbHost := "db"
	dbPort := "5432"
	dbUser := "postgres"
	dbPassword := "password"
	dbName := "postgres"
	sslmode := "disable"

	connectionString := "postgres://" + dbUser + ":" + dbPassword + "@" + dbHost + ":" + dbPort + "/" + dbName + "?sslmode=" + sslmode

	database, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}
	defer database.Close()

	if err = database.Ping(); err != nil {
		log.Fatalf("Could not ping to the database: %v", err)
	}

	userRepository := repository.NewPostgresUserRepository(database)

	userHandler := handler.NewUserHandler(userRepository)

	r := router.SetupRouter(userHandler)

	log.Println("Starting server on port 8080...")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
