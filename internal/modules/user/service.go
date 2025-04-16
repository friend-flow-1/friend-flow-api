package user

type UserService struct {
	Repo Repository
}

func NewUserService(repo Repository) *UserService {
	return &UserService{Repo: repo}
}
