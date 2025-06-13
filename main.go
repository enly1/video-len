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

func processVideoFiles(path string, recursive bool) error {
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("error accessing path %s: %v", path, err)
	}

	if info.IsDir() {
		// Process directory
		entries, err := os.ReadDir(path)
		if err != nil {
			return fmt.Errorf("error reading directory %s: %v", path, err)
		}

		for _, entry := range entries {
			if entry.IsDir() {
				if recursive {
					subPath := filepath.Join(path, entry.Name())
					if err := processVideoFiles(subPath, recursive); err != nil {
						fmt.Printf("Warning: %v\n", err)
					}
				}
				continue
			}

			ext := strings.ToLower(filepath.Ext(entry.Name()))
			if !videoExtensions[ext] {
				continue
			}

			filePath := filepath.Join(path, entry.Name())
			duration, err := getVideoDuration(filePath)
			if err != nil {
				fmt.Printf("Warning: %v\n", err)
				continue
			}

			formattedDuration := formatDuration(duration)
			cleanName := cleanFilename(strings.TrimSuffix(entry.Name(), ext))
			newName := fmt.Sprintf("%s_%s%s", cleanName, formattedDuration, ext)
			newPath := filepath.Join(path, newName)

			if err := os.Rename(filePath, newPath); err != nil {
				fmt.Printf("Error renaming %s: %v\n", filePath, err)
				continue
			}

			fmt.Printf("Renamed: %s -> %s\n", filePath, newPath)
		}
	} else {
		// Process single file
		ext := strings.ToLower(filepath.Ext(path))
		if !videoExtensions[ext] {
			return fmt.Errorf("file %s is not a supported video format", path)
		}

		duration, err := getVideoDuration(path)
		if err != nil {
			return fmt.Errorf("error getting duration for %s: %v", path, err)
		}

		formattedDuration := formatDuration(duration)
		cleanName := cleanFilename(strings.TrimSuffix(filepath.Base(path), ext))
		newName := fmt.Sprintf("%s_%s%s", cleanName, formattedDuration, ext)
		newPath := filepath.Join(filepath.Dir(path), newName)

		if err := os.Rename(path, newPath); err != nil {
			return fmt.Errorf("error renaming %s: %v", path, err)
		}

		fmt.Printf("Renamed: %s -> %s\n", path, newPath)
	}

	return nil
}

func printHelp() {
	fmt.Println("Usage: vidlen [file_or_directory_path]")
	fmt.Println("Options:")
	fmt.Println("  -h, --help    Display this help message")
	fmt.Println("  -R            Process directories recursively")
	fmt.Println("Description:")
	fmt.Println("  Renames video files by appending their duration in [XhYm] or [Ym] format.")
	fmt.Println("  If no path is provided, processes the current directory.")
	fmt.Println("  If a path is provided, processes either the named file or the directory contents.")
}

func main() {
	args := os.Args[1:]
	recursive := false

	for i, arg := range args {
		if arg == "-R" {
			recursive = true
			args = append(args[:i], args[i+1:]...)
			break
		}
	}

	if len(args) > 1 {
		fmt.Println("Usage: vidlen [file_or_directory_path]")
		os.Exit(1)
	}

	if len(args) == 1 {
		if args[0] == "-h" || args[0] == "--help" {
			printHelp()
			os.Exit(0)
		}
		path := args[0]
		if err := processVideoFiles(path, recursive); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	} else {
		if err := processVideoFiles(".", recursive); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	}
}
