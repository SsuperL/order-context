package main

import (
	"fmt"
	"order-context/ohs/remote/controllers"
	"order-context/utils/common"

	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
)

// func main() {
// 	config := common.LoadConfig()
// 	// addr := fmt.Sprintf(":%d", config.Port)
// 	// lis, err := net.Listen("tcp", addr)
// 	// if err != nil {
// 	// 	log.Fatalf("failed to listen: %v", err)
// 	// }

// 	serverOptions := []grpc.ServerOption{grpc.UnaryInterceptor(remote.UnaryInterceptor)}
// 	s := grpc.NewServer(serverOptions...)
// 	orderServer := resources.NewOrderResource()
// 	invoiceServer := resources.NewInvoiceResource()
// 	pl.RegisterOrderServiceServer(s, orderServer)
// 	pl.RegisterInvoiceServiceServer(s, invoiceServer)
// 	reflection.Register(s)

// 	// c := make(chan os.Signal, 1)
// 	// signal.Notify(c, os.Interrupt)
// 	// go func() {
// 	// 	for range c {
// 	// 		log.Println("shutting down gRPC server...")
// 	// 		s.GracefulStop()
// 	// 		// <-ctx.Done()
// 	// 	}
// 	// }()
// 	// log.Println(fmt.Sprintf("server listening at %v", lis.Addr()))
// 	// go func() {
// 	// 	if err := s.Serve(lis); err != nil {
// 	// 		log.Fatalf("failed to serve: %v", err)
// 	// 	}
// 	// }()

// 	dialOptions := []grpc.DialOption{grpc.WithInsecure()}
// 	// start http server
// 	conn, err := grpc.DialContext(
// 		context.Background(),
// 		fmt.Sprintf("localhost:%d", config.Port),
// 		grpc.WithInsecure(),
// 	)
// 	if err != nil {
// 		log.Fatalln("Failed to dial server:", err)
// 	}

// 	gwmux := runtime.NewServeMux(
// 		runtime.WithIncomingHeaderMatcher(customMatcher),
// 	)
// 	// register gateway
// 	err = pl.RegisterOrderServiceHandlerFromEndpoint(context.Background(), gwmux, fmt.Sprintf("localhost:%d", config.Port), dialOptions)
// 	if err != nil {
// 		log.Fatalln("Failed to register gateway:", err)
// 	}
// 	err = pl.RegisterInvoiceServiceHandler(context.Background(), gwmux, conn)
// 	if err != nil {
// 		log.Fatalln("Failed to register gateway:", err)
// 	}

// 	// gwServer := &http.Server{
// 	// 	Addr:    fmt.Sprintf(":%d", config.HTTPPort),
// 	// 	Handler: gwmux,
// 	// }

// 	// log.Println("Serving gRPC-Gateway on http://0.0.0.0" + fmt.Sprintf(":%d", config.Port))
// 	log.Printf("Listening at :%d", config.Port)
// 	err = http.ListenAndServe(fmt.Sprintf(":%d", config.Port), h2c.NewHandler(
// 		httpGrpcRouter(s, gwmux),
// 		&http2.Server{}))
// 	// http.ListenAndServe(fmt.Sprintf(":%d", config.Port), httpGrpcRouter(s, gwmux))
// }

// func httpGrpcRouter(grpcServer *grpc.Server, httpHandler http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
// 			grpcServer.ServeHTTP(w, r)
// 		} else {
// 			httpHandler.ServeHTTP(w, r)
// 		}
// 	})
// }

// func customMatcher(key string) (string, bool) {
// 	switch strings.ToLower(key) {
// 	case "site-code":
// 		return key, true
// 	default:
// 		return runtime.DefaultHeaderMatcher(key)
// 	}
// }

var logger = common.NewLogger()

func main() {
	config := common.LoadConfig()
	fmt.Println(config)

	router := gin.New()
	router.Use(common.LogMiddleWare(), gin.Recovery())
	orderGroup := router.Group("/api")
	orderGroup.Use(common.SideCodeHanler())
	controllers.OrderGroup(orderGroup)

	// srv := &http.Server{
	// 	Addr:    fmt.Sprintf(":%d", config.Port),
	// 	Handler: router,
	// }
	// srv.ListenAndServe()
	endless.ListenAndServe(fmt.Sprintf(":%d", config.Port), router)
}
