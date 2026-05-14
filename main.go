package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

func main() {
	// 1. Define CLI flags
	setDest := flag.String("set-dest", "", "Set the destination folder for daily logs")
	mAgo := flag.Int("m-ago", 0, "Minutes ago to log the entry")
	hAgo := flag.Int("h-ago", 0, "Hours ago to log the entry")

	// Custom usage message
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: log [flags] <activity>\n\nFlags:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExample:\n  log --m-ago 30 went to gym\n  log --set-dest \"C:\\Users\\Name\\Documents\\Obsidian\\Daily\"\n")
	}

	flag.Parse()

	// 2. Resolve Config Path (~/.config/log/config.toml)
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting home directory: %v\n", err)
		os.Exit(1)
	}
	configDir := filepath.Join(homeDir, ".config", "log")
	configFile := filepath.Join(configDir, "config.toml")

	// 3. Handle --set-dest
	if *setDest != "" {
		err := os.MkdirAll(configDir, 0755)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create config dir: %v\n", err)
			os.Exit(1)
		}

		// Simple TOML format for our needs
		configContent := fmt.Sprintf("dest_folder = \"%s\"\n", strings.ReplaceAll(*setDest, "\\", "\\\\"))
		err = os.WriteFile(configFile, []byte(configContent), 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to write config: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Destination folder successfully set to: %s\n", *setDest)
		return
	}

	// 4. Ensure we have an activity to log
	args := flag.Args()
	if len(args) == 0 {
		fmt.Println("Error: No activity provided.")
		flag.Usage()
		os.Exit(1)
	}
	activity := strings.Join(args, " ")

	// 5. Read Destination Folder from Config
	destFolder := readConfig(configFile)
	if destFolder == "" {
		fmt.Println("Error: Destination folder is not set.")
		fmt.Println("Please run: log --set-dest \"C:\\Path\\To\\Your\\Folder\"")
		os.Exit(1)
	}

	// Ensure destination folder exists
	if err := os.MkdirAll(destFolder, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create destination directory: %v\n", err)
		os.Exit(1)
	}

	// 6. Calculate Time
	now := time.Now()
	// Subtract the provided hours and minutes
	logTime := now.Add(-time.Duration(*mAgo) * time.Minute).Add(-time.Duration(*hAgo) * time.Hour)
	
	// Formatting
	dateStr := logTime.Format("2006.01.02") // YYYY.MM.DD
	timeStr := logTime.Format("03:04 PM")   // 12-hour AM/PM
	fileName := fmt.Sprintf("%s.md", dateStr)
	filePath := filepath.Join(destFolder, fileName)

	newEntry := fmt.Sprintf("%s - %s", timeStr, activity)

	// 7. Read existing file and perform Binary Search
	lines := readLogFile(filePath)
	
	// Parse the target time so we can compare chronologically (string sorting fails on AM/PM)
	targetTime, _ := time.Parse("03:04 PM", timeStr)
	
	// Binary search to find the correct chronological insertion point
	// sort.Search returns the smallest index i where the function is true
	insertIdx := sort.Search(len(lines), func(i int) bool {
		line := lines[i]
		if len(line) >= 8 {
			lineTimeStr := line[:8]
			lineTimeParsed, err := time.Parse("03:04 PM", lineTimeStr)
			if err == nil {
				// Return true if the line's time is strictly greater than our target time
				return lineTimeParsed.After(targetTime)
			}
		}
		// If line doesn't start with a timestamp, treat it as something that should stay at the top or bottom depending on structure
		// For safety, assume non-timestamp lines are headers and we go after them
		return false 
	})

	// Insert the new entry at the found index
	lines = append(lines[:insertIdx], append([]string{newEntry}, lines[insertIdx:]...)...)

	// 8. Write back to file with exactly one newline between all lines
	output := strings.Join(lines, "\n")
	
	// Ensure the file ends with a trailing newline
	if !strings.HasSuffix(output, "\n") {
		output += "\n"
	}

	err = os.WriteFile(filePath, []byte(output), 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write to log file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Logged to %s: %s\n", fileName, newEntry)
}

// readConfig parses our simple TOML config to extract the destination folder
func readConfig(path string) string {
	content, err := os.ReadFile(path)
	if err != nil {
		return ""
	}

	scanner := bufio.NewScanner(bytes.NewReader(content))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "dest_folder") {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				// Strip surrounding spaces and quotes
				val := strings.TrimSpace(parts[1])
				val = strings.Trim(val, "\"")
				// Restore escaped windows slashes if present
				val = strings.ReplaceAll(val, "\\\\", "\\")
				return val
			}
		}
	}
	return ""
}

// readLogFile reads the daily file into a slice of strings, removing any empty lines
func readLogFile(path string) []string {
	content, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{} // Return empty slice if file doesn't exist yet
		}
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		return []string{}
	}

	// Normalize windows line endings just in case, then split
	normalized := strings.ReplaceAll(string(content), "\r\n", "\n")
	rawLines := strings.Split(normalized, "\n")
	
	var cleanLines []string
	for _, line := range rawLines {
		if strings.TrimSpace(line) != "" {
			cleanLines = append(cleanLines, line)
		}
	}
	
	return cleanLines
}
