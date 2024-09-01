/*
Copyright © 2024 Macaroni OS Linux
See AUTHORS and LICENSE for the license details and contributors.
*/
package macaroni

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/macaroni-os/gpu-configurator/pkg/kernel"
	"github.com/macaroni-os/gpu-configurator/pkg/logger"
	"github.com/macaroni-os/gpu-configurator/pkg/specs"

	"github.com/macaroni-os/macaronictl/pkg/utils"
)

const (
	NvidiaEnvFileName      = "09nvidia"
	NvidiaPrefixDriverPath = "/opt/nvidia"
)

var (
	KernelModuleSupportedCompression = []string{
		".zst",
		".xz",
	}
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

func (b *MacaroniBackend) GetNVIDIADriverActive() (string, error) {
	ans := ""
	log := logger.GetDefaultLogger()

	drivers, err := b.GetNVIDIADrivers()
	if err != nil {
		return "", err
	}

	envNvidia := filepath.Join(b.GetEnvironmentDir(), NvidiaEnvFileName)

	if !utils.Exists(envNvidia) {
		// POST: there isn't an active nvidia driver.
		return "", nil
	}

	// Read file /etc/env.d/09nvidia and get NVIDIA_DRIVER_VERSION env
	f, err := os.Open(envNvidia)
	if err != nil {
		return "", fmt.Errorf("Error on open file %s: %s", envNvidia, err.Error())
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") {
			continue
		}
		tokens := strings.Split(line, "=")
		if len(tokens) != 2 {
			continue
		}
		if tokens[0] == "NVIDIA_DRIVER_VERSION" {
			ans = strings.ReplaceAll(tokens[1], "\"", "")
			log.DebugC("Found nvidia env version with value", ans)
			break
		}
	}

	if err := scanner.Err(); err != nil {
		return ans, fmt.Errorf("Error on read file %s: %s",
			envNvidia, err.Error())
	}

	// Check if the version is present between
	// drivers else ignore it.
	hasVersion := false
	for _, driver := range *drivers {
		if driver.Version == ans {
			hasVersion = true
			break
		}
	}
	if hasVersion {
		return ans, nil
	}
	if ans != "" {
		log.DebugC("Found nvidia env variable mismatch. Ignored.")
	}
	return "", nil
}

func (b *MacaroniBackend) GetNVIDIAKernelModulesActive(open bool) (*[]*specs.KernelModule, error) {
	modulePath := "/lib/modules"
	ans := []*specs.KernelModule{}
	log := logger.GetDefaultLogger()

	if !utils.Exists(modulePath) {
		return &ans, nil
	}

	dirEntries, err := os.ReadDir(modulePath)
	if err != nil {
		return nil, err
	}

	for _, file := range dirEntries {
		if !file.IsDir() {
			continue
		}

		driverType := ""
		kversion := ""
		kVersion := file.Name()

		nvidiaKmoduleDir := filepath.Join(
			modulePath, kVersion, "video")

		// NOTE: It seems that there are few difference from modinfo
		//       between open and proprietary driver. I use license.

		nvidiaKModule := filepath.Join(nvidiaKmoduleDir, "nvidia.ko")

		if utils.Exists(nvidiaKModule) {
			kversion, _ = kernel.ModinfoField(nvidiaKModule, "version")
			driverType, _ = kernel.ModinfoField(nvidiaKModule, "license")
		} else {
			// Check for the compresses modules
			for _, c := range KernelModuleSupportedCompression {
				nvidiaKModule = filepath.Join(nvidiaKmoduleDir, "nvidia.ko")
				nvidiaKModule += c
				if utils.Exists(nvidiaKModule) {
					kversion, _ = kernel.ModinfoField(nvidiaKModule, "version")
					driverType, _ = kernel.ModinfoField(nvidiaKModule, "license")
					break
				}
			}
		}

		// If driverType is empty I consider the kernel module
		// as NVIDIA driver.
		if open && (driverType == "" || driverType == "NVIDIA") {
			continue
		}

		if !open && (driverType != "" && driverType != "NVIDIA") {

			log.DebugC(fmt.Sprintf("Skipping kernel driver for version %s and kernel %s on path %s.",
				kversion, kVersion, nvidiaKmoduleDir))
			continue
		}

		if kversion != "" {

			if open {
				log.DebugC(fmt.Sprintf("Found open active kernel driver for version %s and kernel %s.",
					kversion, kVersion))
			} else {
				log.DebugC(fmt.Sprintf("Found active kernel driver for version %s and kernel %s.",
					kversion, kVersion))
			}

			lp := &specs.KernelModule{
				Path:          nvidiaKModule,
				KernelVersion: kVersion,
				Name:          "nvidia",
				Fields:        make(map[string]string, 0),
			}
			lp.Fields["version"] = kversion
			lp.Fields["license"] = driverType
			ans = append(ans, lp)
		}

	}

	return &ans, nil
}

func (b *MacaroniBackend) GetNVIDIAKernelModules(open bool) (*[]*specs.KernelModule, error) {
	modulePath := "/lib/modules/nvidia"
	if open {
		modulePath = "/lib/modules/nvidia-open"
	}

	log := logger.GetDefaultLogger()
	ans := []*specs.KernelModule{}

	if !utils.Exists(modulePath) {
		return &ans, nil
	}

	dirEntries, err := os.ReadDir(modulePath)
	if err != nil {
		return nil, err
	}

	// Path used by slotted package follow this pattern
	// /lib/modules/[nvidia|nvidia-open]/<NVIDIA_DRIVER_VERSION>/<KVERSION>/video/*.ko[.zst]

	for _, file := range dirEntries {
		if !file.IsDir() {
			continue
		}

		nvidiaVersion := file.Name()

		nvidiaKVersionPath := filepath.Join(modulePath, nvidiaVersion)
		kernelDirs, err := os.ReadDir(nvidiaKVersionPath)
		if err != nil {
			return nil, err
		}

		for _, kf := range kernelDirs {
			kVersion := kf.Name()
			kversion := ""
			license := ""

			nvidiaKmoduleDir := filepath.Join(
				nvidiaKVersionPath, kVersion, "video")

			nvidiaKModule := filepath.Join(nvidiaKmoduleDir, "nvidia.ko")

			if utils.Exists(nvidiaKModule) {
				kversion, _ = kernel.ModinfoField(nvidiaKModule, "version")
				license, _ = kernel.ModinfoField(nvidiaKModule, "license")
			} else {
				for _, c := range KernelModuleSupportedCompression {
					nvidiaKModule = filepath.Join(nvidiaKmoduleDir, "nvidia.ko")
					nvidiaKModule += c

					if utils.Exists(nvidiaKModule) {
						kversion, _ = kernel.ModinfoField(nvidiaKModule, "version")
						license, _ = kernel.ModinfoField(nvidiaKModule, "license")
						break
					}
				}
			}

			if open {
				log.DebugC(fmt.Sprintf(
					"Found open kernel module %s for nvidia version %s under %s (%s).",
					kVersion, nvidiaVersion, nvidiaKmoduleDir,
					filepath.Base(nvidiaKModule)))
			} else {
				log.DebugC(fmt.Sprintf(
					"Found kernel module %s for nvidia version %s under %s (%s).",
					kVersion, nvidiaVersion, nvidiaKmoduleDir,
					filepath.Base(nvidiaKModule)))
			}

			// TODO: if nvidiaKModule != nvidiaVersion add Warning.

			if kversion != "" {
				lp := &specs.KernelModule{
					Path:          nvidiaKModule,
					KernelVersion: kVersion,
					Name:          "nvidia",
					Fields:        make(map[string]string, 0),
				}
				lp.Fields["version"] = kversion
				lp.Fields["license"] = license
				ans = append(ans, lp)
			}
		}

	}

	return &ans, nil
}

func (b *MacaroniBackend) GetNVIDIADrivers() (*[]*specs.NVIDIADriver, error) {
	ans := []*specs.NVIDIADriver{}
	log := logger.GetDefaultLogger()

	dirPrefix := "nvidia-drivers"
	modulesPath := "/lib/modules/"

	if !utils.Exists(NvidiaPrefixDriverPath) {
		// POST: no nvidia drivers available
		return &ans, nil
	}

	dirEntries, err := os.ReadDir(NvidiaPrefixDriverPath)
	if err != nil {
		return nil, err
	}

	// Retrieve the list of the kernels directories available.
	kernels := []string{}

	if utils.Exists(modulesPath) {
		kernelsDirEntries, err := os.ReadDir(modulesPath)
		if err != nil {
			return nil, err
		}
		for _, kv := range kernelsDirEntries {
			if kv.IsDir() {
				kernels = append(kernels, kv.Name())
			}
		}
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
			Path:    filepath.Join(NvidiaPrefixDriverPath, file.Name()),
			Version: version,
		}

		for _, kv := range kernels {
			nvidiaKmoduleDir := filepath.Join(
				modulesPath, kv, "video")
			nvidiaKModule := filepath.Join(nvidiaKmoduleDir, "nvidia.ko")

			log.DebugC(fmt.Sprintf(
				"Checking kernel version %s and version %s...",
				kv, version))

			kversion := ""
			if utils.Exists(nvidiaKModule) {
				kversion, _ = kernel.ModinfoField(nvidiaKModule, "version")
			} else {
				for _, c := range KernelModuleSupportedCompression {
					nvidiaKModule = filepath.Join(nvidiaKmoduleDir, "nvidia.ko")
					nvidiaKModule += c

					if utils.Exists(nvidiaKModule) {
						kversion, _ = kernel.ModinfoField(nvidiaKModule, "version")
						break
					}
				}
			}

			if version == kversion {
				driverDir.WithKernelModules = true
			} else {
				continue
			}

			log.DebugC(fmt.Sprintf("Found driver %s under %s with kernel module %v",
				version, driverDir.Path, driverDir.WithKernelModules))
		}

		ans = append(ans, driverDir)
	}

	return &ans, nil
}
