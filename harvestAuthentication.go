package main

import (
	"encoding/json"
	"fmt"
	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"log"
	"net/http"
	"time"
)

// Struct to parse incoming API request parameters
type QueryRequest struct {
	MetricName string            `json:"metric_name"`
	Labels     map[string]string `json:"labels"`
	StartTime  string            `json:"start_time"`
	EndTime    string            `json:"end_time"`
}

// Function to build a PromQL query from API request parameters
func buildPromQLQuery(req QueryRequest) string {
	// Base PromQL query with metric name
	query := req.MetricName

	// Append labels to the query
	if len(req.Labels) > 0 {
		query += "{"
		for k, v := range req.Labels {
			query += fmt.Sprintf(`%s="%s",`, k, v)
		}
		query = query[:len(query)-1] + "}" // Remove trailing comma and close label set
	}

	return query
}

// Handler function to process API request and query Prometheus
func queryHandler(w http.ResponseWriter, r *http.Request) {
	// Decode incoming JSON request
	var queryReq QueryRequest
	err := json.NewDecoder(r.Body).Decode(&queryReq)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Build the PromQL query
	query := buildPromQLQuery(queryReq)
	fmt.Printf("Generated PromQL Query: %s\n", query)

	// Parse time range
	startTime, err := time.Parse(time.RFC3339, queryReq.StartTime)
	if err != nil {
		http.Error(w, "Invalid start time format", http.StatusBadRequest)
		return
	}

	endTime, err := time.Parse(time.RFC3339, queryReq.EndTime)
	if err != nil {
		http.Error(w, "Invalid end time format", http.StatusBadRequest)
		return
	}

	// Log the query, start, and end times to ensure they are correct
	log.Printf("Executing PromQL Query: %s\nStart Time: %v\nEnd Time: %v\n", query, startTime, endTime)

	// Connect to Prometheus
	client, err := api.NewClient(api.Config{Address: "http://localhost:9090"})
	if err != nil {
		http.Error(w, "Failed to create Prometheus client", http.StatusInternalServerError)
		return
	}

	v1api := v1.NewAPI(client)
	ctx := r.Context()

	result, warnings, err := v1api.QueryRange(ctx, query, v1.Range{
		Start: startTime,
		End:   endTime,
		Step:  time.Minute,
	})
	if err != nil {
		log.Printf("Error querying Prometheus: %v\n", err)
		http.Error(w, fmt.Sprintf("Failed to query Prometheus: %v", err), http.StatusInternalServerError)
		return
	}

	if len(warnings) > 0 {
		fmt.Printf("Warnings: %v\n", warnings)
	}

	// Print the response to the server log for debugging
	log.Printf("Prometheus response type: %T\n", result)
	log.Printf("Prometheus response: %v\n", result)

	// Respond with the query results
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func main() {
	http.HandleFunc("/query", queryHandler)
	fmt.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
