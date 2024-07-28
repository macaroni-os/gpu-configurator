/*
Copyright © 2024 Macaroni OS Linux
See AUTHORS and LICENSE for the license details and contributors.
*/
package macaroni

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/macaroni-os/gpu-configurator/pkg/specs"
	"github.com/macaroni-os/macaronictl/pkg/utils"
)

func (b *MacaroniBackend) SetNVIDIAVersion(setup *specs.NVIDIASetup, v string) error {
	// NOTE: I want to reset the links and setup every time. This permits
	//       to fix things also when there are bugs on gpu-configurator with
	//       previous versions.

	// Configure NVIDIA version needs:
	// 8. create links under /usr/lib64/xorg/modules/drivers/
	// 9. create links under /usr/lib64/xorgs/modules/extensions/
	// 10. create hardlink to nvidia kernel driver.
	// 11. create /etc/conf.d/* (if doesn't exist)
	// 12. create /etc/ld.so.conf.d/07-nvidia.conf

	// 1. create /etc/env.d/09nvidia file
	err := b.createNvidiaEnvfile(v)
	if err != nil {
		return err
	}

	// 2. create links to /usr/bin
	err = b.createNvidiaBins(v)
	if err != nil {
		return err
	}

	// 3. create files under /etc/init.d
	err = b.createInitd(v)
	if err != nil {
		return err
	}

	// 4. create files under /etc/X11, /etc/sandbox.d, /etc/tmpfiles.d
	err = b.createEtc(v)
	if err != nil {
		return err
	}

	// 5. create links for .desktop
	err = b.createDesktopFile(v)
	if err != nil {
		return err
	}

	// 6. create links for .png
	err = b.createPngFile(v)
	if err != nil {
		return err
	}

	// 7. create links under /usr/share
	err = b.createUsrShare(v)
	if err != nil {
		return err
	}

	return nil
}

func (b *MacaroniBackend) createUsrShare(v string) error {
	driverPath := b.getDriverDir(v)

	// Create /usr/share/vulkan/icd.d/nvidia_icd.json file

	// Create /usr/share/vulkan/implicit_layer.d/nvidia_layers.json

	// Create /usr/share/nvidia/ files

	// Create /usr/share/glvnd/egl_vendor.d/10_nvidia.json
	eglvendorTargetPath := "/usr/share/glvnd/egl_vendor.d"
	eglvendorOriginPath := filepath.Join(
		driverPath, eglvendorTargetPath,
		"10_nvidia.json",
	)

	if utils.Exists(eglvendorOriginPath) {

		if !utils.Exists(eglvendorTargetPath) {
			err := os.MkdirAll(eglvendorTargetPath, os.ModePerm)
			if err != nil {
				return err
			}
		}

		eglvendorTargetFile := filepath.Join(
			eglvendorTargetPath,
			"10_nvidia.json",
		)

		err := os.Symlink(eglvendorOriginPath, eglvendorTargetFile)
		if err != nil {
			return fmt.Errorf("error on linking file %s to %s: %s",
				eglvendorOriginPath, eglvendorTargetFile, err.Error())
		}

	} // else TODO add warning

	// Create /usr/share/X11/xorg.conf.d/nvidia-drm-outputclass.conf
	outputclassTargetPath := "/usr/share/X11/xorg.conf.d"
	outputclassOriginPath := filepath.Join(
		driverPath, outputclassTargetPath,
		"nvidia-drm-outputclass.conf",
	)

	if utils.Exists(outputclassOriginPath) {

		if !utils.Exists(outputclassTargetPath) {
			err := os.MkdirAll(outputclassTargetPath, os.ModePerm)
			if err != nil {
				return err
			}
		}

		outputclassTargetFile := filepath.Join(
			outputclassTargetPath,
			"nvidia-drm-outputclass.conf",
		)

		err := os.Symlink(outputclassOriginPath, outputclassTargetFile)
		if err != nil {
			return fmt.Errorf("error on linking file %s to %s: %s",
				outputclassOriginPath, outputclassTargetFile, err.Error())
		}

	} // else TODO add warning

	// Create /usr/share/dbus-1/system.d/nvidia-dbus.conf
	dbusSystemTargetPath := "/usr/share/dbus-1/system.d"
	dbusSystemOriginPath := filepath.Join(
		driverPath, dbusSystemTargetPath,
		"nvidia-dbus.conf",
	)
	if utils.Exists(dbusSystemOriginPath) {

		if !utils.Exists(dbusSystemTargetPath) {
			err := os.MkdirAll(dbusSystemTargetPath, os.ModePerm)
			if err != nil {
				return err
			}
		}

		dbusTargetFile := filepath.Join(
			dbusSystemTargetPath, "nvidia-dbus.conf",
		)

		err := os.Symlink(dbusSystemOriginPath, dbusTargetFile)
		if err != nil {
			return fmt.Errorf("error on linking file %s to %s: %s",
				dbusSystemOriginPath, dbusTargetFile, err.Error())
		}

	} // else TODO: Add warning

	// Create /usr/share/man/man1/* files
	manPages := []string{
		"nvidia-smi.1",
		"nvidia-cuda-mps-control.1",
		"nvidia-persistenced.1",
		"nvidia-xconfig.1",
		"nvidia-settings.1",
	}
	manTargetPath := "/usr/share/man/man1"
	manOriginPath := filepath.Join(
		driverPath, manTargetPath,
	)

	if !utils.Exists(manTargetPath) {
		err := os.MkdirAll(manTargetPath, os.ModePerm)
		if err != nil {
			return err
		}
	}

	for _, f := range manPages {
		manFile := filepath.Join(manOriginPath, f)
		if !utils.Exists(manFile) {
			// TODO: Add warning
			continue
		}

		targetManFile := filepath.Join(manTargetPath, f)
		err := os.Symlink(manFile, targetManFile)
		if err != nil {
			return fmt.Errorf("error on linking file %s to %s: %s",
				targetManFile, manFile, err.Error())
		}
	}

	return nil
}

func (b *MacaroniBackend) createPngFile(v string) error {
	dirPrefix := "nvidia-drivers"
	pixmapsDir := "/usr/share/pixmaps"
	sourceFile := filepath.Join(NvidiaPrefixDriverPath,
		dirPrefix+"-"+v,
		pixmapsDir,
		"nvidia-settings.png",
	)
	targetFile := filepath.Join(pixmapsDir,
		"nvidia-settings.png",
	)

	if utils.Exists(sourceFile) {
		err := os.Symlink(sourceFile, targetFile)
		if err != nil {
			return fmt.Errorf("error on linking file %s to %s: %s",
				targetFile, sourceFile, err.Error())
		}
	}
	// TODO: Add warning if file doesn't exist

	return nil
}

func (b *MacaroniBackend) createDesktopFile(v string) error {
	dirPrefix := "nvidia-drivers"
	appsDesktopDir := "/usr/share/applications"
	sourceFile := filepath.Join(NvidiaPrefixDriverPath,
		dirPrefix+"-"+v,
		appsDesktopDir,
		"nvidia-settings.desktop",
	)
	targetFile := filepath.Join(appsDesktopDir,
		"nvidia-settings.desktop",
	)

	if utils.Exists(sourceFile) {
		err := os.Symlink(sourceFile, targetFile)
		if err != nil {
			return fmt.Errorf("error on linking file %s to %s: %s",
				targetFile, sourceFile, err.Error())
		}
	}
	// TODO: Add warning if file doesn't exist

	return nil
}

func (b *MacaroniBackend) createEtc(v string) error {
	var etcsandboxd = "/etc/sandbox.d"
	var xinitrcd = "/etc/X11/xinit/xinitrc.d"
	var tmpfilesd = "/etc/tmpfiles.d"
	var nvidiaSettingsFile = filepath.Join(xinitrcd, "95-nvidia-settings")
	var nvidiaFile = filepath.Join(etcsandboxd, "20nvidia")
	var nvidiaTmpfilesd = filepath.Join(tmpfilesd, "nvidia-drivers.conf")

	// Create /etc/sandbox.d/20nvidia
	if !utils.Exists(etcsandboxd) {
		err := os.MkdirAll(etcsandboxd, os.ModePerm)
		if err != nil {
			return err
		}
	}

	err := os.WriteFile(nvidiaFile, []byte(
		`SANDBOX_PREDICT="/dev/nvidiactl:/dev/nvidia-caps:/dev/char"
`), 0644,
	)

	if err != nil {
		return fmt.Errorf("error on write file %s: %s",
			nvidiaFile, err.Error())
	}

	if !utils.Exists(xinitrcd) {
		err := os.MkdirAll(xinitrcd, os.ModePerm)
		if err != nil {
			return err
		}
	}

	// Create /etc/X11/xinit/xinitrc.d/95-nvidia-settings
	err = os.WriteFile(nvidiaSettingsFile, []byte(
		`#!/bin/sh
if [ $(lsmod | grep nvidia | wc -l) != "0" ] ; then
  /usr/bin/nvidia-settings --load-config-only
fi
`), 0644)
	if err != nil {
		return fmt.Errorf("error on write file %s: %s",
			nvidiaSettingsFile, err.Error())
	}

	if !utils.Exists(tmpfilesd) {
		err := os.MkdirAll(tmpfilesd, os.ModePerm)
		if err != nil {
			return err
		}
	}

	// Create /etc/tmpfiles.d/nvidia-drivers.conf
	err = os.WriteFile(nvidiaTmpfilesd, []byte(
		`d /run/nvidia-xdriver 0775 root video -`),
		0644)
	if err != nil {
		return fmt.Errorf("error on write file %s: %s",
			nvidiaTmpfilesd, err.Error())
	}

	// Create link on /etc/OpenCL/vendors/nvidia.icd
	dirPrefix := "nvidia-drivers"
	openCLDir := "/etc/OpenCL/vendors"
	sourceOpenCLFile := filepath.Join(NvidiaPrefixDriverPath,
		dirPrefix+"-"+v,
		openCLDir, "nvidia.icd",
	)
	if !utils.Exists(openCLDir) {
		err = os.MkdirAll(openCLDir, os.ModePerm)
		if err != nil {
			return fmt.Errorf("error on create dir %s: %s",
				openCLDir, err.Error())
		}
	}

	linkOpenCLFile := filepath.Join(
		openCLDir, "nvidia.icd",
	)
	err = os.Symlink(sourceOpenCLFile, linkOpenCLFile)
	if err != nil {
		return fmt.Errorf("error on create link file %s to %s: %s",
			linkOpenCLFile, sourceOpenCLFile, err.Error())
	}

	return nil
}

func (b *MacaroniBackend) createInitd(v string) error {
	dirPrefix := "nvidia-drivers"
	initdDir := filepath.Join(NvidiaPrefixDriverPath,
		dirPrefix+"-"+v,
		"etc/init.d",
	)

	scripts := []string{
		"nvidia-persistenced",
		"nvidia-powerd",
		"nvidia-smi",
	}

	for idx := range scripts {

		f := filepath.Join(initdDir, scripts[idx])
		if !utils.Exists(f) {
			continue
		}

		if !utils.Exists(f) {
			// TODO: Add warning
			continue
		}

		t := filepath.Join("/etc/init.d", scripts[idx])

		// Open destination file (truncate it if exists)
		tfd, err := os.OpenFile(t, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		defer tfd.Close()

		// Open source file
		sourcefd, err := os.Open(f)
		if err != nil {
			return err
		}
		defer sourcefd.Close()

		_, err = io.Copy(tfd, sourcefd)
		if err != nil {
			return err
		}
	}

	return nil
}

func (b *MacaroniBackend) getDriverDir(v string) string {
	dirPrefix := "nvidia-drivers"
	driverDir := filepath.Join(NvidiaPrefixDriverPath,
		dirPrefix+"-"+v,
	)
	return driverDir
}

func (b *MacaroniBackend) createNvidiaBins(v string) error {
	var err error

	dirPrefix := "nvidia-drivers"
	driverDir := filepath.Join(NvidiaPrefixDriverPath,
		dirPrefix+"-"+v,
	)
	driverBinDir := filepath.Join(driverDir, "/bin")

	binariesBin := []string{
		"nvidia-cuda-mps-control",
		"nvidia-cuda-mps-server",
		"nvidia-debugdump",
		"nvidia-settings",
		"nvidia-smi",
		"nvidia-xconfig",
		"nvidia-powerd",
		"nvidia-persistenced",
	}

	for idx := range binariesBin {
		f := filepath.Join(driverBinDir, binariesBin[idx])

		if utils.Exists(f) {
			b := filepath.Join("/usr/bin/", binariesBin[idx])
			err = os.Symlink(f, b)
			if err != nil {
				return fmt.Errorf(
					"error on create symlink %s: %s",
					b, err.Error())
			}
		} // else {
		// TODO: Add warning
	}

	return nil
}

func (b *MacaroniBackend) createNvidiaEnvfile(v string) error {
	envNvidia := filepath.Join(b.GetEnvironmentDir(), NvidiaEnvFileName)

	dirPrefix := "nvidia-drivers"
	libDir := filepath.Join(NvidiaPrefixDriverPath,
		dirPrefix+"-"+v,
		"lib64",
	)

	err := os.WriteFile(envNvidia, []byte(
		fmt.Sprintf(`
# autogenerated file by gpu-configurator
LDPATH="%s"
NVIDIA_DRIVER_VERSION="%s"
`,
			libDir, v)), 0644)

	if err != nil {
		return fmt.Errorf("Error on write env file %s: %s",
			envNvidia, err.Error())
	}
	return nil
}
