package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
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
	silent := flag.Bool("silent", false, "Disable the ticking sound")
	blockList := flag.String("blocklist", "", "The list of sites to block, separated by commas.")
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

	var splitBlockList []string
	if *blockList != "" {
		splitBlockList = strings.Split(*blockList, ",")
	}
	// Configure the Pomodoro session
	cfg := pomodoro.Config{
		WorkDuration:  time.Duration(*timer) * time.Minute,
		BreakDuration: time.Duration(*breakTimer) * time.Minute,
		Intervals:     *intervals,
		BlockList:     splitBlockList,
		HostsFilePath: *hostsFile,
		Silent:        *silent,
	}

	// Run the Pomodoro timer
	if err := pomodoro.Run(ctx, cfg); err != nil {
		log.Fatalf("Error running pomodoro: %v\n", err)
	}
}
