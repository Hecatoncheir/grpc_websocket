# Websocket gateway for gRPC

---

# Protobuf:

````protobuf
syntax = "proto3";
package grpc;

option go_package = "grpc;grpc";
import "google/api/annotations.proto";

message Request {
  string message = 1;
  string details = 2;
}

message Response {
  string message = 1;
  string details = 2;
}

service Service {
  rpc Stream(stream Request) returns (stream Response) {
    option (google.api.http) = {get: "/ws"};
  }
}
````

# Generate rpc stub:

```bash
protoc -I . \
   -I $GOPATH/pkg/mod/github.com/grpc-ecosystem/grpc-gateway/v2@v2.2.0/third_party/googleapis \
   -I $GOPATH/src/github.com/protocolbuffers/protobuf/src \
   --go_out . --go_opt paths=source_relative \
   --go-grpc_out . --go-grpc_opt paths=source_relative \
   socket/grpc/grpc.proto
```

For socket gateway:
```bash
protoc -I . \
   -I $GOPATH/pkg/mod/github.com/grpc-ecosystem/grpc-gateway/v2@v2.2.0/third_party/googleapis \
   -I $GOPATH/src/github.com/protocolbuffers/protobuf/src \
   --grpc-gateway_out . \
   --grpc-gateway_opt logtostderr=true \
   --grpc-gateway_opt paths=source_relative \
   --grpc-gateway_opt generate_unbound_methods=true \
   socket/grpc/grpc.proto
```

# Run local tests:
```
go test ./...
```

# Environment variables

## gRPC: 
 
`GRPCSocketIp` - server ip. Default: `127.0.0.1` <br>
`GRPCSocketPort` - server port. Default: `9001`

`SocketIp` - server ip. Default: `127.0.0.1` <br>
`SocketPort` - server port. Default: `8001`

# Run

```go
// main.go

func main(){
    socketServerConfiguration, err := configuration.New()
    if err != nil {
        return err
    }

    socketServer, err := socket.New(socketServerConfiguration)
    if err != nil {
        return err
    }

    err := socketServer.Run()
    if err != nil {
        log.Fatal(err)
    }

    for client := range socketServer.OnClientAddedChannel() {
    // TODO
        for messages := range client.GetMessagesChan() {
        // TODO
        }
    }
}
````

```bash
go run main.go
```