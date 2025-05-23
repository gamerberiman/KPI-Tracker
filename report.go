package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// handleGenerateReport manages the report generation flow
func handleGenerateReport(scanner *bufio.Scanner) {
	fmt.Println("\n=== Generate Reports ===")
	fmt.Println("1. Monthly Report")
	fmt.Println("2. Quarterly Report")
	fmt.Println("3. Yearly Report")
	fmt.Println("4. Custom Report")
	fmt.Println("5. Export Data to Excel") // This already happens automatically
	fmt.Println("0. Back to Main Menu")

	fmt.Print("\nEnter your choice: ")
	scanner.Scan()
	choice := scanner.Text()

	switch choice {
	case "1":
		generateMonthlyReport(scanner)
	case "2":
		generateQuarterlyReport(scanner)
	case "3":
		generateYearlyReport(scanner)
	case "4":
		generateCustomReport(scanner)
	case "5":
		exportToExcel(scanner)
	case "0":
		return
	default:
		fmt.Println("Invalid choice.")
	}
}

// generateMonthlyReport generates a monthly report
func generateMonthlyReport(scanner *bufio.Scanner) {
	period := selectPeriod(scanner)
	if period.IsZero() {
		return
	}

	// Select report format
	format := selectReportFormat(scanner)
	if format == "" {
		return
	}

	fmt.Printf("\nGenerating monthly report for %s...\n", period.Format("January 2006"))

	report := generateReport(period, period, format)

	// Save the report
	saveReport(report, fmt.Sprintf("Monthly_Report_%s.%s",
		period.Format("Jan2006"), format))
}

// generateQuarterlyReport generates a quarterly report
func generateQuarterlyReport(scanner *bufio.Scanner) {
	fmt.Print("Enter year: ")
	scanner.Scan()
	yearStr := scanner.Text()

	year, err := strconv.Atoi(yearStr)
	if err != nil || year < 2000 || year > 2100 {
		fmt.Println("Invalid year.")
		return
	}

	fmt.Print("Enter quarter (1-4): ")
	scanner.Scan()
	quarterStr := scanner.Text()

	quarter, err := strconv.Atoi(quarterStr)
	if err != nil || quarter < 1 || quarter > 4 {
		fmt.Println("Invalid quarter.")
		return
	}

	// Calculate start and end periods
	startMonth := (quarter-1)*3 + 1
	endMonth := quarter * 3

	startPeriod := time.Date(year, time.Month(startMonth), 1, 0, 0, 0, 0, time.Local)
	endPeriod := time.Date(year, time.Month(endMonth), 1, 0, 0, 0, 0, time.Local)

	// Select report format
	format := selectReportFormat(scanner)
	if format == "" {
		return
	}

	fmt.Printf("\nGenerating quarterly report for Q%d %d...\n", quarter, year)

	report := generateReport(startPeriod, endPeriod, format)

	// Save the report
	saveReport(report, fmt.Sprintf("Quarterly_Report_Q%d_%d.%s",
		quarter, year, format))
}

// generateYearlyReport generates a yearly report
func generateYearlyReport(scanner *bufio.Scanner) {
	fmt.Print("Enter year: ")
	scanner.Scan()
	yearStr := scanner.Text()

	year, err := strconv.Atoi(yearStr)
	if err != nil || year < 2000 || year > 2100 {
		fmt.Println("Invalid year.")
		return
	}

	// Calculate start and end periods
	startPeriod := time.Date(year, 1, 1, 0, 0, 0, 0, time.Local)
	endPeriod := time.Date(year, 12, 1, 0, 0, 0, 0, time.Local)

	// Select report format
	format := selectReportFormat(scanner)
	if format == "" {
		return
	}

	fmt.Printf("\nGenerating yearly report for %d...\n", year)

	report := generateReport(startPeriod, endPeriod, format)

	// Save the report
	saveReport(report, fmt.Sprintf("Yearly_Report_%d.%s", year, format))
}

// generateCustomReport generates a custom report
func generateCustomReport(scanner *bufio.Scanner) {
	fmt.Println("\n=== Custom Report ===")

	// Select start period
	fmt.Println("Select start period:")
	startPeriod := selectPeriod(scanner)
	if startPeriod.IsZero() {
		return
	}

	// Select end period
	fmt.Println("Select end period:")
	endPeriod := selectPeriod(scanner)
	if endPeriod.IsZero() {
		return
	}

	if endPeriod.Before(startPeriod) {
		fmt.Println("End period cannot be before start period.")
		return
	}

	// Select roles to include
	var selectedRoles []Role
	fmt.Println("\nSelect roles to include:")
	fmt.Println("1. All roles")
	fmt.Println("2. Specific roles")

	fmt.Print("Enter your choice: ")
	scanner.Scan()
	roleChoice := scanner.Text()

	if roleChoice == "1" {
		selectedRoles = roles
	} else if roleChoice == "2" {
		fmt.Println("\nSelect roles (comma-separated numbers, e.g., 1,3,5):")

		for i, role := range roles {
			fmt.Printf("%d. %s\n", i+1, role.Name)
		}

		fmt.Print("\nEnter role numbers: ")
		scanner.Scan()
		roleIdxsStr := scanner.Text()

		roleIdxs := strings.Split(roleIdxsStr, ",")
		for _, idxStr := range roleIdxs {
			idx, err := strconv.Atoi(strings.TrimSpace(idxStr))
			if err == nil && idx >= 1 && idx <= len(roles) {
				selectedRoles = append(selectedRoles, roles[idx-1])
			}
		}

		if len(selectedRoles) == 0 {
			fmt.Println("No valid roles selected.")
			return
		}
	} else {
		fmt.Println("Invalid choice.")
		return
	}

	// Select report format
	format := selectReportFormat(scanner)
	if format == "" {
		return
	}

	fmt.Printf("\nGenerating custom report from %s to %s...\n",
		startPeriod.Format("January 2006"), endPeriod.Format("January 2006"))

	// In a real implementation, this would generate a custom report
	// using the selected parameters

	// For demonstration, we'll just call the generic report generator
	report := generateReport(startPeriod, endPeriod, format)

	// Save the report
	saveReport(report, fmt.Sprintf("Custom_Report_%s_to_%s.%s",
		startPeriod.Format("Jan2006"), endPeriod.Format("Jan2006"), format))
}

// exportToExcel exports all KPI data to Excel
func exportToExcel(scanner *bufio.Scanner) {
	fmt.Println("\nExporting all data to Excel...")

	err := saveToExcel()
	if err != nil {
		fmt.Printf("Error exporting to Excel: %v\n", err)
		return
	}

	fmt.Printf("Data successfully exported to Excel: %s\n", getExcelDBPath())

	// Wait for user to press enter
	fmt.Print("\nPress Enter to continue...")
	scanner.Scan()
}

// selectReportFormat allows the user to select a report format
func selectReportFormat(scanner *bufio.Scanner) string {
	fmt.Println("\n=== Select Report Format ===")
	fmt.Println("1. Text (.txt)")
	fmt.Println("2. CSV (.csv)")
	fmt.Println("3. HTML (.html)")
	fmt.Println("0. Cancel")

	fmt.Print("\nEnter your choice: ")
	scanner.Scan()
	choice := scanner.Text()

	switch choice {
	case "1":
		return "txt"
	case "2":
		return "csv"
	case "3":
		return "html"
	case "0":
		return ""
	default:
		fmt.Println("Invalid choice.")
		return ""
	}
}

// generateReport generates a report for the given period range and format
func generateReport(startPeriod, endPeriod time.Time, format string) string {
	// In a real implementation, this would generate a report in the specified format
	// For demonstration, we'll return a placeholder
	fmt.Println("MASUK GENERATEREPORTTTTTTTTTTTTTTTTTTTT")
	var report string

	switch format {
	case "txt":
		report = generateTextReport(startPeriod, endPeriod)
	case "csv":
		report = generateCSVReport(startPeriod, endPeriod)
	case "html":
		fmt.Println("MASUK HTMLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLL")
		report = generateHTMLReport(startPeriod, endPeriod)
	default:
		report = "Unsupported format"
	}

	return report
}

// generateTextReport generates a text report
func generateTextReport(startPeriod, endPeriod time.Time) string {
	report := fmt.Sprintf("KPI REPORT: %s - %s\n",
		startPeriod.Format("January 2006"), endPeriod.Format("January 2006"))
	report += "==========================================\n\n"

	// Add report generation timestamp
	report += fmt.Sprintf("Generated: %s\n\n", time.Now().Format("2006-01-02 15:04:05"))

	// For each role
	for _, role := range roles {
		report += fmt.Sprintf("ROLE: %s\n", role.Name)
		report += "----------------------------------------\n\n"

		roleKPIs := getKPIsByRoleID(role.ID)
		if len(roleKPIs) == 0 {
			report += "No KPIs defined.\n\n"
			continue
		}

		// Quantitative KPIs
		report += "QUANTITATIVE KPIs (70% Weight)\n"
		report += "----------------------------------------\n"

		for _, kpi := range roleKPIs {
			if kpi.Category == "Quantitative" {
				report += fmt.Sprintf("KPI: %s\n", kpi.Name)
				report += fmt.Sprintf("Metric: %s\n", kpi.Metric)
				report += fmt.Sprintf("Target: %s (Weight: %.1f%%)\n", kpi.Target, kpi.Weight)

				// Loop through each month in the period range
				currentPeriod := startPeriod
				for currentPeriod.Before(endPeriod) || currentPeriod.Equal(endPeriod) {
					measurement := getExistingMeasurement(kpi.ID, currentPeriod)
					if measurement != nil {
						achievementPct := calculateAchievement(kpi, measurement)
						report += fmt.Sprintf("  %s: %.2f %s (%.2f%%)\n",
							currentPeriod.Format("Jan 2006"),
							measurement.MetricValue,
							measurement.Unit,
							achievementPct)
					}

					// Move to the next month
					currentPeriod = currentPeriod.AddDate(0, 1, 0)
				}

				report += "\n"
			}
		}

		// Qualitative KPIs
		report += "QUALITATIVE KPIs (30% Weight)\n"
		report += "----------------------------------------\n"

		for _, kpi := range roleKPIs {
			if kpi.Category == "Qualitative" {
				report += fmt.Sprintf("KPI: %s\n", kpi.Name)
				report += fmt.Sprintf("Metric: %s\n", kpi.Metric)
				report += fmt.Sprintf("Target: %s (Weight: %.1f%%)\n", kpi.Target, kpi.Weight)

				// Loop through each month in the period range
				currentPeriod := startPeriod
				for currentPeriod.Before(endPeriod) || currentPeriod.Equal(endPeriod) {
					measurement := getExistingMeasurement(kpi.ID, currentPeriod)
					if measurement != nil {
						achievementPct := calculateAchievement(kpi, measurement)
						report += fmt.Sprintf("  %s: %.2f %s (%.2f%%)\n",
							currentPeriod.Format("Jan 2006"),
							measurement.MetricValue,
							measurement.Unit,
							achievementPct)
					}

					// Move to the next month
					currentPeriod = currentPeriod.AddDate(0, 1, 0)
				}

				report += "\n"
			}
		}

		// Overall score
		report += "OVERALL SCORES\n"
		report += "----------------------------------------\n"

		// Loop through each month in the period range
		currentPeriod := startPeriod
		for currentPeriod.Before(endPeriod) || currentPeriod.Equal(endPeriod) {
			score := calculateOverallScore(roleKPIs, currentPeriod)
			if score > 0 {
				report += fmt.Sprintf("  %s: %.2f%%\n",
					currentPeriod.Format("Jan 2006"), score)
			}

			// Move to the next month
			currentPeriod = currentPeriod.AddDate(0, 1, 0)
		}

		report += "\n"
	}

	return report
}

// generateCSVReport generates a CSV report
func generateCSVReport(startPeriod, endPeriod time.Time) string {
	// This is a simplified CSV generator
	// In a real implementation, this would be more sophisticated

	report := "Role,KPI,Category,Weight,Target"

	// Add month columns
	currentPeriod := startPeriod
	for currentPeriod.Before(endPeriod) || currentPeriod.Equal(endPeriod) {
		report += fmt.Sprintf(",%s", currentPeriod.Format("Jan 2006"))
		currentPeriod = currentPeriod.AddDate(0, 1, 0)
	}

	report += "\n"

	// Add data rows
	for _, role := range roles {
		roleKPIs := getKPIsByRoleID(role.ID)

		for _, kpi := range roleKPIs {
			// Add KPI info
			report += fmt.Sprintf("%s,%s,%s,%.1f%%,%s",
				role.Name, kpi.Name, kpi.Category, kpi.Weight, kpi.Target)

			// Add monthly values
			currentPeriod = startPeriod
			for currentPeriod.Before(endPeriod) || currentPeriod.Equal(endPeriod) {
				measurement := getExistingMeasurement(kpi.ID, currentPeriod)
				if measurement != nil {
					achievementPct := calculateAchievement(kpi, measurement)
					report += fmt.Sprintf(",%.2f%%", achievementPct)
				} else {
					report += ","
				}

				currentPeriod = currentPeriod.AddDate(0, 1, 0)
			}

			report += "\n"
		}

		// Add overall score row
		report += fmt.Sprintf("%s,OVERALL SCORE,,,", role.Name)

		currentPeriod = startPeriod
		for currentPeriod.Before(endPeriod) || currentPeriod.Equal(endPeriod) {
			score := calculateOverallScore(roleKPIs, currentPeriod)
			if score > 0 {
				report += fmt.Sprintf(",%.2f%%", score)
			} else {
				report += ","
			}

			currentPeriod = currentPeriod.AddDate(0, 1, 0)
		}

		report += "\n"
	}

	return report
}

// generateHTMLReport generates an HTML report
func generateHTMLReport(startPeriod, endPeriod time.Time) string {
	// This is a simplified HTML generator
	// In a real implementation, this would be more sophisticated

	report := `<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8">
  <title>KPI Report</title>
  <style>
    body { font-family: Arial, sans-serif; margin: 20px; }
    h1, h2, h3 { color: #333; }
    table { border-collapse: collapse; width: 100%; margin-bottom: 20px; }
    th, td { border: 1px solid #ddd; padding: 8px; text-align: left; }
    th { background-color: #f2f2f2; }
    tr:nth-child(even) { background-color: #f9f9f9; }
    .good { color: green; }
    .warning { color: orange; }
    .bad { color: red; }
  </style>
</head>
<body>
  <h1>KPI Report</h1>
  <p>Period: `

	report += fmt.Sprintf("%s - %s</p>",
		startPeriod.Format("January 2006"), endPeriod.Format("January 2006"))

	report += fmt.Sprintf("<p>Generated: %s</p>", time.Now().Format("2006-01-02 15:04:05"))

	// For each role
	for _, role := range roles {
		report += fmt.Sprintf("<h2>%s</h2>", role.Name)

		roleKPIs := getKPIsByRoleID(role.ID)
		if len(roleKPIs) == 0 {
			report += "<p>No KPIs defined.</p>"
			continue
		}

		// Quantitative KPIs
		report += "<h3>Quantitative KPIs (70% Weight)</h3>"
		report += "<table><tr><th>KPI</th><th>Target</th><th>Weight</th>"

		// Add month columns
		currentPeriod := startPeriod
		for currentPeriod.Before(endPeriod) || currentPeriod.Equal(endPeriod) {
			report += fmt.Sprintf("<th>%s</th>", currentPeriod.Format("Jan 2006"))
			currentPeriod = currentPeriod.AddDate(0, 1, 0)
		}

		report += "</tr>"

		for _, kpi := range roleKPIs {
			if kpi.Category == "Quantitative" {
				report += fmt.Sprintf("<tr><td>%s</td><td>%s</td><td>%.1f%%</td>",
					kpi.Name, kpi.Target, kpi.Weight)

				// Add monthly values
				currentPeriod = startPeriod
				for currentPeriod.Before(endPeriod) || currentPeriod.Equal(endPeriod) {
					measurement := getExistingMeasurement(kpi.ID, currentPeriod)
					if measurement != nil {
						achievementPct := calculateAchievement(kpi, measurement)

						// Add color class based on achievement
						colorClass := "good"
						if achievementPct < 70 {
							colorClass = "bad"
						} else if achievementPct < 90 {
							colorClass = "warning"
						}

						report += fmt.Sprintf("<td class=\"%s\">%.2f%%</td>",
							colorClass, achievementPct)
					} else {
						report += "<td>-</td>"
					}

					currentPeriod = currentPeriod.AddDate(0, 1, 0)
				}

				report += "</tr>"
			}
		}

		report += "</table>"

		// Qualitative KPIs
		report += "<h3>Qualitative KPIs (30% Weight)</h3>"
		report += "<table><tr><th>KPI</th><th>Target</th><th>Weight</th>"

		// Add month columns
		currentPeriod = startPeriod
		for currentPeriod.Before(endPeriod) || currentPeriod.Equal(endPeriod) {
			report += fmt.Sprintf("<th>%s</th>", currentPeriod.Format("Jan 2006"))
			currentPeriod = currentPeriod.AddDate(0, 1, 0)
		}

		report += "</tr>"

		for _, kpi := range roleKPIs {
			if kpi.Category == "Qualitative" {
				report += fmt.Sprintf("<tr><td>%s</td><td>%s</td><td>%.1f%%</td>",
					kpi.Name, kpi.Target, kpi.Weight)

				// Add monthly values
				currentPeriod = startPeriod
				for currentPeriod.Before(endPeriod) || currentPeriod.Equal(endPeriod) {
					measurement := getExistingMeasurement(kpi.ID, currentPeriod)
					if measurement != nil {
						achievementPct := calculateAchievement(kpi, measurement)

						// Add color class based on achievement
						colorClass := "good"
						if achievementPct < 70 {
							colorClass = "bad"
						} else if achievementPct < 90 {
							colorClass = "warning"
						}

						report += fmt.Sprintf("<td class=\"%s\">%.2f%%</td>",
							colorClass, achievementPct)
					} else {
						report += "<td>-</td>"
					}

					currentPeriod = currentPeriod.AddDate(0, 1, 0)
				}

				report += "</tr>"
			}
		}

		report += "</table>"

		// Overall score
		report += "<h3>Overall Scores</h3>"
		report += "<table><tr><th>Period</th><th>Score</th></tr>"

		currentPeriod = startPeriod
		for currentPeriod.Before(endPeriod) || currentPeriod.Equal(endPeriod) {
			score := calculateOverallScore(roleKPIs, currentPeriod)
			if score > 0 {
				// Add color class based on score
				colorClass := "good"
				if score < 70 {
					colorClass = "bad"
				} else if score < 90 {
					colorClass = "warning"
				}

				report += fmt.Sprintf("<tr><td>%s</td><td class=\"%s\">%.2f%%</td></tr>",
					currentPeriod.Format("Jan 2006"), colorClass, score)
			}

			currentPeriod = currentPeriod.AddDate(0, 1, 0)
		}

		report += "</table>"
	}

	report += `</body>
</html>`

	return report
}

// saveReport saves a report to a file
func saveReport(report, filename string) {
	reportsDir := filepath.Join(appSettings.DatabasePath, "reports")

	// Create reports directory if it doesn't exist
	if _, err := os.Stat(reportsDir); os.IsNotExist(err) {
		err = os.MkdirAll(reportsDir, 0755)
		if err != nil {
			fmt.Printf("Error creating reports directory: %v\n", err)
			return
		}
	}

	// Save the report
	reportPath := filepath.Join(reportsDir, filename)
	err := os.WriteFile(reportPath, []byte(report), 0644)
	if err != nil {
		fmt.Printf("Error saving report: %v\n", err)
		return
	}

	fmt.Printf("Report saved to: %s\n", reportPath)
}
