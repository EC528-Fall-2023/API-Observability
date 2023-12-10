package instrumentation

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gofiber/fiber/v2"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"sync"
	"time"
)

var (
	requestCounts sync.Map
	errorCounts   sync.Map
)

var (
	// ... (existing variable declarations)
	wsConn    *websocket.Conn
	connMutex sync.Mutex
)

var (
	influxDBURL string
	token       string
	org         string
	bucket      string
	wsSocketURL string
	measurement string
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
	size       int
}

type Metrics struct {
	InfluxDBURL string                 `json:"influxdb_url"`
	Token       string                 `json:"token"`
	Org         string                 `json:"org"`
	Bucket      string                 `json:"bucket"`
	Measurement string                 `json:"measurement"`
	Tags        map[string]string      `json:"tags"`
	Fields      map[string]interface{} `json:"fields"`
}

func InstrumentEndpoint(routerOrServer interface{}, centralregWSURL string, serviceName string, influxdburl string, Token string, Org string, Bucket string) error {

	wsSocketURL = centralregWSURL + "/metrics"
	influxDBURL = influxdburl
	token = Token
	org = Org
	bucket = Bucket
	measurement = serviceName

	switch r := routerOrServer.(type) {
	case *gin.Engine:
		r.Use(ginMetricsMiddleware())
	case *echo.Echo:
		r.Use(echoMetricsMiddleware)
	case *mux.Router:
		r.Use(gorillaMuxMetricsMiddleware)
	case *http.ServeMux:
		// Wrap the default ServeMux with the net/http middleware
		instrumentedHandler := netHttpMetricsMiddleware(r)
		http.Handle("/", instrumentedHandler)
	case *fiber.App:
		r.Use(fiberMetricsMiddleware)
	// Add additional cases here for other frameworks...
	default:
		return fmt.Errorf("unsupported framework or server type: %T", r)
	}

	return nil
}

func ginMetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		path := c.Request.URL.Path
		userAgent := c.Request.UserAgent()
		ipAddress := c.ClientIP()
		incrementEndpointRequestCount(path)
		currentCount := getEndpointRequestCount(path)
		if len(c.Errors) > 0 {
			incrementEndpointErrorCount(path)
		}
		errorCount := getEndpointErrorCount(path)
		// Continue processing
		c.Next()

		latency := time.Since(startTime)
		statusCode := c.Writer.Status()
		responseSize := c.Writer.Size()

		tags := map[string]string{
			"endpoint":   path,
			"user_agent": userAgent,
			"ip_address": ipAddress,
		}
		fields := map[string]interface{}{
			"request_size":  c.Request.ContentLength,
			"status_code":   statusCode,
			"response_size": responseSize,
			"latency_ms":    latency.Milliseconds(),
			"request_count": currentCount,
			"error_count":   errorCount,
		}

		metrics := Metrics{
			InfluxDBURL: influxDBURL,
			Token:       token,
			Org:         org,
			Bucket:      bucket,
			Measurement: measurement,
			Tags:        tags,
			Fields:      fields,
		}

		// Send metrics
		if err := sendMetrics(metrics); err != nil {
			log.Printf("Error sending metrics: %v\n", err)
		}
	}
}

func echoMetricsMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		startTime := time.Now()
		path := c.Request().URL.Path
		userAgent := c.Request().UserAgent()
		ipAddress := c.RealIP()
		incrementEndpointRequestCount(path)
		currentCount := getEndpointRequestCount(path)
		// Continue processing
		err := next(c)
		if err != nil {
			incrementEndpointErrorCount(path)
		}
		errorCount := getEndpointErrorCount(path)
		latency := time.Since(startTime)
		statusCode := c.Response().Status
		responseSize := c.Response().Size

		tags := map[string]string{
			"endpoint":   path,
			"user_agent": userAgent,
			"ip_address": ipAddress,
		}
		fields := map[string]interface{}{
			"request_size":  c.Request().ContentLength,
			"status_code":   statusCode,
			"response_size": responseSize,
			"latency_ms":    latency.Milliseconds(),
			"request_count": currentCount,
			"error_count":   errorCount,
		}

		metrics := Metrics{
			InfluxDBURL: influxDBURL,
			Token:       token,
			Org:         org,
			Bucket:      bucket,
			Measurement: measurement,
			Tags:        tags,
			Fields:      fields,
		}

		// Send metrics
		if err := sendMetrics(metrics); err != nil {
			log.Printf("Error sending metrics: %v\n", err)
		}

		return err
	}
}

func gorillaMuxMetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		path := r.URL.Path
		userAgent := r.UserAgent()
		ipAddress := r.RemoteAddr // You might want to parse out just the IP
		incrementEndpointRequestCount(path)
		currentCount := getEndpointRequestCount(path)
		// Response writer wrapper to capture the status code and size
		rw := NewResponseWriter(w)
		next.ServeHTTP(rw, r)

		if rw.StatusCode() >= 400 {
			incrementEndpointErrorCount(path)
		}
		errorCount := getEndpointErrorCount(path)
		latency := time.Since(startTime)
		statusCode := rw.StatusCode
		responseSize := rw.Size

		tags := map[string]string{
			"endpoint":   path,
			"user_agent": userAgent,
			"ip_address": ipAddress,
		}
		fields := map[string]interface{}{
			"request_size":  r.ContentLength,
			"status_code":   statusCode,
			"response_size": responseSize,
			"latency_ms":    latency.Milliseconds(),
			"request_count": currentCount,
			"error_count":   errorCount,
		}

		metrics := Metrics{
			InfluxDBURL: influxDBURL,
			Token:       token,
			Org:         org,
			Bucket:      bucket,
			Measurement: measurement,
			Tags:        tags,
			Fields:      fields,
		}

		// Send metrics
		if err := sendMetrics(metrics); err != nil {
			log.Printf("Error sending metrics: %v\n", err)
		}

	})
}

func netHttpMetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		path := r.URL.Path
		userAgent := r.UserAgent()
		ipAddress := r.RemoteAddr // You might want to parse out just the IP
		incrementEndpointRequestCount(path)
		currentCount := getEndpointRequestCount(path)
		// Response writer wrapper to capture the status code and size
		rw := NewResponseWriter(w)
		next.ServeHTTP(rw, r)

		if rw.StatusCode() >= 400 {
			incrementEndpointErrorCount(path)
		}

		latency := time.Since(startTime)
		statusCode := rw.StatusCode
		responseSize := rw.Size
		errorCount := getEndpointErrorCount(path)
		tags := map[string]string{
			"endpoint":   path,
			"user_agent": userAgent,
			"ip_address": ipAddress,
		}
		fields := map[string]interface{}{
			"request_size":  r.ContentLength,
			"status_code":   statusCode,
			"response_size": responseSize,
			"latency_ms":    latency.Milliseconds(),
			"request_count": currentCount,
			"error_count":   errorCount,
		}

		metrics := Metrics{
			InfluxDBURL: influxDBURL,
			Token:       token,
			Org:         org,
			Bucket:      bucket,
			Measurement: measurement,
			Tags:        tags,
			Fields:      fields,
		}

		// Send metrics
		if err := sendMetrics(metrics); err != nil {
			log.Printf("Error sending metrics: %v\n", err)
		}
	})
}

func fiberMetricsMiddleware(c *fiber.Ctx) error {
	startTime := time.Now()
	path := c.OriginalURL()
	userAgent := c.Get(fiber.HeaderUserAgent)
	ipAddress := c.IP()
	incrementEndpointRequestCount(path)
	currentCount := getEndpointRequestCount(path)
	// Continue processing
	err := c.Next()
	if err != nil {
		incrementEndpointErrorCount(path)
	}
	errorCount := getEndpointErrorCount(path)
	latency := time.Since(startTime)
	statusCode := c.Response().StatusCode()
	responseSize := len(c.Response().Body()) // Fiber may have a better way to get this

	tags := map[string]string{
		"endpoint":   path,
		"user_agent": userAgent,
		"ip_address": ipAddress,
	}
	fields := map[string]interface{}{
		"request_size":  c.Request().Header.ContentLength(),
		"status_code":   statusCode,
		"response_size": responseSize,
		"latency_ms":    latency.Milliseconds(),
		"request_count": currentCount,
		"error_count":   errorCount,
	}

	metrics := Metrics{
		InfluxDBURL: influxDBURL,
		Token:       token,
		Org:         org,
		Bucket:      bucket,
		Measurement: measurement,
		Tags:        tags,
		Fields:      fields,
	}

	// Send metrics
	if err := sendMetrics(metrics); err != nil {
		log.Printf("Error sending metrics: %v\n", err)
	}

	return err
}

func incrementEndpointRequestCount(endpoint string) {
	val, _ := requestCounts.LoadOrStore(endpoint, int64(0))
	count := val.(int64)
	requestCounts.Store(endpoint, count+1)
}

// getEndpointRequestCount retrieves the current request count for a given endpoint.
func getEndpointRequestCount(endpoint string) int64 {
	val, _ := requestCounts.Load(endpoint)
	if val == nil {
		return 0
	}
	return val.(int64)
}

func incrementEndpointErrorCount(endpoint string) {
	val, _ := errorCounts.LoadOrStore(endpoint, int64(0))
	count := val.(int64)
	errorCounts.Store(endpoint, count+1)
}

// getEndpointErrorCount retrieves the current error count for a given endpoint.
func getEndpointErrorCount(endpoint string) int64 {
	val, _ := errorCounts.Load(endpoint)
	if val == nil {
		return 0
	}
	return val.(int64)
}

func NewResponseWriter(w http.ResponseWriter) *responseWriter {
	// Default the status code to 200 for HTTP, since if WriteHeader is not called explicitly,
	// the net/http package assumes a "200 OK" response.
	return &responseWriter{w, http.StatusOK, 0}
}

// WriteHeader captures the status code and calls the underlying WriteHeader method
func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// Write captures the size of the response and calls the underlying Write method
func (rw *responseWriter) Write(data []byte) (int, error) {
	size, err := rw.ResponseWriter.Write(data)
	rw.size += size
	return size, err
}

// StatusCode exposes the captured status code
func (rw *responseWriter) StatusCode() int {
	return rw.statusCode
}

// Size exposes the captured response size
func (rw *responseWriter) Size() int {
	return rw.size
}

func sendMetrics(metrics Metrics) error {
	if err := ensureWebSocketConnection(wsSocketURL); err != nil {
		return err
	}

	jsonData, err := json.Marshal(metrics)
	if err != nil {
		return err
	}

	if err := wsConn.WriteMessage(websocket.TextMessage, jsonData); err != nil {
		return fmt.Errorf("failed to write message: %v", err)
	}

	return nil
}

func ensureWebSocketConnection(centralRegisterWSURL string) error {
	connMutex.Lock()
	defer connMutex.Unlock()

	if wsConn != nil {
		return nil // Connection is already established
	}

	var err error
	wsConn, _, err = websocket.DefaultDialer.Dial(centralRegisterWSURL, nil)
	if err != nil {
		return fmt.Errorf("failed to dial WebSocket: %v", err)
	}

	// Start a goroutine to keep the connection alive
	go func() {
		for {
			if _, _, err := wsConn.NextReader(); err != nil {
				wsConn.Close()
				wsConn = nil
				return
			}
		}
	}()

	return nil
}
