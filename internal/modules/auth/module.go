package auth

import (
	"github.com/casbin/casbin/v2"
	"github.com/haxxu/friend-flow-api/internal/modules/user"
	"gorm.io/gorm"
)

type Module struct {
	Service *AuthService
	Handler *AuthHandler
}

func InitModule(db *gorm.DB, enforcer *casbin.Enforcer, jwtSecret string, userService *user.UserService) *Module {
	authService := NewAuthService(userService, enforcer, jwtSecret)
	authHandler := NewAuthHandler(authService)

	return &Module{
		Service: authService,
		Handler: authHandler,
	}
}
