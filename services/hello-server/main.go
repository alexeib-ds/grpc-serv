package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	pb "github.com/alexeib-ds/grpc-serv/services/proto"
)

type server struct {
	pb.UnimplementedGreeterServer

	permClient pb.PermissionClient
}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{
		Message: fmt.Sprintf("Hello, %v!", in.GetName()),
	}, nil
}

func validateRequest(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "metadata is missing")
	}

	auth := md["authorization"]
	if len(auth) < 1 {
		return nil, status.Errorf(codes.Unauthenticated, "token is missing")
	}

	s := info.Server.(*server)
	r, err := s.permClient.ValidateToken(ctx, &pb.ValidateTokenRequest{
		Token: auth[0],
	})

	if err != nil {
		return nil, status.Errorf(codes.Internal, "can't validate token")
	}

	if !r.GetIsValid() {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token")
	}

	return handler(ctx, req)
}

func main() {
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}

	permConn, err := grpc.Dial("localhost:8082", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer(grpc.UnaryInterceptor(validateRequest))
	pb.RegisterGreeterServer(s, &server{
		permClient: pb.NewPermissionClient(permConn),
	})

	if err := s.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
