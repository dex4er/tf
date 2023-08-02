package main

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
)

var version = "dev"

var emptyStringsList = make([]string, 0)

func main() {
	var err error

	if len(os.Args) < 2 {
		err = runTerraformHelp(emptyStringsList)
	} else {
		command := os.Args[1]
		args := os.Args[2:]

		switch command {
		case "apply":
			err = runTerraformApply(args)
		case "init":
			err = runTerraformInit(args)
		case "upgrade":
			err = runTerraformUpgrade(args)
		case "version":
			err = runTerraformVersion(args)
		default:
			err = runTerraformCommand(command, args)
		}
	}

	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			exitCode := exitError.ExitCode()
			os.Exit(exitCode)
		} else {
			fmt.Println("Error:", err)
			os.Exit(2)
		}
	}
}

func runTerraformHelp(args []string) error {
	err := runTerraformCommand("", args)

	fmt.Println()
	fmt.Println("Extra commands:")
	fmt.Println("  upgrade       The same as terraform init -upgrade")

	return err
}

func runTerraformVersion(args []string) error {
	fmt.Println("tf", version)

	return runTerraformCommand("version", args)
}

func runTerraformApply(args []string) error {
	return runTerraformCommandWithProgress("apply", args)
}

func runTerraformInit(args []string) error {
	ignorePattern := "Finding .* versions matching" +
		"|Initializing Terraform" +
		"|Initializing (modules" +
		"|the backend" +
		"|provider plugins)\\.\\.\\." +
		"|Upgrading modules\\.\\.\\." +
		"|Using previously-installed" +
		"|Reusing previous version of" +
		"|from the shared cache directory" +
		"|in \\.terraform/modules/" +
		"|Finding latest version of" +
		"|Partner and community providers are signed by their developers\\." +
		"|If you'd like to know more about provider signing, you can read about it here:" +
		"|https://www\\.terraform\\.io/docs/cli/plugins/signing\\.html" +
		"|Terraform has made some changes to the provider dependency selections recorded" +
		"|in the \\.terraform\\.lock\\.hcl file. Review those changes and commit them to your" +
		"|version control system if they represent changes you intended to make\\."
	endPattern := "Terraform.* has been successfully initialized!"
	return runTerraformCommandWithFilter("init", args, ignorePattern, endPattern)
}

func runTerraformUpgrade(args []string) error {
	return runTerraformInit(append(args, "-upgrade"))
}

func runTerraformCommand(command string, args []string) error {
	cmd := exec.Command("terraform", append([]string{command}, args...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func runTerraformCommandWithProgress(command string, args []string) error {
	return runTerraformCommandWithFilter(command, args, "\\(known after apply\\)", "")
}

func runTerraformCommandWithFilter(command string, args []string, ignorePattern string, endPattern string) error {
	defer fmt.Printf("\033[0m")

	signal.Ignore(syscall.SIGINT)

	cmd := exec.Command("terraform", append([]string{command}, args...)...)

	cmd.Stdin = os.Stdin

	file, err := openLogfile()
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

		if strings.Contains(line, "\x1b[1mEnter a value:\x1b[0m \x1b[0m") || strings.Contains(line, "Enter a value: ") || r == '\n' || isEof {
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

func openLogfile() (*os.File, error) {
	if logFilename := os.Getenv("TF_LOG_FILE"); logFilename != "" {
		file, err := os.OpenFile(logFilename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			return nil, fmt.Errorf("opening the log file: %w", err)
		}
		return file, nil
	}

	return nil, nil
}
