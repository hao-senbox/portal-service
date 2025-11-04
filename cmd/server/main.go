package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"portal/config"
	"portal/internal/bmi"
	"portal/internal/body"
	"portal/internal/drink"
	"portal/internal/ieb"
	"portal/internal/portal"
	"portal/internal/program_planner"
	selectoptions "portal/internal/select_options"
	studypreference "portal/internal/study_preference"
	studyprogram "portal/internal/study_program"
	teacherassign "portal/internal/teacher_assign"
	"portal/internal/term"
	"portal/internal/timer"
	"portal/internal/user"
	"portal/pkg/consul"
	"portal/pkg/uploader"
	"portal/pkg/zap"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
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

	userService := user.NewUserService(consulClient)
	imageService := uploader.NewImageService(consulClient)
	termService := term.NewTermService(consulClient)

	drinkCollection := mongoClient.Database(cfg.MongoDB).Collection("drinks")
	drinkRepository := drink.NewDrinkRepository(drinkCollection)
	drinkService := drink.NewDrinkService(drinkRepository, userService)
	drinkHandler := drink.NewDrinkHandler(drinkService)

	bmiCollection := mongoClient.Database(cfg.MongoDB).Collection("bmis")
	bmiRepository := bmi.NewBMIRepository(bmiCollection)
	bmiService := bmi.NewBMIService(bmiRepository, userService)
	bmiHandler := bmi.NewBMIHandler(bmiService)

	timerCollection := mongoClient.Database(cfg.MongoDB).Collection("timers")
	isTimeCollection := mongoClient.Database(cfg.MongoDB).Collection("is_times")
	timerRepository := timer.NewTimerRepository(timerCollection, isTimeCollection)
	timerService := timer.NewTimerService(timerRepository, userService, imageService)
	timerHandler := timer.NewTimerHandler(timerService)

	bodyCollection := mongoClient.Database(cfg.MongoDB).Collection("bodies")
	bodyRepository := body.NewBodyRepository(bodyCollection)
	bodyService := body.NewBodyService(bodyRepository, userService)
	bodyHandler := body.NewBodyHandler(bodyService)

	portalCollection := mongoClient.Database(cfg.MongoDB).Collection("portals")
	portalRepository := portal.NewPortalRepository(portalCollection)
	portalService := portal.NewPortalService(portalRepository)
	portalHandler := portal.NewPortalHandlers(portalService)

	iebCollection := mongoClient.Database(cfg.MongoDB).Collection("iebs")
	iebRepository := ieb.NewIEBRepository(iebCollection)
	iebService := ieb.NewIEBService(iebRepository)
	iebHandler := ieb.NewIEBHandler(iebService)

	programPlannerCollection := mongoClient.Database(cfg.MongoDB).Collection("program_planners")
	programPlannerRepository := program_planner.NewProgramPlanerRepository(programPlannerCollection)
	programPlannerService := program_planner.NewProgramPlanerService(programPlannerRepository)
	programPlannerHandler := program_planner.NewProgramPlanerHandler(programPlannerService)

	teacherAssignmentCollection := mongoClient.Database(cfg.MongoDB).Collection("teacher_assignments")
	teacherAssignmentRepository := teacherassign.NewTeacherAssignmentRepository(teacherAssignmentCollection)
	teacherAssignmentService := teacherassign.NewTeacherAssignmentService(teacherAssignmentRepository)
	teacherAssignmentHandler := teacherassign.NewTeacherAssignmentHandler(teacherAssignmentService)

	studyProgramCollection := mongoClient.Database(cfg.MongoDB).Collection("study_programs")
	studyProgramRepository := studyprogram.NewStudyProgramRepository(studyProgramCollection)
	studyProgramService := studyprogram.NewStudyProgramService(studyProgramRepository)
	studyProgramHandler := studyprogram.NewStudyProgramHandler(studyProgramService)

	selectOptionsCollection := mongoClient.Database(cfg.MongoDB).Collection("select_options")
	selectOptionsRepository := selectoptions.NewSelectOptionsRepository(selectOptionsCollection)
	selectOptionsService := selectoptions.NewSelectOptionsService(selectOptionsRepository, termService)
	selectOptionsHandler := selectoptions.NewSelectOptionsHandler(selectOptionsService)

	studyPreferenceCollection := mongoClient.Database(cfg.MongoDB).Collection("study_preferences")
	studyPreferenceRepository := studypreference.NewStudyPreferenceRepository(studyPreferenceCollection)
	studyPreferenceService := studypreference.NewStudyPreferenceService(studyPreferenceRepository, termService)
	studyPreferenceHandler := studypreference.NewStudyPreferenceHandler(studyPreferenceService)

	router := gin.Default()

	drink.RegisterRoutes(router, drinkHandler)
	bmi.RegisterRoutes(router, bmiHandler)
	timer.RegisterRoutes(router, timerHandler)
	body.RegisterRoutes(router, bodyHandler)
	portal.RegisterRoutes(router, portalHandler)
	ieb.RegisterRouters(router, iebHandler)
	program_planner.RegisterRoutes(router, programPlannerHandler)
	teacherassign.RegisterRoutes(router, teacherAssignmentHandler)
	studyprogram.RegisterRoutes(router, studyProgramHandler)
	selectoptions.RegisterRoutes(router, selectOptionsHandler)
	studypreference.RegisterRoutes(router, studyPreferenceHandler)

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
