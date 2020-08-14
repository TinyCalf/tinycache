package main

/*
$ curl "http://localhost:9999/api?key=Tom"
630
$ curl "http://localhost:9999/api?key=kkk"
kkk not exist
*/

import (
	//"flag"
	"fmt"
	"log"
	"net/http"

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

func createGroup(name string) *tinycache.Group {
	var f tinycache.GetterFunc = func(key string) ([]byte, error) {
		log.Println("[SlowDB] search key", key)
		if v, ok := dbData[key]; ok {
			return []byte(v), nil
		}
		return nil, fmt.Errorf("%s not exist", key)
	}
	return tinycache.NewGroup(name, 2<<10, f)
}

func startCacheServer(addr string, addrs []string, tiny *tinycache.Group) {
	peers := tinycache.NewHTTPPool(addr)
	peers.Set(addrs...)
	tiny.RegisterPeers(peers)
	log.Println("tinycache is running at", addr)
	log.Fatal(http.ListenAndServe(addr[7:], peers))
}

func startAPIServer(apiAddr string, group *tinycache.Group) {
	http.Handle("/api", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			key := r.URL.Query().Get("key")
			view, err := group.Get(key)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/octet-stream")
			w.Write(view.ByteSlice())

		}))
	log.Println("api server is running at", apiAddr)
	log.Fatal(http.ListenAndServe(apiAddr[7:], nil))

}

func example1() {
	ngroup := createGroup(dbName)
	ngroup.Get("Tom")
	ngroup.Get("Tom")
}

// func example2() {
// 	var port int
// 	var api bool
// 	flag.IntVar(&port, "port", 8001, "TinyCache server port")
// 	flag.BoolVar(&api, "api", false, "Start a api server?")
// 	flag.Parse()

// 	apiAddr := "http://localhost:9999"
// 	addrMap := map[int]string{
// 		8001: "http://localhost:8001",
// 		8002: "http://localhost:8002",
// 		8003: "http://localhost:8003",
// 	}

// 	var addrs []string
// 	for _, v := range addrMap {
// 		addrs = append(addrs, v)
// 	}

// 	tiny := createGroup("scores")
// 	if api {
// 		go startAPIServer(apiAddr, tiny)
// 	}
// 	startCacheServer(addrMap[port], addrs, tiny)
// }

func main() {
	example1()
}
