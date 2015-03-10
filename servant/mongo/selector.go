package mongo

import (
	"github.com/nicholaskh/fae/config"
)

type ServerSelector interface {
	SetServers(servers map[string]*config.ConfigMongodbServer)
	PickServer(pool string, shardId int) (server *config.ConfigMongodbServer,
		err error)
	ServerList() []*config.ConfigMongodbServer
}
