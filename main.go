package main

import (
	"flag"
	
	"./tinycache"
)


//示例用源数据
var (
	dbName = "scores"
	dbData = map[string]string{
		"Tom":  "630",
		"Jack": "589",
		"Sam":  "567",
	}
)

var servers = map[string]string{
	"0": "localhost:10000",	
	"1": "localhost:10001",	
	"2": "localhost:10002",	
}

func main() {
	var peername string
	flag.StringVar(&peername, "peername", "0", "the name of the peer")
	flag.Parse()


	var loader tinycache.Loader = func(key string) (string, error) {
		return dbData[key], nil
	}
	tinycache.CreateCollection("scores", 1024, loader)

	for name, addr := range servers {
		tinycache.RegistPeer(name, addr)
	}
	tinycache.ServeRPC(peername, servers[peername])
}
