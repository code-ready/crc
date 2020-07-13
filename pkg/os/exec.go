package os

import (
	"bytes"
	"os"
	"os/exec"
	"strings"

	"github.com/code-ready/crc/pkg/crc/logging"
)

func runCmd(command string, args []string, env map[string]string) (string, string, error) {
	cmd := exec.Command(command, args...) // #nosec G204
	if len(env) != 0 {
		cmd.Env = os.Environ()
		for key, value := range env {
			cmd.Env = ReplaceOrAddEnv(cmd.Env, key, value)
		}
	}
	stdOut := new(bytes.Buffer)
	stdErr := new(bytes.Buffer)
	cmd.Stdout = stdOut
	cmd.Stderr = stdErr
	err := cmd.Run()
	if err != nil {
		logging.Debugf("Command failed: %v", err)
		logging.Debugf("stdout: %s", stdOut.String())
		logging.Debugf("stderr: %s", stdErr.String())
	}
	return stdOut.String(), stdErr.String(), err
}

func run(command string, args []string, env map[string]string) (string, string, error) {
	logging.Debugf("Running '%s %s'", command, strings.Join(args, " "))
	return runCmd(command, args, env)
}

func runPrivate(command string, args []string, env map[string]string) (string, string, error) {
	logging.Debugf("About to run a hidden command")
	return runCmd(command, args, env)
}

// RunWithPrivilege executes a command using sudo
// provide a reason why root is needed as the first argument
func RunWithPrivilege(reason string, cmdAndArgs ...string) (string, string, error) {
	sudo, err := exec.LookPath("sudo")
	if err != nil {
		return "", "", err
	}
	logging.Infof("Will use root access: %s", reason)
	return run(sudo, cmdAndArgs, map[string]string{})
}

func RunWithDefaultLocale(command string, args ...string) (string, string, error) {
	return RunWithDefaultLocaleAndEnv(command, args, map[string]string{})
}

func RunWithDefaultLocalePrivate(command string, args ...string) (string, string, error) {
	return RunWithDefaultLocalePrivateAndEnv(command, args, map[string]string{})
}

func RunWithDefaultLocaleAndEnv(command string, args []string, extraEnv map[string]string) (string, string, error) {
	env := map[string]string{"LC_ALL": "C", "LANG": "C"}
	for k, v := range extraEnv {
		env[k] = v
	}
	return run(command, args, env)
}

func RunWithDefaultLocalePrivateAndEnv(command string, args []string, extraEnv map[string]string) (string, string, error) {
	env := map[string]string{"LC_ALL": "C", "LANG": "C"}
	for k, v := range extraEnv {
		env[k] = v
	}
	return runPrivate(command, args, env)
}
