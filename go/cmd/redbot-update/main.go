package main

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/exec"

	"github.com/cog-creators/redbot-update-wrapper/go/internal/osutils"
	"github.com/cog-creators/redbot-update-wrapper/go/internal/virtualenv"
)

const (
	DefaultProgramName = "redbot-update"

	envVarNamePrefix   = "REDBOT_UPDATE_WRAPPER_"
	LogDebugEnvVarName = envVarNamePrefix + "LOG_DEBUG"
	LogFileEnvVarName  = envVarNamePrefix + "LOG_FILE"
)

func main() {
	handlerOptions := &slog.HandlerOptions{}
	if os.Getenv(LogDebugEnvVarName) == "1" {
		handlerOptions.Level = slog.LevelDebug
	}
	handlers := []slog.Handler{slog.NewTextHandler(os.Stderr, handlerOptions)}

	if filename := os.Getenv(LogFileEnvVarName); filename != "" {
		logFile, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			panic(err)
		}
		defer logFile.Close()
		handlers = append(handlers, slog.NewTextHandler(logFile, handlerOptions))
	}
	logger := slog.New(slog.NewMultiHandler(handlers...))
	slog.SetDefault(logger)

	exe, err := osutils.GetExecutableWithPreservedSymlinks(DefaultProgramName)
	if err != nil {
		panic(err)
	}
	slog.Debug("Found executable", "executable", exe)

	venv, err := virtualenv.GetVirtualEnv(exe)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
	slog.Debug("Found virtual environment", "venv", venv)

	pythonExe, err := venv.GetPythonExecutable()
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
	slog.Debug("Found Python executable", "python_executable", pythonExe)

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
