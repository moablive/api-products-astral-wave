
package router

import (
	"net/http"

	"apistore/handlers"
)

func NewRouter(authHandler *handlers.AuthHandler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/user/register", authHandler.Register)
	mux.HandleFunc("/api/user/login", authHandler.Login)
	return mux
}
