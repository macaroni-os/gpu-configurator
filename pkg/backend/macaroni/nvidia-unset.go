/*
Copyright © 2024 Macaroni OS Linux
See AUTHORS and LICENSE for the license details and contributors.
*/
package macaroni

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/macaroni-os/gpu-configurator/pkg/specs"
	"github.com/macaroni-os/macaronictl/pkg/utils"
)

func (b *MacaroniBackend) PurgeNVIDIADriver(setup *specs.NVIDIASetup) error {

	// 1. Removing /etc/env.d/09nvidia file
	f := "/etc/env.d/09nvidia"
	if utils.Exists(f) {
		err := os.Remove(f)
		if err != nil {
			return err
		}
	}

	// 2. Removing /usr/bin/ links
	err := b.purgeNvidiaBins()
	if err != nil {
		return err
	}

	// 3. Removing /etc/init.d links
	err = b.purgeNvidiaInitd()
	if err != nil {
		return err
	}

	// 4. Removing file from /etc/X11, /etc/sandbox.d, /etc/tmpfiles.d
	err = b.purgeEtc()
	if err != nil {
		return err
	}

	// 5. remove link for .desktop
	err = b.purgeDesktopfile()
	if err != nil {
		return err
	}

	// 6. remove png file
	err = b.purgePngfile()
	if err != nil {
		return err
	}

	// 7. remove links under /usr/share
	err = b.purgeUsrShare(setup.VersionActive)
	if err != nil {
		return err
	}

	// 8. remove links under /usr/lib64/xorg/modules/drivers/
	err = b.purgeXorgModulesDriver()
	if err != nil {
		return err
	}

	// 9. removing links under /usr/lib64/xorgs/modules/extensions
	err = b.purgeXorgModulesExtension()
	if err != nil {
		return err
	}

	// 10. I avoid to remove file from /etc/conf.d/

	// 11. removing /etc/ld.so.conf.d file
	err = b.purgeLdsoconfdFile()
	if err != nil {
		return err
	}

	return nil
}

func (b *MacaroniBackend) purgeLdsoconfdFile() error {
	targetDir := "/etc/ld.so.conf.d"
	targetFile := filepath.Join(targetDir,
		"07-nvidia",
	)

	if utils.Exists(targetFile) {
		err := os.Remove(targetFile)
		if err != nil {
			return err
		}
	}

	return nil
}

func (b *MacaroniBackend) purgeXorgModulesExtension() error {
	targetPath := "/usr/lib64/xorg/modules/extensions"
	targetFile := filepath.Join(
		targetPath, "nvidia_drv.so",
	)

	if utils.Exists(targetFile) {
		err := os.Remove(targetFile)
		if err != nil {
			return err
		}
	}

	return nil
}

func (b *MacaroniBackend) purgeXorgModulesDriver() error {
	targetPath := "/usr/lib64/xorg/modules/drivers"
	targetFile := filepath.Join(
		targetPath, "nvidia_drv.so",
	)

	if utils.Exists(targetFile) {
		err := os.Remove(targetFile)
		if err != nil {
			return err
		}
	}

	return nil
}

func (b *MacaroniBackend) purgeUsrShare(v string) error {
	nvidiaVulkanIcdTargetPath := "/usr/share/vulkan/icd.d"
	nvidiaVulkanIcdTargetFile := filepath.Join(
		nvidiaVulkanIcdTargetPath,
		"nvidia_icd.json",
	)

	if utils.Exists(nvidiaVulkanIcdTargetFile) {
		err := os.Remove(nvidiaVulkanIcdTargetFile)
		if err != nil {
			return err
		}
	}

	// Removing /usr/share/vulkan/implicit_layer.d/nvidia_layers.json
	nvidiaVulkanLayerTargetPath := "/usr/share/vulkan/implicit_layer.d"
	nvidiaVulkanLayerTargetFile := filepath.Join(
		nvidiaVulkanLayerTargetPath,
		"nvidia_layers.json",
	)

	if utils.Exists(nvidiaVulkanLayerTargetFile) {
		err := os.Remove(nvidiaVulkanLayerTargetFile)
		if err != nil {
			return err
		}
	}

	// Removing /usr/share/nvidia/ files
	shareNvidiaTargetPath := "/usr/share/nvidia"
	for _, f := range shareNvidiaFiles {
		f := strings.ReplaceAll(f, "PV", v)
		targetfile := filepath.Join(
			shareNvidiaTargetPath, f)

		err := os.Remove(targetfile)
		if err != nil {
			return err
		}
	}

	// Removing /usr/share/glvnd/egl_vendor.d/10_nvidia.json
	eglvendorTargetPath := "/usr/share/glvnd/egl_vendor.d"
	eglvendorTargetFile := filepath.Join(
		eglvendorTargetPath,
		"10_nvidia.json",
	)

	if utils.Exists(eglvendorTargetFile) {
		err := os.Remove(eglvendorTargetFile)
		if err != nil {
			return err
		}
	}

	// Removing /usr/share/dbus-1/system.d/nvidia-dbus.conf
	dbusSystemTargetPath := "/usr/share/dbus-1/system.d"
	dbusTargetFile := filepath.Join(
		dbusSystemTargetPath, "nvidia-dbus.conf",
	)
	if utils.Exists(dbusTargetFile) {
		err := os.Remove(dbusTargetFile)
		if err != nil {
			return err
		}
	}

	// Removing /usr/share/man/man1/* files
	manTargetPath := "/usr/share/man/man1"
	for _, f := range manPages {
		manFile := filepath.Join(manTargetPath, f)
		if utils.Exists(manFile) {
			err := os.Remove(manFile)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (b *MacaroniBackend) purgePngfile() error {
	pixmapsDir := "/usr/share/pixmaps"
	pngFile := filepath.Join(pixmapsDir,
		"nvidia-settings.png",
	)

	if utils.Exists(pngFile) {
		err := os.Remove(pngFile)
		if err != nil {
			return err
		}
	}

	return nil
}

func (b *MacaroniBackend) purgeDesktopfile() error {
	appsDesktopDir := "/usr/share/applications"
	desktopFile := filepath.Join(appsDesktopDir,
		"nvidia-settings.desktop",
	)

	if utils.Exists(desktopFile) {
		err := os.Remove(desktopFile)
		if err != nil {
			return err
		}
	}

	return nil
}

func (b *MacaroniBackend) purgeEtc() error {
	var etcsandboxd = "/etc/sandbox.d"
	var xinitrcd = "/etc/X11/xinit/xinitrc.d"
	var tmpfilesd = "/etc/tmpfiles.d"
	var nvidiaSettingsFile = filepath.Join(xinitrcd, "95-nvidia-settings")
	var nvidiaFile = filepath.Join(etcsandboxd, "20nvidia")
	var nvidiaTmpfilesd = filepath.Join(tmpfilesd, "nvidia-drivers.conf")

	if utils.Exists(nvidiaFile) {
		err := os.Remove(nvidiaFile)
		if err != nil {
			return err
		}
	}

	if utils.Exists(nvidiaSettingsFile) {
		err := os.Remove(nvidiaSettingsFile)
		if err != nil {
			return err
		}
	}

	if utils.Exists(nvidiaTmpfilesd) {
		err := os.Remove(nvidiaTmpfilesd)
		if err != nil {
			return err
		}
	}

	openCLDir := "/etc/OpenCL/vendors"
	linkOpenCLFile := filepath.Join(
		openCLDir, "nvidia.icd",
	)
	if utils.Exists(linkOpenCLFile) {
		err := os.Remove(linkOpenCLFile)
		if err != nil {
			return err
		}
	}

	return nil
}

func (b *MacaroniBackend) purgeNvidiaBins() error {
	for idx := range binariesBin {
		f := filepath.Join("/usr/bin/", binariesBin[idx])
		if utils.Exists(f) {
			err := os.Remove(f)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (b *MacaroniBackend) purgeNvidiaInitd() error {
	for idx := range initdscripts {
		f := filepath.Join("/etc/init.d/", initdscripts[idx])
		if utils.Exists(f) {
			err := os.Remove(f)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
