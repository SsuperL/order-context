package main

import (
	"log"
	"net"
	"order-service/ohs/local/pl"
	"order-service/ohs/remote"
	"order-service/ohs/remote/resources"

	"google.golang.org/grpc"
)

func main() {
	var addr = "localhost:50051"
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	serverOptions := []grpc.ServerOption{grpc.UnaryInterceptor(remote.UnaryInterceptor)}
	s := grpc.NewServer(serverOptions...)
	orderServer := resources.NewOrderResource()
	invoiceServer := resources.NewInvoiceResource()
	pl.RegisterOrderServiceServer(s, orderServer)
	pl.RegisterInvoiceServiceServer(s, invoiceServer)
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
