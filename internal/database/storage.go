package database

import (
	"URL-Shortener/internal/database/models"
	"URL-Shortener/internal/database/repositories/recordsrepository"
	"URL-Shortener/internal/database/repositories/usersrepository"
)

func MustGetInstanseOfDatabase(dataSourceName string) *models.DataBase {
	return &models.DataBase{
		Users:   usersrepository.UsersRepository{Path: dataSourceName},
		Records: recordsrepository.RecordsRepository{Path: dataSourceName},
	}
}
