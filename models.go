package main

import (
	"time"
)

// Role represents a job position with associated KPIs
type Role struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// KPI represents a Key Performance Indicator
type KPI struct {
	ID          int     `json:"id"`
	RoleID      int     `json:"role_id"`
	Category    string  `json:"category"` // "Quantitative" or "Qualitative"
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Metric      string  `json:"metric"`
	Unit        string  `json:"unit"` // The unit of measurement
	Target      string  `json:"target"`
	TargetValue float64 `json:"target_value"` // Numerical target value
	Operator    string  `json:"operator"`     // "≤", "≥", "=", etc.
	Weight      float64 `json:"weight"`       // In percentage
}

// Measurement represents an actual KPI measurement
type Measurement struct {
	ID          int       `json:"id"`
	KPIID       int       `json:"kpi_id"`
	MetricValue float64   `json:"metric_value"`
	Unit        string    `json:"unit"` // e.g., "days", "%", "count"
	Period      time.Time `json:"period"`
	Notes       string    `json:"notes"`
	CreatedAt   time.Time `json:"created_at"`
}

// Achievement represents the calculation of a KPI achievement
type Achievement struct {
	KPI         KPI
	Measurement Measurement
	Percentage  float64 // Achievement percentage
	Score       float64 // Achievement score (percentage * weight)
}

// Report represents a KPI report for a specific role and time period
type Report struct {
	Role         Role
	Period       string // e.g., "January 2025", "Q1 2025", "2025"
	Achievements []Achievement
	TotalScore   float64
	GeneratedAt  time.Time
}

// Settings for the application
type Settings struct {
	DatabasePath string `json:"database_path"`
	ExcelDBPath  string `json:"excel_db_path"`
}

// Global variables to store data
var roles []Role
var kpis []KPI
var measurements []Measurement
