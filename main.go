package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"syscall"

	"github.com/samber/lo"
	"github.com/yz89122/go-grpc-extend-error-example/proto"
	"github.com/yz89122/go-grpc-extend-error-example/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func main() {
	var ctx, cancel = context.WithCancel(context.Background())
	defer cancel()

	var serverAddress string
	{
		var tcpListener, err = net.ListenTCP("tcp", nil)
		if err != nil {
			panic(err)
		}

		serverAddress = tcpListener.Addr().String()

		var grpcServer = grpc.NewServer()
		var service = service.NewService()
		proto.RegisterExampleServiceServer(grpcServer, service)

		go func() {
			defer cancel()
			if err := grpcServer.Serve(tcpListener); err != nil {
				fmt.Fprintf(os.Stderr, "listen TCP err: %v\n", err)
				os.Exit(int(syscall.SIGHUP))
			}
		}()

		go func() {
			<-ctx.Done()
			grpcServer.GracefulStop()
		}()
	}

	var grpcClientConn *grpc.ClientConn
	{
		var err error
		grpcClientConn, err = grpc.Dial(
			serverAddress,
			grpc.WithInsecure(),
		)
		if err != nil {
			fmt.Fprintf(os.Stderr, "gRPC dial err: %v\n", err)
			os.Exit(int(syscall.SIGHUP))
		}
	}

	var client = proto.NewExampleServiceClient(grpcClientConn)

	// no err
	{
		fmt.Println("// no err")

		var response, err = client.ExampleMethod(ctx, &proto.ExampleMethodRequest{
			Field1: lo.ToPtr("test message"),
		})
		if err != nil {
			panic(err)
		}

		fmt.Println("echo:", response.GetEchoField1())
	}

	// regular error
	{
		fmt.Println("// regular error")

		var _, err = client.ExampleMethod(ctx, &proto.ExampleMethodRequest{
			ErrorType: proto.ExampleMethodRequest_ERROR_TYPE_REGULAR.Enum(),
			Field1:    lo.ToPtr("test message"),
		})
		if err == nil {
			panic("expect an error")
		}

		var status, ok = status.FromError(err)
		if !ok {
			panic("not status")
		}

		fmt.Println("code:", status.Code())
		fmt.Println("message:", status.Message())
		fmt.Println("len(details):", len(status.Details()))
	}

	// error with details
	{
		fmt.Println("// error with details")

		var _, err = client.ExampleMethod(ctx, &proto.ExampleMethodRequest{
			ErrorType: proto.ExampleMethodRequest_ERROR_TYPE_EXTENDED.Enum().Enum(),
			Field1:    lo.ToPtr("test message"),
		})
		if err == nil {
			panic("expect an error")
		}

		var status, ok = status.FromError(err)
		if !ok {
			panic("not status")
		}

		fmt.Println("code:", status.Code())
		fmt.Println("message:", status.Message())
		fmt.Println("len(details):", len(status.Details()))

		// handle details
		lo.ForEach(status.Details(), func(detailInterface interface{}, _ int) {
			switch detailInterface.(type) {
			case *proto.ExampleErrorDetail:
				var detail = detailInterface.(*proto.ExampleErrorDetail)
				fmt.Println("err echo:", detail.GetEchoField1())
			default:
				panic("unknown type")
			}
		})
	}
}
