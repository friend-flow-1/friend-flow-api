package user

type UserHandler struct {
	UserService *UserService
}

func NewUserHandler(service *UserService) *UserHandler {
	return &UserHandler{UserService: service}
}
