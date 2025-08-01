package main

import (
	"fmt"
	"log"

	"github.com/edynnt/veloras-api/internal/initialize"
	"github.com/edynnt/veloras-api/pkg/response/msg"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Go Drunk Backend API by DDD
// @version 1.0
// @description This is a server for a Go Drunk Backend API, demonstrating DDD principles.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8800
// @BasePath /v1/2025

// @externalDocs.description OpenAPI
// @externalDocs.url https://swagger.io/resources/open-api/
func main() {
	r, port := initialize.Run()

	fmt.Println("Server is running on port " + port)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("%s: %v", msg.FailedToStartServer, err)
	}
}
