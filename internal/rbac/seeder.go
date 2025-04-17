package rbac

import (
	"log"

	"github.com/casbin/casbin/v2"
	"github.com/haxxu/friend-flow-api/internal/modules/user"
)

func SeedPolicies(enforcer *casbin.Enforcer) error {
	// Define roles and permissions
	if ok, _ := enforcer.HasGroupingPolicy(string(user.RoleSuperAdmin), string(user.RoleAdmin)); !ok {
		_, _ = enforcer.AddGroupingPolicy(string(user.RoleSuperAdmin), string(user.RoleAdmin))
	}
	if ok, _ := enforcer.HasGroupingPolicy(string(user.RoleAdmin), string(user.RoleUser)); !ok {
		_, _ = enforcer.AddGroupingPolicy(string(user.RoleAdmin), string(user.RoleUser))
	}

	// Permissions
	if ok, _ := enforcer.HasPolicy(string(user.RoleUser), "/me", "GET"); !ok {
		_, _ = enforcer.AddPolicy(string(user.RoleUser), "/me", "GET") // user can view their profile
	}
	if ok, _ := enforcer.HasPolicy(string(user.RoleAdmin), "/admin", "GET"); !ok {
		_, _ = enforcer.AddPolicy(string(user.RoleAdmin), "/admin", "GET") // admin can view admin panel
	}
	if ok, _ := enforcer.HasPolicy(string(user.RoleSuperAdmin), "/admin", "GET"); !ok {
		_, _ = enforcer.AddPolicy(string(user.RoleSuperAdmin), "/admin", "GET") // super admin can view admin panel
	}

	log.Println("✅ Casbin policies seeded")
	return nil
}
