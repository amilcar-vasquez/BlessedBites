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
