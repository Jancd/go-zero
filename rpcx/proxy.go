package rpcx

import (
	"context"
	"sync"

	"zero/core/rpc"
	"zero/core/syncx"
	"zero/rpcx/auth"

	"google.golang.org/grpc"
)

type RpcProxy struct {
	backend     string
	clients     map[string]*RpcClient
	options     []rpc.ClientOption
	sharedCalls syncx.SharedCalls
	lock        sync.Mutex
}

func NewRpcProxy(backend string, opts ...rpc.ClientOption) *RpcProxy {
	return &RpcProxy{
		backend:     backend,
		clients:     make(map[string]*RpcClient),
		options:     opts,
		sharedCalls: syncx.NewSharedCalls(),
	}
}

func (p *RpcProxy) TakeConn(ctx context.Context) (*grpc.ClientConn, error) {
	cred := auth.ParseCredential(ctx)
	key := cred.App + "/" + cred.Token
	val, err := p.sharedCalls.Do(key, func() (interface{}, error) {
		p.lock.Lock()
		client, ok := p.clients[key]
		p.lock.Unlock()
		if ok {
			return client, nil
		}

		client, err := NewClient(RpcClientConf{
			Server: p.backend,
			App:    cred.App,
			Token:  cred.Token,
		}, p.options...)
		if err != nil {
			return nil, err
		}

		p.lock.Lock()
		p.clients[key] = client
		p.lock.Unlock()
		return client, nil
	})
	if err != nil {
		return nil, err
	}

	conn, ok := val.(*RpcClient).Next()
	if !ok {
		return nil, grpc.ErrServerStopped
	}

	return conn, nil
}
