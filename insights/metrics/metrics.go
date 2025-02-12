package metrics

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/applicationinsights/armapplicationinsights/v2"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/monitor/armmonitor"
	ha "github.com/NateDuff/n8ha"
)

// SiteInfo represents the site related information to be published
type SiteInfo struct {
	RequestCount int
	ErrorCount   int
	AvgPageTime  time.Duration
	Uptime       float64
}

// GetSiteInfo retrieves the site information
func GetSiteInfo(ctx context.Context, siteName string) (SiteInfo, error) {
	return SiteInfo{}, nil
}

func appInsightsIterator(ctx context.Context, cred *azidentity.ChainedTokenCredential, subscriptionID string, fn func(instanceName string, resourceID string) error) error {
	clientFactory, err := armapplicationinsights.NewClientFactory(subscriptionID, cred, nil)
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	pager := clientFactory.NewComponentsClient().NewListPager(nil)

	for pager.More() {
		resp, err := pager.NextPage(ctx)
		if err != nil {
			return fmt.Errorf("failed to get Application Insights instances: %w", err)
		}

		for _, insight := range resp.Value {
			if insight.Name == nil || insight.ID == nil {
				continue
			}

			instanceName := *insight.Name
			resourceID := *insight.ID

			if err := fn(instanceName, resourceID); err != nil {
				return fmt.Errorf("failed to process instance %s: %w", instanceName, err)
			}
		}
	}

	return nil
}

// main queries the total request count metric for all Application Insights instances in the subscription
func main() {
	subscriptionID := os.Getenv("AZURE_SUBSCRIPTION_ID")

	ctx, err := ha.WithUserAuth(context.Background())
	if err != nil {
		log.Fatalf("Failed to get user auth context: %v", err)
	}

	err = appInsightsIterator(ctx, ha.Credentials, subscriptionID, func(instanceName string, resourceID string) error {
		// Query metrics for the instance
		requestCount, err := getRequestCount(ctx, ha.Credentials, subscriptionID, resourceID)
		if err != nil {
			log.Printf("Failed to get request count for %s: %v", instanceName, err)
		}

		fmt.Printf("Instance: %s, Requests (Last 90 Days): %d\n", instanceName, requestCount)
		return nil
	})

	if err != nil {
		log.Fatalf("Failed to iterate Application Insights instances: %v", err)
	}
}

// getRequestCount queries the total request count metric for the given app insights resource
func getRequestCount(ctx context.Context, cred *azidentity.ChainedTokenCredential, subscriptionID, resourceID string) (int64, error) {
	metricsClient, err := armmonitor.NewMetricsClient(subscriptionID, cred, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to create Metrics client: %w", err)
	}

	// Define the time range (last 90 days) and metric to query
	endTime := time.Now().UTC()
	startTime := endTime.AddDate(0, 0, -90)
	metricName := "requests/count"

	resp, err := metricsClient.List(ctx, resourceID, &armmonitor.MetricsClientListOptions{
		Timespan:    toPtr(fmt.Sprintf("%s/%s", startTime.Format(time.RFC3339), endTime.Format(time.RFC3339))),
		Interval:    nil,
		Metricnames: toPtr(metricName),
		Aggregation: toPtr("Total"),
	})
	if err != nil {
		return 0, fmt.Errorf("failed to query metrics: %w", err)
	}

	// Sum the total request counts
	var totalCount int64
	for _, metric := range resp.Value {
		if metric.Timeseries == nil {
			continue
		}
		for _, timeseries := range metric.Timeseries {
			for _, data := range timeseries.Data {
				if data.Total != nil {
					totalCount += int64(*data.Total)
				}
			}
		}
	}

	return totalCount, nil
}

func toPtr[T any](v T) *T {
	return &v
}
