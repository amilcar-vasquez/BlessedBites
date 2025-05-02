// data/users.go
package data

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/amilcar-vasquez/blessed-bites/internal/validator"
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
