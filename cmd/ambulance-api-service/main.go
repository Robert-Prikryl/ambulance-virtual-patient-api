package main

import (
	"log"
	"os"
	"strings"

	"github.com/Robert-Prikryl/ambulance-virtual-patient-api/api"
	"github.com/Robert-Prikryl/ambulance-virtual-patient-api/internal/ambulance_virtual_patient_list"
	"github.com/gin-gonic/gin"

	"context"
	"github.com/Robert-Prikryl/ambulance-virtual-patient-api/internal/db_service"
	"github.com/gin-contrib/cors"
	"time"
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
	// cors middleware
	corsMiddleware := cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "PUT", "POST", "DELETE", "PATCH"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{""},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	})
	engine.Use(corsMiddleware)

	// setup context update  middleware
	dbService := db_service.NewMongoService[ambulance_virtual_patient_list.VirtualPatient](db_service.MongoServiceConfig{})
	defer dbService.Disconnect(context.Background())
	engine.Use(func(ctx *gin.Context) {
		ctx.Set("db_service", dbService)
		ctx.Next()
	})

	// request routings
	handleFunctions := &ambulance_virtual_patient_list.ApiHandleFunctions{
		VirtualPatientListAPI: ambulance_virtual_patient_list.NewPatientApi(),
	}
	ambulance_virtual_patient_list.NewRouterWithGinEngine(engine, *handleFunctions)
	engine.GET("/openapi", api.HandleOpenApi)
	engine.Run(":" + port)
}
