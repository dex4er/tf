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

func commandWithFilter(command string, args []string, ignoreLinePattern string, ignoreNextLinePattern string, ignoreBlockStartPattern string, ignoreBlockEndPattern string, footerPattern string) error {
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

	reIgnoreLine := regexp.MustCompile(ignoreLinePattern)
	reIgnoreNextLine := regexp.MustCompile(ignoreNextLinePattern)
	reIgnoreBlockStart := regexp.MustCompile(ignoreBlockStartPattern)
	reIgnoreBlockEnd := regexp.MustCompile(ignoreBlockEndPattern)
	reFooter := regexp.MustCompile(footerPattern)

	isEof := false
	ignoreNextLine := false
	ignoreBlock := false
	skipHeader := true
	skipFooter := false
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

			if skipFooter {
				line = ""
				continue
			}

			if reIgnoreBlockStart.MatchString(line) {
				ignoreBlock = true
				line = ""
				continue
			}

			if reIgnoreBlockEnd.MatchString(line) {
				ignoreBlock = false
				line = ""
				continue
			}

			if ignoreBlock {
				line = ""
				continue
			}

			if reIgnoreNextLine.MatchString(line) {
				ignoreNextLine = true
				line = ""
				continue
			}

			if ignoreNextLine {
				ignoreNextLine = false
				line = ""
				continue
			}

			if reIgnoreLine.MatchString(line) {
				line = ""
				continue
			}

			if skipHeader && line != "" {
				skipHeader = false
			}

			if skipHeader {
				line = ""
				continue
			}

			if wasEmptyLine && (line == "" || line == "\n" || line == "\r\n") {
				line = ""
				continue
			}

			fmt.Print(line)

			wasEmptyLine = line == "" || line == "\n" || line == "\r\n"

			if footerPattern != "" && reFooter.MatchString(line) {
				skipFooter = true
			}

			if isEof {
				break
			}

			line = ""
		}
	}

	return cmd.Wait()
}
