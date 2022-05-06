package model

import (
	"database/sql"
	"fmt"
	"time"
)

type User struct {
	ID        int        `json:"id" sql:"id"`
	Username  string     `json:"username" validate:"required" sql:"username"`
	Password  string     `json:"password,omitempty" sql:"password"`
	CreatedAt *time.Time `json:"createdat,omitempty" sql:"createdat"`
	UpdatedAt *time.Time `json:"updatedat,omitempty" sql:"updatedat"`
}

func GetUsers(db *sql.DB, start, count int, sort, direction string) ([]User, error) {
	var query string
	if sort != "" {
		query = fmt.Sprintf("SELECT id, username FROM users ORDER BY %s %s LIMIT %d OFFSET %d", sort, direction, count, start)
	} else {
		query = fmt.Sprintf("SELECT id, username FROM users LIMIT %d OFFSET %d", count, start)
	}
	rows, err := db.Query(query)

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

func (u *User) CreateUser(db *sql.DB) error {
	timestamp := time.Now()
	err := db.QueryRow(
		"INSERT INTO users(username, password, createdat, updatedat) VALUES($1, $2, $3, $4) RETURNING id, username, password, createdat, updatedat", u.Username, u.Password, timestamp, timestamp).Scan(&u.ID, &u.Username, &u.Password, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (u *User) GetUserByUserNameAndPassword(db *sql.DB) error {
	return db.QueryRow("SELECT id, username, password, createdat, updatedat FROM users WHERE username=$1 AND password=$2",
		u.Username, u.Password).Scan(&u.ID, &u.Username, &u.Password, &u.CreatedAt, &u.UpdatedAt)
}

func (u *User) DeleteUser(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM users WHERE id=$1", u.ID)

	return err
}

func (u *User) UpdateUser(db *sql.DB) error {
	timestamp := time.Now()
	_, err :=
		db.Exec("UPDATE users SET username=$1, password=$2, updatedat=$3 WHERE id=$4 RETURNING id, username, password, createdat, updatedat", u.Username, u.Password, timestamp, u.ID)

	return err
}
