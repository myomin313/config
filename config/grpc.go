package config

import (
	"context"
	"log"
	"time"

	// "config-service/pb/commonpb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// var GameEngineGrpcClientInstance commonpb.GameEngineServiceClient

func GrpcClient() *grpc.ClientConn {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(
		ctx,
		"kla-klouk-auth-service:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(), // Waits for the connection to be established
	)
	if err != nil {
		log.Fatalf("❌ Failed to connect to GRPC server: %v", err)
	}

	return conn
}

func GameEngineGrpcClient() *grpc.ClientConn {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	log.Println("🔗 Attempting to connect to Game Engine GRPC server at kla-klouk-game-engine-service:50051...")
	conn, err := grpc.DialContext(
		ctx,
		"kla-klouk-game-engine-service:50051", // Corrected port
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		log.Printf("❌ Failed to connect to Game Engine GRPC server: %v", err.Error())
	}

	return conn
}

func InitGameEngineGrpcClient() {
	conn := GameEngineGrpcClient()
	log.Println("✅ Connected to Game Engine GRPC server", conn)
	// GameEngineGrpcClientInstance = commonpb.NewGameEngineServiceClient(conn)
}
