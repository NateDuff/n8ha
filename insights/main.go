package main

// Application to

import (
	"context"
	"encoding/json"
	"log"
	"time"

	ha "github.com/NateDuff/n8ha"
	"github.com/NateDuff/n8ha/insights/metrics"
)

var (
	reportInterval = 5
	siteConfigs    = []siteConfig{
		{
			subscriptionID:  "b45f7ace-a0c3-4bbf-a66e-6ca65f2484ea", // JustPoshEsthetics.com
			appInsightsName: "ai-justposh",
		},
		{
			subscriptionID:  "37b26ee4-9d20-4f0b-9521-f87e6beb4c81", // KineticEnergyPT.com
			appInsightsName: "ai-kineticenergy",
		},
		{
			subscriptionID:  "8c73818d-17aa-49c4-8876-c9a53f09ba11", // DCS
			appInsightsName: "ai-dcs-internal",
		},
		{
			subscriptionID:  "8c73818d-17aa-49c4-8876-c9a53f09ba11", // NateDuff.com
			appInsightsName: "ai-dcs-blog",
		},
		{
			subscriptionID:  "8c73818d-17aa-49c4-8876-c9a53f09ba11", // NicDuff.com
			appInsightsName: "ai-nicmed",
		},
		{
			subscriptionID:  "e9744476-ad56-44c8-a7fe-24b32504a678", // K9ProDogFence.com
			appInsightsName: "ai-k9prodogfence",
		},
		{
			subscriptionID:  "", // KimDuffHomes.com
			appInsightsName: "",
		},
	}
)

type siteConfig struct {
	subscriptionID  string
	appInsightsName string
}

// publishToMQTT publishes the SiteInfo value to MQTT
func publishToMQTT(svc ha.MqttService, topic string, siteInfo metrics.SiteInfo) {
	payload, err := json.Marshal(siteInfo)
	if err != nil {
		log.Printf("Failed to marshal site info: %v", err)
		return
	}

	if token := svc.Client.Publish(topic, 0, false, payload); token.Wait() && token.Error() != nil {
		log.Printf("Failed to publish to MQTT: %v", token.Error())
	}
}

func pullSiteInfo(ctx context.Context, svc *ha.MqttService) {
	sites := getSites()

	for _, siteName := range sites {
		siteInfo, err := metrics.GetSiteInfo(ctx, siteName)
		if err != nil {
			log.Printf("Failed to get site info: %v", err)
			continue
		}

		publishToMQTT(*svc, "homeassistant/sites/"+siteName, siteInfo)
		time.Sleep(time.Duration(reportInterval) * time.Minute)
	}
}

func main() {
	ctx := context.Background()
	svc := ha.NewMqttService()

	for {
		pullSiteInfo(ctx, svc)
	}
}
