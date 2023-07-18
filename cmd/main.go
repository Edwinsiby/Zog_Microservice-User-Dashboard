package main

import (
	"log"
	"net"
	"net/http"
	"service3/pb"
	"service3/pkg/service"

	"google.golang.org/grpc"
)

func main() {
	grpcServer := grpc.NewServer()

	UserService := &service.UserDashboard{}

	pb.RegisterUserDashboardServer(grpcServer, UserService)

	listener, err := net.Listen("tcp", ":5052")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	log.Println("User dashboard service is running...")
	go grpcServer.Serve(listener)

	if err := http.ListenAndServe(":8083", nil); err != nil {
		log.Fatalf("Failed to start health check server: %v", err)
	}
}
