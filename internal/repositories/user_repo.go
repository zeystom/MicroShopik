package repositories

import (
	"MicroShopik/internal/domain"
	"database/sql"
	"errors"
)

type UserRepository interface {
	Create(user *domain.User) error
	LastLoginUpdate(userID int) error
	GetByEmail(email string) (*domain.User, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (a *userRepository) Create(user *domain.User) error {

	query := `INSERT INTO users (username, email, password, created_at) VALUES ($1, $2, $3, $4) RETURNING id`
	return a.db.QueryRow(query, user.Username, user.Email, user.Password, user.CreatedAt).Scan(&user.ID)

}

func (a *userRepository) LastLoginUpdate(userID int) error {
	query := `UPDATE users SET last_login = NOW() WHERE id = $1`
	_, err := a.db.Exec(query, userID)
	return err
}

func (a *userRepository) GetByEmail(email string) (*domain.User, error) {
	var user domain.User
	query := `SELECT id, username, email, password, created_at FROM users WHERE email = $1`
	err := a.db.QueryRow(query, email).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}
