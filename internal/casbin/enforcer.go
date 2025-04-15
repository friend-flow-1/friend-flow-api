package casbin

import (
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/persist"

	"path/filepath"
	"runtime"
)

var Enforcer *casbin.Enforcer

func InitEnforcer(adapter persist.Adapter) error {
	_, b, _, _ := runtime.Caller(0)
	basePath := filepath.Dir(b)
	modelPath := filepath.Join(basePath, "model.conf")
	m, err := model.NewModelFromFile(modelPath)
	if err != nil {
		return err
	}

	e, err := casbin.NewEnforcer(m, adapter)
	if err != nil {
		return err
	}

	e.LoadPolicy()
	Enforcer = e
	return nil
}
