package recordsrepository

import (
	"URL-Shortener/internal/database/models"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type RecordsRepository struct {
	Path string
}

func New(dataSourcePath string) *RecordsRepository {
	db, err := sql.Open("sqlite3", dataSourcePath)
	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS records (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		url TEXT NOT NULL,
		alias TEXT NOT NULL
	)`)
	if err != nil {
		log.Fatal("failed to create table records, error:", err)
	}

	return &RecordsRepository{Path: dataSourcePath}
}

func (rep *RecordsRepository) Get() ([]models.DbRecord, error) {
	db, err := sql.Open("sqlite3", rep.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	defer db.Close()

	rows, err := db.Query("SELECT * FROM records")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	records := make([]models.DbRecord, 0, 10)
	var record models.DbRecord

	for rows.Next() {
		if err := rows.Scan(&record.Id, &record.URL, &record.Alias); err != nil {
			log.Println("cannot parse row")
			continue
		}
		records = append(records, record)
	}

	return records, nil
}

func (rep *RecordsRepository) GetById(id int) (models.DbRecord, error) {
	db, err := sql.Open("sqlite3", rep.Path)
	if err != nil {
		return models.DbRecord{}, fmt.Errorf("failed to connect to database: %w", err)
	}

	defer db.Close()

	row := db.QueryRow("SELECT * FROM records WHERE id = $1;", id)
	var record models.DbRecord
	err = row.Scan(&record.Id, &record.URL, &record.Alias)
	if err != nil {
		return models.DbRecord{}, err
	}

	return record, nil
}

func (rep *RecordsRepository) Insert(obj models.DbRecord) (int, error) {
	db, err := sql.Open("sqlite3", rep.Path)
	if err != nil {
		return 0, fmt.Errorf("failed to connect to database: %w", err)
	}

	defer db.Close()

	res, err := db.Exec("INSERT INTO records (url, alias) VALUES ($1, $2)", obj.URL, obj.Alias)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (rep *RecordsRepository) Update(id int, obj models.DbRecord) error {
	db, err := sql.Open("sqlite3", rep.Path)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()

	_, err = db.Exec("UPDATE records SET url = $1, alias = $2 WHERE id = $3", obj.URL, obj.Alias, id)
	if err != nil {
		return err
	}

	return nil
}

func (rep *RecordsRepository) Delete(id int) (models.DbRecord, error) {
	record, err := rep.GetById(id)
	if err != nil {
		return models.DbRecord{}, err
	}

	db, err := sql.Open("sqlite3", rep.Path)
	if err != nil {
		return models.DbRecord{}, fmt.Errorf("failed to connect to database: %w", err)
	}

	defer db.Close()

	_, err = db.Exec("DELETE FROM records WHERE id = $1", id)
	if err != nil {
		return models.DbRecord{}, err
	}

	return record, nil
}
