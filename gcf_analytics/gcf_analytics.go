//nolint:revive,stylecheck // using underscore in package name for clarity
package gcf_analytics

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"cloud.google.com/go/storage"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/cloudevents/sdk-go/v2/event"
)

// required as per google cloud functions version 2
// https://cloud.google.com/functions/docs/writing/write-http-functions
//
//nolint:gochecknoinits
func init() {
	functions.CloudEvent("DeltaCalAppendAnalytics", deltaCalAppendAnalytics)
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
func deltaCalAppendAnalytics(ctx context.Context, evt event.Event) error {
	log.Println("deltaCalAppendAnalytics Triggered: event triggered by output storage bucket file")
	log.Printf("Event ID: %s", evt.ID())
	log.Printf("Event Type: %s", evt.Type())

	var data StorageObjectData
	if err := evt.DataAs(&data); err != nil {
		return fmt.Errorf("event.DataAs: %w", err)
	}

	log.Printf("Bucket: %s", data.Bucket)
	log.Printf("File: %s", data.Name)
	log.Printf("Metageneration: %d", data.Metageneration)
	log.Printf("Created: %s", data.TimeCreated)
	log.Printf("Updated: %s", data.Updated)

	log.Println("ctx is", ctx)

	// read the file from the bucket
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

func readObjectFromBucket(ctx context.Context, bucketName, objectName string) ([]byte, error) {
	// Creates a client.
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("storage.NewClient: %w", err)
	}

	// Get the bucket.
	bucket := client.Bucket(bucketName)

	// Get the object.
	obj := bucket.Object(objectName)

	// Read the data from the object.
	reader, err := obj.NewReader(ctx)
	if err != nil {
		return nil, fmt.Errorf("obj.NewReader: %w", err)
	}
	defer reader.Close()

	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("ioutil.ReadAll: %w", err)
	}

	return data, nil
}
