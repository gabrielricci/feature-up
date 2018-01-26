package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func RunUnitTests() {
	root := GitGetRepositoryRoot()
	filePath := root + "/.feature-up.tests"

	if _, err := os.Stat(filePath); err == nil {
		fileHandle, _ := os.Open(filePath)
		defer fileHandle.Close()
		fileScanner := bufio.NewScanner(fileHandle)

		for fileScanner.Scan() {
			parts := strings.Split(fileScanner.Text(), " ")
			RunCommand(parts[0], parts[1:], true)
		}

		fmt.Println("All tests passing!")
	} else {
		fmt.Println("No tests are configured for this project, skipping...")
	}

}
