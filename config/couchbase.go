package config

import (
	conf "github.com/nicholaskh/jsconf"
	log "github.com/nicholaskh/log4go"
)

type ConfigCouchbase struct {
	Servers []string
}

func (this *ConfigCouchbase) LoadConfig(cf *conf.Conf) {
	this.Servers = cf.StringList("servers", nil)
	log.Debug("couchbase conf: %+v", *this)
}

func (this *ConfigCouchbase) Enabled() bool {
	return len(this.Servers) > 0
}
