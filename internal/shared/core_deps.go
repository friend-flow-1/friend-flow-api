package shared

import (
	"github.com/bwmarrin/snowflake"
	"github.com/casbin/casbin/v2"
	"github.com/haxxu/friend-flow-api/internal/config"
	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"
)

type CoreDeps struct {
	Config        *config.Config
	PostgresDB    *gorm.DB
	MinioClient   *minio.Client
	EnforcerPG    *casbin.Enforcer
	SnowflakeNode *snowflake.Node
}
