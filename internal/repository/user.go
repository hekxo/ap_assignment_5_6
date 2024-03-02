package repository

import (
	"adv_programming_3_4-main/internal/model"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type UserRegistration struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// UserRepository is an interface that your handler will use to interact with the database.
type UserRepository interface {
	CreateUser(email, hashedPassword string) (confirmationToken string, err error)
	ConfirmUserEmail(token string) error
	GetUserByEmail(email string) (model.User, error)
	CreatePasswordResetToken(email string) (resetToken string, err error)
	ResetPassword(token, newPassword string) error
}

type SQLUserRepository struct {
	DB *sql.DB
}

func NewSQLUserRepository(db *sql.DB) *SQLUserRepository {
	return &SQLUserRepository{DB: db}
}

func (r *SQLUserRepository) CreateUser(email, hashedPassword string) (confirmationToken string, err error) {
	confirmationToken = uuid.New().String()
	tx, err := r.DB.Begin()
	if err != nil {
		return "", err
	}
	defer func(tx *sql.Tx) {
		err := tx.Rollback()
		if err != nil {

		}
	}(tx)
	_, err = tx.Exec("INSERT INTO users (email, password_hash, confirmation_token) VALUES (?, ?, ?)", email, hashedPassword, confirmationToken)
	if err != nil {
		return "", err
	}
	if err := tx.Commit(); err != nil {
		return "", err
	}
	return confirmationToken, nil
}

func (r *SQLUserRepository) GetUserByEmail(email string) (*model.User, error) {
	user := &model.User{}
	err := r.DB.QueryRow("SELECT id, email, password_hash FROM users WHERE email = $1", email).
		Scan(&user.ID, &user.Email, &user.PasswordHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

func (r *SQLUserRepository) ConfirmUserEmail(token string) error {
	result, err := r.DB.Exec("UPDATE users SET is_email_confirmed = TRUE WHERE confirmation_token = $1", token)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}
func (r *SQLUserRepository) CreatePasswordResetToken(email string) (string, error) {
	resetToken := uuid.New().String()
	_, err := r.DB.Exec("UPDATE users SET reset_token = $1, reset_token_expiry = $2 WHERE email = $3",
		resetToken, time.Now().Add(24*time.Hour), email)
	if err != nil {
		return "", err
	}
	return resetToken, nil
}

func (r *SQLUserRepository) ResetPassword(token, newPassword string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	_, err = r.DB.Exec("UPDATE users SET password_hash = $1 WHERE reset_token = $2 AND reset_token_expiry > $3",
		hashedPassword, token, time.Now())
	if err != nil {
		return err
	}
	return nil
}
