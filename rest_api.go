package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

// startRESTServer starts the REST API server
func startRESTServer() {
	router := mux.NewRouter()

	// Define API routes
	// Roles endpoints
	router.HandleFunc("/api/roles", getRoles).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/roles/{id}", getRole).Methods("GET", "OPTIONS")

	// KPIs endpoints
	router.HandleFunc("/api/kpis", getKPIs).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/kpis/{id}", getKPI).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/roles/{id}/kpis", getKPIsByRole).Methods("GET", "OPTIONS")

	// Measurements endpoints
	router.HandleFunc("/api/measurements", getMeasurements).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/measurements", createMeasurement).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/measurements/{id}", updateMeasurement).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/kpis/{id}/measurements", getMeasurementsByKPI).Methods("GET", "OPTIONS")

	// Reports endpoints
	router.HandleFunc("/api/reports/monthly/{year}/{month}", getMonthlyReport).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/reports/quarterly/{year}/{quarter}", getQuarterlyReport).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/reports/yearly/{year}", getYearlyReport).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/reports/custom", getCustomReport).Methods("POST", "OPTIONS")

	// Dashboard endpoints
	router.HandleFunc("/api/dashboard/overview", getDashboardOverview).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/dashboard/trends", getDashboardTrends).Methods("GET", "OPTIONS")

	// Settings endpoints
	router.HandleFunc("/api/settings", getSettings).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/settings", updateSettings).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/settings/reload-excel", reloadExcel).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/settings/save-excel", saveExcelAPI).Methods("POST", "OPTIONS")

	// Setup CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // Allow all origins for development
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	// Wrap the router with CORS middleware
	handler := c.Handler(router)

	// Start the server
	fmt.Println("Starting REST API server on http://localhost:8080")
	http.ListenAndServe(":8080", handler)
}

// Handler functions

// getRoles returns all roles
func getRoles(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(roles)
}

// getRole returns a specific role by ID
func getRole(w http.ResponseWriter, r *http.Request) {
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

// getKPIs returns all KPIs
func getKPIs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(kpis)
}

// getKPI returns a specific KPI by ID
func getKPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid KPI ID", http.StatusBadRequest)
		return
	}

	for _, kpi := range kpis {
		if kpi.ID == id {
			json.NewEncoder(w).Encode(kpi)
			return
		}
	}

	http.Error(w, "KPI not found", http.StatusNotFound)
}

// getKPIsByRole returns all KPIs for a specific role
func getKPIsByRole(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid role ID", http.StatusBadRequest)
		return
	}

	roleKPIs := getKPIsByRoleID(id)
	json.NewEncoder(w).Encode(roleKPIs)
}

// getMeasurements returns all measurements with optional filtering
func getMeasurements(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get query parameters for filtering
	query := r.URL.Query()
	yearStr := query.Get("year")
	monthStr := query.Get("month")

	// If no filters, return all measurements
	if yearStr == "" && monthStr == "" {
		json.NewEncoder(w).Encode(measurements)
		return
	}

	// Apply filters
	var filteredMeasurements []Measurement
	for _, m := range measurements {
		// Filter by year if provided
		if yearStr != "" {
			year, err := strconv.Atoi(yearStr)
			if err != nil {
				http.Error(w, "Invalid year", http.StatusBadRequest)
				return
			}
			if m.Period.Year() != year {
				continue
			}
		}

		// Filter by month if provided
		if monthStr != "" {
			month, err := strconv.Atoi(monthStr)
			if err != nil {
				http.Error(w, "Invalid month", http.StatusBadRequest)
				return
			}
			if int(m.Period.Month()) != month {
				continue
			}
		}

		filteredMeasurements = append(filteredMeasurements, m)
	}

	json.NewEncoder(w).Encode(filteredMeasurements)
}

// createMeasurement adds a new measurement
func createMeasurement(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var measurementRequest struct {
		KPIID       int       `json:"kpi_id"`
		MetricValue float64   `json:"metric_value"`
		Unit        string    `json:"unit"`
		Period      time.Time `json:"period"`
		Notes       string    `json:"notes"`
	}

	// Decode the request body
	err := json.NewDecoder(r.Body).Decode(&measurementRequest)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate the KPI ID
	validKPI := false
	for _, kpi := range kpis {
		if kpi.ID == measurementRequest.KPIID {
			validKPI = true
			break
		}
	}
	if !validKPI {
		http.Error(w, "Invalid KPI ID", http.StatusBadRequest)
		return
	}

	// Save the measurement
	saveMeasurement(
		measurementRequest.KPIID,
		measurementRequest.MetricValue,
		measurementRequest.Unit,
		measurementRequest.Period,
		measurementRequest.Notes,
	)

	// Save to Excel
	err = saveToExcel()
	if err != nil {
		http.Error(w, "Failed to save to Excel: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the newly created measurement
	existingMeasurement := getExistingMeasurement(measurementRequest.KPIID, measurementRequest.Period)
	json.NewEncoder(w).Encode(existingMeasurement)
}

// updateMeasurement updates an existing measurement
func updateMeasurement(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid measurement ID", http.StatusBadRequest)
		return
	}

	var measurementRequest struct {
		MetricValue float64 `json:"metric_value"`
		Notes       string  `json:"notes"`
	}

	// Decode the request body
	err = json.NewDecoder(r.Body).Decode(&measurementRequest)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Find and update the measurement
	found := false
	for i := range measurements {
		if measurements[i].ID == id {
			measurements[i].MetricValue = measurementRequest.MetricValue
			measurements[i].Notes = measurementRequest.Notes
			found = true
			break
		}
	}

	if !found {
		http.Error(w, "Measurement not found", http.StatusNotFound)
		return
	}

	// Save to Excel
	err = saveToExcel()
	if err != nil {
		http.Error(w, "Failed to save to Excel: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the updated measurement
	for i := range measurements {
		if measurements[i].ID == id {
			json.NewEncoder(w).Encode(measurements[i])
			break
		}
	}
}

// getMeasurementsByKPI returns all measurements for a specific KPI
func getMeasurementsByKPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	kpiID, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid KPI ID", http.StatusBadRequest)
		return
	}

	// Get query parameters for filtering
	query := r.URL.Query()
	yearStr := query.Get("year")

	var filteredMeasurements []Measurement
	for _, m := range measurements {
		if m.KPIID == kpiID {
			// Filter by year if provided
			if yearStr != "" {
				year, err := strconv.Atoi(yearStr)
				if err != nil {
					http.Error(w, "Invalid year", http.StatusBadRequest)
					return
				}
				if m.Period.Year() != year {
					continue
				}
			}
			filteredMeasurements = append(filteredMeasurements, m)
		}
	}

	json.NewEncoder(w).Encode(filteredMeasurements)
}

// getMonthlyReport generates a monthly report
func getMonthlyReport(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	year, err := strconv.Atoi(params["year"])
	if err != nil || year < 2000 || year > 2100 {
		http.Error(w, "Invalid year", http.StatusBadRequest)
		return
	}

	month, err := strconv.Atoi(params["month"])
	if err != nil || month < 1 || month > 12 {
		http.Error(w, "Invalid month", http.StatusBadRequest)
		return
	}

	// Create the period
	period := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)

	// Get the report format from query params (default to JSON)
	format := r.URL.Query().Get("format")
	if format == "" {
		format = "json"
	}

	if format == "json" {
		w.Header().Set("Content-Type", "application/json")
		// For JSON, we'll create a structured report
		type KPIAchievement struct {
			KPI         KPI          `json:"kpi"`
			Measurement *Measurement `json:"measurement,omitempty"`
			Achievement float64      `json:"achievement_percent"`
			Score       float64      `json:"score"`
		}

		type RoleReport struct {
			Role         Role             `json:"role"`
			TotalScore   float64          `json:"total_score"`
			Achievements []KPIAchievement `json:"achievements"`
		}

		var report []RoleReport

		for _, role := range roles {
			roleKPIs := getKPIsByRoleID(role.ID)
			if len(roleKPIs) == 0 {
				continue
			}

			roleReport := RoleReport{
				Role:         role,
				Achievements: []KPIAchievement{},
			}

			for _, kpi := range roleKPIs {
				measurement := getExistingMeasurement(kpi.ID, period)

				achievementPct := 0.0
				score := 0.0

				if measurement != nil {
					achievementPct = calculateAchievement(kpi, measurement)
					score = achievementPct * kpi.Weight / 100
				}

				roleReport.Achievements = append(roleReport.Achievements, KPIAchievement{
					KPI:         kpi,
					Measurement: measurement,
					Achievement: achievementPct,
					Score:       score,
				})
			}

			roleReport.TotalScore = calculateOverallScore(roleKPIs, period)
			report = append(report, roleReport)
		}

		json.NewEncoder(w).Encode(report)
	} else {
		// For other formats, generate the report as before
		reportContent := generateReport(period, period, format)

		// Set appropriate content type
		switch format {
		case "txt":
			w.Header().Set("Content-Type", "text/plain")
		case "csv":
			w.Header().Set("Content-Type", "text/csv")
		case "html":
			w.Header().Set("Content-Type", "text/html")
		}

		fmt.Fprint(w, reportContent)
	}
}

// getQuarterlyReport generates a quarterly report
func getQuarterlyReport(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	year, err := strconv.Atoi(params["year"])
	if err != nil || year < 2000 || year > 2100 {
		http.Error(w, "Invalid year", http.StatusBadRequest)
		return
	}

	quarter, err := strconv.Atoi(params["quarter"])
	if err != nil || quarter < 1 || quarter > 4 {
		http.Error(w, "Invalid quarter", http.StatusBadRequest)
		return
	}

	// Calculate start and end periods
	startMonth := (quarter-1)*3 + 1
	endMonth := quarter * 3

	startPeriod := time.Date(year, time.Month(startMonth), 1, 0, 0, 0, 0, time.Local)
	endPeriod := time.Date(year, time.Month(endMonth), 1, 0, 0, 0, 0, time.Local)

	// Get the report format from query params (default to JSON)
	format := r.URL.Query().Get("format")
	if format == "" {
		format = "json"
	}

	if format == "json" {
		w.Header().Set("Content-Type", "application/json")
		// JSON structured report (simplified for brevity)
		response := struct {
			Year      int       `json:"year"`
			Quarter   int       `json:"quarter"`
			StartDate time.Time `json:"start_date"`
			EndDate   time.Time `json:"end_date"`
			// Would include role-based reports here in a real implementation
		}{
			Year:      year,
			Quarter:   quarter,
			StartDate: startPeriod,
			EndDate:   endPeriod,
		}

		json.NewEncoder(w).Encode(response)
	} else {
		// For other formats, generate the report as before
		reportContent := generateReport(startPeriod, endPeriod, format)

		// Set appropriate content type
		switch format {
		case "txt":
			w.Header().Set("Content-Type", "text/plain")
		case "csv":
			w.Header().Set("Content-Type", "text/csv")
		case "html":
			w.Header().Set("Content-Type", "text/html")
		}

		fmt.Fprint(w, reportContent)
	}
}

// getYearlyReport generates a yearly report
func getYearlyReport(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	year, err := strconv.Atoi(params["year"])
	if err != nil || year < 2000 || year > 2100 {
		http.Error(w, "Invalid year", http.StatusBadRequest)
		return
	}

	// Calculate start and end periods
	startPeriod := time.Date(year, 1, 1, 0, 0, 0, 0, time.Local)
	endPeriod := time.Date(year, 12, 1, 0, 0, 0, 0, time.Local)

	// Get the report format from query params (default to JSON)
	format := r.URL.Query().Get("format")
	if format == "" {
		format = "json"
	}

	if format == "json" {
		w.Header().Set("Content-Type", "application/json")
		// JSON structured report (simplified for brevity)
		response := struct {
			Year      int       `json:"year"`
			StartDate time.Time `json:"start_date"`
			EndDate   time.Time `json:"end_date"`
			// Would include role-based reports here in a real implementation
		}{
			Year:      year,
			StartDate: startPeriod,
			EndDate:   endPeriod,
		}

		json.NewEncoder(w).Encode(response)
	} else {
		// For other formats, generate the report as before
		reportContent := generateReport(startPeriod, endPeriod, format)

		// Set appropriate content type
		switch format {
		case "txt":
			w.Header().Set("Content-Type", "text/plain")
		case "csv":
			w.Header().Set("Content-Type", "text/csv")
		case "html":
			w.Header().Set("Content-Type", "text/html")
		}

		fmt.Fprint(w, reportContent)
	}
}

// getCustomReport generates a custom report based on provided parameters
func getCustomReport(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var customReportRequest struct {
		StartPeriod time.Time `json:"start_period"`
		EndPeriod   time.Time `json:"end_period"`
		RoleIDs     []int     `json:"role_ids,omitempty"` // Optional role IDs filter
		Format      string    `json:"format"`
	}

	// Decode the request body
	err := json.NewDecoder(r.Body).Decode(&customReportRequest)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate periods
	if customReportRequest.EndPeriod.Before(customReportRequest.StartPeriod) {
		http.Error(w, "End period cannot be before start period", http.StatusBadRequest)
		return
	}

	// If format is empty, default to JSON
	if customReportRequest.Format == "" {
		customReportRequest.Format = "json"
	}

	if customReportRequest.Format == "json" {
		// JSON structured report (simplified for brevity)
		response := struct {
			StartDate time.Time `json:"start_date"`
			EndDate   time.Time `json:"end_date"`
			RoleIDs   []int     `json:"role_ids,omitempty"`
			// Would include role-based reports here in a real implementation
		}{
			StartDate: customReportRequest.StartPeriod,
			EndDate:   customReportRequest.EndPeriod,
			RoleIDs:   customReportRequest.RoleIDs,
		}

		json.NewEncoder(w).Encode(response)
	} else {
		// For other formats, generate the report as before
		reportContent := generateReport(customReportRequest.StartPeriod, customReportRequest.EndPeriod, customReportRequest.Format)

		// Set appropriate content type
		switch customReportRequest.Format {
		case "txt":
			w.Header().Set("Content-Type", "text/plain")
		case "csv":
			w.Header().Set("Content-Type", "text/csv")
		case "html":
			w.Header().Set("Content-Type", "text/html")
		}

		fmt.Fprint(w, reportContent)
	}
}

// getDashboardOverview returns an overview of KPI achievements for the dashboard
func getDashboardOverview(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get query parameters
	query := r.URL.Query()
	yearStr := query.Get("year")
	monthStr := query.Get("month")

	// Default to current year and month if not provided
	year := time.Now().Year()
	month := int(time.Now().Month())

	if yearStr != "" {
		parsedYear, err := strconv.Atoi(yearStr)
		if err == nil {
			year = parsedYear
		}
	}

	if monthStr != "" {
		parsedMonth, err := strconv.Atoi(monthStr)
		if err == nil && parsedMonth >= 1 && parsedMonth <= 12 {
			month = parsedMonth
		}
	}

	// Create the period
	period := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)

	// Prepare dashboard overview
	type RoleOverview struct {
		RoleID     int     `json:"role_id"`
		RoleName   string  `json:"role_name"`
		TotalScore float64 `json:"total_score"`
		KPICount   int     `json:"kpi_count"`
		Measured   int     `json:"measured_kpis"`
	}

	var overview []RoleOverview

	for _, role := range roles {
		roleKPIs := getKPIsByRoleID(role.ID)
		if len(roleKPIs) == 0 {
			continue
		}

		measured := 0
		for _, kpi := range roleKPIs {
			measurement := getExistingMeasurement(kpi.ID, period)
			if measurement != nil {
				measured++
			}
		}

		totalScore := calculateOverallScore(roleKPIs, period)

		overview = append(overview, RoleOverview{
			RoleID:     role.ID,
			RoleName:   role.Name,
			TotalScore: totalScore,
			KPICount:   len(roleKPIs),
			Measured:   measured,
		})
	}

	// Return the overview
	json.NewEncoder(w).Encode(overview)
}

// getDashboardTrends returns trend data for the dashboard
func getDashboardTrends(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get query parameters
	query := r.URL.Query()
	yearStr := query.Get("year")

	// Default to current year if not provided
	year := time.Now().Year()

	if yearStr != "" {
		parsedYear, err := strconv.Atoi(yearStr)
		if err == nil {
			year = parsedYear
		}
	}

	// Prepare trends data
	type MonthlyTrend struct {
		Month      int                `json:"month"`
		MonthName  string             `json:"month_name"`
		RoleScores map[string]float64 `json:"role_scores"`
	}

	var trends []MonthlyTrend

	// For each month
	for month := 1; month <= 12; month++ {
		period := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)
		monthName := period.Format("January")

		trend := MonthlyTrend{
			Month:      month,
			MonthName:  monthName,
			RoleScores: make(map[string]float64),
		}

		// For each role
		for _, role := range roles {
			roleKPIs := getKPIsByRoleID(role.ID)
			if len(roleKPIs) == 0 {
				continue
			}

			score := calculateOverallScore(roleKPIs, period)
			trend.RoleScores[role.Name] = score
		}

		trends = append(trends, trend)
	}

	// Return the trends
	json.NewEncoder(w).Encode(trends)
}

// getSettings returns the application settings
func getSettings(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(appSettings)
}

// updateSettings updates the application settings
func updateSettings(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var newSettings Settings

	// Decode the request body
	err := json.NewDecoder(r.Body).Decode(&newSettings)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate settings
	if newSettings.DatabasePath == "" {
		http.Error(w, "Database path cannot be empty", http.StatusBadRequest)
		return
	}

	// Update settings
	appSettings.DatabasePath = newSettings.DatabasePath

	// Save settings
	saveSettings()

	// Return the updated settings
	json.NewEncoder(w).Encode(appSettings)
}

// reloadExcel reloads data from Excel
func reloadExcel(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	err := loadFromExcel()
	if err != nil {
		http.Error(w, "Failed to reload from Excel: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return success message
	response := struct {
		Message string `json:"message"`
		Success bool   `json:"success"`
	}{
		Message: "Data reloaded successfully from Excel",
		Success: true,
	}

	json.NewEncoder(w).Encode(response)
}

// saveExcelAPI saves data to Excel
func saveExcelAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	err := saveToExcel()
	if err != nil {
		http.Error(w, "Failed to save to Excel: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return success message
	response := struct {
		Message string `json:"message"`
		Success bool   `json:"success"`
	}{
		Message: "Data saved successfully to Excel",
		Success: true,
	}

	json.NewEncoder(w).Encode(response)
}
