package interceptor

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/protobuf/proto"
	"log"
	"net"
	"time"
)

const (
	serverURL      = "http://35.236.200.122:8086/"
	influxDBToken  = "AxNHAn8hBBhsHz0o6HVJ2iM9gfGqybVWugTx5crw0o2yvkPTURsZqztPjxOXp4YWR2Hy9jiQPZePyilXFh7lcg=="
	influxDBOrg    = "API-Observability"
	influxDBBucket = "combined_metrics"
)

type Metrics struct {
	InfluxDBURL string                 `json:"influxdb_url"`
	Token       string                 `json:"token"`
	Org         string                 `json:"org"`
	Bucket      string                 `json:"bucket"`
	Measurement string                 `json:"measurement"`
	Tags        map[string]string      `json:"tags"`
	Fields      map[string]interface{} `json:"fields"`
}

func MetricsInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()
	resp, err := handler(ctx, req)
	duration := time.Since(start)
	// Measure request and response size (assuming they can be converted to string)
	reqSize := proto.Size(req.(proto.Message))
	respSize := 0
	if resp != nil {
		respSize = proto.Size(resp.(proto.Message))
	}

	// Get method name
	methodName := info.FullMethod
	//statusCode := status.Code(err).String()

	// Extract peer information
	p, ok := peer.FromContext(ctx)
	ipAddress := ""
	if ok && p.Addr != net.Addr(nil) {
		host, _, err := net.SplitHostPort(p.Addr.String())
		if err == nil {
			ipAddress = host
		} else {
			log.Fatalf("Error while parsing peer address: %v", err)
		}
	}

	// Increment request count
	requestCount := 1

	// Error rate
	errorRate := 0
	if err != nil {
		errorRate = 1
	}

	// Extract metadata from context
	md, ok := metadata.FromIncomingContext(ctx)
	userAgent := ""
	if ok {
		// Metadata keys are normalized to lowercase
		if ua, exists := md["user-agent"]; exists && len(ua) > 0 {
			userAgent = ua[0]
		}
	}

	metrics := Metrics{
		InfluxDBURL: serverURL,
		Token:       influxDBToken,
		Org:         influxDBOrg,
		Bucket:      influxDBBucket,
		Measurement: "Student-Info gRPC Service",
		Tags:        map[string]string{"endpoint": methodName, "ip_address": ipAddress, "user_agent": userAgent},
		Fields: map[string]interface{}{
			"duration":      duration.Seconds(),
			"error":         err != nil,
			"request_size":  reqSize,
			"response_size": respSize,
			"request_count": requestCount,
			"error_rate":    errorRate,
		},
	}

	//fmt.Printf("metrics %v", metrics)

	metricsErr := sendMetrics(metrics, "wss://centralreg-necuf5ddgq-ue.a.run.app/metrics")
	if metricsErr != nil {
		fmt.Printf("%v", metricsErr)
		return nil, metricsErr
	}

	// Record metrics to InfluxDB (or print to console, log, etc.)
	//writeMetrics(duration, err, reqSize, respSize, methodName, statusCode, requestCount, errorRate, ipAddress, userAgent)
	return resp, err
}

func writeMetrics(duration time.Duration, err error, reqSize, respSize int, methodName, statusCode string, requestCount, errorRate int, ipAddress, userAgent string) {
	// Create a new InfluxDB client
	client := influxdb2.NewClient(serverURL, influxDBToken)
	defer client.Close()

	// Create a write API (this can be reused)
	writeAPI := client.WriteAPI(influxDBOrg, influxDBBucket)

	// Create a point to write (measurement name is "gRPCMetrics")
	point := influxdb2.NewPoint(
		"gRPCMetrics",
		map[string]string{"endpoint": methodName, "ip_address": ipAddress, "user_agent": userAgent},
		map[string]interface{}{
			"duration":      duration.Seconds(),
			"error":         err != nil,
			"request_size":  reqSize,
			"response_size": respSize,
			"request_count": requestCount,
			"error_rate":    errorRate,
		},
		time.Now(),
	)

	// Write the point
	writeAPI.WritePoint(point)

	// Ensure data is written
	writeAPI.Flush()
}

func sendMetrics(metrics Metrics, centralRegisterWSURL string) error {
	// Serialize the Metrics struct into JSON
	jsonData, err := json.Marshal(metrics)
	if err != nil {
		return err
	}

	// Connect to the WebSocket server
	c, _, err := websocket.DefaultDialer.Dial(centralRegisterWSURL, nil)
	if err != nil {
		log.Println("dial:", err)
		return err
	}
	defer c.Close()

	// Send the JSON data to the WebSocket server
	err = c.WriteMessage(websocket.TextMessage, jsonData)
	if err != nil {
		log.Println("write:", err)
		return err
	}

	return nil
}
