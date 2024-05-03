package main

import (
    "log"
    "os"
    "strings"
    "github.com/gin-gonic/gin"
    "github.com/ImrichDiscantiny/stomatology-webapi/api"
)

func main() {
    log.Printf("Server started")
    port := os.Getenv("STOMATOLOGY_API_PORT")
    if port == "" {
        port = "8080"
    }
    environment := os.Getenv("STOMATOLOGY_API_ENVIRONMENT")
    if !strings.EqualFold(environment, "production") { // case insensitive comparison
        gin.SetMode(gin.DebugMode)
    }
    engine := gin.New()
    engine.Use(gin.Recovery())
    // request routings
    engine.GET("/openapi", api.HandleOpenApi)
    engine.Run(":" + port)
}