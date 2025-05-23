package main

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// handleInputKPI manages the KPI input flow
func handleInputKPI(scanner *bufio.Scanner) {
	// Select role
	selectedRole := selectRole(scanner)
	if selectedRole == nil {
		return
	}

	// Select time period
	period := selectPeriod(scanner)
	if period.IsZero() {
		return
	}

	// Get KPIs for the selected role
	roleKPIs := getKPIsByRoleID(selectedRole.ID)
	if len(roleKPIs) == 0 {
		fmt.Println("No KPIs found for this role.")
		return
	}

	// Input values for each KPI
	fmt.Printf("\n=== Entering KPI values for %s - %s ===\n",
		selectedRole.Name, period.Format("January 2006"))

	// First handle quantitative KPIs
	fmt.Println("\nQUANTITATIVE KPIs (70% Weight)")
	fmt.Println("------------------------------")

	for _, kpi := range roleKPIs {
		if kpi.Category == "Quantitative" {
			inputKPIValue(scanner, kpi, period)
		}
	}

	// Then handle qualitative KPIs
	fmt.Println("\nQUALITATIVE KPIs (30% Weight)")
	fmt.Println("------------------------------")

	for _, kpi := range roleKPIs {
		if kpi.Category == "Qualitative" {
			inputKPIValue(scanner, kpi, period)
		}
	}

	fmt.Println("\nAll KPI values saved successfully!")

	// Save to Excel after all inputs
	err := saveToExcel()
	if err != nil {
		fmt.Printf("Warning: Failed to save to Excel: %v\n", err)
	}
}

// selectRole allows user to select a role
func selectRole(scanner *bufio.Scanner) *Role {
	fmt.Println("\n=== Select Role ===")

	for i, role := range roles {
		fmt.Printf("%d. %s\n", i+1, role.Name)
	}

	fmt.Print("\nEnter role number (0 to cancel): ")
	scanner.Scan()
	choice := scanner.Text()

	index, err := strconv.Atoi(choice)
	if err != nil || index < 0 || index > len(roles) {
		if choice != "0" {
			fmt.Println("Invalid selection.")
		}
		return nil
	}

	if index == 0 {
		return nil
	}

	return &roles[index-1]
}

// selectPeriod allows user to select a time period
func selectPeriod(scanner *bufio.Scanner) time.Time {
	fmt.Println("\n=== Select Period ===")

	currentYear := time.Now().Year()
	fmt.Printf("Year (default: %d): ", currentYear)
	scanner.Scan()
	yearStr := scanner.Text()

	year := currentYear
	if yearStr != "" {
		var err error
		year, err = strconv.Atoi(yearStr)
		if err != nil || year < 2000 || year > 2100 {
			fmt.Println("Invalid year. Using current year.")
			year = currentYear
		}
	}

	fmt.Print("Month (1-12): ")
	scanner.Scan()
	monthStr := scanner.Text()

	month, err := strconv.Atoi(monthStr)
	if err != nil || month < 1 || month > 12 {
		fmt.Println("Invalid month.")
		return time.Time{}
	}

	// Return the first day of the selected month
	return time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)
}

// inputKPIValue handles input for a specific KPI
func inputKPIValue(scanner *bufio.Scanner, kpi KPI, period time.Time) {
	fmt.Printf("\n%s\n", kpi.Name)
	fmt.Printf("Description: %s\n", kpi.Description)
	fmt.Printf("Metric: %s\n", kpi.Metric)
	fmt.Printf("Target: %s (Weight: %.1f%%)\n", kpi.Target, kpi.Weight)

	// Check if there's an existing measurement for this period
	existingMeasurement := getExistingMeasurement(kpi.ID, period)
	if existingMeasurement != nil {
		fmt.Printf("Current value: %.2f %s\n", existingMeasurement.MetricValue, existingMeasurement.Unit)
	}

	// Provide guidance on the expected input
	var unitText string
	switch kpi.Unit {
	case "days":
		unitText = "days"
	case "%":
		unitText = "percentage"
	case "score":
		unitText = "score (0-10)"
	default:
		unitText = kpi.Unit
	}

	fmt.Printf("Enter value (%s) or press enter to skip: ", unitText)
	scanner.Scan()
	valueStr := scanner.Text()

	if valueStr == "" {
		return
	}

	value, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		fmt.Println("Invalid number. Skipping.")
		return
	}

	// Input validation based on metric type
	valid := true
	switch kpi.Unit {
	case "%":
		if value < 0 || value > 100 {
			fmt.Println("Percentage must be between 0 and 100. Skipping.")
			valid = false
		}
	case "score":
		if value < 0 || value > 10 {
			fmt.Println("Score must be between 0 and 10. Skipping.")
			valid = false
		}
	case "days":
		if value < 0 {
			fmt.Println("Days cannot be negative. Skipping.")
			valid = false
		}
	default:
		// For other units, just ensure it's not negative
		if value < 0 && !strings.Contains(kpi.Metric, "error") { // Errors can be negative
			fmt.Println("Value cannot be negative. Skipping.")
			valid = false
		}
	}

	if !valid {
		return
	}

	fmt.Print("Enter notes (optional): ")
	scanner.Scan()
	notes := scanner.Text()

	// Save the measurement
	saveMeasurement(kpi.ID, value, kpi.Unit, period, notes)
}

// getExistingMeasurement retrieves an existing measurement for a KPI and period
func getExistingMeasurement(kpiID int, period time.Time) *Measurement {
	for i, m := range measurements {
		// Check if the measurement is for the same KPI and same month/year
		if m.KPIID == kpiID &&
			m.Period.Year() == period.Year() &&
			m.Period.Month() == period.Month() {
			return &measurements[i]
		}
	}
	return nil
}

// saveMeasurement saves a new KPI measurement
func saveMeasurement(kpiID int, value float64, unit string, period time.Time, notes string) {
	// Check if measurement already exists
	existingMeasurement := getExistingMeasurement(kpiID, period)

	if existingMeasurement != nil {
		// Update existing measurement
		existingMeasurement.MetricValue = value
		existingMeasurement.Unit = unit
		existingMeasurement.Notes = notes
		fmt.Printf("Updated measurement: KPI ID %d, Value %.2f %s\n", kpiID, value, unit)
	} else {
		// Create new measurement
		newID := 1
		if len(measurements) > 0 {
			// Find the highest ID and increment
			for _, m := range measurements {
				if m.ID >= newID {
					newID = m.ID + 1
				}
			}
		}

		measurement := Measurement{
			ID:          newID,
			KPIID:       kpiID,
			MetricValue: value,
			Unit:        unit,
			Period:      period,
			Notes:       notes,
			CreatedAt:   time.Now(),
		}

		// Add to measurements slice
		measurements = append(measurements, measurement)
		fmt.Printf("Saved new measurement: KPI ID %d, Value %.2f %s\n", kpiID, value, unit)
	}
}
