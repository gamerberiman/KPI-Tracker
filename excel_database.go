package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/xuri/excelize/v2"
)

const (
	excelDBFilename   = "kpi_database.xlsx"
	rolesSheet        = "Roles"
	kpisSheet         = "KPIs"
	measurementsSheet = "Measurements"
)

// getExcelDBPath returns the full path to the Excel database file
func getExcelDBPath() string {
	return filepath.Join(appSettings.DatabasePath, excelDBFilename)
}

// initExcelDatabase creates or loads the Excel database
func initExcelDatabase() error {
	excelPath := getExcelDBPath()

	// Check if Excel file exists
	_, err := os.Stat(excelPath)
	if os.IsNotExist(err) {
		// Create new Excel database file with empty sheets
		fmt.Println("Excel database not found. Creating new database...")
		return createEmptyExcelDB()
	} else if err != nil {
		return fmt.Errorf("error checking Excel database: %v", err)
	}

	// Load data from Excel
	fmt.Println("Loading data from Excel database...")
	return loadFromExcel()
}

// createEmptyExcelDB creates a new Excel database with empty sheets
func createEmptyExcelDB() error {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	// Create sheets
	f.NewSheet(rolesSheet)
	f.NewSheet(kpisSheet)
	f.NewSheet(measurementsSheet)
	f.DeleteSheet("Sheet1")

	// Set up headers for Roles sheet
	f.SetSheetRow(rolesSheet, "A1", &[]interface{}{"ID", "Name", "Description"})

	// Set up headers for KPIs sheet
	f.SetSheetRow(kpisSheet, "A1", &[]interface{}{
		"ID", "RoleID", "Category", "Name", "Description",
		"Metric", "Unit", "Target", "TargetValue", "Operator", "Weight",
	})

	// Set up headers for Measurements sheet
	f.SetSheetRow(measurementsSheet, "A1", &[]interface{}{
		"ID", "KPIID", "MetricValue", "Unit", "Period", "Notes", "CreatedAt",
	})

	// Format headers as tables
	formatAsTable(f, rolesSheet, 1, 3)
	formatAsTable(f, kpisSheet, 1, 11)
	formatAsTable(f, measurementsSheet, 1, 7)

	// Save the Excel file
	excelPath := getExcelDBPath()
	if err := f.SaveAs(excelPath); err != nil {
		return fmt.Errorf("failed to create Excel database: %v", err)
	}

	fmt.Printf("Created new Excel database at: %s\n", excelPath)
	return nil
}

// loadFromExcel loads all data from the Excel database
func loadFromExcel() error {
	excelPath := getExcelDBPath()

	fmt.Println("Trying to open Excel file at:", excelPath)
	f, err := excelize.OpenFile(excelPath)
	if err != nil {
		// If we can't open the file, try to create a new one
		if os.IsNotExist(err) {
			fmt.Println("Excel file doesn't exist, creating a new one...")
			err = createEmptyExcelDB()
			if err != nil {
				return fmt.Errorf("failed to create new Excel database: %v", err)
			}
			// Try opening again
			f, err = excelize.OpenFile(excelPath)
			if err != nil {
				return fmt.Errorf("failed to open newly created Excel database: %v", err)
			}
		} else {
			return fmt.Errorf("failed to open Excel database: %v", err)
		}
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	// Write Excel file info to a debug file
	debugData := map[string]interface{}{
		"filePath":  excelPath,
		"fileInfo":  fmt.Sprintf("%+v", f),
		"sheetList": f.GetSheetList(),
	}
	debugJSON, _ := json.MarshalIndent(debugData, "", "  ")
	os.WriteFile(filepath.Join(appSettings.DatabasePath, "excel_debug.json"), debugJSON, 0644)

	// Clear existing data
	roles = []Role{}
	kpis = []KPI{}
	measurements = []Measurement{}

	// Load roles
	rows, err := f.GetRows(rolesSheet)
	if err != nil {
		return fmt.Errorf("failed to read roles sheet: %v", err)
	}

	for i, row := range rows {
		if i == 0 { // Skip header row
			continue
		}
		if len(row) < 3 {
			continue // Skip incomplete rows
		}

		id, err := strconv.Atoi(row[0])
		if err != nil {
			fmt.Printf("Warning: Invalid role ID '%s' in row %d, skipping\n", row[0], i+1)
			continue
		}

		role := Role{
			ID:          id,
			Name:        row[1],
			Description: row[2],
		}

		roles = append(roles, role)
	}

	// Load KPIs
	rows, err = f.GetRows(kpisSheet)
	if err != nil {
		return fmt.Errorf("failed to read KPIs sheet: %v", err)
	}

	for i, row := range rows {
		if i == 0 { // Skip header row
			continue
		}
		if len(row) < 11 {
			continue // Skip incomplete rows
		}

		id, err := strconv.Atoi(row[0])
		if err != nil {
			fmt.Printf("Warning: Invalid KPI ID '%s' in row %d, skipping\n", row[0], i+1)
			continue
		}

		roleID, err := strconv.Atoi(row[1])
		if err != nil {
			fmt.Printf("Warning: Invalid Role ID '%s' in row %d, skipping\n", row[1], i+1)
			continue
		}

		targetValue, err := strconv.ParseFloat(row[8], 64)
		if err != nil {
			targetValue = 0
		}

		weight, err := strconv.ParseFloat(row[10], 64)
		if err != nil {
			weight = 0
		}

		kpi := KPI{
			ID:          id,
			RoleID:      roleID,
			Category:    row[2],
			Name:        row[3],
			Description: row[4],
			Metric:      row[5],
			Unit:        row[6],
			Target:      row[7],
			TargetValue: targetValue,
			Operator:    row[9],
			Weight:      weight,
		}

		kpis = append(kpis, kpi)
	}

	// Load measurements
	rows, err = f.GetRows(measurementsSheet)
	if err != nil {
		return fmt.Errorf("failed to read measurements sheet: %v", err)
	}

	for i, row := range rows {
		if i == 0 { // Skip header row
			continue
		}
		if len(row) < 7 {
			continue // Skip incomplete rows
		}

		id, err := strconv.Atoi(row[0])
		if err != nil {
			fmt.Printf("Warning: Invalid measurement ID '%s' in row %d, skipping\n", row[0], i+1)
			continue
		}

		kpiID, err := strconv.Atoi(row[1])
		if err != nil {
			fmt.Printf("Warning: Invalid KPI ID '%s' in row %d, skipping\n", row[1], i+1)
			continue
		}

		metricValue, err := strconv.ParseFloat(row[2], 64)
		if err != nil {
			fmt.Printf("Warning: Invalid metric value '%s' in row %d, skipping\n", row[2], i+1)
			continue
		}

		// Parse period date
		period, err := time.Parse("2006-01-02", row[4])
		if err != nil {
			fmt.Printf("Warning: Invalid period '%s' in row %d, skipping\n", row[4], i+1)
			continue
		}

		// Parse created date
		createdAt, err := time.Parse("2006-01-02 15:04:05", row[6])
		if err != nil {
			createdAt = time.Now() // Use current time if parsing fails
		}

		measurement := Measurement{
			ID:          id,
			KPIID:       kpiID,
			MetricValue: metricValue,
			Unit:        row[3],
			Period:      period,
			Notes:       row[5],
			CreatedAt:   createdAt,
		}

		measurements = append(measurements, measurement)
	}

	fmt.Printf("Loaded %d roles, %d KPIs, and %d measurements from Excel database\n",
		len(roles), len(kpis), len(measurements))
	return nil
}

// saveToExcel saves all data to the Excel database
func saveToExcel() error {
	fmt.Println("Saving data to Excel...")
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	// Create sheets
	f.NewSheet(rolesSheet)
	f.NewSheet(kpisSheet)
	f.NewSheet(measurementsSheet)
	f.DeleteSheet("Sheet1")

	// Save Roles
	f.SetSheetRow(rolesSheet, "A1", &[]interface{}{"ID", "Name", "Description"})
	for i, role := range roles {
		row := fmt.Sprintf("A%d", i+2)
		f.SetSheetRow(rolesSheet, row, &[]interface{}{role.ID, role.Name, role.Description})
	}

	// Save KPIs
	f.SetSheetRow(kpisSheet, "A1", &[]interface{}{
		"ID", "RoleID", "Category", "Name", "Description",
		"Metric", "Unit", "Target", "TargetValue", "Operator", "Weight",
	})
	for i, kpi := range kpis {
		row := fmt.Sprintf("A%d", i+2)
		f.SetSheetRow(kpisSheet, row, &[]interface{}{
			kpi.ID, kpi.RoleID, kpi.Category, kpi.Name, kpi.Description,
			kpi.Metric, kpi.Unit, kpi.Target, kpi.TargetValue, kpi.Operator, kpi.Weight,
		})
	}

	// Save Measurements
	f.SetSheetRow(measurementsSheet, "A1", &[]interface{}{
		"ID", "KPIID", "MetricValue", "Unit", "Period", "Notes", "CreatedAt",
	})
	for i, m := range measurements {
		row := fmt.Sprintf("A%d", i+2)
		f.SetSheetRow(measurementsSheet, row, &[]interface{}{
			m.ID, m.KPIID, m.MetricValue, m.Unit,
			m.Period.Format("2006-01-02"), m.Notes, m.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	// Format as tables for better viewing
	formatAsTable(f, rolesSheet, len(roles)+1, 3)
	formatAsTable(f, kpisSheet, len(kpis)+1, 11)
	formatAsTable(f, measurementsSheet, len(measurements)+1, 7)

	// Save the Excel file
	excelPath := getExcelDBPath()

	// Make sure the directory exists
	dirPath := filepath.Dir(excelPath)
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			return fmt.Errorf("failed to create directory for Excel database: %v", err)
		}
	}

	// Create a backup of the existing Excel file if it exists
	if _, err := os.Stat(excelPath); err == nil {
		backupDir := filepath.Join(appSettings.DatabasePath, "excel_backups")
		if _, err := os.Stat(backupDir); os.IsNotExist(err) {
			os.MkdirAll(backupDir, 0755)
		}

		timestamp := time.Now().Format("20060102_150405")
		backupPath := filepath.Join(backupDir, fmt.Sprintf("kpi_database_backup_%s.xlsx", timestamp))

		// Copy the file
		data, err := os.ReadFile(excelPath)
		if err == nil {
			err = os.WriteFile(backupPath, data, 0644)
			if err == nil {
				fmt.Printf("Created backup of Excel database: %s\n", backupPath)
			}
		}
	}

	// Save the new version
	if err := f.SaveAs(excelPath); err != nil {
		return fmt.Errorf("failed to save Excel database: %v", err)
	}

	fmt.Printf("Saved data to Excel database: %s\n", excelPath)
	return nil
}

// formatAsTable formats a sheet as a table for better viewing
func formatAsTable(f *excelize.File, sheet string, rows, cols int) {
	// Set header style
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
			Size: 12,
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#C6EFCE"},
			Pattern: 1,
		},
		Border: []excelize.Border{
			{Type: "bottom", Color: "#000000", Style: 2},
		},
	})

	// Apply header style to the first row
	for col := 0; col < cols; col++ {
		colLetter := columnToLetter(col)
		f.SetCellStyle(sheet, fmt.Sprintf("%s1", colLetter), fmt.Sprintf("%s1", colLetter), headerStyle)
	}

	// Set data rows style (alternating colors)
	evenRowStyle, _ := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#F5F5F5"},
			Pattern: 1,
		},
	})

	// Apply alternating styles to data rows
	for row := 2; row <= rows; row++ {
		if row%2 == 0 { // Even rows
			f.SetCellStyle(sheet, fmt.Sprintf("A%d", row), fmt.Sprintf("%s%d", columnToLetter(cols-1), row), evenRowStyle)
		}
	}

	// Auto-fit columns if possible
	for col := 0; col < cols; col++ {
		colLetter := columnToLetter(col)
		f.SetColWidth(sheet, colLetter, colLetter, 20) // Set a reasonable default width
	}
}

// columnToLetter converts a column number to Excel column letter(s)
func columnToLetter(col int) string {
	if col < 26 {
		return string(rune('A' + col))
	}
	return string(rune('A'+col/26-1)) + string(rune('A'+col%26))
}
