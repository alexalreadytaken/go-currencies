package main

import (
	"log"

	_ "github.com/alexalreadytaken/go-currencies/docs"
	"github.com/alexalreadytaken/go-currencies/internal/controllers"
	"github.com/alexalreadytaken/go-currencies/internal/repos"
	"github.com/alexalreadytaken/go-currencies/internal/services"
	"github.com/alexalreadytaken/go-currencies/internal/utils"
	ginpkg "github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Currencies service api
// @version 1.0
// @BasePath /api
// @schemes http
func main() {
	cnf := utils.LoadConfigFromEnv()
	currRepo, err := repos.NewCurrenciesRepo(cnf)
	if err != nil {
		log.Fatalf("error while starting pg repo = %s", err.Error())
	}
	currController := controllers.NewCurrenciesController(currRepo)
	currClient := services.NewCurrenciesClient(cnf)
	currDaemonWorker, err := services.NewCurrencyDaemonWorker(currClient, currRepo)
	if err != nil {
		log.Fatalf("error while curreny daemon worker = %s", err.Error())
	}
	currDaemonWorker.GetBtcUsdtCourseAndStore()
	currDaemonWorker.GetRubToFiatCourseAndStore()

	daemon := cron.New(cron.WithParser(cron.NewParser(cron.Descriptor)))
	defer daemon.Stop()
	daemon.AddFunc("@every 10s", currDaemonWorker.GetBtcUsdtCourseAndStore)
	daemon.AddFunc("@daily", currDaemonWorker.GetRubToFiatCourseAndStore)
	go daemon.Start()

	gin := ginpkg.Default()
	api := gin.Group("/api")
	controllers.AddCurrenciesRoutes(api, currController)
	api.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler)) // /api/docs/index.html
	gin.Run(":2000")
}
