/*
Copyright © 2024 Macaroni OS Linux
See AUTHORS and LICENSE for the license details and contributors.
*/
package analyzer

import (
	"encoding/json"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	bmacaroni "github.com/macaroni-os/gpu-configurator/pkg/backend"
	"github.com/macaroni-os/gpu-configurator/pkg/specs"
)

type Analyzer struct {
	Backend bmacaroni.SystemBackend

	System *specs.System
}

func NewAnalyzer(btype string) (*Analyzer, error) {
	var err error

	ans := &Analyzer{}
	ans.System = specs.NewSystem()

	ans.Backend, err = bmacaroni.NewBackend(btype)
	if err != nil {
		return nil, err
	}

	return ans, nil
}

func (a *Analyzer) GetSystem() *specs.System { return a.System }

func (a *Analyzer) readGbmLibs() error {
	var regexlib = regexp.MustCompile(`.so$|.so.disabled$`)

	gbmlibdir := a.Backend.GetGBMLibDir()
	if gbmlibdir == "" {
		// POST: nothing to do.
		return nil
	}

	dirEntries, err := os.ReadDir(gbmlibdir)
	if err != nil {
		return err
	}

	for _, file := range dirEntries {
		if file.IsDir() {
			continue
		}

		if !regexlib.MatchString(file.Name()) {
			continue
		}

		lib := &specs.Library{
			Name:       file.Name(),
			LinkedFile: "",
		}

		path := filepath.Join(gbmlibdir, file.Name())

		// Check if the library is a link
		finfo, err := os.Lstat(path)
		if err != nil {
			return err
		}
		if finfo.Mode()&os.ModeSymlink != 0 {
			// Resolve link
			linkedLink, err := os.Readlink(path)
			if err != nil {
				return err
			}
			lib.LinkedFile = linkedLink
		}

		if strings.HasSuffix(file.Name(), "disabled") {
			lib.Disabled = true
			lib.Name = file.Name()[0 : len(file.Name())-len(".disabled")]
		}

		a.System.GbmLibraries = append(a.System.GbmLibraries, lib)
	}

	return nil
}

func (a *Analyzer) Read() error {
	var err error
	var regexICD = regexp.MustCompile(`.json$|.json.disable$`)
	var dirs []string

	// Read egl external platforms directories
	dirs, err = a.Backend.GetEglExternalPlatformsDirs()
	if err != nil {
		return err
	}

	for _, dir := range dirs {
		dirEntries, err := os.ReadDir(dir)
		if err != nil {
			// TODO: Add warning
			continue
		}

		// Create directory container
		egldir := specs.NewEglExternalPlatformFiles(dir)
		a.System.EglExtPlatformDirs = append(a.System.EglExtPlatformDirs,
			egldir)

		for _, file := range dirEntries {
			if file.IsDir() {
				continue
			}

			if !regexICD.MatchString(file.Name()) {
				continue
			}

			content, err := os.ReadFile(path.Join(dir, file.Name()))
			if err != nil {
				// TODO: Add warning
				continue
			}

			icdjson, err := specs.NewICDJson(content)
			if err != nil {
				return err
			}

			egldir.Files[file.Name()] = icdjson
		}

	}

	// Read vulkan layers files
	dirs, err = a.Backend.GetVulkanLayersDirs()
	if err != nil {
		return err
	}

	for _, dir := range dirs {
		dirEntries, err := os.ReadDir(dir)
		if err != nil {
			// TODO: Add warning
			continue
		}

		// Create directory container
		vulkandir := specs.NewVulkanLayersFiles(dir)
		a.System.VulkanLayersDirs = append(a.System.VulkanLayersDirs,
			vulkandir)

		for _, file := range dirEntries {
			if file.IsDir() {
				continue
			}

			if !regexICD.MatchString(file.Name()) {
				continue
			}

			content, err := os.ReadFile(path.Join(dir, file.Name()))
			if err != nil {
				// TODO: Add warning
				continue
			}

			fileData := make(map[string]interface{}, 0)
			if err := json.Unmarshal(content, &fileData); err != nil {
				return err
			}

			vulkandir.Files[file.Name()] = fileData
		}

	}

	// Read vulkan icd directory
	dirs, err = a.Backend.GetVulkanICDDirs()
	if err != nil {
		return err
	}

	for _, dir := range dirs {
		dirEntries, err := os.ReadDir(dir)
		if err != nil {
			// TODO: Add warning
			continue
		}

		// Create directory container
		vulkandir := specs.NewEglExternalPlatformFiles(dir)
		a.System.VulkanICDDirs = append(a.System.VulkanICDDirs,
			vulkandir)

		for _, file := range dirEntries {
			if file.IsDir() {
				continue
			}

			if !regexICD.MatchString(file.Name()) {
				continue
			}

			content, err := os.ReadFile(path.Join(dir, file.Name()))
			if err != nil {
				// TODO: Add warning
				continue
			}

			icdjson, err := specs.NewICDJson(content)
			if err != nil {
				return err
			}

			vulkandir.Files[file.Name()] = icdjson
		}

	}

	err = a.readGbmLibs()
	if err != nil {
		return err
	}

	// Retrieve NVIDIA drivers installed
	nvDrivers, err := a.Backend.GetNVIDIADrivers()
	if err != nil {
		return err
	}

	a.System.Nvidia = specs.NewNVIDIASetup()
	a.System.Nvidia.Drivers = *nvDrivers
	a.System.Nvidia.ElaborateVersion()

	nvidiaKModules, err := a.Backend.GetNVIDIAKernelModules()
	if err != nil {
		return err
	}
	a.System.Nvidia.KModuleAvailable = *nvidiaKModules

	return nil
}
