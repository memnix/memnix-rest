package infrastructures

import (
	"context"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/memnix/memnix-rest/config"
	"github.com/memnix/memnix-rest/pkg/env"
)

var (
	influxClient *influxdb2.Client
)

func ConnectInfluxDB(env env.IEnv) error {
	var host string
	var token string
	if config.IsDevelopment() {
		host = env.GetEnv("DEBUG_INFLUXDB_URL")
		token = env.GetEnv("DEBUG_INFLUXDB_TOKEN")
	} else {
		host = env.GetEnv("INFLUXDB_URL")
		token = env.GetEnv("INFLUXDB_TOKEN")
	}
	client := influxdb2.NewClient(host, token)
	influxClient = &client

	_, err := client.Health(context.Background())

	return err

}

func DisconnectInfluxDB() error {
	(*influxClient).Close()

	return nil
}

func GetInfluxDBClient() *influxdb2.Client {
	return influxClient
}
