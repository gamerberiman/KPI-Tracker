package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// APIServer starts the HTTP server for API endpoints
func startAPIServer() {
	router := mux.NewRouter()

	// Routes for roles
	router.HandleFunc("/api/roles", getRolesHandler).Methods("GET")
	router.HandleFunc("/api/roles/{id}", getRoleHandler).Methods("GET")

	// Routes for KPIs
	// router.HandleFunc("/api/kpis", getKPIsHandler).Methods("GET")
	// router.HandleFunc("/api/kpis/{id}", getKPIHandler).Methods("GET")
	// router.HandleFunc("/api/roles/{id}/kpis", getRoleKPIsHandler).Methods("GET")

	// // Routes for measurements
	// router.HandleFunc("/api/measurements", getMeasurementsHandler).Methods("GET")
	// router.HandleFunc("/api/measurements", createMeasurementHandler).Methods("POST")
	// router.HandleFunc("/api/measurements/{id}", updateMeasurementHandler).Methods("PUT")
	// router.HandleFunc("/api/kpis/{id}/measurements", getKPIMeasurementsHandler).Methods("GET")

	// // Routes for reports
	// router.HandleFunc("/api/reports/monthly", generateMonthlyReportHandler).Methods("GET")
	// router.HandleFunc("/api/reports/quarterly", generateQuarterlyReportHandler).Methods("GET")
	// router.HandleFunc("/api/reports/yearly", generateYearlyReportHandler).Methods("GET")
	// router.HandleFunc("/api/reports/custom", generateCustomReportHandler).Methods("GET")

	// // Routes for charts data
	// router.HandleFunc("/api/charts/trends", getTrendsDataHandler).Methods("GET")
	// router.HandleFunc("/api/charts/yearly", getYearlyDataHandler).Methods("GET")

	// Add CORS middleware
	corsMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}

	// Apply middleware
	corsRouter := corsMiddleware(router)

	// Start server
	fmt.Println("API server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", corsRouter))
}

// Implement handler functions below
func getRolesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(roles)
}

func getRoleHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid role ID", http.StatusBadRequest)
		return
	}

	for _, role := range roles {
		if role.ID == id {
			json.NewEncoder(w).Encode(role)
			return
		}
	}

	http.Error(w, "Role not found", http.StatusNotFound)
}

// Implement remaining handlers...

func getKPIsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(kpis)
}

func getRoleKPIsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid role ID", http.StatusBadRequest)
		return
	}

	var roleKPIs []KPI
	for _, kpi := range kpis {
		if kpi.RoleID == id {
			roleKPIs = append(roleKPIs, kpi)
		}
	}

	json.NewEncoder(w).Encode(roleKPIs)
}

func createMeasurementHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var measurement Measurement
	err := json.NewDecoder(r.Body).Decode(&measurement)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Generate new ID
	newID := 1
	if len(measurements) > 0 {
		for _, m := range measurements {
			if m.ID >= newID {
				newID = m.ID + 1
			}
		}
	}

	measurement.ID = newID
	measurement.CreatedAt = time.Now()

	// Check if measurement already exists
	existing := getExistingMeasurement(measurement.KPIID, measurement.Period)
	if existing != nil {
		existing.MetricValue = measurement.MetricValue
		existing.Unit = measurement.Unit
		existing.Notes = measurement.Notes
		json.NewEncoder(w).Encode(existing)
	} else {
		// Add new measurement
		measurements = append(measurements, measurement)
		json.NewEncoder(w).Encode(measurement)
	}

	// Save to Excel
	saveToExcel()
}

// More handlers would be implemented here...
