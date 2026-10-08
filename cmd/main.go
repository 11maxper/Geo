package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"project/config"
	"project/internal/infrastructure/DB"
	"project/internal/infrastructure/cache"
	"project/internal/infrastructure/cript"
	"project/internal/infrastructure/dadata"
	"project/internal/infrastructure/httpserver"
	"project/internal/infrastructure/jwt"
	"project/internal/infrastructure/logger"
	"project/internal/infrastructure/migrations"
	authController "project/internal/modules/auth/controller"
	authService "project/internal/modules/auth/service"
	authStorage "project/internal/modules/auth/storage"
	geoController "project/internal/modules/geo/controller"
	geoService "project/internal/modules/geo/service"
	geoStorage "project/internal/modules/geo/storage"
	"project/internal/router"
)

// @title Geoservise API
// @version 1.0
// @description Документация API для геосервиса, предоставляющего функции поиска адресов и геокодирования.
// @host localhost:8080
// @BasePath /api/

// @securityDefinitions.apiKey BearerAuth
// @in header
// @name Authorization
// @description Введите токен в формате: **Bearer {your_jwt_token}**

func main() {
	logger := logger.NewLogger()
	ctx := context.Background()

	// Конфиг
	config, err := config.Init(".env")
	if err != nil {
		log.Fatal(err)
	}

	// JWT
	jwt := jwt.NewJwt(config.JWTKey)

	// Xеширование
	hasher := cript.NewBcript()

	// Api dadata
	api := dadata.NewDadata(config.Dadata)

	//Cache
	cacheClient := cache.NewCacheClient(ctx, config.Cache, logger)
	cache := cache.NewCache(cacheClient)

	//DB
	dbx := DB.NewDB(config.DB, logger)
	//Миграции
	migrator := migrations.NewMigration(dbx, logger)
	migrator.Run()

	// Storage
	geoStorage := geoStorage.NewStorage(api, &config.Dadata)
	authStorage := authStorage.NewAuthStorage(dbx)

	// Service
	geoservice := geoService.NewGeoservice(geoStorage, cache, logger)
	authService := authService.NewAuthService(jwt, hasher, authStorage, logger)

	// Controller
	geoController := geoController.NewGeoController(geoservice)
	authController := authController.NewAuthController(authService)

	// Роутер
	r := router.NewRouter(authController, geoController, jwt)

	// Конфигурация и запуск сервера
	server := httpserver.NewServer(r, config.HTTP)
	server.Serve()

	//Shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigChan
	log.Printf("\nReceive signal: %s", sig)

	server.Close()

	log.Println("Server stopped gracefully")

}
