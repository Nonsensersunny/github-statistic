package main

import (
	"github.com/gin-gonic/gin"
	"github_statistics/internal/config"
	"github_statistics/internal/log"
	"github_statistics/pkg/db/influxdb"
	"github_statistics/pkg/rest/v1"
	"strconv"
	"github.com/gin-contrib/cors"
)

func Init() {
	log.Info("Initializing APP")
	influxdb.InitInfluxClient()
}

func main() {
	Init()

	r := gin.Default()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = config.NewServerConfig().Http.AllowOrigin
	corsConfig.AllowCredentials = true
	r.Use(cors.New(corsConfig))

	v1.Router(r)

	r.Run(":" + strconv.Itoa(config.NewServerConfig().Http.Port))
}