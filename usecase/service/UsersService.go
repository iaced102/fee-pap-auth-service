package service

import (
	"aBet/library"
	"aBet/model"
	"aBet/usecase/repository"
)

type UsersService interface {
	AddUsers(Users model.Users) (model.Users, error)
	EditUsers(Users model.Users) (model.Users, error)
	DeleteUsers(Users model.Users) (model.Users, error)
	GetByIDUsers(Id string) ([]model.Users, error)
	LoginUserAccount(userName string, password string) (model.Users, error)
	GetUsersByName(userName string) (model.Users, error)
}

type usersService struct {
	usersRepository repository.UsersRepository
}

// DeleteUserById implements UsersService

func (uS *usersService) GetUsersByName(userName string) (model.Users, error) {
	userInfo, e := uS.usersRepository.GetUsersByName(userName)
	return userInfo, e
}
func (uS *usersService) LoginUserAccount(userName string, password string) (model.Users, error) {
	userInfo, e := uS.usersRepository.LoginUserAccount(userName, password)
	return userInfo, e
}

// AddUsers implements UsersService
func (uS *usersService) AddUsers(Users model.Users) (model.Users, error) {
	Users.CreatedAt = library.GetCurrentTimeStamp()
	Users.UpdatedAt = library.GetCurrentTimeStamp()
	userInfo, e := uS.usersRepository.AddUsers(Users)

	return userInfo, e
}

// DeleteUsers implements UsersService
func (uS *usersService) DeleteUsers(Users model.Users) (model.Users, error) {
	userInfo, e := uS.usersRepository.DeleteUsersById(Users)
	return userInfo, e
}

// EditUsers implements UsersService
func (uS *usersService) EditUsers(Users model.Users) (model.Users, error) {
	Users.UpdatedAt = library.GetCurrentTimeStamp()
	userInfo, e := uS.usersRepository.EditUsers(Users)
	return userInfo, e
}

// GetByIDUsers implements UsersService
func (uS *usersService) GetByIDUsers(Id string) ([]model.Users, error) {
	UserInfo, e := uS.usersRepository.GetByIDUsers(Id)
	return UserInfo, e
}

func NewUsersService(us repository.UsersRepository) UsersService {
	return &usersService{
		usersRepository: us,
	}
}
