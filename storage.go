package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Database paths
const (
	dbDir        = "./data"
	settingsFile = "settings.json"
)

var appSettings Settings

// initDatabase initializes the database directory
func initDatabase() error {
	// Create data directory if it doesn't exist
	if _, err := os.Stat(dbDir); os.IsNotExist(err) {
		err = os.MkdirAll(dbDir, 0755)
		if err != nil {
			return fmt.Errorf("failed to create data directory: %v", err)
		}
	}

	// Load settings
	loadSettings()

	return nil
}

// loadSettings loads application settings from file
func loadSettings() {
	settingsPath := filepath.Join(dbDir, settingsFile)

	// Default settings
	appSettings = Settings{
		DatabasePath: dbDir,
	}

	// Try to load existing settings
	data, err := os.ReadFile(settingsPath)
	if err == nil {
		err = json.Unmarshal(data, &appSettings)
		if err != nil {
			fmt.Printf("Error parsing settings file: %v\n", err)
			// Use default settings
		}
	}

	// Ensure the database directory exists
	if _, err := os.Stat(appSettings.DatabasePath); os.IsNotExist(err) {
		err = os.MkdirAll(appSettings.DatabasePath, 0755)
		if err != nil {
			fmt.Printf("Error creating database directory: %v\n", err)
		}
	}

	// Save settings (creates the file if it doesn't exist)
	saveSettings()
}

// saveSettings saves application settings to file
func saveSettings() {
	settingsPath := filepath.Join(dbDir, settingsFile)

	data, err := json.MarshalIndent(appSettings, "", "  ")
	if err != nil {
		fmt.Printf("Error serializing settings: %v\n", err)
		return
	}

	err = os.WriteFile(settingsPath, data, 0644)
	if err != nil {
		fmt.Printf("Error saving settings: %v\n", err)
	}
}
