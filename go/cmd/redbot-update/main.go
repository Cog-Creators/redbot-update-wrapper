package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/cog-creators/redbot-update-wrapper/go/internal/osutils"
	"github.com/cog-creators/redbot-update-wrapper/go/internal/virtualenv"
)

const DefaultProgramName = "redbot-update"

func main() {
	exe, err := osutils.GetExecutableWithPreservedSymlinks(DefaultProgramName)
	if err != nil {
		panic(err)
	}

	venv, err := virtualenv.GetVirtualEnv(exe)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	pythonExe, err := venv.GetPythonExecutable()
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	args := append([]string{"-m", "redbot._update"}, os.Args[1:]...)
	cmd := exec.Command(pythonExe, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		if exitError, ok := errors.AsType[*exec.ExitError](err); ok {
			os.Exit(exitError.ExitCode())
		} else {
			fmt.Printf("Unexpected error occurred while running internal update command:\n%v\n", err)
			os.Exit(1)
		}
	}
}
