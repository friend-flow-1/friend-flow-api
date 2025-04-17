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

var EnforcerPG *casbin.Enforcer

func InitEnforcerPostgres(db *gorm.DB) error {
	_, b, _, _ := runtime.Caller(0)
	modelPath := filepath.Join(filepath.Dir(b), "model.conf")

	m, err := model.NewModelFromFile(modelPath)
	if err != nil {
		return err
	}

	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		return err
	}

	e, err := casbin.NewEnforcer(m, adapter)
	if err != nil {
		return err
	}

	if err := e.LoadPolicy(); err != nil {
		return err
	}

	// Seed policies if they don't exist yet (can add check if policies exist)
	err = SeedPolicies(e)
	if err != nil {
		return err
	}

	EnforcerPG = e
	log.Println("✅ Casbin Enforcer (PostgreSQL) initialized")
	return nil
}
