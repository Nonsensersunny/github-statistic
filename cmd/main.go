package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github_statistics/internal/config"
	"github_statistics/internal/log"
	"github_statistics/pkg/business"
	"github_statistics/pkg/db/influxdb"
	"github_statistics/pkg/db/sqlite"
	"github_statistics/pkg/rest/v1"
	"strconv"
)

func Init() {
	log.Info("Initializing APP")
	influxdb.InitInfluxClient()
	sqlite.InitSqliteClient(business.MigrateDeveloper)
}

func main() {
	Init()

	gin.SetMode(gin.DebugMode)

	r := gin.Default()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = config.NewServerConfig().Http.AllowOrigin
	corsConfig.AllowCredentials = true
	r.Use(cors.New(corsConfig))

	v1.Router(r)

	r.Run(":" + strconv.Itoa(config.NewServerConfig().Http.Port))
}