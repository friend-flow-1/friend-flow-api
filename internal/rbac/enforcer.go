package rbac

import (
	"log"
	"path/filepath"
	"runtime"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	"github.com/gocql/gocql"
	cassandraadapter "github.com/haxxu/friend-flow-api/pkg/cassandraadapter"
)

var Enforcer *casbin.Enforcer

func InitEnforcer(session *gocql.Session) error {
	// Get the path of the current file
	_, b, _, _ := runtime.Caller(0)
	modelPath := filepath.Join(filepath.Dir(b), "model.conf")

	// Load Casbin model from the model file
	m, err := model.NewModelFromFile(modelPath)
	if err != nil {
		return err
	}

	// Initialize the Cassandra adapter
	adapter := cassandraadapter.NewAdapter(session)

	// Create the Casbin enforcer with the model and adapter
	e, err := casbin.NewEnforcer(m, adapter)
	if err != nil {
		return err
	}

	// Load the policy from the database
	if err := e.LoadPolicy(); err != nil {
		return err
	}

	// Set the global Enforcer
	Enforcer = e
	log.Println("✅ Casbin enforcer initialized")

	return nil
}
