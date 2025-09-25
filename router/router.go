
package router

import (
	"net/http"

	"apistore/handlers"
)

func NewRouter(authHandler *handlers.AuthHandler) *http.ServeMux {
    mux := http.NewServeMux()
    mux.HandleFunc("/user/register", authHandler.Register)
    mux.HandleFunc("/user/login", authHandler.Login)
    return mux
}
