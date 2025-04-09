// data/users.go
package data

import (
	"database/sql"
)

type User struct {
	ID       int
	Email    string
	FullName string
	PhoneNo  string
	Password string
}

type UserModel struct {
	DB *sql.DB
}

func (u UserModel) Insert(user User) error {
	query := `INSERT INTO users (email, full_name, phone_no, password) VALUES ($1, $2, $3, $4)`
	_, err := u.DB.Exec(query, user.Email, user.FullName, user.PhoneNo, user.Password)
	return err
}

func (u UserModel) GetByEmail(email string) (*User, error) {
	query := `SELECT id, email, full_name, phone_no, password FROM users WHERE email = $1`
	row := u.DB.QueryRow(query, email)

	var user User
	err := row.Scan(&user.ID, &user.Email, &user.FullName, &user.PhoneNo, &user.Password)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func NewUserModel(db *sql.DB) UserModel {
	return UserModel{DB: db}
}
