package user

import "gorm.io/gorm"

type Module struct {
	Repo    *repository
	Service *UserService
	Handler *UserHandler
}

func InitModule(db *gorm.DB) *Module {
	repo := NewUserRepository(db)
	service := NewUserService(repo)
	handler := NewUserHandler(service)

	return &Module{
		Repo:    repo,
		Service: service,
		Handler: handler,
	}
}
