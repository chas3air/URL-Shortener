package usersrepository

import (
	"URL-Shortener/internal/database/models"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type UsersRepository struct {
	Path string
}

func New(dataSourcePath string) *UsersRepository {
	db, err := sql.Open("sqlite3", dataSourcePath)
	if err != nil {
		log.Fatal("failed to connect to database, error:", err)
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		login TEXT NOT NULL,
		password TEXT NOT NULL
	)`)
	if err != nil {
		log.Fatal("failed to create table, error:", err)
	}

	return &UsersRepository{Path: dataSourcePath}
}

func (u *UsersRepository) Get() ([]models.User, error) {
	db, err := sql.Open("sqlite3", u.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	defer db.Close()

	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]models.User, 0, 10)
	var user models.User

	for rows.Next() {
		if err := rows.Scan(&user.Id, &user.Login, &user.Password); err != nil {
			log.Println("cannot parse row")
			continue
		}
		users = append(users, user)
	}

	return users, nil
}

func (u *UsersRepository) GetById(id int) (models.User, error) {
	db, err := sql.Open("sqlite3", u.Path)
	if err != nil {
		return models.User{}, fmt.Errorf("failed to connect to database: %w", err)
	}

	defer db.Close()

	row := db.QueryRow("SELECT * FROM users WHERE id = $1;", id)
	var user models.User
	err = row.Scan(&user.Id, &user.Login, &user.Password)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (u *UsersRepository) Insert(obj models.User) (int, error) {
	db, err := sql.Open("sqlite3", u.Path)
	if err != nil {
		return 0, fmt.Errorf("failed to connect to database: %w", err)
	}

	defer db.Close()

	res, err := db.Exec("INSERT INTO users (login, password) VALUES ($1, $2)", obj.Login, obj.Password)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (u *UsersRepository) Update(id int, obj models.User) error {
	db, err := sql.Open("sqlite3", u.Path)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	defer db.Close()

	_, err = db.Exec("UPDATE users SET login = $1, password = $2 WHERE id = $3", obj.Login, obj.Password, id)
	if err != nil {
		return err
	}

	return nil
}

func (u *UsersRepository) Delete(id int) (models.User, error) {
	user, err := u.GetById(id)
	if err != nil {
		return models.User{}, err
	}

	db, err := sql.Open("sqlite3", u.Path)
	if err != nil {
		return models.User{}, fmt.Errorf("failed to connect to database: %w", err)
	}

	defer db.Close()

	_, err = db.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}
