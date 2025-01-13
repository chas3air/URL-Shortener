package interfaces

import "github.com/chas3air/URL-Shortener/BLL/pkg/models"

type AuthService interface {
	SignUp(string, string) (models.User, error)
	SignIn(string, string) (models.User, error)
}

type UsersService interface {
	GetUsers() ([]models.User, error)
	GetUserById(int) (models.User, error)
	InsertUser(models.User) error
	UpdateUser(models.User) error
	DeleteUser(int) (models.User, error)
}

type RecordsService interface {
	GetRecords() ([]models.DbRecord, error)
	GetRecordById(int) (models.DbRecord, error)
	InsertRecord(models.DbRecord) error
	UpdateRecord(models.DbRecord) error
	DeleteRecord(int) (models.DbRecord, error)
}

type Repository[T any] interface {
	Get() ([]T, error)
	GetById(int) (T, error)
	Insert(T) error
	Update(T) error
	Delete(int) (T, error)
}