// main.go

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"apistore/handlers"
	"apistore/middleware"
	"apistore/router"
	"apistore/store"

	"github.com/joho/godotenv"
)

func main() {

	if os.Getenv("GO_ENV") != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Println("Aviso: Erro ao carregar arquivo .env. Usando variáveis de ambiente do sistema.")
		}
	}

	jwtKey := []byte(os.Getenv("JWT_SECRET"))
	db, err := store.NewDBConnection()
	if err != nil {
		log.Fatalf("Erro ao conectar DB: %v", err)
	}

	defer db.Close()
	fmt.Println("Conectado ao PostgreSQL!")

	userStore := store.NewPostgresStore(db)
	authHandler := handlers.NewAuthHandler(userStore, jwtKey)

	r := router.NewRouter(authHandler)
	porta := ":8080"

	fmt.Printf("Servidor na porta %s\n", porta)

	log.Fatal(http.ListenAndServe(porta, middleware.CorsMiddleware(r)))
}
