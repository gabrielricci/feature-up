package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func RunCommand(cmdName string, args []string, strict bool) (string, error) {
	var err error

	stdOut := &bytes.Buffer{}
	stdErr := &bytes.Buffer{}

	cmd := exec.Command(cmdName, args...)
	cmd.Stdout = stdOut
	cmd.Stderr = stdErr

	fmt.Println(" - Running", cmdName, strings.Join(args, " "))

	if err = cmd.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "There was an error running the command: ", err)
		fmt.Fprintln(os.Stderr, " - ", string(stdErr.Bytes()))
		if strict {
			os.Exit(1)
		}
		return "", err
	}

	return string(stdOut.Bytes()), nil
}
