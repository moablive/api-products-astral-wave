
package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"apistore/models"
)

var ErrUserNotFound = errors.New("usuário não encontrado")
var ErrUserExists = errors.New("usuário ou email já existe")

type UserStore interface {
	CreateUser(ctx context.Context, req models.RegisterRequest, hashedPassword string) (int, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	UserExists(ctx context.Context, username, email string) (bool, error)
}

type PostgresStore struct {
	DB *sql.DB
}

func NewPostgresStore(db *sql.DB) UserStore {
	return &PostgresStore{DB: db}
}

func (s *PostgresStore) UserExists(ctx context.Context, username, email string) (bool, error) {
	var count int
	err := s.DB.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM public.usuarios WHERE nome_usuario = $1 OR email = $2
	`, username, email).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("erro ao checar usuário: %w", err)
	}
	return count > 0, nil
}

func (s *PostgresStore) CreateUser(ctx context.Context, req models.RegisterRequest, hashedPassword string) (int, error) {
	exists, err := s.UserExists(ctx, req.NomeUsuario, req.Email)
	if err != nil {
		return 0, err
	}
	if exists {
		return 0, ErrUserExists
	}
	nivelAcessoID := req.NivelAcessoID
	if nivelAcessoID == 0 {
		nivelAcessoID = 1 
	}
	var newUserID int
	err = s.DB.QueryRowContext(ctx, `
		INSERT INTO public.usuarios (nome_usuario, email, senha, nivel_acesso_id)
		VALUES ($1, $2, $3, $4) RETURNING id
	`, req.NomeUsuario, req.Email, hashedPassword, nivelAcessoID).Scan(&newUserID)
	if err != nil {
		return 0, fmt.Errorf("erro ao criar usuário: %w", err)
	}
	return newUserID, nil
}

func (s *PostgresStore) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	err := s.DB.QueryRowContext(ctx, `
		SELECT id, nome_usuario, email, senha, nivel_acesso_id
		FROM public.usuarios 
		WHERE nome_usuario = $1
	`, username).Scan(&user.ID, &user.NomeUsuario, &user.Email, &user.SenhaHash, &user.NivelAcessoID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("erro ao buscar usuário: %w", err)
	}
	return &user, nil
}
