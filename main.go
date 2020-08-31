package main

/*
$ curl "http://localhost:9999/api?key=Tom"
630
$ curl "http://localhost:9999/api?key=kkk"
kkk not exist
*/

import (
	//"flag"
	// "fmt"
	// "log"
	// "net/http"

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

func main() {
	var loader tinycache.Loader = func(key string) (string, error) {
		return dbData[key], nil
	}
	tinycache.CreateCollection("scores", 1024, loader)
	go tinycache.ServeRPC("localhost:8081")
	tinycache.RegistPeer(dbName, "localhost:8081")
}
