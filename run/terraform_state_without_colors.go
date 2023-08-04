package run

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/dex4er/tf/util"
)

func terraformStateWithoutColors(command string, args []string) error {
	newArgs := []string{}

	for _, arg := range args {
		if util.StartsWith(arg, '-') {
			newArgs = append(newArgs, arg)
		} else {
			newArgs = append(newArgs, util.AddQuotes(arg))
		}
	}

	signal.Ignore(syscall.SIGINT)

	cmd := exec.Command("terraform", append([]string{"state", command}, newArgs...)...)

	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	cmdStdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("creating stdout pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("starting the command: %w", err)
	}

	scanner := bufio.NewScanner(cmdStdout)

	for scanner.Scan() {
		line := scanner.Text()
		line = util.RemoveColors(line)
		fmt.Println(line)
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return cmd.Wait()
}
