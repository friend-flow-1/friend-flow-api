package auth

import (
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/haxxu/friend-flow-api/internal/config"
	authsession "github.com/haxxu/friend-flow-api/internal/modules/auth/session"
	authtoken "github.com/haxxu/friend-flow-api/internal/modules/auth/token"
	"github.com/haxxu/friend-flow-api/internal/modules/user"
	"gorm.io/gorm"
)

type Module struct {
	Service            *AuthService
	AuthSessionService *authsession.AuthSessionService
	TokenService       *authtoken.TokenService
	Handler            *AuthHandler
}

func InitModule(db *gorm.DB, enforcer *casbin.Enforcer, cfg *config.Config, userService *user.UserService) *Module {
	authSessionService := authsession.NewAuthSessionService(authsession.NewGormAuthSessionRepository(db))
	tokenService := authtoken.NewTokenService(
		cfg.AccessTokenSecret,
		cfg.RefreshTokenSecret,
		time.Minute*15,  // 15m
		time.Hour*24*30, // 30d
	)
	authService := NewAuthService(userService, authSessionService, tokenService, enforcer, cfg)
	authHandler := NewAuthHandler(authService)

	return &Module{
		Service:            authService,
		AuthSessionService: authSessionService,
		TokenService:       tokenService,
		Handler:            authHandler,
	}
}
