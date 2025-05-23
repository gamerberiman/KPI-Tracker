package main

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// handleViewKPIs manages the KPI viewing flow
func handleViewKPIs(scanner *bufio.Scanner) {
	// View options
	fmt.Println("\n=== View KPIs ===")
	fmt.Println("1. View by Role")
	fmt.Println("2. View by Month")
	fmt.Println("3. View Year-to-Date")
	fmt.Println("4. View Trends")
	fmt.Println("0. Back to Main Menu")

	fmt.Print("\nEnter your choice: ")
	scanner.Scan()
	choice := scanner.Text()

	switch choice {
	case "1":
		viewByRole(scanner)
	case "2":
		viewByMonth(scanner)
	case "3":
		viewYearToDate(scanner)
	case "4":
		viewTrends(scanner)
	case "0":
		return
	default:
		fmt.Println("Invalid choice.")
	}
}

// viewByRole shows KPIs for a specific role
func viewByRole(scanner *bufio.Scanner) {
	role := selectRole(scanner)
	if role == nil {
		return
	}

	period := selectPeriod(scanner)
	if period.IsZero() {
		return
	}

	roleKPIs := getKPIsByRoleID(role.ID)
	if len(roleKPIs) == 0 {
		fmt.Println("No KPIs found for this role.")
		return
	}

	fmt.Printf("\n=== KPI Values for %s - %s ===\n\n",
		role.Name, period.Format("January 2006"))

	// Display KPIs by category
	fmt.Println("QUANTITATIVE KPIs (70% Weight)")
	fmt.Println("------------------------------")
	displayKPIsByCategory(roleKPIs, "Quantitative", period)

	fmt.Println("\nQUALITATIVE KPIs (30% Weight)")
	fmt.Println("------------------------------")
	displayKPIsByCategory(roleKPIs, "Qualitative", period)

	// Calculate and display overall score
	totalScore := calculateOverallScore(roleKPIs, period)
	fmt.Printf("\nOVERALL SCORE: %.2f%%\n", totalScore)

	// Wait for user to press enter
	fmt.Print("\nPress Enter to continue...")
	scanner.Scan()
}

// displayKPIsByCategory shows KPIs of a specific category
func displayKPIsByCategory(kpis []KPI, category string, period time.Time) {
	fmt.Printf("%-40s %-20s %-15s %-15s %-10s %-10s\n",
		"KPI Name", "Metric", "Target", "Actual", "Achievement", "Score")
	fmt.Println(strings.Repeat("-", 115))

	for _, kpi := range kpis {
		if kpi.Category == category {
			// Get the actual measurement
			measurement := getExistingMeasurement(kpi.ID, period)

			// Calculate achievement percentage
			achievementPct := calculateAchievement(kpi, measurement)

			// Calculate score (achievement * weight)
			score := achievementPct * kpi.Weight / 100

			// Display the KPI
			actualValue := "-"
			if measurement != nil {
				actualValue = fmt.Sprintf("%.2f %s", measurement.MetricValue, measurement.Unit)
			}

			fmt.Printf("%-40s %-20s %-15s %-15s %-10.2f%% %-10.2f\n",
				kpi.Name, kpi.Metric, kpi.Target, actualValue, achievementPct, score)
		}
	}
}

// viewByMonth shows all KPIs for a specific month
func viewByMonth(scanner *bufio.Scanner) {
	period := selectPeriod(scanner)
	if period.IsZero() {
		return
	}

	fmt.Printf("\n=== All KPIs for %s ===\n\n", period.Format("January 2006"))

	for _, role := range roles {
		fmt.Printf("== %s ==\n", role.Name)
		roleKPIs := getKPIsByRoleID(role.ID)

		if len(roleKPIs) == 0 {
			fmt.Println("No KPIs defined.")
			continue
		}

		totalScore := calculateOverallScore(roleKPIs, period)
		fmt.Printf("Overall Score: %.2f%%\n\n", totalScore)
	}

	// Wait for user to press enter
	fmt.Print("\nPress Enter to continue...")
	scanner.Scan()
}

// viewYearToDate shows year-to-date performance
func viewYearToDate(scanner *bufio.Scanner) {
	fmt.Print("Enter year: ")
	scanner.Scan()
	yearStr := scanner.Text()

	year, err := strconv.Atoi(yearStr)
	if err != nil || year < 2000 || year > 2100 {
		fmt.Println("Invalid year.")
		return
	}

	currentMonth := time.Now().Month()
	if year == time.Now().Year() {
		currentMonth = time.Now().Month()
	} else {
		currentMonth = 12
	}

	fmt.Printf("\n=== Year-to-Date KPI Performance for %d ===\n\n", year)

	for _, role := range roles {
		fmt.Printf("== %s ==\n", role.Name)
		roleKPIs := getKPIsByRoleID(role.ID)

		if len(roleKPIs) == 0 {
			fmt.Println("No KPIs defined.")
			continue
		}

		// Display monthly scores
		fmt.Printf("%-15s", "Month")
		for i := 1; i <= int(currentMonth); i++ {
			monthName := time.Month(i).String()[:3]
			fmt.Printf("%-8s", monthName)
		}
		fmt.Println("Average")

		fmt.Println(strings.Repeat("-", 15+8*int(currentMonth)+8))

		fmt.Printf("%-15s", "Score (%)")

		var totalYearScore float64
		var monthsWithData int

		for i := 1; i <= int(currentMonth); i++ {
			period := time.Date(year, time.Month(i), 1, 0, 0, 0, 0, time.Local)
			score := calculateOverallScore(roleKPIs, period)

			if score > 0 {
				fmt.Printf("%-8.2f", score)
				totalYearScore += score
				monthsWithData++
			} else {
				fmt.Printf("%-8s", "-")
			}
		}

		// Calculate average
		avgScore := 0.0
		if monthsWithData > 0 {
			avgScore = totalYearScore / float64(monthsWithData)
		}

		fmt.Printf("%-8.2f\n\n", avgScore)
	}

	// Wait for user to press enter
	fmt.Print("\nPress Enter to continue...")
	scanner.Scan()
}

// viewTrends shows KPI trends over time
func viewTrends(scanner *bufio.Scanner) {
	role := selectRole(scanner)
	if role == nil {
		return
	}

	fmt.Print("Enter year: ")
	scanner.Scan()
	yearStr := scanner.Text()

	year, err := strconv.Atoi(yearStr)
	if err != nil || year < 2000 || year > 2100 {
		fmt.Println("Invalid year.")
		return
	}

	fmt.Print("Enter KPI number (or 0 for overall score): ")
	scanner.Scan()
	kpiIdxStr := scanner.Text()

	roleKPIs := getKPIsByRoleID(role.ID)
	if len(roleKPIs) == 0 {
		fmt.Println("No KPIs found for this role.")
		return
	}

	// Display KPIs with their numbers for selection
	fmt.Println("\nAvailable KPIs:")
	for i, kpi := range roleKPIs {
		fmt.Printf("%d. %s\n", i+1, kpi.Name)
	}

	kpiIdx, err := strconv.Atoi(kpiIdxStr)
	if err != nil || kpiIdx < 0 || kpiIdx > len(roleKPIs) {
		fmt.Println("Invalid KPI number.")
		return
	}

	currentMonth := time.Now().Month()
	if year == time.Now().Year() {
		currentMonth = time.Now().Month()
	} else {
		currentMonth = 12
	}

	if kpiIdx == 0 {
		// Show overall score trend
		fmt.Printf("\n=== Overall Score Trend for %s in %d ===\n\n", role.Name, year)

		fmt.Printf("%-15s", "Month")
		for i := 1; i <= int(currentMonth); i++ {
			monthName := time.Month(i).String()[:3]
			fmt.Printf("%-8s", monthName)
		}
		fmt.Println()

		fmt.Println(strings.Repeat("-", 15+8*int(currentMonth)))

		fmt.Printf("%-15s", "Score (%)")

		for i := 1; i <= int(currentMonth); i++ {
			period := time.Date(year, time.Month(i), 1, 0, 0, 0, 0, time.Local)
			score := calculateOverallScore(roleKPIs, period)

			if score > 0 {
				fmt.Printf("%-8.2f", score)
			} else {
				fmt.Printf("%-8s", "-")
			}
		}

		fmt.Println("\n")

		// Display ASCII chart
		displayASCIIChart(roleKPIs, year, 0)
	} else {
		// Show specific KPI trend
		kpi := roleKPIs[kpiIdx-1]
		fmt.Printf("\n=== Trend for %s in %d ===\n\n", kpi.Name, year)

		fmt.Printf("%-15s", "Month")
		for i := 1; i <= int(currentMonth); i++ {
			monthName := time.Month(i).String()[:3]
			fmt.Printf("%-8s", monthName)
		}
		fmt.Println()

		fmt.Println(strings.Repeat("-", 15+8*int(currentMonth)))

		fmt.Printf("%-15s", "Achievement (%)")

		for i := 1; i <= int(currentMonth); i++ {
			period := time.Date(year, time.Month(i), 1, 0, 0, 0, 0, time.Local)
			measurement := getExistingMeasurement(kpi.ID, period)
			achievement := calculateAchievement(kpi, measurement)

			if achievement > 0 {
				fmt.Printf("%-8.2f", achievement)
			} else {
				fmt.Printf("%-8s", "-")
			}
		}

		fmt.Println("\n")

		// Display ASCII chart
		displayASCIIChart(roleKPIs, year, kpi.ID)
	}

	// Wait for user to press enter
	fmt.Print("\nPress Enter to continue...")
	scanner.Scan()
}

// displayASCIIChart shows a simple ASCII chart of KPI trends
func displayASCIIChart(kpis []KPI, year int, kpiID int) {
	// Simple ASCII chart representation
	const chartHeight = 10
	const chartSymbol = "*"

	// Get data points
	var dataPoints []float64
	currentMonth := time.Now().Month()
	if year != time.Now().Year() {
		currentMonth = 12
	}

	for i := 1; i <= int(currentMonth); i++ {
		period := time.Date(year, time.Month(i), 1, 0, 0, 0, 0, time.Local)

		if kpiID == 0 {
			// Overall score
			score := calculateOverallScore(kpis, period)
			dataPoints = append(dataPoints, score)
		} else {
			// Specific KPI
			var selectedKPI KPI
			for _, kpi := range kpis {
				if kpi.ID == kpiID {
					selectedKPI = kpi
					break
				}
			}

			measurement := getExistingMeasurement(kpiID, period)
			achievement := calculateAchievement(selectedKPI, measurement)
			dataPoints = append(dataPoints, achievement)
		}
	}

	// Find max value for scaling
	maxValue := 0.0
	for _, value := range dataPoints {
		if value > maxValue {
			maxValue = value
		}
	}

	if maxValue == 0 {
		fmt.Println("No data available for chart.")
		return
	}

	// Print chart (upside down, will flip when displaying)
	chart := make([][]string, chartHeight)
	for i := range chart {
		chart[i] = make([]string, len(dataPoints))
		for j := range chart[i] {
			chart[i][j] = " "
		}
	}

	// Plot data points
	for i, value := range dataPoints {
		if value > 0 {
			// Scale value to chart height
			scaledValue := int((value / maxValue) * float64(chartHeight-1))
			for j := 0; j <= scaledValue; j++ {
				chart[j][i] = chartSymbol
			}
		}
	}

	// Display chart (flipped vertically)
	fmt.Println("Chart:")
	fmt.Printf("%.0f%% +", maxValue)
	fmt.Println(strings.Repeat("-", len(dataPoints)*2))

	for i := chartHeight - 1; i >= 0; i-- {
		fmt.Print("|")
		for j := 0; j < len(dataPoints); j++ {
			fmt.Printf("%s ", chart[i][j])
		}
		fmt.Println()
	}

	fmt.Print("0%  +")
	fmt.Println(strings.Repeat("-", len(dataPoints)*2))

	fmt.Print("    ")
	for i := 1; i <= len(dataPoints); i++ {
		fmt.Printf("%s ", time.Month(i).String()[:1])
	}
	fmt.Println()
}

// getKPIsByRoleID returns KPIs for a specific role
func getKPIsByRoleID(roleID int) []KPI {
	var result []KPI
	for _, kpi := range kpis {
		if kpi.RoleID == roleID {
			result = append(result, kpi)
		}
	}
	return result
}

// calculateAchievement calculates achievement percentage for a KPI
func calculateAchievement(kpi KPI, measurement *Measurement) float64 {
	if measurement == nil {
		return 0
	}

	// Get the value
	value := measurement.MetricValue

	// Calculate achievement based on operator and target
	var achievement float64

	switch kpi.Operator {
	case "≤", "<=":
		// Lower is better (e.g., days to process, error rate)
		if value <= kpi.TargetValue {
			achievement = 100 // Full achievement
		} else {
			// Proportional achievement (inverse relationship)
			achievement = (kpi.TargetValue / value) * 100
		}
	case "≥", ">=":
		// Higher is better (e.g., completion rate, satisfaction score)
		achievement = (value / kpi.TargetValue) * 100
	case "=":
		// Exact match is best (e.g., 100% compliance)
		if value >= kpi.TargetValue {
			achievement = 100
		} else {
			achievement = (value / kpi.TargetValue) * 100
		}
	default:
		// Default behavior (assume higher is better)
		achievement = (value / kpi.TargetValue) * 100
	}

	// Cap at 100% for simplicity
	if achievement > 100 {
		achievement = 100
	}

	return achievement
}

// calculateOverallScore calculates the overall score for a set of KPIs
func calculateOverallScore(kpis []KPI, period time.Time) float64 {
	var totalScore float64
	var totalWeight float64

	for _, kpi := range kpis {
		measurement := getExistingMeasurement(kpi.ID, period)
		if measurement != nil {
			achievementPct := calculateAchievement(kpi, measurement)
			score := achievementPct * kpi.Weight / 100
			totalScore += score
			totalWeight += kpi.Weight
		}
	}

	if totalWeight == 0 {
		return 0
	}

	// Normalize to 100%
	return (totalScore / totalWeight) * 100
}
