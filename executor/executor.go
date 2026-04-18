package executor

import (
	"context"
	"os"
	"os/exec"
)

type cmdExecutor struct{}

// ----------------------------------
//
//	Returns a new cmdExecutor instance.
//
// ----------------------------------
func CmdExecutor() *cmdExecutor {
	return &cmdExecutor{}
}

// ----------------------------------
//
//	Builds an exec.Cmd from args, setting its working directory to
//	executionPath and routing both stdout and stderr to the terminal's
//	stderr so that command output reaches the user without polluting the
//	stdout channel reserved for the rift shell eval.
//	Returns nil when args is empty.
//
// ----------------------------------
func (c *cmdExecutor) RunCmd(args []string, executionPath string, envs []string) *exec.Cmd {
	var argName string
	var argsArray []string
	if len(args) > 1 {
		argName = args[0]
		argsArray = args[1:]
	} else if len(args) == 1 {
		argName = args[0]
	} else {
		return nil
	}

	cmd := exec.Command(argName, argsArray...)
	cmd.Dir = executionPath
	cmd.Env = append(cmd.Env, envs...)
	// stdout is reserved for the rift shell eval (cd command only);
	// route all command output through stderr so it reaches the terminal directly.
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	return cmd
}

// ----------------------------------
//
//	RunCmdWithContext creates a command that respects context cancellation.
//	When the context is cancelled, the command will be terminated automatically.
//
// ----------------------------------
func (c *cmdExecutor) RunCmdWithContext(ctx context.Context, args []string, executionPath string, envs []string) *exec.Cmd {
	var argName string
	var argsArray []string
	if len(args) > 1 {
		argName = args[0]
		argsArray = args[1:]
	} else if len(args) == 1 {
		argName = args[0]
	} else {
		return nil
	}

	cmd := exec.CommandContext(ctx, argName, argsArray...)
	cmd.Dir = executionPath
	cmd.Env = append(cmd.Env, envs...)
	// stdout is reserved for the rift shell eval (cd command only);
	// route all command output through stderr so it reaches the terminal directly.
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	return cmd
}
