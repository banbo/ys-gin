package id

import "fmt"

type idWorker interface {
	Generate() int64
}

var IdWorker idWorker

func NewIdWorker(node int64) {
	// snowflake实现
	if node < 0 || node > nodeMax {
		panic(fmt.Sprintf("snowflake节点必须在0-%d之间", node))
	}

	snowflakeIns := &snowflake{
		timestamp: 0,
		node:      node,
		step:      0,
	}

	IdWorker = snowflakeIns
}
