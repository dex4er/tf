package run

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"regexp"
	"syscall"

	"github.com/dex4er/tf/util"
)

func terraformWithoutColors(command string, args []string, patternIgnoreFooter string) error {
	reIgnoreFooter := regexp.MustCompile(patternIgnoreFooter)

	signal.Ignore(syscall.SIGINT)

	cmd := exec.Command("terraform", append([]string{command}, args...)...)

	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	cmdStdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("creating stdout pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("starting the command: %w", err)
	}

	ignoreFooter := false

	scanner := bufio.NewScanner(cmdStdout)

	for scanner.Scan() {
		if ignoreFooter {
			continue
		}

		line := scanner.Text()
		line = util.RemoveColors(line)

		// ignore from this line to the end
		if patternIgnoreFooter != "" && reIgnoreFooter.MatchString(line) {
			ignoreFooter = true
			continue
		}

		fmt.Println(line)
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return cmd.Wait()
}
