package run

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/dex4er/tf/util"
)

func terraformWithoutColors(command string, args []string) error {
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

	ignoreFooter := false

	scanner := bufio.NewScanner(cmdStdout)

	for scanner.Scan() {
		if ignoreFooter {
			continue
		}

		line := scanner.Text()
		line = util.RemoveColors(line)

		fmt.Println(line)
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return cmd.Wait()
}
