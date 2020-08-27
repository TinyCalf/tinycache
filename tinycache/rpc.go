package tinycache

import (
	"fmt"
	"log"
	"github.com/valyala/gorpc"
)

//ServeRPC 启动RPC服务
func ServeRPC(port string) {
	dispatcher := gorpc.NewDispatcher()
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
	
	s := &gorpc.Server{
		Addr: "localhost:" + port,
		Handler: dispatcher.NewHandlerFunc(),
	}

	if err := s.Serve(); err != nil {
		log.Fatalf("Cannot start rpc server: %s", err)
	}	
}





