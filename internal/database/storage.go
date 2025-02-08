package database

import (
	"URL-Shortener/internal/database/models"
	recordsrepository "URL-Shortener/internal/database/repositories/records"
	usersrepository "URL-Shortener/internal/database/repositories/users"
)

func MustGetInstanseOfDatabase(dataSourceName string) *models.DataBase {
	return &models.DataBase{
		Users:   usersrepository.New(dataSourceName),
		Records: recordsrepository.New(dataSourceName),
	}
}
