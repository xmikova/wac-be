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
    "github.com/xmikova/ambulance-webapi/internal/ambulance_wl"
    "github.com/xmikova/ambulance-webapi/internal/db_service"
    "github.com/xmikova/ambulance-webapi/internal/pharmacy"
)

func main() {
    log.Printf("Server started")
    port := os.Getenv("AMBULANCE_API_PORT")
    if port == "" {
        port = "8080"
    }
    environment := os.Getenv("AMBULANCE_API_ENVIRONMENT")
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
        MaxAge: 12 * time.Hour,
    })
    engine.Use(corsMiddleware)
    // setup context update middleware
    dbService := db_service.NewMongoService[ambulance_wl.Ambulance](db_service.MongoServiceConfig{})
    defer dbService.Disconnect(context.Background())

    pharmacyDbService := db_service.NewMongoService[pharmacy.PharmacyStore](db_service.MongoServiceConfig{
        Collection: "pharmacy",
    })
    defer pharmacyDbService.Disconnect(context.Background())

    engine.Use(func(ctx *gin.Context) {
        ctx.Set("db_service", dbService)
        ctx.Set("pharmacy_db_service", pharmacyDbService)
        ctx.Next()
    })

    // ambulance routes (generated)
    handleFunctions := &ambulance_wl.ApiHandleFunctions{
        AmbulanceConditionsAPI:  ambulance_wl.NewAmbulanceConditionsApi(),
        AmbulanceWaitingListAPI: ambulance_wl.NewAmbulanceWaitingListApi(),
        AmbulancesAPI:           ambulance_wl.NewAmbulancesApi(),
    }
    ambulance_wl.NewRouterWithGinEngine(engine, *handleFunctions)

    // pharmacy routes
    engine.GET("/api/pharmacy/:pharmacyId/medicines", pharmacy.GetMedicines)
    engine.POST("/api/pharmacy/:pharmacyId/medicines", pharmacy.CreateMedicine)
    engine.GET("/api/pharmacy/:pharmacyId/medicines/:medicineId", pharmacy.GetMedicine)
    engine.PUT("/api/pharmacy/:pharmacyId/medicines/:medicineId", pharmacy.UpdateMedicine)
    engine.DELETE("/api/pharmacy/:pharmacyId/medicines/:medicineId", pharmacy.DeleteMedicine)

    engine.GET("/openapi", api.HandleOpenApi)
    engine.Run(":" + port)
}