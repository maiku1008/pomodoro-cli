package pomodoro

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/maiku1008/pomodoro-cli/internal/hosts"
	"github.com/maiku1008/pomodoro-cli/internal/sound"
)

// Config holds the configuration for a Pomodoro session
type Config struct {
	WorkDuration  time.Duration
	BreakDuration time.Duration
	Intervals     int
	BlockList     []string
	HostsFilePath string
	Silent        bool
}

// Run executes the Pomodoro timer with the given configuration
func Run(ctx context.Context, cfg Config) error {
	// create string template to add to the hosts file
	blockTemplate := hosts.BlockTemplate(cfg.BlockList)

	// open the hosts file
	hostsFile, err := os.OpenFile(cfg.HostsFilePath, os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("failed to open hosts file: %w", err)
	}
	defer hostsFile.Close()

	// Setup cleanup to always unblock sites and clean up temp files when exiting
	defer func() {
		sound.Cleanup()
		if err := hosts.Unblock(blockTemplate, hostsFile); err != nil {
			log.Printf("Error during cleanup: %v\n", err)
		}
	}()

	// Run multiple pomodoro cycles
	fmt.Printf("🍅 Starting %d Pomodoro cycle(s)\n\n", cfg.Intervals)

	for i := 1; i <= cfg.Intervals; i++ {
		fmt.Printf("═══════════════════════════════════════\n")
		fmt.Printf("🍅 Pomodoro %d of %d\n", i, cfg.Intervals)
		fmt.Printf("═══════════════════════════════════════\n\n")

		// Phase 1: Work time - block sites
		fmt.Printf("⏰ Work session (%.0f minutes)\n", cfg.WorkDuration.Minutes())
		fmt.Println("Blocking distracting sites...")
		err = hosts.Block(blockTemplate, hostsFile)
		if err != nil {
			return fmt.Errorf("failed to block sites: %w", err)
		}

		// Play windup sound and start ticking
		sound.PlaySound(sound.Windup)
		workCtx, workCancel := context.WithCancel(ctx)
		if !cfg.Silent {
			sound.StartTickingSound(workCtx)
		}

		// Wait for either the work timer to finish or cancellation
		if waitWithCountdown(ctx, cfg.WorkDuration, "🍅") {
			workCancel() // Stop ticking
			sound.PlaySound(sound.Ding)
			fmt.Println("\n✅ Work session complete!")
		} else {
			workCancel() // Stop ticking
			fmt.Println("\n❌ Pomodoro cancelled")
			return nil
		}

		// Phase 2: Break time - unblock sites
		fmt.Println("Unblocking sites for break...")
		err = hosts.Unblock(blockTemplate, hostsFile)
		if err != nil {
			return fmt.Errorf("failed to unblock sites: %w", err)
		}

		fmt.Printf("\n☕ Break time! (%.0f minutes)\n", cfg.BreakDuration.Minutes())
		fmt.Println("Sites are now unblocked. Take a break!")

		// Wait for either the break timer to finish or cancellation
		if waitWithCountdown(ctx, cfg.BreakDuration, "☕") {
			fmt.Println("\n⏰ Break finished!")
		} else {
			fmt.Println("\n❌ Break cancelled")
			return nil
		}

		fmt.Printf("\n✨ Pomodoro %d complete!\n\n", i)
	}

	fmt.Println("═══════════════════════════════════════")
	fmt.Printf("🎉 All %d Pomodoro cycles complete!\n", cfg.Intervals)
	fmt.Println("═══════════════════════════════════════")
	return nil
}

// waitWithCountdown waits for the specified duration while displaying a countdown with progress bar
// Returns true if the duration completed, false if context was cancelled
func waitWithCountdown(ctx context.Context, duration time.Duration, label string) bool {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	endTime := time.Now().Add(duration)
	totalSeconds := duration.Seconds()

	for {
		remaining := time.Until(endTime)
		if remaining <= 0 {
			fmt.Print("\r" + strings.Repeat(" ", 80) + "\r") // Clear the line
			return true
		}

		// Calculate progress
		elapsed := totalSeconds - remaining.Seconds()
		progress := elapsed / totalSeconds
		barWidth := 30
		filled := int(progress * float64(barWidth))

		// Build progress bar
		bar := strings.Repeat("█", filled) + strings.Repeat("░", barWidth-filled)

		// Format time as MM:SS
		minutes := int(remaining.Minutes())
		seconds := int(remaining.Seconds()) % 60

		fmt.Printf("\r%s [%s] %02d:%02d remaining", label, bar, minutes, seconds)

		select {
		case <-ctx.Done():
			fmt.Print("\r" + strings.Repeat(" ", 80) + "\r") // Clear the line
			return false
		case <-ticker.C:
			// Continue loop to update display
		}
	}
}
