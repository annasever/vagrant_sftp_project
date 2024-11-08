package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sync"
)

// LogData structure to hold individual log data
type LogData struct {
	Filename   string `json:"filename"`
	Result     int    `json:"result"`
	SFTPServer string `json:"sftp_server"`
	Timestamp  string `json:"timestamp"`
}

// Data storage and mutex for synchronization
var (
	logStorage = make([]LogData, 0)
	mu         sync.Mutex
)

// Template for the HTML report
var reportTemplate = `
<!DOCTYPE html>
<html>
<head>
    <title>Logs Report</title>
    <style>
        body {
            font-family: Arial, sans-serif;
        }
        .container {
            width: 80%;
            margin: auto;
        }
        table {
            width: 100%;
            border-collapse: collapse;
            margin-bottom: 20px;
        }
        table, th, td {
            border: 1px solid #ddd;
        }
        th, td {
            padding: 8px;
            text-align: left;
        }
        th {
            background-color: #f2f2f2;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>Logs Report</h1>

        <h2>Total Logs Per Server</h2>
        <table id="totalLogsTable">
            <thead>
                <tr>
                    <th>Server</th>
                    <th>Total Logs</th>
                </tr>
            </thead>
            <tbody>
                {{range $server, $count := .TotalLogs}}
                <tr>
                    <td>{{$server}}</td>
                    <td>{{$count}}</td>
                </tr>
                {{end}}
            </tbody>
        </table>

        <h2>Logs Details</h2>
        <table id="logsTable">
            <thead>
                <tr>
                    <th>Server</th>
                    <th>Filename</th>
                </tr>
            </thead>
            <tbody>
                {{range .Logs}}
                <tr>
                    <td>{{.SFTPServer}}</td>
                    <td>{{.Filename}}</td>
                </tr>
                {{end}}
            </tbody>
        </table>
    </div>
</body>
</html>
`

// Struct to hold the report data
type ReportData struct {
	TotalLogs map[string]int
	Logs      []LogData
}

// receiveData handles POST requests to receive log data
func receiveData(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed", http.StatusMethodNotAllowed)
		return
	}

	var logData LogData
	if err := json.NewDecoder(r.Body).Decode(&logData); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Save log data with synchronization to avoid race conditions
	mu.Lock()
	logStorage = append(logStorage, logData)
	mu.Unlock()

	fmt.Fprintf(w, "Log data for file %s received successfully\n", logData.Filename)
}

// showData handles GET requests to display all collected log data in JSON format
func showData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	mu.Lock()
	defer mu.Unlock()
	if err := json.NewEncoder(w).Encode(logStorage); err != nil {
		http.Error(w, "Unable to encode data", http.StatusInternalServerError)
	}
}

// report generates an HTML report showing the count of logs for each SFTP server and detailed logs
func report(w http.ResponseWriter, r *http.Request) {
	// Create a map to count logs per SFTP server
	serverCounts := make(map[string]int)
	var allLogs []LogData

	mu.Lock()
	for _, logData := range logStorage {
		serverCounts[logData.SFTPServer]++
		allLogs = append(allLogs, logData)
	}
	mu.Unlock()

	// Prepare the data for the report template
	reportData := ReportData{
		TotalLogs: serverCounts,
		Logs:      allLogs,
	}

	// Parse and execute the HTML template
	tmpl, err := template.New("report").Parse(reportTemplate)
	if err != nil {
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	if err := tmpl.Execute(w, reportData); err != nil {
		http.Error(w, "Error generating report", http.StatusInternalServerError)
	}
}

// homePage displays a welcome message and available endpoints
func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to the Central Log Server!")
	fmt.Fprintln(w, "Available endpoints:")
	fmt.Fprintln(w, "/receive_data - POST endpoint for agents to send log data")
	fmt.Fprintln(w, "/logs - GET endpoint to view collected log data")
	fmt.Fprintln(w, "/report - GET endpoint to view SFTP server log counts in an HTML report")
}

func main() {
	http.HandleFunc("/", homePage)                // Base URL
	http.HandleFunc("/receive_data", receiveData) // Log data receiving endpoint
	http.HandleFunc("/logs", showData)            // Log display endpoint
	http.HandleFunc("/report", report)            // Report display endpoint

	fmt.Println("Central server is running on 10.0.0.204:5000...")
	if err := http.ListenAndServe("10.0.0.204:5000", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

