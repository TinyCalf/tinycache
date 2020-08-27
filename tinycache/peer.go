package tinycache

import (
	"fmt"
	"github.com/valyala/gorpc"
)

var dispatcher = gorpc.NewDispatcher()

func init() {
	dispatcher.AddFunc("Get", func (collectionName, key string) (string, LoadType, error) {
		collection := GetCollection(collectionName)
		if collection == nil {
			return "", 0, fmt.Errorf("collection not found")
		}
		value, lt, err := collection.Get(key)
		if lt == FromLocalCache {
			lt = FromRemoteCache
		}
		if lt == FromLocalLoader {
			lt = FromRemoteLoader
		}
		return value, lt, err
	})
}

var peers = make(map[string]*peer)

//以后可以实现 gorpc mqtt http等多种形式的peer
type peer interface {
	name() string
	addr() string
	register(name, addr string)
	get(collectionName, key string) (string, LoadType, error)
}

type gorpcPeer struct{
	name string
	addr 
}


func (p *peer) get(collectionName, key string) (string, LoadType, error) {
	resp, err := dispatcher.NewFuncClient(p.rpcCli).Call("Get", collectionName, key)
}

// RegistPeer 注册存在的节点
func RegistPeer(name string, addr string) {
	cli := gorpc.NewTCPClient(addr)
	cli.Start()
	p := &peer{
		name, 
		addr, 
		cli,
	}
	peers[name] = p
}

