package main

import (
	"context"
	"log"
	"net"
	"os"
	"time"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"jsantdev.com/grpc_sm_api/internals/api/handlers"
	"jsantdev.com/grpc_sm_api/internals/repositories/mongodb"
	pb "jsantdev.com/grpc_sm_api/proto/gen"
)

func init() {

}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Unable to load .env file: %v", err)
	}
	cert := "../../cert/cert.pem"
	key := "../../cert/key.pem"
	createMongoDBClient()
	runGRPCServer(cert, key)

}

func createMongoDBClient() {
	log.Println("Connecting to mongodb...")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	client, err := mongodb.CreateMongoClient(ctx)
	if err != nil {
		log.Fatalln("Unable to connect to mongodb:", err)
	}

	//PR from dev branch
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Fatalln("Unable to disconnect to mongodb:", err)
		}
	}()

	log.Println("Connected to mongodb")
}

func runGRPCServer(certFile, keyFile string) {

	port := os.Getenv("SERVER_PORT")
	listen, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal("Failed to listen:", err)
	}

	creds, err := credentials.NewServerTLSFromFile(certFile, keyFile)
	if err != nil {
		log.Fatalln("Unable to load credentials:", err)
	}

	grpcServer := grpc.NewServer(grpc.Creds(creds))

	pb.RegisterTeacherServiceServer(grpcServer, &handlers.Server{})
	pb.RegisterStudentServiceServer(grpcServer, &handlers.Server{})
	pb.RegisterExecServiceServer(grpcServer, &handlers.Server{})

	log.Println("gRPC server is running on port", port)
	err = grpcServer.Serve(listen)
	if err != nil {
		log.Fatal("Unable to start server:", err)
	}
}
