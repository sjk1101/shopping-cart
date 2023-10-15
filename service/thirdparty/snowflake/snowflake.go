package snowflake

import "github.com/bwmarrin/snowflake"

var (
	node *snowflake.Node
)

func NewIDGenerator() *snowflake.Node {
	_node, err := snowflake.NewNode(1)
	if err != nil {
		panic(err)
	}
	node = _node
	return node
}