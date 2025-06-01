package rbac

import (
	"log"
	"path/filepath"
	"runtime"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

func InitEnforcerPostgres(db *gorm.DB) (*casbin.Enforcer, error) {
	_, b, _, _ := runtime.Caller(0)
	modelPath := filepath.Join(filepath.Dir(b), "model.conf")

	m, err := model.NewModelFromFile(modelPath)
	if err != nil {
		return nil, err
	}

	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		return nil, err
	}

	e, err := casbin.NewEnforcer(m, adapter)
	if err != nil {
		return nil, err
	}

	if err := e.LoadPolicy(); err != nil {
		return nil, err
	}

	if err := SeedPolicies(e); err != nil {
		return nil, err
	}

	log.Println("✅ Casbin Enforcer (PostgreSQL) initialized")
	return e, nil
}
