package recordsrepository

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/chas3air/URL-Shortener/DAL/pkg/models"
)

type RecordsRepository struct {
	Path string
}

func (rep RecordsRepository) Get() ([]models.DbRecord, error) {
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

func (rep RecordsRepository) GetById(id int) (models.DbRecord, error) {
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

func (rep RecordsRepository) Insert(obj models.DbRecord) error {
	db, err := sql.Open("sqlite3", rep.Path)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	defer db.Close()

	_, err = db.Exec("INSERT INTO records (url, alias) VALUES ($1, $2)", obj.URL, obj.Alias)
	if err != nil {
		return err
	}

	return nil
}

func (rep RecordsRepository) Update(obj models.DbRecord) error {
	db, err := sql.Open("sqlite3", rep.Path)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()

	_, err = db.Exec("UPDATE records SET url = $1, alias = $2 WHERE id = $3", obj.URL, obj.Alias, obj.Id)
	if err != nil {
		return err
	}

	return nil
}

func (rep RecordsRepository) Delete(id int) (models.DbRecord, error) {
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
