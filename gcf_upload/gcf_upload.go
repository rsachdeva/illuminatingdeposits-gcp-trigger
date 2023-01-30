package interestcal

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

func init() {
	functions.HTTP("UploadHTTP", uploadHTTP)
}

type NewDeposit struct {
	Account     string  `json:"account,omitempty"`
	AccountType string  `json:"account_type,omitempty"`
	Apy         float64 `json:"apy,omitempty"`
	Years       float64 `json:"years,omitempty"`
	Amount      float64 `json:"amount,omitempty"`
}

type NewBank struct {
	Name        string
	NewDeposits []*NewDeposit `json:"new_deposits,omitempty"`
}

type CreateInterestRequest struct {
	NewBanks []*NewBank `json:"new_banks,omitempty"`
}

// helloHTTP is an HTTP Cloud Function with a request parameter.
func uploadHTTP(writer http.ResponseWriter, r *http.Request) {
	var req CreateInterestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Request not sucessfuly submitted. could not decode json %v", err)
		respondWithError(writer, http.StatusBadRequest, "Invalid request payload")
		return
	}
	printDecodedReq(req)
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
