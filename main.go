package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/maiku1008/pomodoro-cli/internal/pomodoro"
)

func main() {
	// Parse command-line flags
	timer := flag.Int("timer", 25, "The timer duration in minutes")
	breakTimer := flag.Int("break", 5, "The break duration in minutes")
	intervals := flag.Int("interval", 1, "The number of pomodoros to complete")
	hostsFile := flag.String("hosts", "/etc/hosts", "The file to modify")
	flag.Parse()

	// Setup context for cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup signal handling for graceful shutdown
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		fmt.Println("\nReceived interrupt signal")
		cancel()
	}()

	// Configure the Pomodoro session
	cfg := pomodoro.Config{
		WorkDuration:  time.Duration(*timer) * time.Minute,
		BreakDuration: time.Duration(*breakTimer) * time.Minute,
		Intervals:     *intervals,
		BlockList: []string{
			"reddit.com",
			"facebook.com",
			"linkedin.com",
			"bbc.com",
			"timesofmalta.com",
			"nintendolife.com",
			"kotaku.com",
			"polygon.com",
		},
		HostsFilePath: *hostsFile,
		WindupSound:   "sounds/windup.wav",
		TickingSound:  "sounds/ticking.wav",
		DingSound:     "sounds/ding.wav",
	}

	// Run the Pomodoro timer
	if err := pomodoro.Run(ctx, cfg); err != nil {
		log.Fatalf("Error running pomodoro: %v\n", err)
	}
}
