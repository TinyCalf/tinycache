package tinycache

import (
	"testing"
	"log"
)

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

func startServer(peername string) {
	var loader Loader = func(key string) (string, error) {
		return dbData[key], nil
	}
	CreateCollection("scores", 1024, loader)

	for name, addr := range servers {
		RegistPeer(name, addr)
	}
	
	ServeRPC(peername, servers[peername])
}

func TestServer0(t *testing.T) {
	log.Println("starting server 0")
	go StartTerminal("localhost:9000")
	startServer("0");
}

func TestServer1(t *testing.T) {
	log.Println("starting server 1")
	startServer("1");
}

func TestServer2(t *testing.T) {
	log.Println("starting server 2")
	startServer("2");
}

