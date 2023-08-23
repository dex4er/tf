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

	"github.com/dex4er/tf/console"
	"github.com/dex4er/tf/progress"
	"github.com/dex4er/tf/util"
)

var TF_IN_AUTOMATION = os.Getenv("TF_IN_AUTOMATION")
var TF_PLAN_FORMAT = os.Getenv("TF_PLAN_FORMAT")
var TF_PROGRESS_FORMAT = os.Getenv("TF_PROGRESS_FORMAT")

func terraformWithProgress(command string, args []string) error {
	patternIgnoreLine := `  Prepared .* for import` +
		`|: Drift detected` +
		`|: Refresh complete after ` +
		`|Acquiring state lock\. This may take a few moments\.\.\.` +
		`|actions to undo or respond to these changes\.` +
		`|and found no differences, so no changes are needed\.` +
		`|Apply complete! Resources: 0 added, 0 changed, 0 destroyed\.` +
		`|at .* \(if it is present\),` +
		`|at .*\.terraformignore (if it is present),` +
		`|excluding files or directories as defined by a \.terraformignore file` +
		`|guarantee to take exactly these actions if you run "terraform apply" now\.` +
		"|guarantee to take exactly these actions if you run `terraform apply` now\\." +
		`|in order to capture the filesystem context the remote workspace expects:` +
		`|Note: You .* use the -out option to save this plan, so Terraform` +
		`|plan. Resource actions are indicated with the following symbols:` +
		`|Preparing the remote plan\.\.\.` +
		`|relative to the target repository.` +
		`|Releasing state lock\. This may take a few moments\.\.\.` +
		`|relevant attributes using ignore_changes, the following plan may include` +
		`|Running plan in Terraform Cloud\. Output will stream here\. Pressing Ctrl-C` +
		`|state, without changing any real infrastructure\.` +
		`|Terraform has checked that the real remote objects still match the result of` +
		`|Terraform has compared your real infrastructure against your configuration` +
		`|Terraform used the selected providers to generate the following execution` +
		`|Terraform will destroy all your managed infrastructure, as shown above\.` +
		`|Terraform will perform the actions described above\.` +
		`|Terraform will perform the following actions:` +
		`|The remote workspace is configured to work with configuration at` +
		`|The resources that were imported are shown above\. These resources are now in` +
		`|There is no undo. Only 'yes' will be accepted to confirm\.` +
		`|these\. If you were expecting these changes then you can apply this plan to` +
		`|This is a refresh-only plan, so Terraform will not take any actions to undo` +
		`|To perform exactly these actions, run the following command to apply:` +
		`|To see the full warning notes, run Terraform without -compact-warnings\.` +
		`|Unless you have made equivalent changes to your configuration, or ignored the` +
		`|Waiting for the plan to start\.\.\.` +
		`|Warnings:` +
		`|will stop streaming the logs, but will not stop the plan running remotely\.` +
		`|You can apply this plan to save these new output values to the Terraform` +
		`|your most recent changes, and found no differences\.` +
		`|your Terraform state and will henceforth be managed by Terraform\.` +
		`|─────────────────────────────────────────────────────────────────────────────`

	patternIgnoreNextLine := ` Experimental feature .* is active` +
		`|record the updated values in the Terraform state without changing any remote` +
		`|record the updated values in the Terraform state without changing any remote` +
		`|Saved the plan to: terraform\.tfplan`

	patternIgnoreBlockStart := `This plan was saved to: terraform.tfplan` +
		`|Warning:.*Applied changes may be incomplete` +
		`|Warning:.*Resource targeting is in effect`

	patternIgnoreBlockEnd := `exceptional situations such as recovering from errors or mistakes` +
		`|suggests to use it as part of an error message` +
		"|terraform apply `terraform\\.tfplan`"

	patternIgnoreShortFormat := `= \(known after apply\)` +
		`|\(\d+ unchanged \w+ hidden\)` +
		`|\(config refers to values not yet known\)` +
		`|read \(data resources\)`

	patternIgnoreShortBlockStart := ` will be read during apply`

	patternIgnoreShortBlockEnd := `^    }`

	patternIgnoreCompactFormat := `^\s\s[\s+~-]` +
		`|<=.* data ".*?" ".*?" \{` +
		`|\(config refers to values not yet known\)` +
		`|\(depends on a resource or a module with changes pending\)` +
		`|read \(data resources\)` +
		`|will be read during apply`

	patternRefreshing := `(?:.\[0m.\[1m)?(.*?): (.)(?:efreshing(?: state)?)\.\.\..*?(?:\r?\n|$)`
	patternStartOperation := `(?:.\[0m.\[1m)?(.*?): (.)(?:eading|reating|estroying|odifying)\.\.\..*?(?:\r?\n|$)`
	patternStillOperation := `(?:.\[0m.\[1m)?(.*?): Still (.).*ing\.\.\..*?(?:\r?\n|$)`
	patternStopOperation := `(?:.\[0m.\[1m)?(.*?): (.)(?:ead|reation|estruction|odifications) complete after.*?(?:\r?\n|$)`

	patternIgnoreOutputs := `^(Changes to )?Outputs:(\n|$)`

	reIgnoreLine := regexp.MustCompile(patternIgnoreLine)
	reIgnoreNextLine := regexp.MustCompile(patternIgnoreNextLine)
	reIgnoreBlockStart := regexp.MustCompile(patternIgnoreBlockStart)
	reIgnoreBlockEnd := regexp.MustCompile(patternIgnoreBlockEnd)
	reIgnoreShortFormat := regexp.MustCompile(patternIgnoreShortFormat)
	reIgnoreShortBlockStart := regexp.MustCompile(patternIgnoreShortBlockStart)
	reIgnoreShortBlockEnd := regexp.MustCompile(patternIgnoreShortBlockEnd)
	reIgnoreCompactFormat := regexp.MustCompile(patternIgnoreCompactFormat)
	reRefreshing := regexp.MustCompile(patternRefreshing)
	reStartOperation := regexp.MustCompile(patternStartOperation)
	reStillOperation := regexp.MustCompile(patternStillOperation)
	reStopOperation := regexp.MustCompile(patternStopOperation)
	reIgnoreOutputs := regexp.MustCompile(patternIgnoreOutputs)

	planFormat := "short"
	progressFormat := "counters"
	noColor := false
	noOutputs := false

	if TF_IN_AUTOMATION == "1" {
		progressFormat = "verbose"
	}

	if TF_PLAN_FORMAT != "" {
		planFormat = TF_PLAN_FORMAT
	}

	if TF_PROGRESS_FORMAT != "" {
		progressFormat = TF_PROGRESS_FORMAT
	}

	newArgs := []string{}

	for _, arg := range args {
		switch util.ReplaceFirstTwoDashes(arg) {
		case "-compact":
			planFormat = "compact"
		case "-counter":
			progressFormat = "counters"
		case "-counters":
			progressFormat = "counters"
		case "-dot":
			progressFormat = "dots"
		case "-dots":
			progressFormat = "dots"
		case "-fan":
			progressFormat = "fan"
		case "-full":
			planFormat = "full"
		case "-no-color":
			noColor = true
			console.NoColor = true
			newArgs = append(newArgs, arg)
		case "-no-colors":
			noColor = true
			console.NoColor = true
			newArgs = append(newArgs, "-no-color")
		case "-no-output":
			noOutputs = true
		case "-no-outputs":
			noOutputs = true
		case "-quiet":
			progressFormat = "quiet"
		case "-short":
			planFormat = "short"
		case "-verbatim":
			planFormat = "full"
			progressFormat = "verbatim"
		case "-verbose":
			progressFormat = "verbose"
		default:
			newArgs = append(newArgs, arg)
		}
	}

	// clear color even after errors
	if !noColor {
		defer fmt.Print("\x1b[0m")
	}

	// original terraform still handles ctrl-c
	signal.Ignore(syscall.SIGINT)

	cmd := execTerraformCommand(append([]string{command}, newArgs...)...)

	cmd.Stdin = os.Stdin

	// errors are written as-is to `TF_OUTPUT_PATH`
	outputFile, err := util.OpenOutputFile()
	if err != nil {
		return err
	}
	if outputFile != nil {
		defer outputFile.Close()
		cmd.Stderr = io.MultiWriter(os.Stderr, outputFile)
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

	reader := bufio.NewReader(cmdStdout)

	isEof := false
	ignoreNextLine := false
	ignoreBlock := false
	skipOutputs := false
	wasEmptyLine := false

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
			if outputFile != nil {
				fmt.Fprint(outputFile, line)
			}

			// skip after "Outputs:" line
			if skipOutputs {
				goto NEXT
			}

			// verbatim progress format is not processed or ignored
			if progressFormat != "verbatim" {
				if m := reRefreshing.FindStringSubmatch(line); m != nil {
					line := m[0]
					line = strings.TrimSuffix(line, "\n")
					line = strings.TrimSuffix(line, "\r")
					progress.Refresh(progressFormat, line, m[1], m[2])
					goto NEXT
				}

				if m := reStartOperation.FindStringSubmatch(line); m != nil {
					line := m[0]
					line = strings.TrimSuffix(line, "\n")
					line = strings.TrimSuffix(line, "\r")
					progress.Start(progressFormat, line, m[1], m[2])
					goto NEXT
				}

				if m := reStillOperation.FindStringSubmatch(line); m != nil {
					line := m[0]
					line = strings.TrimSuffix(line, "\n")
					line = strings.TrimSuffix(line, "\r")
					progress.Still(progressFormat, line, m[1], m[2])
					goto NEXT
				}

				if m := reStopOperation.FindStringSubmatch(line); m != nil {
					line := m[0]
					line = strings.TrimSuffix(line, "\n")
					line = strings.TrimSuffix(line, "\r")
					progress.Stop(progressFormat, line, m[1], m[2])
					goto NEXT
				}
			}

			// dot format trims EOL then we need one extra just before "Apply complete!"
			if progressFormat == "dots" && strings.HasPrefix(line, "Apply complete!") {
				fmt.Println()
			}

			// handles block to ignore: it has start and end pattern
			if reIgnoreBlockStart.MatchString(line) {
				ignoreBlock = true
				goto NEXT
			}

			if planFormat == "short" && reIgnoreShortBlockStart.MatchString(line) {
				ignoreBlock = true
				goto NEXT
			}

			if reIgnoreBlockEnd.MatchString(line) {
				ignoreBlock = false
				goto NEXT
			}

			if planFormat == "short" && reIgnoreShortBlockEnd.MatchString(line) {
				ignoreBlock = false
				goto NEXT
			}

			if ignoreBlock {
				goto NEXT
			}

			// handles pattern that causes ignoring this and next line
			if reIgnoreNextLine.MatchString(line) {
				ignoreNextLine = true
				goto NEXT
			}

			if ignoreNextLine {
				ignoreNextLine = false
				goto NEXT
			}

			// ignore just this line
			if reIgnoreLine.MatchString(line) {
				goto NEXT
			}

			// starts ignoring the outputs
			if noOutputs && reIgnoreOutputs.MatchString(line) {
				skipOutputs = true
				goto NEXT
			}

			// handles different plan formats
			if planFormat == "short" && reIgnoreShortFormat.MatchString(line) {
				goto NEXT
			}

			if planFormat == "compact" && reIgnoreCompactFormat.MatchString(line) {
				goto NEXT
			}

			// skip another empty line but preserve color codes
			if wasEmptyLine && util.IsEmptyLine(line) {
				line = strings.TrimSuffix(line, "\n")
				line = strings.TrimSuffix(line, "\r")
				line = strings.ReplaceAll(line, "╵", "")
				line = strings.ReplaceAll(line, "╷", "")
			}

			// in CI do not add spaces clearing progress indicator
			if TF_IN_AUTOMATION == "1" {
				fmt.Print(line)
			} else {
				console.Print(line)
			}

			// mark if current line was empty for next loop iteration
			wasEmptyLine = util.IsEmptyLine(line)

		NEXT:
			// empty line buffer before next loop iteration
			line = ""
		}
	}

	return cmd.Wait()
}
