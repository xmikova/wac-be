package main

import (
	"context"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/xmikova/ambulance-webapi/api"
	"github.com/xmikova/ambulance-webapi/internal/db_service"
	"github.com/xmikova/ambulance-webapi/internal/pharmacy"
	"github.com/xmikova/ambulance-webapi/internal/pharmacy_orders"
)

func main() {
	log.Printf("Server started")
	port := os.Getenv("PHARMACY_API_PORT")
	if port == "" {
		port = "8080"
	}
	environment := os.Getenv("PHARMACY_API_ENVIRONMENT")
	if !strings.EqualFold(environment, "production") { // case insensitive comparison
		gin.SetMode(gin.DebugMode)
	}
	engine := gin.New()
	engine.Use(gin.Recovery())
	corsMiddleware := cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "PUT", "POST", "DELETE", "PATCH"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{""},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	})
	engine.Use(corsMiddleware)
	pharmacyDbService := db_service.NewMongoService[pharmacy.PharmacyStore](db_service.MongoServiceConfig{
		Collection: "pharmacy",
	})
	defer pharmacyDbService.Disconnect(context.Background())

	engine.Use(func(ctx *gin.Context) {
		ctx.Set("pharmacy_db_service", pharmacyDbService)
		ctx.Next()
	})

	handleFunctions := &pharmacy_orders.ApiHandleFunctions{
		PharmacyOrdersAPI:      pharmacy_orders.NewPharmacyOrdersApi(),
		PharmacyMedicinesAPI:   pharmacy_orders.NewPharmacyMedicinesApi(),
		PharmacyDispensingsAPI: pharmacy_orders.NewPharmacyDispensingsApi(),
	}
	pharmacy_orders.NewRouterWithGinEngine(engine, *handleFunctions)

	engine.GET("/openapi", api.HandleOpenApi)
	engine.Run(":" + port)
}
