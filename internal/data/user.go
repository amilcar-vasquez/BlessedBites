// data/users.go
package data

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/amilcar-vasquez/blessed-bites/internal/validator"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type User struct {
	ID        int64  `json:"id"`
	Email     string `json:"email"`
	FullName  string `json:"fullname"`
	PhoneNo   string `json:"phoneNo"`
	Password  string `json:"password_hash"`
	Role      string `json:"role"`
	CreatedAt string `json:"created_at"`
}

type UserModel struct {
	DB *sql.DB
}

// ValidateUser validates the user data
func ValidateUser(v *validator.Validator, user *User) {
	fmt.Println("Validating user data")
	v.Check(validator.NotBlank(user.Email), "email", "Email must be provided")
	v.Check(validator.IsEmail(user.Email), "email", "Email must be a valid email address")

	v.Check(validator.NotBlank(user.FullName), "fullname", "Full name must be provided")
	v.Check(validator.MaxLength(user.FullName, 100), "fullname", "Full name must not exceed 100 characters")
	v.Check(validator.MinLength(user.FullName, 5), "fullname", "Full name must be at least 2 characters long")

	v.Check(validator.NotBlank(user.Password), "password", "Password must be provided")
	v.Check(validator.MinLength(user.Password, 8), "password", "Password must be at least 8 characters long")
	v.Check(validator.MaxLength(user.Password, 100), "password", "Password must not exceed 100 characters")
}

func (u *UserModel) Insert(user *User) error {
	query := `INSERT INTO users (email, full_name, phone_no, password_hash, role)
          		VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return u.DB.QueryRowContext(
		ctx,
		query,
		user.Email,
		user.FullName,
		user.PhoneNo,
		user.Password,
		user.Role, // now expected
	).Scan(&user.ID, &user.CreatedAt)

}

// GetByID retrieves a user by ID
func (u *UserModel) GetByID(id int64) (*User, error) {
	query := `SELECT id, email, full_name, phone_no, password_hash, role, created_at FROM users WHERE id=$1`
	row := u.DB.QueryRow(query, id)

	var user User
	err := row.Scan(&user.ID, &user.Email, &user.FullName, &user.PhoneNo, &user.Password, &user.Role, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update updates a user in the database
func (u *UserModel) Update(user *User) error {
	query := `UPDATE users 
			  SET email = $1, full_name = $2, phone_no = $3, password_hash = $4, role = $5 
			  WHERE id = $6`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := u.DB.ExecContext(
		ctx,
		query,
		user.Email,
		user.FullName,
		user.PhoneNo,
		user.Password,
		user.Role,
		user.ID,
	)
	return err
}

// delete user by id
func (u *UserModel) Delete(id int64) error {
	query := `DELETE FROM users WHERE id=$1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := u.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserModel) GetByEmail(email string) (*User, error) {
	query := `SELECT id, email, full_name, phone_no, password_hash, role, created_at FROM users WHERE email=$1`
	row := u.DB.QueryRow(query, email)

	var user User
	err := row.Scan(&user.ID, &user.Email, &user.FullName, &user.PhoneNo, &user.Password, &user.Role, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetAll retrieves all users from the database
func (u *UserModel) GetAll() ([]*User, error) {
	query := `SELECT id, email, full_name, phone_no, password_hash, role, created_at FROM users`
	rows, err := u.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Email, &user.FullName, &user.PhoneNo, &user.Password, &user.Role, &user.CreatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	return users, nil
}

func ValidateLogin(v *validator.Validator, users *User) {
	v.Check(validator.NotBlank(users.Email), "email", "Email must be provided")
	v.Check(validator.NotBlank(users.Password), "password", "Password must be provided")
}

// validate email
func ValidateEmail(v *validator.Validator, email string) {
	v.Check(validator.NotBlank(email), "email", "Email must be provided")
	v.Check(validator.IsEmail(email), "email", "Email must be a valid email address")
}

// validate password
func ValidatePassword(v *validator.Validator, password string) {
	v.Check(validator.NotBlank(password), "password", "Password must be provided")
	v.Check(validator.MinLength(password, 8), "password", "Password must be at least 8 characters long")
	v.Check(validator.MaxLength(password, 100), "password", "Password must not exceed 100 characters")
}

// password reset token helper function
func GenerateResetToken() (plain string, hash string, expiry time.Time, err error) {
	bytes := make([]byte, 32)
	_, err = rand.Read(bytes)
	if err != nil {
		return "", "", time.Time{}, err
	}

	plain = hex.EncodeToString(bytes)
	hashBytes := sha256.Sum256([]byte(plain))
	hash = hex.EncodeToString(hashBytes[:])
	expiry = time.Now().UTC().Add(1 * time.Hour)
	fmt.Printf("Debug: Now (UTC): %s | Expiry (UTC): %s\n", time.Now().UTC().Format(time.RFC3339), expiry.Format(time.RFC3339))

	return plain, hash, expiry, nil
}

// Save reset token for a user
func (u *UserModel) SetResetToken(email, tokenHash string, expiry time.Time) error {
	query := `UPDATE users SET reset_token_hash = $1, reset_token_expiry = $2 WHERE email = $3`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := u.DB.ExecContext(ctx, query, tokenHash, expiry, email)
	return err
}

// Verify token and get user
func (u *UserModel) GetUserByResetToken(token string) (*User, error) {
	tokenHash := sha256.Sum256([]byte(token))
	hash := hex.EncodeToString(tokenHash[:])

	query := `
		SELECT id, email, full_name, phone_no, password_hash, role, created_at, reset_token_expiry
		FROM users 
		WHERE reset_token_hash = $1
	`
	row := u.DB.QueryRow(query, hash)

	var user User
	var expiry time.Time

	err := row.Scan(&user.ID, &user.Email, &user.FullName, &user.PhoneNo,
		&user.Password, &user.Role, &user.CreatedAt, &expiry)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("reset token not found")
		}
		return nil, fmt.Errorf("error querying user by reset token: %w", err)
	}

	now := time.Now().UTC()
	if now.After(expiry) {
		return nil, fmt.Errorf("reset token expired")
	}

	return &user, nil
}

// Clear reset token
func (u *UserModel) ClearResetToken(userID int64) error {
	query := `UPDATE users SET reset_token_hash = NULL, reset_token_expiry = NULL WHERE id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := u.DB.ExecContext(ctx, query, userID)
	return err
}

// Check if user exists and generate token
func (u *UserModel) InitiatePasswordReset(email string) (string, error) {
	user, err := u.GetByEmail(email)
	if err != nil {
		return "", fmt.Errorf("email not found")
	}

	token, tokenHash, expiry, err := GenerateResetToken()
	if err != nil {
		return "", err
	}

	err = u.SetResetToken(user.Email, tokenHash, expiry)
	if err != nil {
		return "", err
	}

	return token, nil
}

// Finalize password reset
func (u *UserModel) FinalizePasswordReset(token, newPassword string) error {
	// Log the received token for debugging
	fmt.Printf("Debug: Received reset token: %s\n", token)

	user, err := u.GetUserByResetToken(token)
	if err != nil {
		if err.Error() == "reset token expired" {
			return fmt.Errorf("reset token has expired")
		}
		if err.Error() == "reset token not found" {
			return fmt.Errorf("reset token is invalid")
		}
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("error hashing password: %w", err)
	}

	user.Password = string(hashedPassword)
	if err := u.Update(user); err != nil {
		return err
	}

	return u.ClearResetToken(user.ID)
}
