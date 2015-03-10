package config

import (
	conf "github.com/nicholaskh/jsconf"
	log "github.com/nicholaskh/log4go"
	"time"
)

type ConfigGame struct {
	NamegenLength   int
	LockMaxItems    int
	LockExpires     time.Duration
	RedisServerPool string
	ShardSplit      ConfigGameShardSplit
}

func (this *ConfigGame) LoadConfig(cf *conf.Conf) {
	this.RedisServerPool = cf.String("redis_server_pool", "")
	if this.RedisServerPool == "" {
		panic("empty redis_server_pool in game section")
	}
	this.NamegenLength = cf.Int("namegen_length", 3)
	this.LockMaxItems = cf.Int("lock_max_items", 1<<20)
	this.LockExpires = cf.Duration("lock_expires", time.Second*10)
	section, err := cf.Section("shard_split_strategy")
	if err != nil {
		panic("empty shard_split_strategy")
	}
	this.ShardSplit.loadConfig(section)

	log.Debug("game conf: %+v", *this)
}

type ConfigGameShardSplit struct {
	Kingdom  int
	User     int
	Alliance int // how many alliances per shard
	Chat     int
}

func (this *ConfigGameShardSplit) loadConfig(cf *conf.Conf) {
	this.Kingdom = cf.Int("kingdom", 18000)
	this.User = cf.Int("user", 200000)
	this.Chat = cf.Int("chat", 2000000)
	this.Alliance = cf.Int("alliance", 200000/50) // TODO
}
