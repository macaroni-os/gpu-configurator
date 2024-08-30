/*
Copyright © 2024 Macaroni OS Linux
See AUTHORS and LICENSE for the license details and contributors.
*/
package macaroni

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/macaroni-os/gpu-configurator/pkg/kernel"
	"github.com/macaroni-os/gpu-configurator/pkg/logger"
	"github.com/macaroni-os/gpu-configurator/pkg/specs"
	"github.com/macaroni-os/macaronictl/pkg/utils"
)

var (
	extraNvidiaModules = []string{
		"nvidia-drm.ko",
		"nvidia-modeset.ko",
		"nvidia-uvm.ko",
	}
)

func (b *MacaroniBackend) ActiveNVIDIAKernelDriver(setup *specs.NVIDIASetup,
	nvidiaVersion string, kernelVersion string, open bool) error {
	log := logger.GetDefaultLogger()

	log.DebugC(fmt.Sprintf(
		"Trying to active nvidia kernel driver for kernel %s and version %s (open = %v)",
		kernelVersion, nvidiaVersion, open))

	// NOTE: I active a kernel selected only if it matches with
	//       the driver type (open/proprietary).

	// Retrieve list of kernel modules with the selected version
	kmodules := setup.GetKernelModulesAvailable(nvidiaVersion, open)
	if len(*kmodules) == 0 {
		log.InfoC(fmt.Sprintf(
			"No kernel modules available for nvidia version %s",
			nvidiaVersion))
		return nil
	}

	// POST: Is not present an already active kernel module
	//       with selected filter.
	var kmod2link *specs.KernelModule = nil

	for _, k := range *kmodules {
		if k.KernelVersion == kernelVersion {
			kmod2link = k
			break
		}
	}

	if kmod2link == nil {
		return fmt.Errorf(
			"No available kernel module with version %s for nvidia version %s.",
			kernelVersion, nvidiaVersion)
	}

	kactive := setup.GetKernelModulesActive(nvidiaVersion, kernelVersion)

	if kactive != nil {
		if kactive.IsOpen() == open && open == true {
			log.InfoC(fmt.Sprintf(
				"Open kernel module %s for nvidia version %s already present.",
				kernelVersion, nvidiaVersion,
			))
		} else if (!kactive.IsOpen()) && !open && open == false {
			log.InfoC(fmt.Sprintf(
				"Proprietary kernel module %s for nvidia version %s already present.",
				kernelVersion, nvidiaVersion,
			))
		} else {
			log.InfoC(fmt.Sprintf(
				"Kernel module %s for nvidia version %s (%v) already "+
					"present by with a different type.",
				kernelVersion, nvidiaVersion, open,
			))
		}
	} else {
		// Create hardlinks
		err := b.activeKernelModule(setup, kmod2link)
		if err != nil {
			return err
		}
	}

	return nil
}

func (b *MacaroniBackend) activeKernelModule(setup *specs.NVIDIASetup, km *specs.KernelModule) error {
	log := logger.GetDefaultLogger()
	modulePath := "/lib/modules"
	nvidiaKmoduleDir := filepath.Join(modulePath, km.KernelVersion, "video")
	nvidiakmName := filepath.Base(km.Path)
	targetFile := filepath.Join(nvidiaKmoduleDir, nvidiakmName)

	if !utils.Exists(nvidiaKmoduleDir) {
		log.DebugC(fmt.Sprintf(
			"Creating directory %s...",
			nvidiaKmoduleDir))
		err := os.MkdirAll(nvidiaKmoduleDir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	log.DebugC(fmt.Sprintf(
		"Creating hardlink of %s to %s...",
		km.Path, targetFile))
	err := os.Link(km.Path, targetFile)
	if err != nil {
		return err
	}

	ext := filepath.Ext(km.Path)
	sourceDir := filepath.Dir(km.Path)

	// Create hardlink of all modules
	for _, em := range extraNvidiaModules {

		link := filepath.Join(nvidiaKmoduleDir, em)
		f := filepath.Join(sourceDir, em)
		if ext != ".ko" {
			f += ext
			link += ext
		}

		if utils.Exists(f) {
			log.DebugC(fmt.Sprintf(
				"Creating hardlink of %s to %s...",
				f, link))
			err := os.Link(f, link)
			if err != nil {
				return err
			}
		} else {
			log.DebugC(fmt.Sprintf(
				"Kernel module %s not found.",
				f))
		}

	}

	// After that the driver is been installed we need
	// to rebuild the kernel symbols with depmod
	return kernel.Depmod(km.KernelVersion, []string{})
}

func (b *MacaroniBackend) PurgeNVIDIAKernelDriverActive(setup *specs.NVIDIASetup,
	nvidiaVersion, kernelVersion, driverType string) error {
	log := logger.GetDefaultLogger()

	log.DebugC(fmt.Sprintf(
		"Trying to purge nvidia kernel driver for kernel %s and version %s",
		kernelVersion, nvidiaVersion))

	// Check driver between kernel active modules
	if (driverType == "all" || driverType == "proprietary") &&
		len(setup.KModuleActive) > 0 {

		err := b.purgeKernelDriverFromList(&setup.KModuleActive,
			nvidiaVersion, kernelVersion)
		if err != nil {
			return err
		}
	}

	// Check driver between kernel active open modules
	if (driverType == "all" || driverType == "open") &&
		len(setup.KOpenModuleActive) > 0 {

		err := b.purgeKernelDriverFromList(&setup.KOpenModuleActive,
			nvidiaVersion, kernelVersion)
		if err != nil {
			return err
		}
	}

	return nil
}

func (b *MacaroniBackend) purgeKernelDriverFromList(listRef *[]*specs.KernelModule,
	nvidiaVersion string, kernelVersion string) error {
	log := logger.GetDefaultLogger()

	purgedDrivers := 0

	for _, kmodule := range *listRef {
		if kmodule.KernelVersion != kernelVersion &&
			(kernelVersion != "" && kernelVersion != "*") {
			log.DebugC(fmt.Sprintf(
				"Ignoring module %s (%s)",
				kmodule.KernelVersion, kernelVersion))
			continue
		}

		kv, _ := kmodule.Fields["version"]
		if nvidiaVersion != "" && nvidiaVersion != "*" {
			// POST: Check if version match
			if kv != nvidiaVersion {
				log.DebugC(fmt.Sprintf(
					"Ignoring module with a mismatch version %s != %s",
					kv, nvidiaVersion))

				continue
			}
		}

		// POST: nvidiaVersion == "" || nvidiaVersion == "*" || nvidiaVersion matches

		log.InfoC(fmt.Sprintf(
			"Removing kernel driver for kernel %s and NVIDIA version %s...",
			kmodule.KernelVersion, kv))

		err := b.removeFileIfExist(kmodule.Path)
		if err != nil {
			return err
		}

		kernelDir := filepath.Dir(kmodule.Path)

		for _, em := range extraNvidiaModules {
			// NOTE: I try to remove all files to cleanup stalled condition.
			f := filepath.Join(kernelDir, em)
			err := b.removeFileIfExist(f)
			if err != nil {
				return err
			}

			for _, c := range KernelModuleSupportedCompression {
				if strings.HasSuffix(kmodule.Path, c) {
					fc := f + c
					err := b.removeFileIfExist(fc)
					if err != nil {
						return err
					}
				}
			}

		}

		purgedDrivers++

		// After that the driver is been removed we need
		// to rebuild the kernel symbols with depmod
		err = kernel.Depmod(kmodule.KernelVersion, []string{})
		if err != nil {
			return err
		}
	}

	if purgedDrivers == 0 {
		if kernelVersion != "" && kernelVersion != "*" {
			log.InfoC(fmt.Sprintf(
				"No kernel driver found for version %s and NVIDIA version %s.",
				kernelVersion, nvidiaVersion,
			))
		} else {
			log.InfoC(fmt.Sprintf(
				"No kernel driver found for NVIDIA version %s.",
				nvidiaVersion,
			))
		}
	}

	return nil
}
