package main

import (
	"config-service/config"
	"config-service/docs"
	_ "config-service/docs" // generated docs
	"config-service/internal/nacos"
	"config-service/migrations"
	"config-service/pb/commonpb"
	"config-service/routes"
	"log"
	"net"
	"os"
	"os/signal"
	"runtime/debug"
	"strconv"
	"syscall"
	"time"

	"config-service/Services"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
	fiberSwagger "github.com/swaggo/fiber-swagger"
	"google.golang.org/grpc"
)

// @title Kla Klouk configuration service
// @version 1.0
// @description Kla Klouk configuration service API documentation
// @Host os.getenv("APP_URL")
// @BasePath /configuration/api
func main() {
	// Load .env file

	err := godotenv.Load("/app/.env")
	if err != nil {
		log.Println("⚠️ Could not load .env")
	}
	config.ErrorLog.Printf(
		"info | %v",
		"Starting Kla Klouk Admin Service...",
	)
	app := fiber.New(fiber.Config{
		ErrorHandler: config.ErrorHandler,
	})

	//log the panic message
	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
		StackTraceHandler: func(c *fiber.Ctx, e interface{}) {

			stack := debug.Stack()

			config.ErrorLog.Printf(`
				[PANIC]
				Method: %s
				URL: %s
				Error: %v
				Stack Trace:
				%s
				`,
				c.Method(),
				c.OriginalURL(),
				e,
				string(stack),
			)
		},
	}))

	app.Use(logger.New(config.AccessLogger()))

	docs.SwaggerInfo.Host = os.Getenv("APP_URL")

	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	if os.Getenv("ENABLE_NACOS") == "true" {
		if err := nacos.Init(); err != nil {
			log.Fatal(err)
		}

		portr, _ := strconv.Atoi(os.Getenv("APP_PORT"))
		// Register service

		err := nacos.RegisterService(
			os.Getenv("APP_NAME"),
			os.Getenv("SERVICE_IP"), // container IP if running in Docker
			uint64(portr),
		)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("✅ Registered service with Nacos: %s\n", os.Getenv("APP_NAME"))
	}
	redisConn := config.NewRedisClient()
	publisher := Services.NewRedisPublisher(redisConn, 300*time.Millisecond)
	Services.GlobalVariablePublisher = publisher

	// config.ConnectRabbitMQ()
	// defer config.RabbitConn.Close()
	// defer config.RabbitChannel.Close()

	config.ConnectRabbitMQ()

	defer func() {
		if config.RabbitChannel != nil {
			config.RabbitChannel.Close()
		}

		if config.RabbitConn != nil {
			config.RabbitConn.Close()
		}
	}()

	// Connect DB
	config.Connect()

	// Initialize gRPC clients
	config.InitGameEngineGrpcClient()

	//run migraion if any
	config.Migration()
	// Auto migrate tables

	// Create super admin
	migrations.Seeders()
	config.InitMinio()

	// Routes
	routes.Setup(app, redisConn)

	log.Println("✅ Routes initialized")

	// Graceful shutdown
	go listenForShutdown(app)

	// Start listening on port 3000
	go startGRPCServer()

	// Start listening on port 3000
	port := ":3000"
	log.Printf("🚀 Server starting on port %s\n", port)
	if err := app.Listen(port); err != nil {
		log.Printf("❌ Server error: %v\n", err)
	}

}

func listenForShutdown(app *fiber.App) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Println("🛑 Shutting down server...")
	nacos.DeregisterService()

	if err := app.Shutdown(); err != nil {
		log.Printf("❌ Server shutdown error: %v\n", err)
	}
	log.Println("✅ Server stopped")
}

func startGRPCServer() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Printf("❌ GRPC Server error: %v\n", err)
	}

	grpcServer := grpc.NewServer()

	adminService := &Services.AdminService{DB: config.DB}
	commonpb.RegisterAdminServiceServer(grpcServer, adminService)

	grpcuserService := &Services.GrpcUserService{DB: config.DB}
	commonpb.RegisterPlatformServiceServer(grpcServer, grpcuserService)

	playTableActiveService := &Services.PlayTableActiveService{DB: config.DB}
	commonpb.RegisterGameTableActiveServiceServer(grpcServer, playTableActiveService)

	grpcPlayService := &Services.GrpcPlayService{DB: config.DB}
	commonpb.RegisterGameServiceServer(grpcServer, grpcPlayService)

	log.Println("gRPC server running on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
