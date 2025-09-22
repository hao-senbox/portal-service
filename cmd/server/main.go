package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"portal/config"
	"portal/internal/bmi"
	"portal/internal/body"
	"portal/internal/drink"
	"portal/internal/portal"
	"portal/internal/timer"
	"portal/internal/user"
	"portal/pkg/consul"
	"portal/pkg/zap"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	cfg := config.LoadConfig()

	logger, err := zap.New(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	consulConn := consul.NewConsulConn(logger, cfg)
	consulClient := consulConn.Connect()
	defer consulConn.Deregister()

	mongoClient, err := connectToMongoDB(cfg.MongoURI)
	if err != nil {
		logger.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	if err := waitPassing(consulClient, "go-main-service", 60*time.Second); err != nil {
		logger.Fatalf("Dependency not ready: %v", err)
	}

	userService := user.NewUserService(consulClient)
	drinkCollection := mongoClient.Database(cfg.MongoDB).Collection("drinks")
	drinkRepository := drink.NewDrinkRepository(drinkCollection)
	drinkService := drink.NewDrinkService(drinkRepository, userService)
	drinkHandler := drink.NewDrinkHandler(drinkService)

	bmiCollection := mongoClient.Database(cfg.MongoDB).Collection("bmis")
	bmiRepository := bmi.NewBMIRepository(bmiCollection)
	bmiService := bmi.NewBMIService(bmiRepository, userService)
	bmiHandler := bmi.NewBMIHandler(bmiService)

	timerCollection := mongoClient.Database(cfg.MongoDB).Collection("timers")
	timerRepository := timer.NewTimerRepository(timerCollection)
	timerService := timer.NewTimerService(timerRepository, userService)
	timerHandler := timer.NewTimerHandler(timerService)

	bodyCollection := mongoClient.Database(cfg.MongoDB).Collection("bodies")
	bodyRepository := body.NewBodyRepository(bodyCollection)
	bodyService := body.NewBodyService(bodyRepository, userService)
	bodyHandler := body.NewBodyHandler(bodyService)

	portalCollection := mongoClient.Database(cfg.MongoDB).Collection("portals")
	portalRepository := portal.NewPortalRepository(portalCollection)
	portalService := portal.NewPortalService(portalRepository)
	portalHandler := portal.NewPortalHandlers(portalService)

	router := gin.Default()

	drink.RegisterRoutes(router, drinkHandler)
	bmi.RegisterRoutes(router, bmiHandler)
	timer.RegisterRoutes(router, timerHandler)
	body.RegisterRoutes(router, bodyHandler)
	portal.RegisterRoutes(router, portalHandler)
	defer func() {
		if err := mongoClient.Disconnect(context.Background()); err != nil {
			logger.Fatal(err)
		}
	}()

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	// Run server in a separate goroutine
	go func() {
		logger.Infof("Server running on port %s", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Error starting server: %v", err)
		}
	}()

	// Set up graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logger.Fatalf("Error shutting down server: %v", err)
	}
	logger.Info("Server stopped")
}

func connectToMongoDB(uri string) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Println("Failed to connect to MongoDB")
		return nil, err
	}

	// Check connection
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		log.Println("Failed to ping to MongoDB")
		return nil, err
	}

	log.Println("Successfully connected to MongoDB")
	return client, nil
}

func waitPassing(cli *consulapi.Client, name string, timeout time.Duration) error {
	dl := time.Now().Add(timeout)
	for time.Now().Before(dl) {
		entries, _, err := cli.Health().Service(name, "", true, nil)
		if err == nil && len(entries) > 0 {
			return nil
		}
		time.Sleep(2 * time.Second)
	}
	return fmt.Errorf("%s not ready in consul", name)
}
