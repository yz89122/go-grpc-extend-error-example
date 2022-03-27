package service

import (
	"context"

	spb "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/yz89122/go-grpc-extend-error-example/proto"
)

func NewService() proto.ExampleServiceServer {
	return &service{}
}

type service struct {
	proto.UnimplementedExampleServiceServer
}

func (s *service) ExampleMethod(ctx context.Context, request *proto.ExampleMethodRequest) (*proto.ExampleMethodResponse, error) {
	switch request.GetErrorType() {
	case proto.ExampleMethodRequest_ERROR_TYPE_REGULAR:
		return nil, status.Error(codes.FailedPrecondition, "regular error")

	case proto.ExampleMethodRequest_ERROR_TYPE_EXTENDED:
		var detail, err = anypb.New(&proto.ExampleErrorDetail{
			EchoField1: request.Field1,
		})
		if err != nil {
			panic(err)
		}

		return nil, status.FromProto(&spb.Status{
			Code:    int32(codes.FailedPrecondition),
			Message: "extended error",
			Details: []*anypb.Any{detail},
		}).Err()
	}

	return &proto.ExampleMethodResponse{
		EchoField1: request.Field1,
	}, nil
}
