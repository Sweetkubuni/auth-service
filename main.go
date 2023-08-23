package main

import (
	"auth-service/config"
	"auth-service/controllers"
	"auth-service/database"
	"auth-service/docs"
	"auth-service/middlewares"
	"flag"
	"fmt"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

//	@title			Auth Server
//	@version		1.0
//	@description	Auth service for jwt authentication
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://demothesoftwarepls.com/support
//	@contact.email	mail@demothesoftwarepls.com

//	@host		localhost:9000
//	@BasePath	/api

//	@securityDefinitions.apikey	Authentication
//	@in							header
//	@name						Bearer

//	@externalDocs.description	OpenAPI
//	@externalDocs.url			https://swagger.io/resources/open-api
func main() {
	// Get Config
	envPtr := flag.String("new", "development", "defines enviroment default:production")
	flag.Parse()

	appCfg, err := config.GetConfig(*envPtr)
	if err != nil {
		log.Fatalf("could not get app config: %v", err.Error())
	}

	// Database
	dbHost := appCfg.GetAppConfig().DbHost
	dbPort := appCfg.GetAppConfig().DbPort
	dbName := appCfg.GetAppConfig().DbName
	dbUser := appCfg.GetAppConfig().DbUser
	dbPassword := appCfg.GetAppConfig().DbPassword

	// Initialize Database
	database.Connect(fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPassword, dbHost, dbPort, dbName))
	database.Migrate()
	// Initialize Router
	router := initRouter(appCfg)
	router.Run(":9000")
}

func initRouter(appCfg config.IConfig) *gin.Engine {
	router := gin.Default()
	app := controllers.App{
		Config: appCfg,
	}

	router.Use(cors.New(cors.Config{
		AllowOrigins:     app.Config.GetAppConfig().TrustedOrigins,
		AllowMethods:     []string{"PUT", "PATCH", "GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type", "Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	router.SetTrustedProxies(nil)

	docs.SwaggerInfo.BasePath = "/api"

	api := router.Group("/api")
	{
		api.POST("/token", app.GenerateToken)
		api.POST("/user/register", app.RegisterUser)
		secured := api.Group("/secured").Use(middlewares.Auth(app.Config.GetAppConfig().JwtSecret))
		{
			secured.GET("/ping", app.Ping)
		}
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	fmt.Printf("This %v", app.Config.GetAppConfig().TrustedOrigins)
	return router
}
