//nolint:revive,stylecheck // using underscore in package name for clarity
package gcf_analytics

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/cloudevents/sdk-go/v2/event"
)

// required as per google cloud functions version 2
// https://cloud.google.com/functions/docs/writing/write-http-functions
//
//nolint:gochecknoinits
func init() {
	functions.CloudEvent("DeltaCalAnalytics", deltaCalAnalytics)
}

// StorageObjectData contains metadata of the Cloud Storage object.
type StorageObjectData struct {
	Bucket         string    `json:"bucket,omitempty"`
	Name           string    `json:"name,omitempty"`
	Metageneration int64     `json:"metageneration,string,omitempty"`
	TimeCreated    time.Time `json:"timeCreated,omitempty"`
	Updated        time.Time `json:"updated,omitempty"`
}

// analyzeDeltaStorage consumes a CloudEvent message for Interest calculation with Delta uploaded to output bucket.
func deltaCalAnalytics(ctx context.Context, evt event.Event) error {
	log.Println("deltaCalAnalytics Triggered: event triggered by output storage bucket file")
	log.Printf("Event ID: %s", evt.ID())
	log.Printf("Event Type: %s", evt.Type())

	var data StorageObjectData
	if err := evt.DataAs(&data); err != nil {
		return fmt.Errorf("event.DataAs: %w", err)
	}

	log.Printf("Bucket: %s", data.Bucket)
	log.Printf("File: %s", data.Name)
	// log.Printf("Metageneration: %d", data.Metageneration)
	// log.Printf("Created: %s", data.TimeCreated)
	// log.Printf("Updated: %s", data.Updated)
	log.Println("ctx is", ctx)
	log.Println("QUERY_MODE_ONLY_DEVELOPMENT is", os.Getenv("QUERY_MODE_ONLY_DEVELOPMENT"))

	// if queryOnlyMode is true, then  no data will be imported to big query
	// only query will be executed and published to pubsub
	queryOnlyMode, err := strconv.ParseBool(os.Getenv("QUERY_MODE_ONLY_DEVELOPMENT"))
	if err != nil {
		return fmt.Errorf("strconv.ParseBool for QUERY_MODE_ONLY_DEVELOPMENT : %w", err)
	}

	log.Println("queryOnlyMode bool got is", queryOnlyMode)

	// read the file from the bucket
	if err := readAndAppendToBigQuery(ctx, queryOnlyMode, data); err != nil {
		return fmt.Errorf("readAndAppendToBigQuery: %w", err)
	}

	if err := queryAndPublishAnalytics(ctx); err != nil {
		return fmt.Errorf("queryAndPublishAnalytics: %w", err)
	}

	return nil
}

func readAndAppendToBigQuery(ctx context.Context, queryOnlyMode bool, data StorageObjectData) error {
	if queryOnlyMode {
		return nil
	}

	respData, err := readObjectFromBucket(ctx, data.Bucket, data.Name)
	if err != nil {
		return fmt.Errorf("readObjectFromBucket: %w", err)
	}

	log.Printf("File Data read is: %v", string(respData))

	err = appendCal(ctx, respData)
	if err != nil {
		log.Println("import data to big query error is", err)

		return fmt.Errorf("appendCal err : %w", err)
	}

	return nil
}
