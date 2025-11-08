package pomodoro

import (
	"context"
	"fmt"
	"log"
	"os"
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
	WindupSound   string
	TickingSound  string
	DingSound     string
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

	// Setup cleanup to always unblock sites when exiting
	defer func() {
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
		sound.Play(cfg.WindupSound)
		workCtx, workCancel := context.WithCancel(ctx)
		sound.StartTicking(workCtx, cfg.TickingSound)

		// Wait for either the work timer to finish or cancellation
		select {
		case <-time.After(cfg.WorkDuration):
			workCancel() // Stop ticking
			sound.Play(cfg.DingSound)
			fmt.Println("\n✅ Work session complete!")
		case <-ctx.Done():
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
		select {
		case <-time.After(cfg.BreakDuration):
			fmt.Println("\n⏰ Break finished!")
		case <-ctx.Done():
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
