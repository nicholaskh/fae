package servant

import (
	"fmt"
	"github.com/nicholaskh/fae/config"
	"github.com/nicholaskh/fae/servant/gen-go/fun/rpc"
	log "github.com/nicholaskh/log4go"
	"time"
)

// profiler and auditter
type profiler struct {
	on bool
	t0 time.Time // start of each session
	t1 time.Time // start of each call
}

func (this *profiler) do(callName string, ctx *rpc.Context, format string,
	args ...interface{}) {
	elapsed := time.Since(this.t1)
	slow := elapsed > config.Engine.Servants.CallSlowThreshold
	if !(slow || this.on) {
		return
	}

	body := fmt.Sprintf(format, args...)
	if slow {
		svtStats.incCallSlow()

		header := fmt.Sprintf("SLOW=%s/%s Q=%s ",
			elapsed, time.Since(this.t0), callName)
		log.Warn(header + this.truncatedStr(body))
	} else if this.on {
		header := fmt.Sprintf("T=%s/%s Q=%s ",
			elapsed, time.Since(this.t0), callName)
		log.Trace(header + this.truncatedStr(body))
	}

}

func (this *profiler) truncatedStr(val string) string {
	if len(val) < config.Engine.Servants.ProfilerMaxBodySize {
		return val
	}

	return val[:config.Engine.Servants.ProfilerMaxBodySize] + "..."
}
