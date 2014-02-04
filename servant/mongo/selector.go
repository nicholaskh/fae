package mongo

import (
	"fmt"
	"github.com/funkygao/fae/config"
	"net"
)

type ServerSelector interface {
	SetServers(servers ...string) error
	PickServer(key string) (net.Addr, error)
}

func (this *MongodbPool) lookupDbName(shardKey string, shardId int) string {
	n := (shardId / config.Servants.Mongodb.ShardBaseNum) + 1
	return fmt.Sprintf("db%s", n)
}
