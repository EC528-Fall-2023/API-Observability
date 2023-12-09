package main

import (
	pb "centralReg/service_reg"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

type SystemMetrics struct {
	CPUUsage    []float64
	MemoryUsage *mem.VirtualMemoryStat
	DiskUsage   []disk.UsageStat
	NetIO       []net.IOCountersStat
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

type ApiService struct {
	pb.APIService
	WSConn *websocket.Conn // WebSocket connection for REST services
}

var registry = struct {
	sync.RWMutex
	services map[string]*ApiService
}{services: make(map[string]*ApiService)}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {
	go monitorServices()
	setupEndpoints()

	clientWrapperSystem := NewClientWrapper("http://35.236.200.122:8086/", "CentralReg", "API-Observability", "AxNHAn8hBBhsHz0o6HVJ2iM9gfGqybVWugTx5crw0o2yvkPTURsZqztPjxOXp4YWR2Hy9jiQPZePyilXFh7lcg==")
	defer clientWrapperSystem.client.Close()

	// Run metric collection in a separate goroutine
	go startMetricCollection(clientWrapperSystem, 15*time.Second)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8090"
	}

	log.Printf("Service Registry is running on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func startMetricCollection(clientWrapper *ClientWrapper, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		metrics, err := collectSystemMetrics()
		if err != nil {
			fmt.Println("Error collecting system metrics:", err)
			continue
		}

		if err := clientWrapper.WriteMetric("system_metrics", nil, metrics); err != nil {
			fmt.Println("Error writing to InfluxDB:", err)
		}
	}
}

func collectSystemMetrics() (map[string]interface{}, error) {
	var metrics = make(map[string]interface{})

	// Collect CPU usage
	cpuUsage, err := cpu.Percent(time.Second, false) // false to get aggregated percentage across all CPUs
	if err != nil {
		return nil, fmt.Errorf("error collecting CPU metrics: %w", err)
	}
	if len(cpuUsage) > 0 {
		metrics["cpu_usage"] = cpuUsage[0] // Assuming you want the aggregated CPU usage
	}

	// Collect Memory usage
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		return nil, fmt.Errorf("error collecting memory metrics: %w", err)
	}
	metrics["memory_total"] = vmStat.Total
	metrics["memory_used"] = vmStat.Used
	metrics["memory_free"] = vmStat.Free

	// Collect Disk usage
	// Assuming you want to collect data of the root partition '/'
	diskStat, err := disk.Usage("/")
	if err != nil {
		return nil, fmt.Errorf("error collecting disk metrics: %w", err)
	}
	metrics["disk_total"] = diskStat.Total
	metrics["disk_used"] = diskStat.Used
	metrics["disk_free"] = diskStat.Free

	// Collect Network I/O stats
	netIO, err := net.IOCounters(false) // false to get stats summed across all network interfaces
	if err != nil {
		return nil, fmt.Errorf("error collecting network IO metrics: %w", err)
	}
	if len(netIO) > 0 {
		metrics["net_bytes_sent"] = netIO[0].BytesSent
		metrics["net_bytes_recv"] = netIO[0].BytesRecv
		metrics["net_packets_sent"] = netIO[0].PacketsSent
		metrics["net_packets_recv"] = netIO[0].PacketsRecv
	}

	return metrics, nil
}

func setupEndpoints() {
	http.HandleFunc("/register", registerService)
	http.HandleFunc("/services", listServices)
	http.HandleFunc("/metrics", metricsEndpoint)
}

func monitorServices() {
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for {
		<-ticker.C
		registry.RLock()
		for _, service := range registry.services {
			if service.WSConn != nil {
				err := service.WSConn.WriteMessage(websocket.PingMessage, nil)
				if err != nil {
					log.Printf("Failed to send heartbeat to service %s: %v", service.Name, err)
					service.Status = "Down"
				} else {
					service.Status = "Up"
					log.Printf("Successful heartbeat to service %s", service.Name)
				}
			}
		}
		registry.RUnlock()
	}
}

func registerService(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade to websocket failed:", err)
		return
	}

	_, msg, err := conn.ReadMessage()
	if err != nil {
		log.Println("Error reading registration message:", err)
		return
	}

	var reg pb.Registration
	if err := json.Unmarshal(msg, &reg); err != nil {
		log.Printf("Error decoding registration data: %s", err)
		return
	}

	registry.Lock()
	if existingService, exists := registry.services[reg.Name]; exists {
		existingService.Type = reg.Type
		existingService.WSConn = conn
	} else {
		registry.services[reg.Name] = &ApiService{
			APIService: pb.APIService{
				Name: reg.Name,
				Type: reg.Type,
			},
			WSConn: conn,
		}
	}
	registry.Unlock()
}

func listServices(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade to websocket failed:", err)
		return
	}
	defer conn.Close()

	registry.RLock()
	defer registry.RUnlock()

	services := make([]*ApiService, 0, len(registry.services))
	for _, service := range registry.services {
		services = append(services, service)
	}

	err = conn.WriteJSON(services)
	if err != nil {
		log.Println("Error sending services over WebSocket:", err)
	}
}

func metricsEndpoint(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade to websocket failed:", err)
		return
	}
	defer conn.Close()

	handleMetrics(conn)
}

func handleMetrics(conn *websocket.Conn) {
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading metrics message:", err)
			break
		}

		var metrics Metrics
		if err := json.Unmarshal(msg, &metrics); err != nil {
			log.Printf("Error decoding metrics data: %s", err)
			continue
		}

		if err := sendToInfluxDB(metrics); err != nil {
			log.Printf("Error sending to InfluxDB: %s", err)
		}
	}
}

func sendToInfluxDB(metrics Metrics) error {
	clientWrapper := NewClientWrapper(metrics.InfluxDBURL, metrics.Bucket, metrics.Org, metrics.Token)
	defer clientWrapper.client.Close()
	err := clientWrapper.WriteMetric(metrics.Measurement, metrics.Tags, metrics.Fields)
	if err != nil {
		return err
	}
	return nil
}

type ClientWrapper struct {
	client influxdb2.Client
	org    string
	bucket string
}

func NewClientWrapper(url, bucket, org, token string) *ClientWrapper {
	client := influxdb2.NewClient(url, token)
	return &ClientWrapper{
		client: client,
		org:    org,
		bucket: bucket,
	}
}

func (cw *ClientWrapper) WriteMetric(measurement string, tags map[string]string, fields map[string]interface{}) error {
	writeAPI := cw.client.WriteAPIBlocking(cw.org, cw.bucket)
	p := influxdb2.NewPoint(measurement, tags, fields, time.Now())
	return writeAPI.WritePoint(context.Background(), p)
}
