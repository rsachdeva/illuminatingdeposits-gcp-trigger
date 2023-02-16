//nolint:revive,stylecheck // using underscore in package name for clarity
package gcf_interestcal

import (
	"context"
	"encoding/json"
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
	functions.CloudEvent("InterestCalStorage", interestCalStorage)
}

// StorageObjectData contains metadata of the Cloud Storage object.
type StorageObjectData struct {
	Bucket         string    `json:"bucket,omitempty"`
	Name           string    `json:"name,omitempty"`
	Metageneration int64     `json:"metageneration,string,omitempty"`
	TimeCreated    time.Time `json:"timeCreated,omitempty"`
	Updated        time.Time `json:"updated,omitempty"`
}

// helloStorage consumes a CloudEvent message and logs details about the changed object.
func interestCalStorage(ctx context.Context, evt event.Event) error {
	log.Println("interestCalStorage Triggered: event triggered by input storage bucket file upload")
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
	reqData, err := readObjectFromBucket(ctx, data.Bucket, data.Name)
	if err != nil {
		return fmt.Errorf("readObjectFromBucket: %w", err)
	}

	log.Printf("File Data read is: %v", string(reqData))

	// calculate the delta
	resp, err := calculateDelta(ctx, reqData)
	if err != nil {
		log.Println("calculate delta error is", err)

		return fmt.Errorf("calculateDelta: %w", err)
	}

	printDeltaResponseToSaveToBucket(resp)

	// save response with delta calculations to the bucket
	bucketName := "illuminating_upload_json_bucket_output"
	objectName := "interestresponse.json"

	if err := writeObjectToBucket(ctx, resp, bucketName, objectName); err != nil {
		return fmt.Errorf("write object to bucket for input request %w", err)
	}

	return nil
}

func printDeltaResponseToSaveToBucket(resp *CreateInterestResponse) {
	inJSONResp, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		log.Println("err marshalling to json is", err)
	}

	log.Printf("Response Data with Delta is: %v", string(inJSONResp))
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

func writeObjectToBucket(ctx context.Context, resp *CreateInterestResponse, bucketName, objectName string) error {
	// debug.PrintStack()
	// create a client for Google Cloud Storage
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("could not create storage client %w", err)
	}
	defer client.Close()

	// upload the temporary file to the specified bucket name from env variable.
	bucket := client.Bucket(bucketName)
	bucketObj := bucket.Object(objectName)
	bucketObjWriter := bucketObj.NewWriter(ctx)
	bucketObjWriter.ObjectAttrs.ContentType = "application/json"

	// marshall go value req to json
	respWithDelta, err := json.Marshal(resp)
	if err != nil {
		return fmt.Errorf("could not marshal request %w", err)
	}
	// https://cloud.google.com/go/docs/reference/cloud.google.com/go/storage/latest#hdr-Objects
	if _, err := fmt.Fprint(bucketObjWriter, string(respWithDelta)); err != nil {
		return fmt.Errorf("could not write response with delta data %w", err)
	}

	if err := bucketObjWriter.Close(); err != nil {
		return fmt.Errorf("could not close bucket object writer %w", err)
	}

	log.Println("Delta calculated response body uploaded to file : gs://" + bucketName + "/" + objectName)

	return nil
}
