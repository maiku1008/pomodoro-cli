package sound

import (
	"context"
	"log"
	"os/exec"
)

// Play plays a sound file once using afplay (macOS)
func Play(soundFile string) {
	cmd := exec.Command("afplay", soundFile)
	if err := cmd.Start(); err != nil {
		log.Printf("Warning: Could not play sound %s: %v\n", soundFile, err)
	}
}

// StartTicking starts playing the ticking sound in a loop
// The ticking will stop when the context is cancelled
func StartTicking(ctx context.Context, soundFile string) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				cmd := exec.CommandContext(ctx, "afplay", soundFile)
				_ = cmd.Run() // Will be interrupted when context is cancelled
			}
		}
	}()
}
