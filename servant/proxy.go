package servant

import (
	"github.com/nicholaskh/fae/servant/gen-go/fun/rpc"
	"github.com/nicholaskh/fae/servant/proxy"
)

// proxy mode, dispatching(routing) call
func (this *FunServantImpl) peerServantByKey(ctx *rpc.Context, key string) (
	*proxy.FunServantPeer, error) {
	svt, err := this.proxy.ServantByKey(key)
	if err != nil {
		return nil, err

	}

	if svt == nil {
		// should never happen
		return nil, ErrProxyNotFound
	}

	svt.HijackContext(ctx)
	return svt, nil
}

func (this *FunServantImpl) peerServantRand(ctx *rpc.Context) (
	*proxy.FunServantPeer, error) {
	svt, err := this.proxy.RandServant()
	if err != nil {
		return nil, err

	}

	if svt == nil {
		// should never happen
		return nil, ErrProxyNotFound
	}

	svt.HijackContext(ctx)
	return svt, nil
}
