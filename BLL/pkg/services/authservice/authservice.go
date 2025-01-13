package authservice

import (
	"fmt"

	"github.com/chas3air/URL-Shortener/BLL/pkg/interfaces"
	"github.com/chas3air/URL-Shortener/BLL/pkg/models"
)

type AuthService struct {
	UserRepo interfaces.Repository[models.User]
}

// TODO: нужно переделать репозитории так чтобы они возвращали при добавдении объекта пользователя
func (as *AuthService) SignUp(login, password string) (models.User, error) {
	users, err := as.UserRepo.Get()
	if err != nil {
		return models.User{}, err
	}

	for _, user := range users {
		if user.Login == login && user.Password == password {
			return models.User{}, fmt.Errorf("user alreadu exists")
		}
	}

	err = as.UserRepo.Insert(models.User{
		Login:    login,
		Password: password,
	})
	if err != nil {
		return models.User{}, err
	}

	return models.User{}, nil
}

func (as *AuthService) SignIn(login, password string) (models.User, error) {
	return models.User{}, nil
}
