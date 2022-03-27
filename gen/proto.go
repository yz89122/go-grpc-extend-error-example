package gen

//go:generate env GOBIN=$PWD/../bin go install google.golang.org/protobuf/cmd/protoc-gen-go
//go:generate env GOBIN=$PWD/../bin go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
//go:generate mkdir -p ../proto
//go:generate env GOBIN=$PWD/../bin PATH=$PWD/../bin:$PATH protoc -I ../api/proto --go_out=../proto --go_opt=paths=source_relative --go-grpc_out=../proto --go-grpc_opt=paths=source_relative ../api/proto/rpc.proto
