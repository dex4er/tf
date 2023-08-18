package run

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"

	"github.com/mitchellh/colorstring"

	"github.com/dex4er/tf/util"
)

func terraformWithFilter(command string, args []string, patternIgnoreLine string, patternIgnoreFooter string) error {
	reIgnoreLine := regexp.MustCompile(patternIgnoreLine)
	reIgnoreFooter := regexp.MustCompile(patternIgnoreFooter)

	noColor := false

	newArgs := []string{}

	for _, arg := range args {
		switch util.ReplaceFirstTwoDashes(arg) {
		case "-no-color":
			noColor = true
			newArgs = append(newArgs, arg)
		case "-no-colors":
			noColor = true
			newArgs = append(newArgs, "-no-color")
		default:
			newArgs = append(newArgs, arg)
		}
	}

	// clear color even after errors
	if !noColor {
		defer fmt.Print("\x1b[0m")
	}

	signal.Ignore(syscall.SIGINT)

	cmd := execTerraformCommand(append([]string{command}, newArgs...)...)

	cmd.Stdin = os.Stdin

	file, err := util.OpenOutputFile()
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

	isEof := false
	ignoreBlock := false
	ignoreFooter := false
	wasEmptyLine := false

	reader := bufio.NewReader(cmdStdout)

	// buffer for the current line
	line := ""

	// Token scanner cannot be used here because of interactive prompts from
	// terraform. The prompt doesn't end with EOL then Stdout must be read rune
	// by rune.

	for {
		// stream was ended in previous iteration of the loop
		if isEof {
			break
		}

		r, _, err := reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				// still we need to process the line and end this loop in the next iteration
				isEof = true
			} else {
				return fmt.Errorf("reading command output: %w", err)
			}
		} else {
			line = line + string(r)
		}

		// tokens that triggers processing of the line
		if strings.Contains(line, colorstring.Color("[bold]Enter a value:[reset] ")) ||
			strings.Contains(line, "Enter a value: ") ||
			r == '\n' || isEof {

			// verbatim output to the log file
			if file != nil {
				fmt.Fprint(file, line)
			}

			if ignoreFooter {
				goto NEXT
			}

			if ignoreBlock {
				goto NEXT
			}

			// ignore just this line
			if patternIgnoreLine != "" && reIgnoreLine.MatchString(line) {
				goto NEXT
			}

			// skip another empty line but preserve color codes
			if wasEmptyLine && util.IsEmptyLine(line) {
				line = strings.TrimSuffix(line, "\n")
				line = strings.TrimSuffix(line, "\r")
			}

			// print current line buffer
			fmt.Print(line)

			// mark if current line was empty for next loop iteration
			wasEmptyLine = util.IsEmptyLine(line)

			// ignore from this line to the end
			if patternIgnoreFooter != "" && reIgnoreFooter.MatchString(line) {
				ignoreFooter = true
			}

		NEXT:
			// empty line buffer before next loop iteration
			line = ""
		}
	}

	return cmd.Wait()
}
