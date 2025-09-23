
package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"apistore/models"
	"apistore/store"
	"apistore/utils"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	Store  store.UserStore
	JwtKey []byte
}

func NewAuthHandler(store store.UserStore, jwtKey []byte) *AuthHandler {
	return &AuthHandler{Store: store, JwtKey: jwtKey}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if r.Method != http.MethodPost {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Método não permitido")
		return
	}
	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "JSON inválido")
		return
	}
	if req.NomeUsuario == "" || req.Email == "" || req.Senha == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Campos obrigatórios faltando")
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Senha), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Erro ao hash senha: %v", err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Erro ao processar senha")
		return
	}
	newUserID, err := h.Store.CreateUser(ctx, req, string(hashedPassword))
	if err != nil {
		log.Printf("Erro ao criar usuário: %v", err)
		if errors.Is(err, store.ErrUserExists) {
			utils.RespondWithError(w, http.StatusConflict, "Usuário ou email já existe")
			return
		}
		utils.RespondWithError(w, http.StatusInternalServerError, "Erro ao criar usuário")
		return
	}
	utils.RespondWithJSON(w, http.StatusCreated, map[string]interface{}{
		"message": "Usuário criado",
		"user_id": newUserID,
	})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if r.Method != http.MethodPost {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Método não permitido")
		return
	}
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "JSON inválido")
		return
	}
	user, err := h.Store.GetUserByUsername(ctx, req.NomeUsuario)
	if err != nil {
		if errors.Is(err, store.ErrUserNotFound) {
			utils.RespondWithError(w, http.StatusUnauthorized, "Credenciais inválidas")
			return
		}
		log.Printf("Erro ao buscar usuário: %v", err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Erro interno")
		return
	}
	if err = bcrypt.CompareHashAndPassword([]byte(user.SenhaHash), []byte(req.Senha)); err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Credenciais inválidas")
		return
	}
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &models.Claims{
		UserID:        user.ID,
		Username:      user.NomeUsuario,
		NivelAcessoID: user.NivelAcessoID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(h.JwtKey)
	if err != nil {
		log.Printf("Erro ao gerar token: %v", err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Erro ao gerar token")
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"token": tokenString})
}
