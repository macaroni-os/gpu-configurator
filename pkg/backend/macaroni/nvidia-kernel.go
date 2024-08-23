/*
Copyright © 2024 Macaroni OS Linux
See AUTHORS and LICENSE for the license details and contributors.
*/
package macaroni

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/macaroni-os/gpu-configurator/pkg/logger"
	"github.com/macaroni-os/gpu-configurator/pkg/specs"
)

func (b *MacaroniBackend) PurgeNVIDIAKernelDriverActive(setup *specs.NVIDIASetup,
	nvidiaVersion string, kernelVersion string) error {
	log := logger.GetDefaultLogger()

	log.DebugC(fmt.Sprintf(
		"Trying to purge nvidia kernel driver for kernel %s and version %s",
		kernelVersion, nvidiaVersion))

	// Check driver between kernel active modules
	if len(setup.KModuleActive) > 0 {

		err := b.purgeKernelDriverFromList(&setup.KModuleActive,
			nvidiaVersion, kernelVersion)
		if err != nil {
			return err
		}
	}

	// Check driver between kernel active open modules
	if len(setup.KOpenModuleActive) > 0 {

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

	extraModules := []string{
		"nvidia-drm.ko",
		"nvidia-modeset.ko",
		"nvidia-uvm.ko",
	}

	log := logger.GetDefaultLogger()
	for _, kmodule := range *listRef {
		if kmodule.KernelVersion != kernelVersion {
			continue
		}

		kv, _ := kmodule.Fields["version"]
		if nvidiaVersion != "" && nvidiaVersion != "*" {
			// POST: Check if version match
			if kv != nvidiaVersion {
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

		for _, em := range extraModules {
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

	}

	return nil
}
