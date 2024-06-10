package run

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"regexp"
	"syscall"

	"github.com/dex4er/tf/util"
)

func terraformWithoutColors(command string, noOutputs bool, args []string) error {
	patternIgnoreOutputs := `^Outputs:(\n|$)`
	reIgnoreOutputs := regexp.MustCompile(patternIgnoreOutputs)

	signal.Ignore(syscall.SIGINT)

	cmd := execTerraformCommand(append([]string{command}, args...)...)

	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	cmdStdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("creating stdout pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("starting the command: %w", err)
	}

	skipOutputs := false

	scanner := bufio.NewScanner(cmdStdout)

	for scanner.Scan() {
		if skipOutputs {
			continue
		}

		line := scanner.Text()
		line = util.RemoveColors(line)

		// starts ignoring the outputs
		if noOutputs && reIgnoreOutputs.MatchString(line) {
			skipOutputs = true
			continue
		}

		fmt.Println(line)
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return cmd.Wait()
}
