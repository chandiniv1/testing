package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

type Confidence struct {
	Block                uint32  `json:"status"`
	Confidence           float64 `json:"confidence"`
	SerialisedConfidence *string `json:"serialised_confidence,omitempty"`
}

const LightClientURL = "http://localhost:8000"

func GetConfidenceHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the block number from the URL path parameter
	vars := mux.Vars(r)
	blockStr := vars["blockNumbers"]

	blockNumber, err := strconv.ParseUint(blockStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid block number", http.StatusBadRequest)
		return
	}

	confidenceURL := fmt.Sprintf("%s/v2/confidence/%d", LightClientURL, blockNumber)
	response, err := http.Get(confidenceURL)
	if err != nil {
		http.Error(w, "Failed to retrieve confidence", http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		var confidence Confidence
		err := json.NewDecoder(response.Body).Decode(&confidence)
		if err != nil {
			http.Error(w, "Failed to decode response", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(confidence)
	} else {
		http.Error(w, "Failed to retrieve confidence", response.StatusCode)
	}
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/v1/confidence/{blockNumber}", GetConfidenceHandler).Methods("GET")

	log.Fatal(http.ListenAndServe(":7000", router))
}
