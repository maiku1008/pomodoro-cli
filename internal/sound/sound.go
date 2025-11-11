package sound

import (
	"context"
	"embed"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

//go:embed sounds/*.wav
var soundFiles embed.FS

// SoundType represents the type of sound to play
type SoundType int

const (
	Windup SoundType = iota
	Ticking
	Ding
)

var (
	tempFiles     = make(map[SoundType]string)
	tempFilesLock sync.Mutex
	initOnce      sync.Once
)

// soundFileNames maps SoundType to embedded file names
var soundFileNames = map[SoundType]string{
	Windup:  "sounds/windup.wav",
	Ticking: "sounds/ticking.wav",
	Ding:    "sounds/ding.wav",
}

// Initialize sets up temporary files for all embedded sounds
// This is called automatically on first use
func initialize() error {
	tempFilesLock.Lock()
	defer tempFilesLock.Unlock()

	for soundType, fileName := range soundFileNames {
		// Read embedded file
		data, err := soundFiles.ReadFile(fileName)
		if err != nil {
			return fmt.Errorf("failed to read embedded sound %s: %w", fileName, err)
		}

		// Create temp file
		tmpFile, err := os.CreateTemp("", "pomodoro-sound-*.wav")
		if err != nil {
			return fmt.Errorf("failed to create temp file for %s: %w", fileName, err)
		}

		// Write data to temp file
		if _, err := tmpFile.Write(data); err != nil {
			tmpFile.Close()
			os.Remove(tmpFile.Name())
			return fmt.Errorf("failed to write temp file for %s: %w", fileName, err)
		}

		if err := tmpFile.Close(); err != nil {
			os.Remove(tmpFile.Name())
			return fmt.Errorf("failed to close temp file for %s: %w", fileName, err)
		}

		tempFiles[soundType] = tmpFile.Name()
	}

	return nil
}

// getTempFilePath returns the temporary file path for a given sound type
func getTempFilePath(soundType SoundType) (string, error) {
	var initErr error
	initOnce.Do(func() {
		initErr = initialize()
	})
	if initErr != nil {
		return "", initErr
	}

	tempFilesLock.Lock()
	defer tempFilesLock.Unlock()

	path, exists := tempFiles[soundType]
	if !exists {
		return "", fmt.Errorf("temp file not found for sound type %d", soundType)
	}
	return path, nil
}

// Cleanup removes all temporary sound files
// Should be called when the application exits
func Cleanup() {
	tempFilesLock.Lock()
	defer tempFilesLock.Unlock()

	for _, path := range tempFiles {
		if err := os.Remove(path); err != nil {
			log.Printf("Warning: Could not remove temp file %s: %v\n", filepath.Base(path), err)
		}
	}
	tempFiles = make(map[SoundType]string)
}

// PlaySound plays a sound once using afplay (macOS)
func PlaySound(soundType SoundType) {
	path, err := getTempFilePath(soundType)
	if err != nil {
		log.Printf("Warning: Could not get sound file: %v\n", err)
		return
	}

	cmd := exec.Command("afplay", path)
	if err := cmd.Start(); err != nil {
		log.Printf("Warning: Could not play sound: %v\n", err)
	}
}

// StartTickingSound starts playing the ticking sound in a continuous loop
// The ticking will stop when the context is cancelled
// To avoid gaps between loops, we start the next play slightly before the current one finishes
func StartTickingSound(ctx context.Context) {
	path, err := getTempFilePath(Ticking)
	if err != nil {
		log.Printf("Warning: Could not get ticking sound file: %v\n", err)
		return
	}

	go func() {
		// Get the duration of the sound file
		duration, err := getSoundDuration(path)
		if err != nil {
			log.Printf("Warning: Could not get sound duration, falling back to sequential play: %v\n", err)
			// Fallback to simple sequential play
			for {
				select {
				case <-ctx.Done():
					return
				default:
					cmd := exec.CommandContext(ctx, "afplay", path)
					_ = cmd.Run()
				}
			}
		}

		// Start overlapping playback to avoid gaps
		// We start the next play 200ms before the current one finishes
		overlapTime := duration - 200
		if overlapTime < 0 {
			overlapTime = duration / 2
		}

		for {
			select {
			case <-ctx.Done():
				return
			default:
				// Start playing
				cmd := exec.CommandContext(ctx, "afplay", path)
				if err := cmd.Start(); err != nil {
					log.Printf("Warning: Could not start ticking sound: %v\n", err)
					return
				}

				// Wait for overlap time, then start the next one
				select {
				case <-ctx.Done():
					return
				case <-time.After(time.Duration(overlapTime) * time.Millisecond):
					// Continue to next iteration to start the next play
				}
			}
		}
	}()
}

// getSoundDuration returns the duration of an audio file in milliseconds
// Uses afinfo command available on macOS
func getSoundDuration(soundFile string) (int64, error) {
	cmd := exec.Command("afinfo", soundFile)
	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	// Parse the output to find "estimated duration: X.XXX sec"
	var duration float64
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "estimated duration:") {
			// Format: "estimated duration: 33.123456 sec"
			parts := strings.Fields(line)
			if len(parts) >= 3 {
				fmt.Sscanf(parts[2], "%f", &duration)
				break
			}
		}
	}

	if duration == 0 {
		return 0, fmt.Errorf("could not parse duration from afinfo output")
	}

	return int64(duration * 1000), nil // Convert to milliseconds
}
