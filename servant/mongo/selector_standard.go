package mongo

import (
	"fmt"
	"github.com/nicholaskh/fae/config"
	"strings"
)

type StandardServerSelector struct {
	shardBaseNum int
	servers      map[string]*config.ConfigMongodbServer // key is pool
}

func NewStandardServerSelector(baseNum int) *StandardServerSelector {
	return &StandardServerSelector{shardBaseNum: baseNum}
}

func (this *StandardServerSelector) PickServer(pool string,
	shardId int) (server *config.ConfigMongodbServer, err error) {
	const SHARD_POOL_PREFIX = "db"

	var bucket string
	if !strings.HasPrefix(pool, SHARD_POOL_PREFIX) {
		bucket = pool
	} else {
		bucket = fmt.Sprintf("db%d", (shardId/this.shardBaseNum)+1)
	}

	var present bool
	server, present = this.servers[bucket]
	if !present {
		err = ErrServerNotFound
	}

	return
}

func (this *StandardServerSelector) SetServers(servers map[string]*config.ConfigMongodbServer) {
	this.servers = servers
}

func (this *StandardServerSelector) ServerList() (servers []*config.ConfigMongodbServer) {
	for _, s := range this.servers {
		servers = append(servers, s)
	}
	return
}
