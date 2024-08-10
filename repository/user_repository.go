package repository

import "go-rest/model"

type IUserRepository interface {
	GeetUserByEmail(user *model.User, email string) error
	CreateUser(user *model.User) error
}
