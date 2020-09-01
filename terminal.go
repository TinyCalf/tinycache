package tinycache

import (
    "net"
    "log"
	"bufio"
	"io"
	"strings"
)

// StartTerminal 开启终端
// Terminal是一个简单的tcp连接，可用于查看管理缓存状态
func StartTerminal(addr string) {
    listener, err := net.Listen("tcp", addr)
    if err != nil {
        log.Fatalf("listen err: %v\n", err)
    }
    defer listener.Close()
    log.Printf("bind: %s, start listening...\n", addr)

    for {
        conn, err := listener.Accept()
        if err != nil {
			log.Fatalf("accept err: %v\n", err)
        }
        go handle(conn)
    }
}

func handle(conn net.Conn) {
    reader := bufio.NewReader(conn)
    for {
        msg, err := reader.ReadString('\n')
        if err != nil {
            if err == io.EOF {
                log.Println("client disconnected!")
            } else {
                log.Println(err)
            }
            return
		}
		params := strings.Fields(msg)
		var bytes []byte = []byte("nil\n")
		if len(params) == 3 && params[0] == "get" {
			collection := GetCollection(params[1])
			if collection != nil {
				value, err := collection.Get(params[2])
				if err == nil && len([]byte(value)) > 0 {
					bytes = []byte(value + "\n")
				}
			}
		}
        conn.Write(bytes)
    }
}
