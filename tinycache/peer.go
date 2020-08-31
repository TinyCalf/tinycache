package tinycache

import (
	"fmt"
	"log"
	"github.com/valyala/gorpc"
)

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

var peers = make(map[string]peer)

//以后可以实现 gorpc mqtt http等多种形式的peer
type peer interface {
	getName() string
	getAddr() string
	get(collectionName, key string) (string, error)
}

type gorpcPeer struct{
	name string
	addr string
	cli *gorpc.Client
}

var _ peer = (*gorpcPeer)(nil)

func (p gorpcPeer) getName() string{
	return p.name;
}

func (p gorpcPeer) getAddr() string {
	return p.addr;
}

type funcGetParams struct {
	CollectionName string
	Key string
}

func (p *gorpcPeer) get(collectionName, key string) (string, error) {
	params := &funcGetParams{collectionName, key}
	res, err := dispatcher.NewFuncClient(p.cli).Call("Get", params)
	log.Println("func get:",res, err)
	return res.(string), err
}

// RegistPeer 注册存在的节点
func RegistPeer(name string, addr string) {
	cli := gorpc.NewTCPClient(addr)
	cli.Start()
	//有个问题，Start会不会掉，或者选择每次访问的时候Start()
	p := &gorpcPeer{ name, addr, cli}
	peers[name] = peer(p)

	//test
	value, _ := p.get("scores", "Tom")
	log.Println("rpc response:", value)
}

//ServeRPC 启动RPC服务
func ServeRPC(addr string) {	
	s := &gorpc.Server{
		Addr: addr,
		Handler: dispatcher.NewHandlerFunc(),
	}
	if err := s.Serve(); err != nil {
		log.Fatalf("Cannot start rpc server: %s", err)
	}	
}

