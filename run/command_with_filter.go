package run

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"regexp"
	"strings"
	"syscall"

	"github.com/mitchellh/colorstring"

	"github.com/dex4er/tf/util"
)

func commandWithFilter(command string, args []string, ignorePattern string, endPattern string) error {
	defer colorstring.Printf("[reset]")

	signal.Ignore(syscall.SIGINT)

	cmd := exec.Command("terraform", append([]string{command}, args...)...)

	cmd.Stdin = os.Stdin

	file, err := util.OpenLogfile()
	if err != nil {
		return err
	}
	if file != nil {
		defer file.Close()
		cmd.Stderr = io.MultiWriter(os.Stderr, file)
	} else {
		cmd.Stderr = os.Stderr
	}

	cmdStdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("creating stdout pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("starting the command: %w", err)
	}

	reIgnore := regexp.MustCompile(ignorePattern)
	reEnd := regexp.MustCompile(endPattern)

	isEof := false
	skipBegin := true
	skipEnd := false
	wasEmptyLine := false

	reader := bufio.NewReader(cmdStdout)

	line := ""

	for {
		if isEof {
			break
		}

		r, _, err := reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				isEof = true
			} else {
				return fmt.Errorf("reading command output: %w", err)
			}
		} else {
			line = line + string(r)
		}

		if strings.Contains(line, colorstring.Color("[bold]Enter a value:[reset] ")) || strings.Contains(line, "Enter a value: ") || r == '\n' || isEof {
			if file != nil {
				fmt.Fprintln(file, line)
			}

			if skipEnd {
				line = ""
				continue
			}

			if reIgnore.MatchString(line) {
				line = ""
				continue
			}

			if skipBegin && line != "" {
				skipBegin = false
			}

			if skipBegin {
				line = ""
				continue
			}

			if wasEmptyLine && (line == "" || line == "\n" || line == "\r\n") {
				line = ""
				continue
			}

			fmt.Print(line)

			wasEmptyLine = line == "" || line == "\n" || line == "\r\n"

			if endPattern != "" && reEnd.MatchString(line) {
				skipEnd = true
			}

			if isEof {
				break
			}

			line = ""
		}
	}

	return cmd.Wait()
}
