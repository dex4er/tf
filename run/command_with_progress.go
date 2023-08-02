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

	"github.com/dex4er/tf/util"
	"github.com/mitchellh/colorstring"
)

func commandWithProgress(command string, args []string) error {
	patternIgnoreLine := "Terraform used the selected providers to generate the following execution" +
		"|Preparing the remote plan\\.\\.\\." +
		"|Running plan in Terraform Cloud\\. Output will stream here\\. Pressing Ctrl-C" +
		"|will stop streaming the logs, but will not stop the plan running remotely\\." +
		"|The remote workspace is configured to work with configuration at" +
		"|relative to the target repository." +
		"|excluding files or directories as defined by a \\.terraformignore file" +
		"|at .*\\.terraformignore (if it is present)," +
		"|in order to capture the filesystem context the remote workspace expects:" +
		"|plan. Resource actions are indicated with the following symbols:" +
		"|Waiting for the plan to start\\.\\.\\." +
		"|Terraform will perform the following actions:" +
		"|Terraform will perform the actions described above\\." +
		"|Terraform has compared your real infrastructure against your configuration" +
		"|and found no differences, so no changes are needed\\." +
		"|Unless you have made equivalent changes to your configuration, or ignored the" +
		"|relevant attributes using ignore_changes, the following plan may include" +
		"|actions to undo or respond to these changes\\." +
		"|This is a refresh-only plan, so Terraform will not take any actions to undo" +
		"|these\\. If you were expecting these changes then you can apply this plan to" +
		"|Terraform has checked that the real remote objects still match the result of" +
		"|your most recent changes, and found no differences\\." +
		"|To perform exactly these actions, run the following command to apply:" +
		"|To see the full warning notes, run Terraform without -compact-warnings\\." +
		"|Acquiring state lock\\. This may take a few moments\\.\\.\\." +
		"|Releasing state lock\\. This may take a few moments\\.\\.\\." +
		"|Warnings:" +
		"|Note: You .* use the -out option to save this plan, so Terraform" +
		"|guarantee to take exactly these actions if you run \"terraform apply\" now\\." +
		"|Apply complete! Resources: 0 added, 0 changed, 0 destroyed\\." +
		"|─────────────────────────────────────────────────────────────────────────────"

	patternIgnoreNextLine := "record the updated values in the Terraform state without changing any remote" +
		"|record the updated values in the Terraform state without changing any remote" +
		"| Experimental feature .* is active" +
		"|Saved the plan to: terraform\\.tfplan"

	patternIgnoreBlockStart := "Warning:.*Applied changes may be incomplete" +
		"|Warning:.*Resource targeting is in effect" +
		"|This plan was saved to: terraform.tfplan"

	patternIgnoreBlockEnd := "suggests to use it as part of an error message" +
		"|exceptional situations such as recovering from errors or mistakes" +
		"|terraform apply \"terraform\\.tfplan\""

	patternIgnoreShortFormat := "= \\(known after apply\\)" +
		"|\\(\\d+ unchanged \\w+ hidden\\)" +
		"|\\(config refers to values not yet known\\)"

	patternIgnoreCompactFormat := "^\\s\\s[\\s+~-]" +
		"|\\(config refers to values not yet known\\)"

	// patternStartRefreshing := ": Refreshing state\\.\\.\\." +
	// 	"|: Refreshing\\.\\.\\."
	// patternStopRefreshing := ": Drift detected"

	patternStartReading := "(?:.\\[0m.\\[1m)?(.*?): Reading\\.\\.\\."
	patternStopReading := "(?:.\\[0m.\\[1m)?(.*?): Read complete after"

	// patternStartCreating := ": Creating\\.\\.\\."
	// patternStopCreating := ": Creation complete after"

	// patternStartDestroying := ": Destroying\\.\\.\\."
	// patternStopDestroying := ": Destruction complete after"

	// patternStartModyfying := ": Modifying\\.\\.\\."
	// patternStopModifying := ": Modifications complete after"

	// patternStillProcessing := ": Still .*ing\\.\\.\\."

	reIgnoreLine := regexp.MustCompile(patternIgnoreLine)
	reIgnoreNextLine := regexp.MustCompile(patternIgnoreNextLine)
	reIgnoreBlockStart := regexp.MustCompile(patternIgnoreBlockStart)
	reIgnoreBlockEnd := regexp.MustCompile(patternIgnoreBlockEnd)
	reIgnoreShortFormat := regexp.MustCompile(patternIgnoreShortFormat)
	reIgnoreCompactFormat := regexp.MustCompile(patternIgnoreCompactFormat)
	reStartReading := regexp.MustCompile(patternStartReading)
	reStopReading := regexp.MustCompile(patternStopReading)

	format := "short"
	progress := "fan"
	noOutputs := false

	if os.Getenv("TF_IN_AUTOMATION") == "1" {
		progress = "verbose"
	}

	newArgs := []string{}

	for _, arg := range args {
		switch arg {
		case "-compact":
			format = "compact"
		case "-dot":
			progress = "dot"
		case "-fan":
			progress = "fan"
		case "-full":
			format = "full"
		case "-no-outputs":
			noOutputs = true
		case "-quiet":
			progress = "quiet"
		case "-short":
			format = "short"
		case "-verbose":
			progress = "verbose"
		default:
			if util.StartsWith(arg, '-') {
				newArgs = append(newArgs, arg)
			} else {
				newArgs = append(newArgs, fmt.Sprintf("-target=%s", util.AddQuotes(arg)))
			}
		}
	}

	fmt.Println(progress)
	fmt.Println(noOutputs)

	defer fmt.Print(util.ColorReset)

	signal.Ignore(syscall.SIGINT)

	cmd := exec.Command("terraform", append([]string{command}, newArgs...)...)

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

	isEof := false
	ignoreNextLine := false
	ignoreBlock := false
	skipHeader := true
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

			if m := reStartReading.FindStringSubmatch(line); m != nil {
				r := m[1]
				colorstring.Printf("[cyan]0/1[reset] [green]0/0[reset] [yellow]0/0[reset] [red]0/0[reset] %s\r", r)
				line = ""
				continue
			}

			if m := reStopReading.FindStringSubmatch(line); m != nil {
				r := m[1]
				colorstring.Printf("[cyan]1/1[reset] [green]0/0[reset] [yellow]0/0[reset] [red]0/0[reset] %s\r", r)
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

			if format == "short" && reIgnoreShortFormat.MatchString(line) {
				line = ""
				continue
			}

			if format == "compact" && reIgnoreCompactFormat.MatchString(line) {
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

			if wasEmptyLine && util.IsEmptyLine(line) {
				line = ""
				continue
			}

			if strings.HasSuffix(line, "\n") {
				line = strings.TrimSuffix(line, "\n")
				fmt.Print(line)
				fmt.Println(strings.Repeat(" ", 79-len(line)))
			} else {
				fmt.Print(line)
			}

			wasEmptyLine = util.IsEmptyLine(line)

			if isEof {
				break
			}

			line = ""
		}
	}

	return cmd.Wait()
}
