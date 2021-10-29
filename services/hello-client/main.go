package main

import (
	"context"
	"log"

	pb "github.com/alexeib-ds/grpc-serv/services/proto"
	"google.golang.org/grpc"
)

type credsProvider struct {
	token string
}

func (cr credsProvider) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"Authorization": cr.token,
	}, nil
}

func (cr credsProvider) RequireTransportSecurity() bool {
	return false
}

func main() {
	creds := credsProvider{
		token: "JoinTheDarkSide",
	}
	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure(), grpc.WithBlock(), grpc.WithPerRPCCredentials(creds))
	if err != nil {
		log.Fatal(err)
	}

	client := pb.NewGreeterClient(conn)
	resp, err := client.SayHello(context.Background(), &pb.HelloRequest{
		Name: "Darth Vader",
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Response: %v\n", resp.GetMessage())
}
