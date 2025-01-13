package storage

import (
	"github.com/chas3air/URL-Shortener/DAL/internal/repositories/recordsrepository"
	"github.com/chas3air/URL-Shortener/DAL/internal/repositories/usersrepository"
	"github.com/chas3air/URL-Shortener/DAL/pkg/interfaces"
	"github.com/chas3air/URL-Shortener/DAL/pkg/models"
)

func GetInstanseOfRecordsRepository(dataSourceName string) (interfaces.Repository[models.DbRecord], error) {
	repo := recordsrepository.RecordsRepository{Path: dataSourceName}
	return repo, nil
}

func GetInstanseOfUsersRepository(dataSourceName string) (interfaces.Repository[models.User], error) {
	repo := usersrepository.UsersRepository{Path: dataSourceName}
	return repo, nil
}
