package servant

import (
	"github.com/nicholaskh/fae/config"
	"github.com/nicholaskh/fae/servant/gen-go/fun/rpc"
	"github.com/nicholaskh/golib/sampling"
	log "github.com/nicholaskh/log4go"
	"sync/atomic"
	"time"
)

type session struct {
	ctx      *rpc.Context // will stay the same during a session
	profiler *profiler
}

func (this *FunServantImpl) getSession(ctx *rpc.Context) *session {
	const DIGIT_REPLACED_WITH = "?"
	s, present := this.sessions.Get(ctx.Rid)
	if !present {
		atomic.AddInt64(&svtStats.sessionN, 1)
		s = &session{ctx: ctx}
		this.sessions.Set(ctx.Rid, s)

		normalizedReason := this.digitNormalizer.ReplaceAll(
			[]byte(ctx.Reason), []byte(DIGIT_REPLACED_WITH))
		this.phpReasonPercent.Inc(string(normalizedReason), 1)

		log.Debug("new session {uid^%d rid^%s reason^%s}", this.extractUid(ctx),
			ctx.Rid, ctx.Reason)
	}

	return s.(*session)
}

func (this *session) startProfiler() (*profiler, error) {
	if this.profiler == nil {
		if this.ctx.Rid == "" || this.ctx.Reason == "" {
			log.Error("Invalid context: %s", this.ctx.String())
			return nil, ErrInvalidContext
		}

		this.profiler = &profiler{}
		// TODO 某些web server需要100%采样
		this.profiler.on = sampling.SampleRateSatisfied(config.Engine.Servants.ProfilerRate)
		this.profiler.t0 = time.Now()
		this.profiler.t1 = this.profiler.t0
	}

	this.profiler.t1 = time.Now()
	return this.profiler, nil
}
