package influxdb

import (
	"context"
	"github.com/influxdata/influxdb-client-go/v2"
	"time"
)

const (
	url    = "http://34.86.236.100"
	bucket = "testing"
	org    = "API-Observability"
	token  = "AxNHAn8hBBhsHz0o6HVJ2iM9gfGqybVWugTx5crw0o2yvkPTURsZqztPjxOXp4YWR2Hy9jiQPZePyilXFh7lcg=="
)

func InitInfluxDB() influxdb2.Client {
	return influxdb2.NewClient(url, token)
}

func WriteMetric(c influxdb2.Client, measurement string, tags map[string]string, fields map[string]interface{}) error {
	// Use a background context
	writeAPI := c.WriteAPIBlocking(org, bucket)
	p := influxdb2.NewPoint(measurement, tags, fields, time.Now())
	return writeAPI.WritePoint(context.Background(), p)
}
