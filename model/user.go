package model

import (
	"database/sql"
	"time"
)

type User struct {
	ID        int        `json:"id" sql:"id"`
	Username  string     `json:"username" validate:"required" sql:"username"`
	Password  string     `json:"password,omitempty" sql:"password"`
	CreatedAt *time.Time `json:"createdat,omitempty" sql:"createdat"`
	UpdatedAt *time.Time `json:"updatedat,omitempty" sql:"updatedat"`
}

func GetUsers(db *sql.DB, start, count int) ([]User, error) {
	rows, err := db.Query(
		"SELECT id, username FROM users LIMIT $1 OFFSET $2",
		count, start)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []User{}

	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Username); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}

func (u *User) GetUserByUserName(db *sql.DB) error {
	return db.QueryRow("SELECT id FROM users WHERE username=$1",
		u.Username).Scan(&u.ID)
}

func (u *User) GetUser(db *sql.DB) error {
	return db.QueryRow("SELECT username, password, createdat, updatedat FROM users WHERE id=$1",
		u.ID).Scan(&u.Username, &u.Password, &u.CreatedAt, &u.UpdatedAt)
}
