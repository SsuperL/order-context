package main

import (
	"context"
	"fmt"
	"log"
	"order-service/ohs/local/pl"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func main() {
	addr := "localhost:50051"
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatal("did not connect: %v", err)
	}

	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("failed to close connection: %v", err)
		}
	}()
	c := pl.NewOrderServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	res, err := c.GetOrderDetail(ctx, &pl.GetOrderDetailRequest{Id: "test"})
	if err != nil {
		s := status.Convert(err)
		// fmt.Printf("%T", s.Code())
		fmt.Println(s.Message())
		os.Exit(1)
	}
	fmt.Println(res)
}
