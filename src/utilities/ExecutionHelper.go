package utilities

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"sync"
)

// TerminalArgs defines the arguments for executing a command on the terminal.
type TerminalArgs struct {
	Command string
	Args    []string
	EnvArgs []string
}

// ExecuteOnTerminalArgs executes a terminal command with the specified arguments and environment variables.
//
// Parameters:
//   - terminalArgs: TerminalArgs struct containing command, args, and envArgs.
//
// Returns:
//   - error: An error if the execution of the command fails.
func ExecuteOnTerminalArgs(terminalArgs *TerminalArgs) error {
	cmd := generateCmdWithTerminalArgs(terminalArgs)
	return getOutput(cmd)
}

// ExecuteOnTerminal executes a terminal command with the specified command and arguments.
//
// Parameters:
//   - command: The command to be executed.
//   - args: Command-line arguments.
//
// Returns:
//   - error: An error if the execution of the command fails.
func ExecuteOnTerminal(command string, args ...string) error {
	cmd := generateCmd(command, args...)
	return getOutput(cmd)
}

// generateCmdWithTerminalArgs creates an exec.Cmd with the specified terminal arguments.
func generateCmdWithTerminalArgs(terminalArgs *TerminalArgs) *exec.Cmd {
	cmd := exec.Command(terminalArgs.Command, terminalArgs.Args...)
	cmd.Env = append(os.Environ(), terminalArgs.EnvArgs...)
	return cmd
}

// generateCmd creates an exec.Cmd with the specified command and arguments.
func generateCmd(command string, args ...string) *exec.Cmd {
	return exec.Command(command, args...)
}

// getOutput runs the command and captures its output to stdout and stderr.
//
// Steps:
// 1. Create pipes for capturing stdout and stderr.
// 2. Start the command.
// 3. Launch goroutines to capture and print the output from stdout and stderr.
// 4. Wait for the command to complete.
//
// Parameters:
//   - cmd: The exec.Cmd representing the command to be executed.
//
// Returns:
//   - error: An error if the execution of the command fails.
func getOutput(cmd *exec.Cmd) error {
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("error in creating stdout pipe: %v", err)
	}

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("error creating StderrPipe: %v", err)
	}

	err = cmd.Start()
	if err != nil {
		return fmt.Errorf("error starting command: %v", err)
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()

		scanner := bufio.NewScanner(stdoutPipe)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
	}()

	go func() {
		defer wg.Done()

		scanner := bufio.NewScanner(stderrPipe)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
	}()

	err = cmd.Wait()
	if err != nil {
		return fmt.Errorf("the Command finished with error: %v", err)
	}

	wg.Wait()
	return nil
}
