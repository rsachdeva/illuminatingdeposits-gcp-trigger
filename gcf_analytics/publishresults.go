//nolint:revive,stylecheck // using underscore in package name for clarity
package gcf_analytics

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"cloud.google.com/go/pubsub"
)

func queryAndPublishAnalytics(ctx context.Context) error {
	var strBuilder strings.Builder

	err := queryDepositHighestDelta(ctx, &strBuilder)
	if err != nil {
		log.Println("queryDepositHighestDelta error is", err)

		return fmt.Errorf("queryDepositHighestDelta err : %w", err)
	}

	analyticResultsStr, err := analyticsResultsJson(strBuilder)
	if err != nil {
		return fmt.Errorf("analyticsResultsJson error: %w", err)
	}

	topicId := os.Getenv("PUBSUB_DELTA_ANALYTICS_TOPIC")
	log.Println("topicId for pubsub to push results is", topicId)

	err = publishAnalytics(ctx, log.Writer(), topicId, analyticResultsStr)
	if err != nil {
		return fmt.Errorf("publishAnalytics error: %w", err)
	}

	return nil
}

func analyticsResultsJson(strBuilder strings.Builder) (string, error) {
	queryHighestDeltaDepositsByDate := strBuilder.String()
	log.Println("queryHighestDeltaDepositsByDate i", queryHighestDeltaDepositsByDate)

	type AnalyticsResults struct {
		//nolint:tagliatelle // json tag is required; keeping snake case for protocol buffer compatibility
		HighestDeltaDepositsByDate string `json:"highest_delta_deposits_by_date,omitempty"`
	}

	analyticsResults := AnalyticsResults{
		HighestDeltaDepositsByDate: queryHighestDeltaDepositsByDate,
	}

	// convert analyticsResults to json
	analyticsResultsJson, err := json.Marshal(analyticsResults)
	if err != nil {
		return "", fmt.Errorf("json.Marshal for analyticsResultsJson: %w", err)
	}

	log.Println("string(analyticsResultsJson) is", string(analyticsResultsJson))

	return string(analyticsResultsJson), nil
}

func publishAnalytics(ctx context.Context, writer io.Writer, topicID, msg string) error {
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("pubsub: NewClient: %writer", err)
	}
	defer client.Close()

	topic := client.Topic(topicID)

	//nolint:exhaustruct // pubsub.Message does not to fill all fields
	pubMsg := pubsub.Message{
		// ID:              "",
		Data: []byte(msg),
		// Attributes:      nil,
		// PublishTime:     time.Time{},
		// DeliveryAttempt: nil,
		// OrderingKey:     "",
	}
	result := topic.Publish(ctx, &pubMsg)
	// Block until the result is returned and a server-generated
	// ID is returned for the published message.
	id, err := result.Get(ctx)
	if err != nil {
		return fmt.Errorf("pubsub: result.Get: %writer", err)
	}

	fmt.Fprintf(writer, "Published a message; msg ID: %v\n", id)

	return nil
}
