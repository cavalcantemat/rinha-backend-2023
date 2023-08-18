package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/streadway/amqp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"github.com/grpc-ecosystem/go-grpc-middleware"
)

// Implemente as definições de serviços, mensagens e interceptores aqui

func asyncUnaryServerInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {
	// Implement the logic of the interceptor here
	// ...

	// Call the actual handler to process the request
	resp, err = handler(ctx, req)

	// Additional logic can be placed here after the handler is called
	// ...

	return
}

type personServer struct {
	personCache map[uint32]string
	rabbitConn  *amqp.Connection
	rabbitChan  *amqp.Channel
}

// Implement service methods here
// ...

func main() {
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	rabbitConn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("failed to connect to RabbitMQ: %v", err)
	}
	defer rabbitConn.Close()

	rabbitChan, err := rabbitConn.Channel()
	if err != nil {
		log.Fatalf("failed to open a RabbitMQ channel: %v", err)
	}
	defer rabbitChan.Close()

	creds, err := credentials.NewServerTLSFromFile("server.crt", "server.key")
	if err != nil {
		log.Fatalf("failed to load TLS certificates: %v", err)
	}

	s := grpc.NewServer(
		grpc.Creds(creds), // Use TLS credentials
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			asyncUnaryServerInterceptor,
		)),
	)

	// ... (Registre seus serviços no servidor gRPC aqui)

	fmt.Println("gRPC server started on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
