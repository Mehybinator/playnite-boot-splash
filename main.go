package main

import (
	_ "embed"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
)

//go:embed splash.mp4
var videoContent []byte

func main() {
	// Initialize log file
	logFile, err := initLogFile("logFile.log")
	if err != nil {
		handleError(fmt.Errorf("failed to initialize log file: %w", err))
	}
	defer logFile.Close()

	// Check if VLC is available
	if _, err := exec.LookPath("vlc"); err != nil {
		handleError(fmt.Errorf("VLC not found in PATH. Ensure VLC is installed and accessible: %w", err))
	}

	// Get current user
	currentUser, err := user.Current()
	if err != nil {
		handleError(fmt.Errorf("failed to retrieve current user: %w", err))
	}

	// Create a temporary video file for the splash
	videoPath, err := extractSplashVideo()
	if err != nil {
		handleError(fmt.Errorf("failed to create temporary splash video file: %w", err))
	}

	// Prepare VLC command
	playVideoCmd := prepareVLCCommand(videoPath)

	// Prepare Playnite command
	playniteExe := filepath.Join(currentUser.HomeDir, "AppData", "Local", "Playnite", "Playnite.fullscreenapp.exe")
	launchPlayniteCmd := exec.Command(playniteExe, "--hidesplashscreen")

	// Start VLC to play the video
	log.Println("Starting VLC to play splash video...")
	if err := playVideoCmd.Start(); err != nil {
		handleError(fmt.Errorf("failed to start VLC: %w", err))
	}

	// Start Playnite
	log.Println("Launching Playnite...")
	if err := launchPlayniteCmd.Start(); err != nil {
		handleError(fmt.Errorf("failed to start Playnite: %w", err))
	}
}

// initLogFile initializes the log file for logging application activity.
func initLogFile(fileName string) (*os.File, error) {
	logFile, err := os.OpenFile(fileName, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	log.SetOutput(logFile)
	log.Println("Log file initialized")
	return logFile, nil
}

// handleError logs an error message and exits the application.
func handleError(err error) {
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
}

// extractSplashVideo creates a temporary file for the splash video and writes its content.
func extractSplashVideo() (string, error) {
	tempDir, err := os.MkdirTemp("", "splash")
	if err != nil {
		return "", fmt.Errorf("failed to create temporary directory: %w", err)
	}

	splashPath := filepath.Join(tempDir, "splash.mp4")
	if err := os.WriteFile(splashPath, videoContent, 0644); err != nil {
		return "", fmt.Errorf("failed to write splash video to temporary file: %w", err)
	}

	log.Printf("Temporary splash video created at: %s", splashPath)
	return splashPath, nil
}

// prepareVLCCommand constructs the VLC command with necessary arguments.
func prepareVLCCommand(videoPath string) *exec.Cmd {
	args := []string{
		"--fullscreen",
		"--video-on-top",  // Keep the video window on top
		"--play-and-exit", // Close VLC after playback
		"--intf", "dummy", // Use dummy interface (no GUI)
		"--dummy-quiet", // Suppress logs from VLC
		"--no-osd",      // Disable On-Screen Display
		videoPath,       // Path to the video file
	}
	return exec.Command("vlc", args...)
}
