package main

import (
	"auth-service/config"
	"auth-service/controllers"
	"auth-service/database"
	"auth-service/middlewares"
	"flag"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

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

	api := router.Group("/api")
	{
		api.POST("/token", app.GenerateToken)
		api.POST("/user/register", app.RegisterUser)
		secured := api.Group("/secured").Use(middlewares.Auth(app.Config.GetAppConfig().JwtSecret))
		{
			secured.GET("/ping", app.Ping)
		}
	}
	return router
}
