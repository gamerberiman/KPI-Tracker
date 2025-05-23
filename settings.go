package main

import (
	"bufio"
	"fmt"
	"os"
)

// handleSettings manages the settings menu
func handleSettings(scanner *bufio.Scanner) {
	fmt.Println("\n=== Settings ===")
	fmt.Println("1. Change Database Path")
	fmt.Println("2. Reload Data from Excel")
	fmt.Println("3. Force Save Data to Excel")
	fmt.Println("0. Back to Main Menu")

	fmt.Print("\nEnter your choice: ")
	scanner.Scan()
	choice := scanner.Text()

	switch choice {
	case "1":
		changeDatabasePath(scanner)
	case "2":
		reloadFromExcel(scanner)
	case "3":
		forceSaveToExcel(scanner)
	case "0":
		return
	default:
		fmt.Println("Invalid choice.")
	}
}

// changeDatabasePath allows changing the database path
func changeDatabasePath(scanner *bufio.Scanner) {
	fmt.Printf("Current database path: %s\n", appSettings.DatabasePath)
	fmt.Print("Enter new database path (or Enter to keep current): ")
	scanner.Scan()
	newPath := scanner.Text()

	if newPath == "" {
		return
	}

	// Check if path exists
	if _, err := os.Stat(newPath); os.IsNotExist(err) {
		fmt.Print("Path does not exist. Create it? (y/n): ")
		scanner.Scan()
		confirm := scanner.Text()

		if confirm != "y" && confirm != "Y" {
			fmt.Println("Canceled.")
			return
		}

		// Create directory
		err = os.MkdirAll(newPath, 0755)
		if err != nil {
			fmt.Printf("Error creating directory: %v\n", err)
			return
		}
	}

	// Save current Excel file to new location if it exists
	oldExcelPath := getExcelDBPath()
	if _, err := os.Stat(oldExcelPath); err == nil {
		// Update path
		appSettings.DatabasePath = newPath

		// Create new Excel path
		newExcelPath := getExcelDBPath()

		// Copy Excel file to new location
		data, err := os.ReadFile(oldExcelPath)
		if err == nil {
			err = os.WriteFile(newExcelPath, data, 0644)
			if err == nil {
				fmt.Printf("Copied Excel database to new location: %s\n", newExcelPath)
			} else {
				fmt.Printf("Warning: Failed to copy Excel database: %v\n", err)
			}
		}
	} else {
		// Update path
		appSettings.DatabasePath = newPath
	}

	// Save settings
	saveSettings()

	fmt.Println("Database path updated.")
}

// reloadFromExcel reloads all data from Excel
func reloadFromExcel(scanner *bufio.Scanner) {
	fmt.Print("This will discard any unsaved changes. Continue? (y/n): ")
	scanner.Scan()
	confirm := scanner.Text()

	if confirm != "y" && confirm != "Y" {
		fmt.Println("Canceled.")
		return
	}

	err := loadFromExcel()
	if err != nil {
		fmt.Printf("Error reloading from Excel: %v\n", err)
		return
	}

	fmt.Println("Data reloaded successfully from Excel.")
}

// forceSaveToExcel forces saving to Excel
func forceSaveToExcel(scanner *bufio.Scanner) {
	err := saveToExcel()
	if err != nil {
		fmt.Printf("Error saving to Excel: %v\n", err)
		return
	}

	fmt.Println("Data saved successfully to Excel.")
}
