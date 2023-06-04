package registry

import (
	"aBet/adapters/controller"
	uR "aBet/adapters/repository"
	"aBet/usecase/repository"
	uSv "aBet/usecase/service"
)

func (r *registry) NewAuthController() controller.AuthController {
	return controller.NewAuthController(r.NewUsersService())
}

func (r *registry) NewUsersService() uSv.UsersService {
	return uSv.NewUsersService(r.NewUsersRepository())
}
func (r *registry) NewUsersRepository() repository.UsersRepository {
	return uR.NewUsersRepository(r.db)
}
