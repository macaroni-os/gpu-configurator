/*
Copyright © 2024 Macaroni OS Linux
See AUTHORS and LICENSE for the license details and contributors.
*/
package kernel

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/macaroni-os/gpu-configurator/pkg/logger"
	"github.com/macaroni-os/macaronictl/pkg/utils"
)

func Depmod(kernelVersion string, flags []string) error {
	log := logger.GetDefaultLogger()
	depmodBin := utils.TryResolveBinaryAbsPath("depmod")
	args := []string{
		depmodBin, "-a",
	}

	if len(flags) > 0 {
		args = append(args, flags...)
	}

	args = append(args, kernelVersion)

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	log.InfoC(fmt.Sprintf("Running %s...",
		strings.Join(args, " ")))

	err := cmd.Start()
	if err != nil {
		return err
	}

	err = cmd.Wait()
	if err != nil {
		return err
	}

	if cmd.ProcessState.ExitCode() != 0 {
		return fmt.Errorf("depmod exiting with %s",
			cmd.ProcessState.ExitCode())
	}

	return nil
}

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
