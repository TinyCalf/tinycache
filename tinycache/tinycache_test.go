package tinycache

import (
	"fmt"
	"log"
	"testing"
)

var (
	dbName = "scores"
	dbData = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
	}
)

func createGroup() *Group {
	var f GetterFunc = func(key string) ([]byte, error) {
		log.Println("[SlowDB] search key", key)
		if v, ok := dbData[key]; ok {
			return []byte(v), nil
		}
		return nil, fmt.Errorf("%s not exist", key)
	}
	return NewGroup(dbName, 2<<10, f)
}

/**
数据组测试1
缓存为空时通过回调获取源数据
有缓存时直接返回缓存
*/
func TestGroupCase1(t *testing.T) {
	g := createGroup()
	key := "Tom"
	if res, _ := g.Get(key); res.String() != dbData[key] {
		t.Fatalf("expect %s but get %s ", key, res)
	}
	if res, _ := g.Get(key); res.String() != dbData[key] {
		t.Fatalf("expect %s but get %s ", key, res)
	}
}