package main

import (
	"fmt"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

var videoExtensions = map[string]bool{
	".mp4":  true,
	".mov":  true,
	".avi":  true,
	".mkv":  true,
	".wmv":  true,
	".flv":  true,
	".webm": true,
}

func getVideoDuration(filePath string) (float64, error) {
	cmd := exec.Command("ffprobe",
		"-v", "error",
		"-show_entries", "format=duration",
		"-of", "default=noprint_wrappers=1:nokey=1",
		filePath)

	output, err := cmd.Output()
	if err != nil {
		return 0, fmt.Errorf("error getting duration for %s: %v", filePath, err)
	}

	duration, err := strconv.ParseFloat(strings.TrimSpace(string(output)), 64)
	if err != nil {
		return 0, fmt.Errorf("error parsing duration for %s: %v", filePath, err)
	}

	return duration, nil
}

func formatDuration(duration float64) string {
	totalMinutes := int(math.Ceil(duration / 60))
	hours := totalMinutes / 60
	minutes := totalMinutes % 60

	if hours > 0 {
		return fmt.Sprintf("[%dh%dm]", hours, minutes)
	}
	return fmt.Sprintf("[%dm]", minutes)
}

func cleanFilename(filename string) string {
	// Remove content within square brackets
	re := regexp.MustCompile(`\[.*?\]`)
	cleaned := re.ReplaceAllString(filename, "")
	// Remove extra underscores and spaces
	return strings.Trim(cleaned, " _")
}

func processVideoFiles() error {
	entries, err := os.ReadDir(".")
	if err != nil {
		return fmt.Errorf("error reading directory: %v", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		ext := strings.ToLower(filepath.Ext(entry.Name()))
		if !videoExtensions[ext] {
			continue
		}

		duration, err := getVideoDuration(entry.Name())
		if err != nil {
			fmt.Printf("Warning: %v\n", err)
			continue
		}

		formattedDuration := formatDuration(duration)
		cleanName := cleanFilename(strings.TrimSuffix(entry.Name(), ext))
		newName := fmt.Sprintf("%s_%s%s", cleanName, formattedDuration, ext)

		if err := os.Rename(entry.Name(), newName); err != nil {
			fmt.Printf("Error renaming %s: %v\n", entry.Name(), err)
			continue
		}

		fmt.Printf("Renamed: %s -> %s\n", entry.Name(), newName)
	}

	return nil
}

func main() {
	if err := processVideoFiles(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
} 