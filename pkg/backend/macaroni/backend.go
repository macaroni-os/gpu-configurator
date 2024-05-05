/*
Copyright © 2024 Macaroni OS Linux
See AUTHORS and LICENSE for the license details and contributors.
*/
package backend

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/macaroni-os/gpu-configurator/pkg/kernel"
	"github.com/macaroni-os/gpu-configurator/pkg/specs"

	"github.com/macaroni-os/macaronictl/pkg/utils"
)

type MacaroniBackend struct {
	Name string
}

func NewMacaroniBackend() (*MacaroniBackend, error) {
	return &MacaroniBackend{
		Name: "macaroni",
	}, nil
}

func (b *MacaroniBackend) GetName() string {
	return b.Name
}

func (b *MacaroniBackend) GetEglExternalPlatformsDirs() ([]string, error) {
	return []string{
		"/usr/share/egl/egl_external_platform.d",
	}, nil
}

func (b *MacaroniBackend) GetVulkanLayersDirs() ([]string, error) {
	return []string{
		"/usr/share/vulkan/explicit_layer.d",
		"/usr/share/vulkan/implicit_layer.d",
	}, nil
}

func (b *MacaroniBackend) GetVulkanICDDirs() ([]string, error) {
	return []string{
		"/usr/share/vulkan/icd.d",
		"/etc/vulkan/icd.d",
	}, nil
}

func (b *MacaroniBackend) GetEnvironmentDir() string { return "/etc/env.d" }

func (b *MacaroniBackend) GetGBMLibDir() string { return "/usr/lib64/gbm" }

func (b *MacaroniBackend) GetNVIDIAEglWaylandLibDir() string { return "/usr/lib64" }
func (b *MacaroniBackend) GetNVIDIAEglGbmLibDir() string     { return "/usr/lib64" }

func (b *MacaroniBackend) GetNVIDIAKernelModules() (*[]*specs.KernelModule, error) {
	modulePath := "/lib/modules"
	ans := []*specs.KernelModule{}

	dirEntries, err := os.ReadDir(modulePath)
	if err != nil {
		return nil, err
	}

	for _, file := range dirEntries {
		if !file.IsDir() {
			continue
		}

		kVersion := file.Name()

		nvidiaKmoduleDir := filepath.Join(
			modulePath, kVersion, "video")
		nvidiaKModule := filepath.Join(nvidiaKmoduleDir, "nvidia.ko.zst")

		kversion := ""
		if utils.Exists(nvidiaKModule) {
			kversion, _ = kernel.ModinfoField(nvidiaKModule, "version")
		} else {
			nvidiaKModule := filepath.Join(nvidiaKmoduleDir, "nvidia.ko")
			if utils.Exists(nvidiaKModule) {
				kversion, _ = kernel.ModinfoField(nvidiaKModule, "version")
			}
		}

		if kversion != "" {
			lp := &specs.KernelModule{
				Path:          nvidiaKModule,
				KernelVersion: kVersion,
				Name:          "nvidia",
				Fields:        make(map[string]string, 0),
			}
			lp.Fields["version"] = kversion
			ans = append(ans, lp)
		}

	}

	return &ans, nil
}

func (b *MacaroniBackend) GetNVIDIADrivers() (*[]*specs.NVIDIADriver, error) {
	ans := []*specs.NVIDIADriver{}

	prefixDriverPath := "/opt/nvidia"
	dirPrefix := "nvidia-drivers"

	if !utils.Exists(prefixDriverPath) {
		// POST: no nvidia drivers available
		return &ans, nil
	}

	dirEntries, err := os.ReadDir(prefixDriverPath)
	if err != nil {
		return nil, err
	}

	// Retrieve current kernel version
	kVersion, err := kernel.GetRuntimeKernelVersion()
	if err != nil {
		return nil, err
	}

	for _, file := range dirEntries {
		if !file.IsDir() {
			continue
		}

		if !strings.HasPrefix(file.Name(), dirPrefix) {
			continue
		}

		version := file.Name()[len(dirPrefix)+1:]
		driverDir := &specs.NVIDIADriver{
			Path:    filepath.Join(prefixDriverPath, file.Name()),
			Version: version,
		}

		nvidiaKmoduleDir := filepath.Join(
			"/lib/modules/", kVersion, "video")
		nvidiaKModule := filepath.Join(nvidiaKmoduleDir, "nvidia.ko.zst")

		kversion := ""
		if utils.Exists(nvidiaKModule) {
			kversion, _ = kernel.ModinfoField(nvidiaKModule, "version")
		} else {
			nvidiaKModule := filepath.Join(nvidiaKmoduleDir, "nvidia.ko")
			if utils.Exists(nvidiaKModule) {
				kversion, _ = kernel.ModinfoField(nvidiaKModule, "version")
			}
		}

		if version == kversion {
			driverDir.WithKernelModules = true
		}

		ans = append(ans, driverDir)
	}

	return &ans, nil
}
