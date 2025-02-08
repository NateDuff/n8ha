package metrics

import (
	"encoding/json"
	"log"
	"math"

	ha "github.com/NateDuff/n8ha"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
)

// SystemMetrics represents the system metrics to be published
type SystemMetrics struct {
	CPUUsage    float64 `json:"cpu_usage"`
	MemoryUsage float64 `json:"memory_usage"`
	DiskUsage   float64 `json:"disk_usage"`
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

// PublishMetrics publishes the system metrics to the specified MQTT topic
func PublishMetrics(svc ha.MqttService, topic string) {
	metrics, err := getSystemMetrics()
	if err != nil {
		log.Printf("Error retrieving system metrics: %v", err)
		return
	}

	payload, err := json.Marshal(metrics)
	if err != nil {
		log.Printf("Error encoding JSON: %v", err)
		return
	}

	token := svc.Client.Publish(topic, 0, false, payload)
	token.Wait()
	log.Printf("Published metrics: %s", payload)
}
