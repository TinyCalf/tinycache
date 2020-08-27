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

func main() {
	
}
