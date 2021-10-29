package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"

	pb "github.com/alexeib-ds/grpc-serv/services/proto"
)

type server struct {
	pb.UnimplementedPermissionServer
}

func validateToken(token string) bool {
	log.Printf("Validating token: %v\n", token)
	return token == "JoinTheDarkSide"
}

func (s *server) ValidateToken(ctx context.Context, in *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	return &pb.ValidateTokenResponse{
		IsValid: validateToken(in.GetToken()),
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":8082")
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()
	pb.RegisterPermissionServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
