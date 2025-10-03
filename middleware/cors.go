// middleware/cors.go
package middleware

import (
	"net/http"
)

// CorsMiddleware adiciona os headers de CORS a cada resposta.
func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Define os headers. Em produção, seja mais específico com a origem.
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		// Se for uma requisição OPTIONS (pre-flight), apenas retorne com os headers.
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Caso contrário, prossiga para o próximo handler na cadeia.
		next.ServeHTTP(w, r)
	})
}
