package tinycache

import (
	"fmt"
	"log"
	"github.com/valyala/gorpc"

	"./consistenthash"
)

// RegistPeer 注册节点
func RegistPeer(name string, addr string) {
	cli := gorpc.NewTCPClient(addr)
	cli.Start()
	p := &peer{ name, addr, cli}
	globalPeerManager.addPeer(p);
}

//ServeRPC 启动RPC服务
func ServeRPC(name string, addr string) {	
	globalPeerManager.selfName = name
	s := &gorpc.Server{
		Addr: addr,
		Handler: dispatcher.NewHandlerFunc(),
	}
	if err := s.Serve(); err != nil {
		log.Fatalf("Cannot start rpc server: %s", err)
	}	
}

var globalPeerManager = &peerManager {
	"",
	make(map[string]*peer),
	consistenthash.New(5, nil),
}

type peerManager struct {
	selfName string
	peers map[string]*peer
	chMap *consistenthash.Map
}

func (pm *peerManager) addPeer(p *peer) {
	pm.peers[p.name] = p
	pm.chMap.Add(p.name)
}

func (pm *peerManager) selectPeer(key string) (p *peer, isSelf bool) {
	name := pm.chMap.Get(key)
	if name == pm.selfName {
		isSelf = true
	}
	p = pm.peers[name]
	return
}

type peer struct{
	name string
	addr string
	cli *gorpc.Client
}

var dispatcher = gorpc.NewDispatcher()

func init() {
	dispatcher.AddFunc("Get", func (params *funcGetParams) (string, error) {
		collection := GetCollection(params.CollectionName)
		if collection == nil {
			return "", fmt.Errorf("collection not found")
		}
		return collection.Get(params.Key)
	})
}

type funcGetParams struct {
	CollectionName string
	Key string
}

func (p *peer) get(collectionName, key string) (string, error) {
	params := &funcGetParams{collectionName, key}
	res, err := dispatcher.NewFuncClient(p.cli).Call("Get", params)
	return res.(string), err
}







