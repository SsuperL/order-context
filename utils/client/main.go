package main

import (
	"context"
	"fmt"
	"log"
	"order-context/ohs/local/pl"
	"os"
	"time"

	"github.com/gogo/status"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// func unaryCallWithMetadata(c pl.OrderServiceClient,siteCode string){
// 	md:=metadata.Pairs()
// }
func main() {
	addr := "127.0.0.1:9080"
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
	ctx, cancel := context.WithTimeout(ctx, time.Minute*3)
	defer cancel()

	res, err := c.GetOrderDetail(ctx, &pl.GetOrderDetailRequest{Id: "11"})
	// res, err := c.GetOrderList(ctx, &pl.GetOrderListRequest{SpaceId: "", Status: 19})
	// res, err := c.CreateOrder(ctx, &pl.CreateOrderRequest{Status: 1, SpaceId: "space1"})
	if err != nil {
		s := status.Convert(err)
		fmt.Println(s)
		fmt.Println(s.Details()[0])
		os.Exit(1)
	}
	// fmt.Println(len(res.Data))
	fmt.Println(res)
}
