package style

import (
	"github.com/charmbracelet/lipgloss"
)

// Color definitions
var (
	// Primary colors
	PrimaryColor   = lipgloss.Color("#00A86B") // Green
	SecondaryColor = lipgloss.Color("#0066CC") // Blue
	
	// Status colors
	SuccessColor = lipgloss.Color("#00C851") // Success green
	ErrorColor   = lipgloss.Color("#FF4444") // Error red
	WarningColor = lipgloss.Color("#FFBB33") // Warning amber
	InfoColor    = lipgloss.Color("#33B5E5") // Info blue
	
	// Neutral colors
	MutedColor = lipgloss.Color("#999999") // Muted gray
	BoldColor  = lipgloss.Color("#FFFFFF") // White/bold
)

// Style definitions
var (
	// Header styles
	TitleStyle = lipgloss.NewStyle().
		Foreground(PrimaryColor).
		Bold(true).
		Margin(1, 0)
	
	SubtitleStyle = lipgloss.NewStyle().
		Foreground(SecondaryColor).
		Bold(true)
	
	// Message styles
	SuccessStyle = lipgloss.NewStyle().
		Foreground(SuccessColor).
		Bold(true)
	
	ErrorStyle = lipgloss.NewStyle().
		Foreground(ErrorColor).
		Bold(true)
	
	WarningStyle = lipgloss.NewStyle().
		Foreground(WarningColor).
		Bold(true)
	
	InfoStyle = lipgloss.NewStyle().
		Foreground(InfoColor)
	
	// File operation styles
	FileCreatedStyle = lipgloss.NewStyle().
		Foreground(SuccessColor)
	
	FileMergedStyle = lipgloss.NewStyle().
		Foreground(InfoColor)
	
	// Summary styles
	SummaryStyle = lipgloss.NewStyle().
		Foreground(PrimaryColor).
		Bold(true).
		Margin(1, 0).
		Padding(0, 1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(PrimaryColor)
	
	// Helper text styles
	MutedStyle = lipgloss.NewStyle().
		Foreground(MutedColor)
)

// Helper functions for common operations
func Success(text string) string {
	return SuccessStyle.Render("✓ " + text)
}

func Error(text string) string {
	return ErrorStyle.Render("✗ " + text)
}

func Warning(text string) string {
	return WarningStyle.Render("⚠ " + text)
}

func Info(text string) string {
	return InfoStyle.Render("ℹ " + text)
}

func FileCreated(filename string) string {
	return FileCreatedStyle.Render("✓ Created: " + filename)
}

func FileMerged(filename string) string {
	return FileMergedStyle.Render("📄 Merged: " + filename)
}

func Summary(text string) string {
	return SummaryStyle.Render(text)
}

func Muted(text string) string {
	return MutedStyle.Render(text)
}

func Title(text string) string {
	return TitleStyle.Render(text)
}

func Subtitle(text string) string {
	return SubtitleStyle.Render(text)
}