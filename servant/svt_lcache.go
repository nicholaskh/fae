/*
local cache key:string, value:[]byte.
*/
package servant

import (
	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/nicholaskh/fae/servant/gen-go/fun/rpc"
	"github.com/nicholaskh/golib/cache"
	log "github.com/nicholaskh/log4go"
)

func (this *FunServantImpl) onLcLruEvicted(key cache.Key, value interface{}) {
	// Can't use LruCache public api
	// Because that will lead to nested LruCache RWMutex lock, dead lock
	// TODO
	log.Debug("lru[%v] evicted", key)
}

func (this *FunServantImpl) LcSet(ctx *rpc.Context,
	key string, value []byte) (r bool, ex error) {
	const IDENT = "lc.set"

	svtStats.inc(IDENT)

	profiler, err := this.getSession(ctx).startProfiler()
	if err != nil {
		ex = err
		svtStats.incErr()
		return
	}

	this.lc.Set(key, value)
	r = true
	profiler.do(IDENT, ctx,
		"{key^%s val^%s} {r^%v}", key, value, r)

	return
}

func (this *FunServantImpl) LcGet(ctx *rpc.Context, key string) (r []byte,
	miss *rpc.TCacheMissed, ex error) {
	const IDENT = "lc.get"

	svtStats.inc(IDENT)

	profiler, err := this.getSession(ctx).startProfiler()
	if err != nil {
		ex = err
		svtStats.incErr()
		return
	}

	result, ok := this.lc.Get(key)
	if !ok {
		miss = rpc.NewTCacheMissed()
		miss.Message = thrift.StringPtr("lcache missed: " + key) // optional
	} else {
		r = result.([]byte)
	}

	profiler.do(IDENT, ctx,
		"{key^%s} {miss^%v r^%s}", key, miss, r)

	return
}

func (this *FunServantImpl) LcDel(ctx *rpc.Context, key string) (ex error) {
	const IDENT = "lc.del"

	svtStats.inc(IDENT)

	profiler, err := this.getSession(ctx).startProfiler()
	if err != nil {
		ex = err
		svtStats.incErr()
		return
	}

	this.lc.Del(key)
	profiler.do(IDENT, ctx, "{key^%s}", key)
	return
}
