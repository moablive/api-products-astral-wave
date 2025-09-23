
package models

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	UserID        int    `json:"user_id"`
	Username      string `json:"username"`
	NivelAcessoID int    `json:"nivel_acesso_id"`
	jwt.RegisteredClaims
}

type RegisterRequest struct {
	NomeUsuario   string `json:"nome_usuario"`
	Email         string `json:"email"`
	Senha         string `json:"senha"`
	NivelAcessoID int    `json:"nivel_acesso_id,omitempty"`
}

type LoginRequest struct {
	NomeUsuario string `json:"nome_usuario"`
	Senha       string `json:"senha"`
}

type User struct {
	ID            int
	NomeUsuario   string
	Email         string
	SenhaHash     string
	NivelAcessoID int
}
