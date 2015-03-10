package main

import (
	"github.com/nicholaskh/golib/locking"
	log "github.com/nicholaskh/log4go"
	"os"
)

func cleanup() {
	if options.lockFile != "" {
		locking.UnlockInstance(options.lockFile)
		log.Debug("Cleanup lock %s", options.lockFile)
	}
}

func shutdown() {
	log.Info("unregistering etcd")
	engineRunner.UnregisterEtcd()

	cleanup()

	log.Info("Terminated")

	os.Exit(0)
}
