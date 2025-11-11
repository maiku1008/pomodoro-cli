package sound

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"
)

// Play plays a sound file once using afplay (macOS)
func Play(soundFile string) {
	cmd := exec.Command("afplay", soundFile)
	if err := cmd.Start(); err != nil {
		log.Printf("Warning: Could not play sound %s: %v\n", soundFile, err)
	}
}

// StartTickingSound starts playing the ticking sound in a continuous loop
// The ticking will stop when the context is cancelled
// To avoid gaps between loops, we start the next play slightly before the current one finishes
func StartTickingSound(ctx context.Context, soundFile string) {
	go func() {
		// Get the duration of the sound file
		duration, err := getSoundDuration(soundFile)
		if err != nil {
			log.Printf("Warning: Could not get sound duration, falling back to sequential play: %v\n", err)
			// Fallback to simple sequential play
			for {
				select {
				case <-ctx.Done():
					return
				default:
					cmd := exec.CommandContext(ctx, "afplay", soundFile)
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
				cmd := exec.CommandContext(ctx, "afplay", soundFile)
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
