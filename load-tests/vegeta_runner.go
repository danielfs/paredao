package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"
)

// VegetaReport represents the JSON output from Vegeta
type VegetaReport struct {
	Latencies struct {
		Total       int64   `json:"total"`
		Mean        int64   `json:"mean"`
		P50         int64   `json:"50th"`
		P90         int64   `json:"90th"`
		P95         int64   `json:"95th"`
		P99         int64   `json:"99th"`
		Max         int64   `json:"max"`
		Min         int64   `json:"min"`
		StdDev      int64   `json:"stddev"`
		MeanByCount float64 `json:"mean_by_count"`
	} `json:"latencies"`
	BytesIn struct {
		Total int     `json:"total"`
		Mean  float64 `json:"mean"`
	} `json:"bytes_in"`
	BytesOut struct {
		Total int     `json:"total"`
		Mean  float64 `json:"mean"`
	} `json:"bytes_out"`
	Earliest    time.Time      `json:"earliest"`
	Latest      time.Time      `json:"latest"`
	End         time.Time      `json:"end"`
	Duration    int64          `json:"duration"`
	Wait        int64          `json:"wait"`
	Requests    int            `json:"requests"`
	Rate        float64        `json:"rate"`
	Success     float64        `json:"success"`
	StatusCodes map[string]int `json:"status_codes"`
	Errors      []string       `json:"errors"`
}

func main() {
	// Define command line flags
	rate := flag.Int("rate", 2000, "Requests per second")
	duration := flag.String("duration", "30s", "Test duration")
	endpoint := flag.String("endpoint", "votos", "API endpoint to test (votos, participantes, votacoes, estatisticas)")
	outputDir := flag.String("output", "./results", "Output directory")
	threshold := flag.Int("threshold", 2000, "Minimum acceptable requests per second")
	flag.Parse()

	// Create output directory if it doesn't exist
	if err := os.MkdirAll(*outputDir, 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	// Run the load test
	fmt.Printf("Starting load test at %d requests per second for %s...\n", *rate, *duration)
	fmt.Printf("Target: %s\n", *endpoint)

	// Determine the targets file based on the endpoint
	targetsFile := filepath.Join(*outputDir, fmt.Sprintf("targets_%s.txt", *endpoint))
	createTargetsFile(targetsFile, *endpoint)

	// Run Vegeta attack
	resultsFile := filepath.Join(*outputDir, fmt.Sprintf("results_%s.bin", *endpoint))
	attackCmd := exec.Command("vegeta", "attack",
		"-rate", strconv.Itoa(*rate),
		"-duration", *duration,
		"-targets", targetsFile,
	)

	resultsOutput, err := os.Create(resultsFile)
	if err != nil {
		log.Fatalf("Failed to create results file: %v", err)
	}
	defer resultsOutput.Close()

	attackCmd.Stdout = resultsOutput
	attackCmd.Stderr = os.Stderr

	if err := attackCmd.Run(); err != nil {
		log.Fatalf("Failed to run Vegeta attack: %v", err)
	}

	// Generate JSON report
	jsonReportFile := filepath.Join(*outputDir, fmt.Sprintf("report_%s.json", *endpoint))
	reportCmd := exec.Command("vegeta", "report", "-type", "json", resultsFile)

	jsonOutput, err := os.Create(jsonReportFile)
	if err != nil {
		log.Fatalf("Failed to create JSON report file: %v", err)
	}
	defer jsonOutput.Close()

	reportCmd.Stdout = jsonOutput
	reportCmd.Stderr = os.Stderr

	if err := reportCmd.Run(); err != nil {
		log.Fatalf("Failed to generate JSON report: %v", err)
	}

	// Generate text report for display
	textReportFile := filepath.Join(*outputDir, fmt.Sprintf("report_%s.txt", *endpoint))
	textReportCmd := exec.Command("vegeta", "report", "-type", "text", resultsFile)

	textOutput, err := os.Create(textReportFile)
	if err != nil {
		log.Fatalf("Failed to create text report file: %v", err)
	}
	defer textOutput.Close()

	textReportCmd.Stdout = textOutput
	textReportCmd.Stderr = os.Stderr

	if err := textReportCmd.Run(); err != nil {
		log.Fatalf("Failed to generate text report: %v", err)
	}

	// Read and parse the JSON report
	jsonData, err := os.ReadFile(jsonReportFile)
	if err != nil {
		log.Fatalf("Failed to read JSON report: %v", err)
	}

	var report VegetaReport
	if err := json.Unmarshal(jsonData, &report); err != nil {
		log.Fatalf("Failed to parse JSON report: %v", err)
	}

	// Display results
	fmt.Println("\nLoad Test Results:")
	fmt.Printf("Requests: %d\n", report.Requests)
	fmt.Printf("Rate: %.2f requests/second\n", report.Rate)
	fmt.Printf("Success Rate: %.2f%%\n", report.Success*100)
	fmt.Printf("Mean Latency: %.2f ms\n", float64(report.Latencies.Mean)/1000000)
	fmt.Printf("P95 Latency: %.2f ms\n", float64(report.Latencies.P95)/1000000)
	fmt.Printf("P99 Latency: %.2f ms\n", float64(report.Latencies.P99)/1000000)

	fmt.Println("\nStatus Codes:")
	for code, count := range report.StatusCodes {
		fmt.Printf("  %s: %d\n", code, count)
	}

	if len(report.Errors) > 0 {
		fmt.Println("\nErrors:")
		for _, err := range report.Errors {
			fmt.Printf("  %s\n", err)
		}
	}

	// Check if the test passed
	fmt.Println("\nTest Results:")
	if report.Rate >= float64(*threshold) {
		fmt.Printf("✅ SUCCESS: The API handled %.2f requests per second (threshold: %d)\n", report.Rate, *threshold)
	} else {
		fmt.Printf("❌ FAILURE: The API only handled %.2f requests per second (threshold: %d)\n", report.Rate, *threshold)
	}

	if report.Success < 1.0 {
		fmt.Printf("⚠️  WARNING: Not all requests were successful. Success rate: %.2f%%\n", report.Success*100)
	}
}

// createTargetsFile creates a targets file for Vegeta based on the endpoint
func createTargetsFile(filePath, endpoint string) {
	var content string

	switch endpoint {
	case "votos":
		content = `POST http://localhost:8080/votos
Content-Type: application/json
@./payload1.json

POST http://localhost:8080/votos
Content-Type: application/json
@./payload2.json

POST http://localhost:8080/votos
Content-Type: application/json
@./payload3.json
`
	case "participantes":
		content = `GET http://localhost:8080/participantes
Accept: application/json
`
	case "votacoes":
		content = `GET http://localhost:8080/votacoes
Accept: application/json
`
	case "estatisticas":
		content = `GET http://localhost:8080/estatisticas/votacoes/1/total
Accept: application/json
`
	default:
		log.Fatalf("Unknown endpoint: %s", endpoint)
	}

	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		log.Fatalf("Failed to create targets file: %v", err)
	}
}
