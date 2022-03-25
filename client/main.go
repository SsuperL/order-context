package main

import (
	"context"
	"fmt"
	"log"
	"order-context/ohs/local/pl"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// func unaryCallWithMetadata(c pl.OrderServiceClient,siteCode string){
// 	md:=metadata.Pairs()
// }
func main() {
	addr := "localhost:50051"
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("failed to close connection: %v", err)
		}
	}()
	c := pl.NewOrderServiceClient(conn)
	siteCode := "001001"
	md := metadata.Pairs("site-code", siteCode)
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	// ctx := context.Background()
	// ctx, cancel := context.WithTimeout(ctx, time.Minute*3)
	// defer cancel()

	// res, err := c.GetOrderDetail(ctx, &pl.GetOrderDetailRequest{Id: "11"})
	// res, err := c.GetOrderList(ctx, &pl.GetOrderListRequest{SpaceId: "", Status: 19})
	res, err := c.CreateOrder(ctx, &pl.CreateOrderRequest{Status: 1, SpaceId: "space1"})
	if err != nil {
		s := status.Convert(err)
		fmt.Println(int(s.Code()))
		fmt.Println(s.Message())
		os.Exit(1)
	}
	// fmt.Println(len(res.Data))
	fmt.Println(res)
}
