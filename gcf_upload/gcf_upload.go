package interestcal

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/storage"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/go-playground/validator/v10"
)

// required as per google cloud functions version 2
// https://cloud.google.com/functions/docs/writing/write-http-functions
//
//nolint:gochecknoinits
func init() {
	functions.HTTP("UploadHTTP", uploadHTTP)
}

// use a single instance of Validate, it caches struct info
var validate *validator.Validate //nolint:gochecknoglobals

type NewDeposit struct {
	Account     string  `json:"account,omitempty"`
	AccountType string  `json:"account_type" validate:"required"` //nolint:tagliatelle
	Apy         float64 `json:"apy" validate:"gte=0"`
	Years       float64 `json:"years" validate:"required"`
	Amount      float64 `json:"amount" validate:"required"`
}

type NewBank struct {
	Name        string        `json:"name" validate:"required"`
	NewDeposits []*NewDeposit `json:"new_deposits" validate:"required,dive"` //nolint:tagliatelle
}

type CreateInterestRequest struct {
	NewBanks []*NewBank `json:"new_banks" validate:"required,dive"` //nolint:tagliatelle
}

// uploadHTTP is an HTTP Cloud Function with a request parameter.
func uploadHTTP(writer http.ResponseWriter, httpReq *http.Request) {
	log.Printf("Request body received is: %v", httpReq.Body)

	ctx := httpReq.Context()
	log.Println("ctx is", ctx)

	var req CreateInterestRequest
	if err := json.NewDecoder(httpReq.Body).Decode(&req); err != nil {
		log.Printf("Request not successfully submitted, could not decode json %v\n", err)
		respondWithError(writer, http.StatusBadRequest, fmt.Sprintf("Decode request payload %v", err))

		return
	}

	printDecodedReq(req)

	validate = validator.New()
	if err := validate.Struct(req); err != nil {
		log.Printf("Invalid request. Could not validate json %v\n", err)
		respondWithError(writer, http.StatusBadRequest, fmt.Sprintf("Invalid request payload %v", err))

		return
	}

	if err := write(ctx, req); err != nil {
		log.Println(err)
	}

	respondWithSuccess(writer, http.StatusOK, "Request submitted successfully")
}

// print nested struct
func printDecodedReq(req CreateInterestRequest) {
	for _, nb := range req.NewBanks {
		log.Printf("new bank name is %#v", nb.Name)

		for _, nd := range nb.NewDeposits {
			log.Printf("for this bank new deposit is %#v", nd)
		}
	}
}

// respond with http status and message

func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJSON(w, code, map[string]string{"error": msg})
}

func respondWithSuccess(w http.ResponseWriter, code int, msg string) {
	respondWithJSON(w, code, map[string]string{"message": msg})
}

func respondWithJSON(respWriter http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		code = http.StatusInternalServerError
		response = []byte("HTTP 500: Internal Server Error")

		log.Printf("could not marshal response %v", err)
	}

	respWriter.Header().Set("Content-Type", "application/json")
	respWriter.WriteHeader(code)

	_, err = respWriter.Write(response)
	if err != nil {
		log.Printf("could not write response %v", err)
	}
}

// save Go value instance to a json file in google cloud platform storage bucket
func write(ctx context.Context, req CreateInterestRequest) error {
	tempFile, err := createTempFileForUpload(req)
	if err != nil {
		return err
	}

	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			log.Printf("could not remove temp file %v", err)
		}
	}(tempFile.Name()) // delete the file when done

	// create a client for Google Cloud Storage
	client, err := storage.NewClient(ctx)

	if err != nil {
		return fmt.Errorf("could not create storage client %w", err)
	}
	defer client.Close()

	// open the temporary file
	file, err := os.Open(tempFile.Name())
	if err != nil {
		return fmt.Errorf("could not open temp file %w", err)
	}
	defer file.Close()

	// upload the temporary file to the specified bucket name from env variable?
	bucketName := "illuminating_upload_json_bucket"
	objectName := "inputrequest.json"
	bucketObj := client.Bucket(bucketName).Object(objectName)
	bucketObjWriter := bucketObj.NewWriter(ctx)

	if _, err := io.Copy(bucketObjWriter, file); err != nil {
		return fmt.Errorf("could not copy temp file to bucket object writer %w", err)
	}

	if err := bucketObjWriter.Close(); err != nil {
		return fmt.Errorf("could not close bucket object writer %w", err)
	}

	log.Println("Submitted request body uploaded to file : gs://" + bucketName + "/" + objectName)

	return nil
}

func createTempFileForUpload(req CreateInterestRequest) (*os.File, error) {
	// create the temporary file
	tempFile, err := os.CreateTemp("", "temp")
	if err != nil {
		return nil, fmt.Errorf("could not create temp file %w", err)
	}

	// marshall go value req to json
	submitted, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("could not marshal request %w", err)
	}
	// write the data to the file
	_, err = tempFile.Write(submitted)
	if err != nil {
		return nil, fmt.Errorf("could not write to temp file %w", err)
	}

	return tempFile, nil
}
