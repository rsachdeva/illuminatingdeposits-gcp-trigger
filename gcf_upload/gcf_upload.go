package interestcal

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/go-playground/validator/v10"
)

func init() {
	functions.HTTP("UploadHTTP", uploadHTTP)
}

// use a single instance of Validate, it caches struct info
var validate *validator.Validate

type NewDeposit struct {
	Account     string  `json:"account,omitempty"`
	AccountType string  `json:"account_type" validate:"required"`
	Apy         float64 `json:"apy" validate:"gte=0"`
	Years       float64 `json:"years" validate:"required"`
	Amount      float64 `json:"amount" validate:"required"`
}

type NewBank struct {
	Name        string        `json:"name" validate:"required"`
	NewDeposits []*NewDeposit `json:"new_deposits" validate:"required,dive"`
}

type CreateInterestRequest struct {
	NewBanks []*NewBank `json:"new_banks" validate:"required,dive"`
}

// uploadHTTP is an HTTP Cloud Function with a request parameter.
func uploadHTTP(writer http.ResponseWriter, r *http.Request) {
	log.Printf("Request received is: %v", r)

	var req CreateInterestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Request not sucessfuly submitted. could not decode json %v\n", err)
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

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
