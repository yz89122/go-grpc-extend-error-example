# Go gRPC extend error example

An example code to demonstrate how to extend the error response of the standard gRPC error in Go.

For details, see [Richer error model](https://www.grpc.io/docs/guides/error/#richer-error-model).

## Run

```bash
go run ./main.go
```

### Expected output

```plain
// no err
echo: test message
// regular error
code: FailedPrecondition
message: regular error
len(details): 0
// error with details
code: FailedPrecondition
message: extended error
len(details): 1
err echo: test message
```

### Proto

see [rpc.proto](./api/proto/rpc.proto)

### Service implementation

see [service.go](./service/service.go)

### Client handle

see [main.go](./main.go)
