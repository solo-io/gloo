package helpers

import (
	"bytes"
	"io"
	"os"
	"os/exec"

	"github.com/pkg/errors"
)

func RunCommand(verbose bool, args ...string) error {
	_, err := RunCommandOutput(verbose, args...)
	return err
}

func RunCommandOutput(verbose bool, args ...string) (string, error) {
	return RunCommandInputOutput("", verbose, args...)
}

func RunCommandInput(input string, verbose bool, args ...string) error {
	_, err := RunCommandInputOutput(input, verbose, args...)
	return err
}

func RunCommandInputOutput(input string, verbose bool, args ...string) (string, error) {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = GlooDir()
	cmd.Env = os.Environ()
	if len(input) > 0 {
		cmd.Stdin = bytes.NewBuffer([]byte(input))
	}
	buf := &bytes.Buffer{}
	var out io.Writer
	if verbose {
		out = io.MultiWriter(buf, os.Stdout, os.Stderr)
	} else {
		out = buf
	}
	cmd.Stdout = out
	cmd.Stderr = out
	if err := cmd.Run(); err != nil {
		return "", errors.Wrapf(err, "%v failed: %s", cmd.Args, buf.String())
	}

	return buf.String(), nil
}
