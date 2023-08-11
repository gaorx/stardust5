package sdsnowflake

import (
	"github.com/samber/lo"
)

var (
	defaultNode *Node
	zeroNode    *Node
)

func init() {
	defaultNode = lo.Must(NewFromIP())
	zeroNode = lo.Must(New(0))
}

func GenerateIP() int64 {
	return defaultNode.Generate()
}

func GenerateZero() int64 {
	return zeroNode.Generate()
}
