package sys

import (
	"net/rpc"
	"sync"
	"zim/common"
)

type rpcNode struct {
	sp   *sync.Pool
	host string
	port string
}

type rpcPool struct {
	pool map[string]*rpcNode
}

var RpcPool *rpcPool

func NewRpcPool() *rpcPool {
	return &rpcPool{
		pool: make(map[string]*rpcNode, 0),
	}
}

func (r *rpcPool) GetClient(host, port string) (c *rpc.Client) {
	dbn := common.Md5Str(host + port)
	p, ok := r.pool[dbn]
	if ok && p != nil {
		c = p.sp.Get().(*rpc.Client)
		if c != nil {
			return
		}
	}
	n := &sync.Pool{
		New: func() interface{} {
			c, err := rpc.DialHTTP("tcp", host+":"+port)
			if err != nil {
				return nil
			}
			return c
		},
	}
	node := new(rpcNode)
	node.host = host
	node.port = port
	node.sp = n
	r.pool[dbn] = node
	return r.pool[dbn].sp.Get().(*rpc.Client)
}

func (r *rpcPool) PutClient(host, port string, c *rpc.Client) {
	dbn := common.Md5Str(host + port)
	p, ok := r.pool[dbn]
	if ok && p != nil {
		r.pool[dbn].sp.Put(c)
	}
}
