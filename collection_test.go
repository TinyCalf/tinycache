package tinycache

import (
	"testing"
)

//测试1 连续请求同一个key，一次从loader获取，一次从cache获取
func TestCase1(t *testing.T) {
	var loader Loader = func(key string) (string, error) {
		return dbData[key], nil
	}
	c := CreateCollection("score", 10, loader)
	key := "Tom"
	if res, _ := c.Get(key); res != dbData[key] {
		t.Fatalf("res: %s",res);
	}
	if res, _ := c.Get(key); res != dbData[key]{
		t.Fatalf("res: %s",res);
	}
}

//测试2 缓存容量很低的时候 连续请求同一个key都是从loader获取
func TestCase2(t *testing.T) {
	var loader Loader = func(key string) (string, error) {
		return dbData[key], nil
	}
	c := CreateCollection("score", 1, loader)
	key := "Tom"
	if res, _ := c.Get(key); res != dbData[key] {
		t.Fatalf("res: %s",res);
	}
	if res, _ := c.Get(key); res != dbData[key]{
		t.Fatalf("res: %s",res);
	}
}