package influxdb

import (
	"fmt"
	"github_statistics/internal/config"
	"github_statistics/internal/log"
	"strconv"
	"time"

	"github.com/influxdata/influxdb/client/v2"
)

func Write(ps ...*client.Point) error {
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database: config.NewServerConfig().InfluxDB.DBName,
	})
	if err != nil {
		return err
	}
	for _, p := range ps {
		bp.AddPoint(p)
	}
	return clt.Write(bp)
}

var clt client.Client

func InitInfluxClient() {
	c := config.NewServerConfig().InfluxDB
	host := c.Host
	port := c.Port
	addr := "http://" + host + ":" + strconv.Itoa(port)
	dbName := c.DBName
	username := c.Username
	password := c.Password
	log.Infof("trying connecting influxdb:%s %s", addr, username)
	var err error
	clt, err = client.NewHTTPClient(client.HTTPConfig{
		Addr:     addr,
		Username: username,
		Password: password,
	})
	if err != nil {
		panic(err)
	}
	_, _, err = clt.Ping(time.Second)
	if err != nil {
		panic(err)
	}
	_, err = clt.Query(client.Query{
		Command: fmt.Sprintf("create database %s", dbName),
	})
}

func NewInfluxClient() (ct client.Client) {
	return clt
}
