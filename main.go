package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Setup signal handling for graceful shutdown
	setupSignalHandling()

	// Create or load database directory
	err := initDatabase()
	if err != nil {
		fmt.Printf("Error initializing database: %v\n", err)
		return
	}
	fmt.Println("Database directory initialized at:", appSettings.DatabasePath)

	// Initialize KPI data only if Excel database doesn't exist
	excelPath := getExcelDBPath()
	if _, err := os.Stat(excelPath); os.IsNotExist(err) {
		fmt.Println("Excel database not found, initializing default KPI data...")
		initializeKPIData()
	}

	// Initialize Excel database and load data
	err = initExcelDatabase()
	if err != nil {
		fmt.Printf("Error initializing Excel database: %v\n", err)
		return
	}
	fmt.Println("Excel database loaded successfully")

	// Start the REST API server in a separate goroutine
	go startRESTServer()
	fmt.Println("REST API server started on http://localhost:8080")

	// Main program loop
	scanner := bufio.NewScanner(os.Stdin)
	for {
		displayMainMenu()

		fmt.Print("Enter your choice (1-5): ")
		scanner.Scan()
		choice := scanner.Text()

		switch choice {
		case "1":
			handleInputKPI(scanner)
		case "2":
			handleViewKPIs(scanner)
		case "3":
			handleGenerateReport(scanner)
		case "4":
			handleSettings(scanner)
		case "5":
			fmt.Println("Exiting program...")
			err := saveToExcel()
			if err != nil {
				fmt.Printf("Error saving data to Excel: %v\n", err)
			}
			return
		default:
			fmt.Println("Invalid choice. Please try again.")
		}
	}
}

// setupSignalHandling sets up handlers for system signals to ensure data is saved
// on unexpected termination
func setupSignalHandling() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		fmt.Println("\nReceived termination signal. Saving data before exit...")
		err := saveToExcel()
		if err != nil {
			fmt.Printf("Error saving data to Excel: %v\n", err)
		}
		os.Exit(0)
	}()
}

func displayMainMenu() {
	fmt.Println("\n==== KPI Tracking System ====")
	fmt.Println("1. Input KPI Achievement")
	fmt.Println("2. View Current KPIs")
	fmt.Println("3. Generate Reports")
	fmt.Println("4. Settings")
	fmt.Println("5. Exit")
}
