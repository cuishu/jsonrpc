## JSON RPC 服务端，实现了JSON RPC 2.0协议

使用方法:
```
go get -u github.com/cuishu/jsonrpc
```

example:
```go
package main

import (
	"github.com/cuishu/jsonrpc"
)

func abc() string {
	return "I am function abc :)"
}

type Server struct{}

func (s *Server) Abc() string {
	return "I am method Abc :)"
}

func main() {
	server := jsonrpc.NewJRPCServer()
	server.RegistMethod("abc", abc)
	server.RegistMethod("s", Server{})
	server.Run("127.0.0.1:8080")
}
```
