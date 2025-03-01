package urlrepository

import (
	"URL-Shortener/internal/models"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type URLRepository struct {
	Path string
}

func New(dataSourcePath string) *URLRepository {
	db, err := sql.Open("sqlite3", dataSourcePath)
	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS url (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		url TEXT NOT NULL,
		alias TEXT NOT NULL UNIQUE
	)`)
	if err != nil {
		log.Fatal("failed to create table url, error:", err)
	}

	return &URLRepository{Path: dataSourcePath}
}

func (rep *URLRepository) Get() ([]models.URL, error) {
	db, err := sql.Open("sqlite3", rep.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	defer db.Close()

	rows, err := db.Query("SELECT * FROM url")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	url := make([]models.URL, 0, 10)
	var record models.URL

	for rows.Next() {
		if err := rows.Scan(&record.Id, &record.URL, &record.Alias); err != nil {
			log.Println("cannot parse row")
			continue
		}
		url = append(url, record)
	}

	return url, nil
}

func (rep *URLRepository) GetById(id int) (models.URL, error) {
	db, err := sql.Open("sqlite3", rep.Path)
	if err != nil {
		return models.URL{}, fmt.Errorf("failed to connect to database: %w", err)
	}

	defer db.Close()

	row := db.QueryRow("SELECT * FROM url WHERE id = $1;", id)
	var record models.URL
	err = row.Scan(&record.Id, &record.URL, &record.Alias)
	if err != nil {
		return models.URL{}, err
	}

	return record, nil
}

func (rep *URLRepository) Insert(obj models.URL) (int, error) {
	db, err := sql.Open("sqlite3", rep.Path)
	if err != nil {
		return 0, fmt.Errorf("failed to connect to database: %w", err)
	}

	defer db.Close()

	res, err := db.Exec("INSERT INTO url (url, alias) VALUES ($1, $2)", obj.URL, obj.Alias)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (rep *URLRepository) Delete(alias string) error {
	db, err := sql.Open("sqlite3", rep.Path)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	defer db.Close()

	_, err = db.Exec("DELETE FROM url WHERE alias = $1", alias)
	if err != nil {
		return err
	}

	return nil
}
