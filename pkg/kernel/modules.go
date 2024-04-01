/*
Copyright © 2024 Macaroni OS Linux
See AUTHORS and LICENSE for the license details and contributors.
*/
package kernel

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/macaroni-os/macaronictl/pkg/utils"
)

func ModinfoField(fpath, field string) (string, error) {
	var errBuffer bytes.Buffer
	var outBuffer bytes.Buffer

	modinfoBin := utils.TryResolveBinaryAbsPath("modinfo")
	args := []string{
		modinfoBin, "-F", field, fpath,
	}

	cmd := exec.Command(args[0], args[1:]...)

	cmd.Stdout = utils.NewNopCloseWriter(&outBuffer)
	cmd.Stderr = utils.NewNopCloseWriter(&errBuffer)

	err := cmd.Start()
	if err != nil {
		return "", err
	}

	err = cmd.Wait()
	if err != nil {
		return "", err
	}

	if cmd.ProcessState.ExitCode() != 0 {
		return "", fmt.Errorf("modinfo exiting with %s: %s",
			cmd.ProcessState.ExitCode(),
			errBuffer.String())
	}

	ans := strings.TrimSpace(
		strings.ReplaceAll(outBuffer.String(), "\n", ""),
	)

	return ans, nil
}

func GetRuntimeKernelVersion() (string, error) {
	var errBuffer bytes.Buffer
	var outBuffer bytes.Buffer

	binary := utils.TryResolveBinaryAbsPath("uname")
	args := []string{
		binary, "-r",
	}

	cmd := exec.Command(args[0], args[1:]...)

	cmd.Stdout = utils.NewNopCloseWriter(&outBuffer)
	cmd.Stderr = utils.NewNopCloseWriter(&errBuffer)

	err := cmd.Start()
	if err != nil {
		return "", err
	}

	err = cmd.Wait()
	if err != nil {
		return "", err
	}

	if cmd.ProcessState.ExitCode() != 0 {
		return "", fmt.Errorf("%s exiting with %s: %s",
			binary, cmd.ProcessState.ExitCode(),
			errBuffer.String())
	}

	ans := strings.TrimSpace(
		strings.ReplaceAll(outBuffer.String(), "\n", ""),
	)

	return ans, nil
}
