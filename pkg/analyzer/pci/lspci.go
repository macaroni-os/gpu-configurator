/*
Copyright © 2024 Macaroni OS Linux
See AUTHORS and LICENSE for the license details and contributors.
*/
package pci

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/macaroni-os/macaronictl/pkg/utils"
)

type PCIDevice struct {
	BusId             string   `json:"bus_id,omitempty" yaml:"bus_id,omitempty"`
	ClassName         string   `json:"class_name,omitempty" yaml:"class_name,omitempty"`
	ClassId           string   `json:"class_id,omitempty" yaml:"class_id,omitempty"`
	Name              string   `json:"name,omitempty" yaml:"name,omitempty"`
	Id                string   `json:"id,omitempty" yaml:"id,omitempty"`
	Subsystem         string   `json:"subsystem,omitempty" yaml:"subsystem,omitempty"`
	DeviceName        string   `json:"device_name,omitempty" yaml:"device_name,omitempty"`
	KernelDriverInUse string   `json:"kernel_driver_inuse,omitempty" yaml:"kernel_driver_inuse,omitempty"`
	KernelModules     []string `json:"kernel_modules,omitempty" yaml:"kernel_modules,omitempty"`
}

type SystemDevices []*PCIDevice

func GetDevices() (*SystemDevices, error) {
	var errBuffer bytes.Buffer
	var outBuffer bytes.Buffer

	lspciBin := utils.TryResolveBinaryAbsPath("lspci")
	args := []string{
		lspciBin, "-nn", "-k",
	}

	cmd := exec.Command(args[0], args[1:]...)

	cmd.Stdout = utils.NewNopCloseWriter(&outBuffer)
	cmd.Stderr = utils.NewNopCloseWriter(&errBuffer)

	err := cmd.Start()
	if err != nil {
		return nil, err
	}

	err = cmd.Wait()
	if err != nil {
		return nil, err
	}

	if cmd.ProcessState.ExitCode() != 0 {
		return nil, fmt.Errorf("lspci exiting with %s: %s",
			cmd.ProcessState.ExitCode(),
			errBuffer.String())
	}

	devices, err := parseLspciOutput(outBuffer.String())
	if err != nil {
		return nil, fmt.Errorf("error on parsing lspci output: %s",
			err.Error())
	}

	ans := SystemDevices(*devices)

	return &ans, nil
}

func parseLspciOutput(output string) (*[]*PCIDevice, error) {
	ans := []*PCIDevice{}
	var lastPCIDevice *PCIDevice

	lines := []string{}
	// Convert output in lines
	scanner := bufio.NewScanner(strings.NewReader(output))

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	// Example of block output
	// 05:00.0 VGA compatible controller [0300]: Advanced Micro Devices, Inc. [AMD/ATI] Picasso [1002:15d8] (rev c1)
	//	Subsystem: ASUSTeK Computer Inc. Device [1043:18f1]
	//	Kernel driver in use: amdgpu
	//	Kernel modules: amdgpu
	for _, line := range lines {
		words := strings.Split(line, " ")
		if line[0] == '\t' {
			// POST: we reading data about previous bus/device

			if strings.HasPrefix(line, "\tSubsystem") {
				lastPCIDevice.Subsystem = line[len("Subsystem:")+2:]
			} else if strings.HasPrefix(line, "\tKernel driver") {
				lastPCIDevice.KernelDriverInUse = words[4]
			} else if strings.HasPrefix(line, "\tDeviceName") {
				lastPCIDevice.DeviceName = strings.TrimSpace(line[len("DeviceName:")+2:])
			} else {
				// POST: parse kernel modules
				modules := ""
				lenm := len(words)
				for i := 2; i < lenm; i++ {
					modules += words[i]
				}
				lastPCIDevice.KernelModules = strings.Split(modules, ",")
			}

		} else {
			// POST: reading line with bus id

			lastPCIDevice = &PCIDevice{
				BusId:         words[0],
				KernelModules: []string{},
			}

			lastPCIDevice.ClassName =
				strings.TrimSpace(
					line[len(words[0])+1 : strings.Index(line, "[")],
				)

			classNameWords := strings.Split(lastPCIDevice.ClassName, " ")
			pos := len(classNameWords) + 1
			lastPCIDevice.ClassId = strings.ReplaceAll(
				strings.ReplaceAll(
					strings.ReplaceAll(words[pos], ":", ""),
					"]", ""), "[", "")

			posId := len(words) - 1
			// Retriev id pos
			if strings.Index(words[posId], ")") > 0 {
				posId -= 2
				// POST: last words is related to (rev XX) word
			}
			lastPCIDevice.Id = strings.ReplaceAll(
				strings.ReplaceAll(words[posId], "]", ""),
				"[", "")

			pos++
			for i := 0; pos < posId; i++ {
				if i == 0 {
					lastPCIDevice.Name = words[pos]
				} else {
					lastPCIDevice.Name = lastPCIDevice.Name + " " + words[pos]
				}
				pos++
			}

			ans = append(ans, lastPCIDevice)
		}
	}

	return &ans, nil
}
