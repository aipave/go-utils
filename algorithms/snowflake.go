package algorithms

import (
    "math/rand"

    "github.com/bwmarrin/snowflake"
    "github.com/sirupsen/logrus"
)

var node *snowflake.Node

type snowflakeMgr struct {
}

var Snowflake snowflakeMgr

func init() {
    var err error
    machineId := rand.Int63() % 1024
    node, err = snowflake.NewNode(machineId)
    if err != nil {
        logrus.Fatalf("init snowflake node error: %v", err)
    }
}

func (snowflakeMgr) NextOptStreamID() int64 {
    return node.Generate().Int64()
}

func (snowflakeMgr) NextRoundID() int64 {
    return node.Generate().Int64()
}

func (snowflakeMgr) NextRpcSeqId() int64 {
    return node.Generate().Int64()
}
