package idgen

import (
	"fmt"

	"github.com/bwmarrin/snowflake"
)

func InitSnowflakeNode(nodeID int64) (*snowflake.Node, error) {
	node, err := snowflake.NewNode(nodeID)
	if err != nil {
		return nil, fmt.Errorf("failed to init snowflake node: %w", err)
	}
	return node, nil
}
