package executor

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
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
	cmd.Env = buildEnv(envs)
	// stdout is reserved for the rift shell eval (cd command only);
	// route all command output through stderr so it reaches the terminal directly.
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd
}

// ----------------------------------
//
//	ExecWithPadding runs args and prefixes every output line with padding on
//	stderr. Falls back to plain Run when padding is empty or pipe fails.
//
// ----------------------------------
func (c *cmdExecutor) ExecWithPadding(args []string, executionPath string, envs []string, padding string) {
	cmd := c.RunCmd(args, executionPath, envs)
	if cmd == nil {
		return
	}
	if padding == "" {
		cmd.Run()
		return
	}
	pr, pw, err := os.Pipe()
	if err != nil {
		cmd.Run()
		return
	}
	defer pr.Close()
	cmd.Stdout = pw
	cmd.Stderr = pw
	if err := cmd.Start(); err != nil {
		pw.Close()
		return
	}
	pw.Close()
	scanner := bufio.NewScanner(pr)
	for scanner.Scan() {
		fmt.Fprintf(os.Stderr, "%s%s\n", padding, scanner.Text())
	}
	cmd.Wait()
}

// buildEnv returns a clean environment for a child process: strips the
// inherited RIFT_RUNE_DEPTH, reads its value (default 0), increments it by 1,
// injects the new value, then appends any caller-supplied overrides.
func buildEnv(envs []string) []string {
	depth := 0
	base := make([]string, 0, len(os.Environ())+1)
	for _, e := range os.Environ() {
		if val, ok := strings.CutPrefix(e, "RIFT_RUNE_DEPTH="); ok {
			if n, err := strconv.Atoi(val); err == nil {
				depth = n
			}
			continue
		}
		base = append(base, e)
	}
	base = append(base, fmt.Sprintf("RIFT_RUNE_DEPTH=%d", depth+1))
	return append(base, envs...)
}
