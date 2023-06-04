package repository

import "aBet/model"

type UsersRepository interface {
	AddUsers(Users model.Users) (model.Users, error)
	EditUsers(Users model.Users) (model.Users, error)
	DeleteUsersById(Users model.Users) (model.Users, error)
	GetByIDUsers(Id string) ([]model.Users, error)
	LoginUserAccount(userName string, password string) (model.Users, error)
	GetUsersByName(userName string) (model.Users, error)
}
