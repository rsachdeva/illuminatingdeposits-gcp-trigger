//nolint:revive,stylecheck // using underscore in package name for clarity
package gcf_analytics

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"time"

	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/civil"
	"cloud.google.com/go/storage"
	"github.com/go-playground/validator/v10"
)

const (
	projectID = "illuminatingdeposits-gcp"
	datasetID = "gcfdeltaanalytics"
	tableID   = "delta_calculations"
)

type DeltaCalculations struct {
	CreateInterestResponse
	CreatedAt civil.DateTime `bigquery:"created_at"`
}

type CreateInterestResponse struct {
	Banks []*Bank `json:"banks,omitempty"`
	Delta float64 `json:"delta,omitempty"`
}

type Deposit struct {
	Account     string  `json:"account,omitempty"`
	AccountType string  `json:"account_type,omitempty" bigquery:"account_type"` //nolint:tagliatelle
	Apy         float64 `json:"apy,omitempty"`
	Years       float64 `json:"years,omitempty"`
	Amount      float64 `json:"amount,omitempty"`
	Delta       float64 `json:"delta,omitempty"`
}

type Bank struct {
	Name     string     `json:"name,omitempty"`
	Deposits []*Deposit `json:"deposits,omitempty"`
	Delta    float64    `json:"delta,omitempty"`
}

// use a single instance of Validate, it caches struct info.
var validate *validator.Validate //nolint:gochecknoglobals

type LoadService struct{}

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

func appendCal(ctx context.Context, respData []byte) error {
	var (
		loadSvc LoadService
		resp    CreateInterestResponse
	)

	if err := json.Unmarshal(respData, &resp); err != nil {
		return fmt.Errorf("decode interest response clculation payload %w", err)
	}

	printDecodedResp(resp)

	validate = validator.New()
	if err := validate.Struct(resp); err != nil {
		return fmt.Errorf("invalid interest response clculation payload %w", err)
	}

	err := loadSvc.addToBigQueryTable(ctx, &resp)
	if err != nil {
		return fmt.Errorf("load service import to big query  err: %w", err)
	}

	log.Println("append to big query for the delta interest response was successful with created at datetime")

	return nil
}

// print nested struct.
func printDecodedResp(resp CreateInterestResponse) {
	for _, b := range resp.Banks {
		log.Printf("bank name is %#v", b.Name)

		for _, d := range b.Deposits {
			log.Printf("for this bank deposit is %#v", d)
		}
	}

	log.Printf("overall delta is %v\n", resp.Delta)
}

func (svc LoadService) addToBigQueryTable(ctx context.Context, ciresp *CreateInterestResponse) error {
	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("bigquery.NewClient err: %w", err)
	}
	defer client.Close()

	table := client.Dataset(datasetID).Table(tableID)

	name, err := table.Identifier(bigquery.StandardSQLID)
	if err != nil {
		return fmt.Errorf("table identifier err: %w", err)
	}

	log.Printf("table fully qualified name is: %#v\n", name)
	// SELECT
	// cir.delta
	// FROM
	// `illuminatingdeposits-gcp.gcfdeltaanalytics.create_interest_responses` cir
	// ORDER BY cir.created_at DESC
	// ciresp.CreatedAt = civil.DateTimeOf(time.Now())

	deltaCalc := DeltaCalculations{
		CreateInterestResponse: *ciresp,
		CreatedAt:              civil.DateTimeOf(time.Now()),
	}

	log.Printf("DeltaCal to be added to table is %#v", deltaCalc)

	err = table.Inserter().Put(ctx, deltaCalc)
	if err != nil {
		return fmt.Errorf("table nserter put err: %w", err)
	}

	return nil
}
