package hosts

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// BlockTemplate creates a string template to add to the hosts file
func BlockTemplate(blockList []string) string {
	var template strings.Builder
	template.WriteString("\n### Pomodoro CLI - Begin Blocked sites ###\n")
	for _, block := range blockList {
		template.WriteString(fmt.Sprintf("127.0.0.1 %s\n127.0.0.1 www.%s\n", block, block))
	}
	template.WriteString("### Pomodoro CLI - End Blocked sites ###\n")
	return template.String()
}

// Block adds the block template to the hosts file
func Block(blockTemplate string, hostsFile *os.File) error {
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

// Unblock removes the block template from the hosts file
func Unblock(blockTemplate string, hostsFile *os.File) error {
	// Seek to beginning of file to read
	_, err := hostsFile.Seek(0, 0)
	if err != nil {
		return err
	}

	// read the file
	hostContents, err := io.ReadAll(hostsFile)
	if err != nil {
		return err
	}

	// check if the block template is in the hosts file
	if strings.Contains(string(hostContents), blockTemplate) {
		// Remove the block template
		newHosts := strings.Replace(string(hostContents), blockTemplate, "", -1)

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
