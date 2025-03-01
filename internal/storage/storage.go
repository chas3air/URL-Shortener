package storage

import "URL-Shortener/internal/models"

type URLRepository interface {
	Get() ([]models.URL, error)
	GetById(int) (models.URL, error)
	Insert(models.URL) (int, error)
	Delete(string) error
}

//go:generate mockgen -destination=mock/mock_repository.go -package=mock . URLRepository
