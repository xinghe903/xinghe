package main

import (
	"log"
	"net/http"

	// authclient "xgateway/client/auth"
	authbakclient "xgateway/client/authbak"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

func startGw() {
	gwmux := runtime.NewServeMux()
	// authclient.Register(gwmux)
	authbakclient.Register(gwmux)
	gwServer := &http.Server{
		Addr:    ":8090",
		Handler: gwmux,
	}
	err := gwServer.ListenAndServe()
	if err != nil {
		log.Fatalf("failed to start gateway server: %v", err)
	}
}

func main() {
	startGw()
	// lis, err := net.Listen("tcp", ":8080")
	// if err != nil {
	// 	log.Fatalln("Failed to listen:", err)
	// }

	// // Create a gRPC server object
	// s := grpc.NewServer()
	// // Attach the Greeter service to the server
	// helloworldpb.RegisterGreeterServer(s, &server{})
	// // Serve gRPC server
	// log.Println("Serving gRPC on 0.0.0.0:8080")
	// go func() {
	// 	log.Fatalln(s.Serve(lis))
	// }()

	// // Create a client connection to the gRPC server we just started
	// // This is where the gRPC-Gateway proxies the requests
	// conn, err := grpc.NewClient(
	// 	"0.0.0.0:8080",
	// 	grpc.WithTransportCredentials(insecure.NewCredentials()),
	// )
	// if err != nil {
	// 	log.Fatalln("Failed to dial server:", err)
	// }

	// gwmux := runtime.NewServeMux()
	// // Register Greeter
	// err = helloworldpb.RegisterGreeterHandler(context.Background(), gwmux, conn)
	// if err != nil {
	// 	log.Fatalln("Failed to register gateway:", err)
	// }

	// gwServer := &http.Server{
	// 	Addr:    ":8090",
	// 	Handler: gwmux,
	// }

	// log.Println("Serving gRPC-Gateway on http://0.0.0.0:8090")
	// log.Fatalln(gwServer.ListenAndServe())
}
