package main

import (
	"amass/internal/handler"
	"amass/internal/infra/database"
	"amass/internal/infra/logger"
	"amass/internal/middlewares"
	"amass/internal/repository"
	"amass/internal/service"
	"amass/internal/utils"
	"context"
	"log"
	"net"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func init() {
	// Force IPv4 to prevent "socket not connected" errors
	net.DefaultResolver = &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: 30 * time.Second,
			}
			// Force TCP4 instead of TCP6
			return d.DialContext(ctx, "tcp4", address)
		},
	}
}

// loadEnv loads and logs environment variables
func loadEnv() (env, version, appName, dbUser, dbPassword, dbHost, dbPort, dbName string, debug bool) {
	env = os.Getenv("ENV")
	version = os.Getenv("VERSION")
	appName = os.Getenv("NAME")

	dbUser = os.Getenv("DB_USER")
	dbPassword = os.Getenv("DB_PASSWORD")
	dbHost = os.Getenv("DB_HOST")
	dbPort = os.Getenv("DB_PORT")
	dbName = os.Getenv("DB_NAME")

	var err error
	debug, err = strconv.ParseBool(os.Getenv("DEBUG"))
	if err != nil {
		debug = false
	}
	return
}

func main() {
	env, version, appName, dbUser, dbPassword, dbHost, dbPort, dbName, debug := loadEnv()

	// Initialize logger
	if err := logger.InitialzeLoggerSystem(); err != nil {
		log.Fatal(err)
	}
	defer logger.Sync()

	logger.Debug(env, version, appName, debug)

	ctx := context.Background()

	// ✅ Database - Connect in goroutine
	var db *pgxpool.Pool
	var wg sync.WaitGroup
	errChan := make(chan error, 1)

	wg.Add(1)
	go func() {
		defer wg.Done()
		logger.Info("starting database connection...")

		var err error
		db, err = database.Connect(ctx, database.Config{
			User:     dbUser,
			Password: dbPassword,
			Host:     dbHost,
			Port:     dbPort,
			Name:     dbName,
		})
		if err != nil {
			errChan <- err
			return
		}
		logger.Info("database connection established")
	}()

	// Wait for database connection to complete
	wg.Wait()
	close(errChan)

	if err := <-errChan; err != nil {
		log.Fatalf("database connection failed: %v", err)
	}

	defer func() {
		logger.Info("closing database connection")
		db.Close()
	}()

	userRepository, err := repository.NewUserRepository(ctx, db)
	if err != nil {
		log.Fatal(err)
	}

	// --- DI ---
	repo := service.AllRepository{
		IUserRepository: userRepository,
	}

	utilsService := utils.NewUtilsService(ctx, debug)
	userService := service.NewUserService(ctx, debug, repo, utilsService)

	svc := service.AllService{
		IUserService:  userService,
		IUtilsService: utilsService,
	}

	userHandler := handler.NewUserHandler(debug, svc)

	// --- Echo ---
	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{
			echo.GET, echo.POST, echo.PUT, echo.PATCH, echo.DELETE,
		},
	}))

	api := e.Group("/api/v1")
	api.GET("/health", userHandler.Health)

	api.POST("/auth/register", userHandler.Register)
	api.POST("/auth/login", userHandler.Login)

	protected := api.Group("/users")
	protected.Use(middlewares.AuthMiddleware())
	protected.GET("/me", userHandler.GetUser)

	// --- Start ---
	e.Logger.Fatal(e.Start(":3001"))
}
