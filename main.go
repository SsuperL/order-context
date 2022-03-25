package main

import (
	"log"
	"net"
	"order-context/ohs/local/pl"
	"order-context/ohs/remote"
	"order-context/ohs/remote/resources"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	var addr = "localhost:50051"
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// init tables
	// dbInstance, err := db.ConnectDB()
	// if err != nil {
	// 	log.Fatalf("failed to connect to database: %v", err)
	// 	os.Exit(1)
	// }

	// defer db.DisconnectDB()

	// db.InitTables(dbInstance)

	serverOptions := []grpc.ServerOption{grpc.UnaryInterceptor(remote.UnaryInterceptor)}
	s := grpc.NewServer(serverOptions...)
	orderServer := resources.NewOrderResource()
	invoiceServer := resources.NewInvoiceResource()
	pl.RegisterOrderServiceServer(s, orderServer)
	pl.RegisterInvoiceServiceServer(s, invoiceServer)
	reflection.Register(s)
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
