package metrics

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"

	ha "github.com/NateDuff/n8ha"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
)

const gpuStatsFile = "/srv/gpu_stats.txt"

// SystemMetrics represents the system metrics to be published
type SystemMetrics struct {
	CPUUsage    float64 `json:"cpu_usage"`
	MemoryUsage float64 `json:"memory_usage"`
	DiskUsage   float64 `json:"disk_usage"`
}

// GPUMetrics represents the GPU metrics to be published
type GPUMetrics struct {
	GPUTemperature float64 `json:"gpu_temperature"`
	GPUUsage       float64 `json:"gpu_usage"`
}

// stringsToFloat converts a string to a float64
func stringsToFloat(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		fmt.Printf("Error converting string to float: %v\n", err)
	}
	return f
}

// getGPUMetrics retrieves the GPU metrics
func getGPUMetrics() (*GPUMetrics, error) {
	file, err := os.Open(gpuStatsFile)
	if err != nil {
		fmt.Printf("Error reading GPU stats file: %v\n", err)
		return nil, err
	}
	defer file.Close()

	tempInC := 0.0
	utilization := 0.0
	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		data := strings.TrimSpace(scanner.Text())
		fields := strings.Split(data, ", ")
		if len(fields) != 2 {
			fmt.Println("Unexpected file format")
			return nil, fmt.Errorf("unexpected file format")
		}

		tempInC = stringsToFloat(fields[0])
		utilization = stringsToFloat(fields[1])
	} else {
		fmt.Println("No data found in GPU stats file")
	}

	tempInFehrenheit := (1.8 * tempInC) + 32

	metrics := &GPUMetrics{
		GPUTemperature: math.Round(tempInFehrenheit*100) / 100,
		GPUUsage:       utilization,
	}

	return metrics, nil
}

// getSystemMetrics retrieves the system metrics
func getSystemMetrics() (*SystemMetrics, error) {
	// Get CPU Usage
	cpuPercentages, err := cpu.Percent(0, false)
	if err != nil {
		return nil, err
	}

	// Get Memory Usage
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}

	// Get Disk Usage (Root Partition)
	diskStat, err := disk.Usage("/")
	if err != nil {
		return nil, err
	}

	metrics := &SystemMetrics{
		CPUUsage:    math.Round(cpuPercentages[0]*100) / 100,
		MemoryUsage: math.Round(vmStat.UsedPercent*100) / 100,
		DiskUsage:   math.Round(diskStat.UsedPercent*100) / 100,
	}

	return metrics, nil
}

// PublishMetrics publishes both system and GPU metrics to the specified MQTT topics
func PublishMetrics(svc ha.MqttService, systemTopic string, gpuTopic string) {
	// Get and publish system metrics
	systemMetrics, err := getSystemMetrics()
	if err != nil {
		log.Printf("Error retrieving system metrics: %v", err)
	} else {
		systemPayload, err := json.Marshal(systemMetrics)
		if err != nil {
			log.Printf("Error encoding system metrics JSON: %v", err)
		} else {
			token := svc.Client.Publish(systemTopic, 0, false, systemPayload)
			token.Wait()
			log.Printf("Published system metrics: %s", systemPayload)
		}
	}

	// Get and publish GPU metrics
	gpuMetrics, err := getGPUMetrics()
	if err != nil {
		log.Printf("Error retrieving GPU metrics: %v", err)
	} else {
		gpuPayload, err := json.Marshal(gpuMetrics)
		if err != nil {
			log.Printf("Error encoding GPU metrics JSON: %v", err)
		} else {
			token := svc.Client.Publish(gpuTopic, 0, false, gpuPayload)
			token.Wait()
			log.Printf("Published GPU metrics: %s", gpuPayload)
		}
	}
}
