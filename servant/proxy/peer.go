package proxy

import (
	"github.com/nicholaskh/fae/servant/gen-go/fun/rpc"
	"github.com/nicholaskh/golib/pool"
	log "github.com/nicholaskh/log4go"
	"strconv"
)

// A single rpc client connection with remote peer
// When ex encountered, must Close/Recyle the servant conn from pool
type FunServantPeer struct {
	id uint64
	pool.Resource
	*rpc.FunServantClient

	pool *funServantPeerPool
}

func newFunServantPeer(id uint64, p *funServantPeerPool,
	c *rpc.FunServantClient) *FunServantPeer {
	this := new(FunServantPeer)
	this.FunServantClient = c
	this.pool = p
	this.id = id
	this.Resource = this
	return this
}

func (this *FunServantPeer) Close() {
	log.Debug("peer[%s] conn txn:%d closed", this.pool.peerAddr, this.Id())
	this.Transport.Close()
	this.Resource = nil
}

func (this *FunServantPeer) Id() uint64 {
	return this.id
}

func (this *FunServantPeer) IsOpen() bool {
	return this.Transport.IsOpen()
}

func (this *FunServantPeer) Recycle() {
	if this.Transport.IsOpen() {
		this.pool.pool.Put(this)
	} else {
		this.pool.pool.Kill(this)
		this.pool.pool.Put(nil)
	}
}

func (this *FunServantPeer) NewContext(reason string, uid *int64) *rpc.Context {
	ctx := rpc.NewContext()
	ctx.Rid = strconv.FormatInt(this.pool.nextTxn(), 10)
	ctx.Reason = reason
	ctx.Uid = uid
	ctx.Host = this.pool.myIp

	return ctx
}

// append my transaction id and my host ip to ctx
func (this *FunServantPeer) HijackContext(ctx *rpc.Context) {
	ctx.Host = ctx.Host + ":" + this.pool.myIp
	ctx.Sticky = new(bool)
	*ctx.Sticky = true // tells peer it's from fae
}

func (this *FunServantPeer) Addr() string {
	return this.pool.peerAddr // peers in the pool share the remote peer addr
}
