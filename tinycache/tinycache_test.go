package tinycache

import (
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

//测试1 连续请求同一个key，一次从loader获取，一次从cache获取
func TestGroupCase1(t *testing.T) {
	var loader Loader = func(key string) (string, error) {
		return dbData[key], nil
	}
	c := CreateCollection("score", 10, loader)
	key := "Tom"
	if res, lt, _ := c.Get(key); res != dbData[key] || lt != FromLocalLoader {
		t.Fatalf("res: %s, loadType: %d",res, lt );
	}
	if res, lt, _ := c.Get(key); res != dbData[key] || lt != FromLocalCache{
		t.Fatalf("res: %s, loadType: %d",res, lt );
	}
}

//测试2 缓存容量很低的时候 连续请求同一个key都是从loader获取
func TestGroupCase2(t *testing.T) {
	var loader Loader = func(key string) (string, error) {
		return dbData[key], nil
	}
	c := CreateCollection("score", 1, loader)
	key := "Tom"
	if res, lt, _ := c.Get(key); res != dbData[key] || lt != FromLocalLoader {
		t.Fatalf("res: %s, loadType: %d",res, lt );
	}
	if res, lt, _ := c.Get(key); res != dbData[key] || lt != FromLocalLoader{
		t.Fatalf("res: %s, loadType: %d",res, lt );
	}
}