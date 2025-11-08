package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func main() {
	blockList := []string{
		"reddit.com",
		"facebook.com",
		"linkedin.com",
	}

	// set a configurable timer for pomodoro
	timer := flag.Int("timer", 25, "The timer duration in minutes")
	breakTimer := flag.Int("break", 5, "The break duration in minutes")
	intervalTimer := flag.Int("interval", 4, "The interval duration in minutes")
	// TODO: implement intervals
	_ = intervalTimer
	flag.Parse()

	// create string template to add to the hosts file
	blockTemplate := blockTemplate(blockList)

	// open the file /etc/hosts
	hostsFile, err := os.OpenFile("hosts.txt", os.O_RDWR, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer hostsFile.Close()

	// set cancel context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup cleanup to always unblock sites when exiting
	// (in case of cancellation during work phase)
	defer func() {
		// Only try to unblock if we're exiting unexpectedly
		// The unblockSites function handles the case where sites are already unblocked
		if err := unblockSites(blockTemplate, hostsFile); err != nil {
			log.Printf("Error during cleanup: %v\n", err)
		}
	}()

	// cancel the context when we receive a signal to stop
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		fmt.Println("\nReceived interrupt signal")
		cancel()
	}()

	// Phase 1: Work time - block sites
	fmt.Printf("🍅 Starting Pomodoro work session (%d minutes)\n", *timer)
	fmt.Println("Blocking distracting sites...")
	err = blockSites(blockTemplate, hostsFile)
	if err != nil {
		log.Fatal(err)
	}

	// Play windup sound and start ticking
	playSound("windup.wav")
	workCtx, workCancel := context.WithCancel(ctx)
	startTickingSound(workCtx)

	// Wait for either the work timer to finish or cancellation
	select {
	case <-time.After(time.Duration(*timer) * time.Minute):
		workCancel() // Stop ticking
		playSound("ding.wav")
		fmt.Println("\n✅ Work session complete!")
	case <-ctx.Done():
		workCancel() // Stop ticking
		fmt.Println("\n❌ Work session cancelled")
		return // defer will unblock sites
	}

	// Phase 2: Break time - unblock sites
	fmt.Println("Unblocking sites for break...")
	err = unblockSites(blockTemplate, hostsFile)
	if err != nil {
		log.Printf("Error unblocking sites: %v\n", err)
		return
	}

	fmt.Printf("\n☕ Break time! (%d minutes)\n", *breakTimer)
	fmt.Println("Sites are now unblocked. Take a break!")

	// Play windup sound and start ticking for break
	playSound("windup.wav")
	breakCtx, breakCancel := context.WithCancel(ctx)
	startTickingSound(breakCtx)

	// Wait for either the break timer to finish or cancellation
	select {
	case <-time.After(time.Duration(*breakTimer) * time.Minute):
		breakCancel() // Stop ticking
		playSound("ding.wav")
		fmt.Println("\n⏰ Break finished!")
	case <-ctx.Done():
		breakCancel() // Stop ticking
		fmt.Println("\n❌ Break cancelled")
		return
	}

	fmt.Println("Pomodoro session complete! 🎉")
}

// playSound plays a sound file once using afplay (macOS)
func playSound(soundFile string) {
	cmd := exec.Command("afplay", soundFile)
	if err := cmd.Start(); err != nil {
		log.Printf("Warning: Could not play sound %s: %v\n", soundFile, err)
	}
}

// startTickingSound starts playing the ticking sound in a loop
// Returns a cancel function to stop the ticking
func startTickingSound(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				cmd := exec.CommandContext(ctx, "afplay", "ticking.wav")
				_ = cmd.Run() // Will be interrupted when context is cancelled
			}
		}
	}()
}

// create a string template to add to the hosts file
func blockTemplate(blockList []string) string {
	var blockTemplate strings.Builder
	blockTemplate.WriteString("\n### Pomodoro CLI - Begin Blocked sites ###\n")
	for _, block := range blockList {
		blockTemplate.WriteString(fmt.Sprintf("127.0.0.1 %s\n127.0.0.1 www.%s\n", block, block))
	}
	blockTemplate.WriteString("### Pomodoro CLI - End Blocked sites ###\n")
	return blockTemplate.String()
}

// block the sites in the hosts file
func blockSites(blockTemplate string, hostsFile *os.File) error {
	// Seek to beginning of file to read
	_, err := hostsFile.Seek(0, 0)
	if err != nil {
		return err
	}

	// read the file
	hosts, err := io.ReadAll(hostsFile)
	if err != nil {
		return err
	}

	// check if the block template is in the hosts file
	if strings.Contains(string(hosts), blockTemplate) {
		fmt.Println("Block template already exists in hosts file")
		return nil
	}

	// add the block template to the hosts file
	_, err = hostsFile.WriteString(blockTemplate)
	if err != nil {
		return err
	}

	fmt.Println("Block template added to hosts file")
	return nil
}

// unblock the sites in the hosts file
func unblockSites(blockTemplate string, hostsFile *os.File) error {
	// Seek to beginning of file to read
	_, err := hostsFile.Seek(0, 0)
	if err != nil {
		return err
	}

	// read the file
	hosts, err := io.ReadAll(hostsFile)
	if err != nil {
		return err
	}

	// check if the block template is in the hosts file
	if strings.Contains(string(hosts), blockTemplate) {
		// Remove the block template
		newHosts := strings.Replace(string(hosts), blockTemplate, "", -1)

		// Truncate the file and seek to beginning
		err = hostsFile.Truncate(0)
		if err != nil {
			return err
		}
		_, err = hostsFile.Seek(0, 0)
		if err != nil {
			return err
		}

		// Write the updated content
		_, err = hostsFile.WriteString(newHosts)
		if err != nil {
			return err
		}

		fmt.Println("Block template removed from hosts file")
	} else {
		fmt.Println("Block template not found in hosts file")
	}

	return nil
}
