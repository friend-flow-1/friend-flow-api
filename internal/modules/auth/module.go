package auth

import (
	"github.com/casbin/casbin/v2"
	authsession "github.com/haxxu/friend-flow-api/internal/modules/auth/session"
	authtoken "github.com/haxxu/friend-flow-api/internal/modules/auth/token"
	"github.com/haxxu/friend-flow-api/internal/modules/user"
	"gorm.io/gorm"
)

type Module struct {
	Service            *AuthService
	AuthSessionService *authsession.AuthSessionService
	Handler            *AuthHandler
}

func InitModule(db *gorm.DB, enforcer *casbin.Enforcer, jwtSecret string, userService *user.UserService) *Module {
	authSessionService = authsession.NewAuthSessionService()
	tokenService * authtoken.TokenService
	authService := NewAuthService(userService, authSessionService, tokenService, enforcer, jwtSecret)
	authHandler := NewAuthHandler(authService)

	return &Module{
		Service: authService,
		Handler: authHandler,
	}
}
