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

	"github.com/macaroni-os/gpu-configurator/pkg/logger"
	"github.com/macaroni-os/gpu-configurator/pkg/specs"
	"github.com/macaroni-os/macaronictl/pkg/utils"
)

func (b *MacaroniBackend) PurgeNVIDIADriver(system *specs.System) error {

	setup := system.GetNvidia()
	log := logger.GetDefaultLogger()

	if setup.VersionActive != "" {
		log.InfoC(fmt.Sprintf("Purging version %s...",
			setup.VersionActive))
	} else {
		log.InfoC("Force purge system from nvidia files...")
	}

	// 1. Removing /etc/env.d/09nvidia file
	f := "/etc/env.d/09nvidia"
	err := b.removeFileIfExist(f)
	if err != nil {
		return err
	}

	// 2. Removing /usr/bin/ links
	err = b.purgeNvidiaBins()
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

	// 12. Remove gbm links
	err = b.PurgeGBMLinks(system, []string{gbmlibName, gbmlibNameShort})
	if err != nil {
		return err
	}

	return nil
}

func (b *MacaroniBackend) removeFileIfExist(f string) error {
	log := logger.GetDefaultLogger()

	if utils.Exists(f) {
		err := os.Remove(f)
		if err != nil {
			return err
		}
		log.DebugC(fmt.Sprintf(
			"File %s removed.", f))
	} else {
		log.DebugC(fmt.Sprintf(
			"File %s not present.", f))
	}

	return nil
}

func (b *MacaroniBackend) purgeLdsoconfdFile() error {
	targetDir := "/etc/ld.so.conf.d"
	targetFile := filepath.Join(targetDir,
		"07-nvidia",
	)

	return b.removeFileIfExist(targetFile)
}

func (b *MacaroniBackend) purgeXorgModulesExtension() error {
	targetPath := "/usr/lib64/xorg/modules/extensions"
	targetFile := filepath.Join(
		targetPath, "libglxserver_nvidia.so",
	)

	return b.removeFileIfExist(targetFile)
}

func (b *MacaroniBackend) purgeXorgModulesDriver() error {
	targetPath := "/usr/lib64/xorg/modules/drivers"
	targetFile := filepath.Join(
		targetPath, "nvidia_drv.so",
	)

	return b.removeFileIfExist(targetFile)
}

func (b *MacaroniBackend) purgeUsrShare(v string) error {
	nvidiaVulkanIcdTargetPath := "/usr/share/vulkan/icd.d"
	nvidiaVulkanIcdTargetFile := filepath.Join(
		nvidiaVulkanIcdTargetPath,
		"nvidia_icd.json",
	)

	err := b.removeFileIfExist(nvidiaVulkanIcdTargetFile)
	if err != nil {
		return err
	}

	// Removing /usr/share/vulkan/implicit_layer.d/nvidia_layers.json
	nvidiaVulkanLayerTargetPath := "/usr/share/vulkan/implicit_layer.d"
	nvidiaVulkanLayerTargetFile := filepath.Join(
		nvidiaVulkanLayerTargetPath,
		"nvidia_layers.json",
	)

	err = b.removeFileIfExist(nvidiaVulkanLayerTargetFile)
	if err != nil {
		return err
	}

	// Removing /usr/share/nvidia/ files
	shareNvidiaTargetPath := "/usr/share/nvidia"
	for _, f := range shareNvidiaFiles {

		if v == "" && strings.Index(f, "PV") > 0 {
			// POST: If v is empty i will try to search file later.
			continue
		}

		f := strings.ReplaceAll(f, "PV", v)
		targetFile := filepath.Join(
			shareNvidiaTargetPath, f)

		err = b.removeFileIfExist(targetFile)
		if err != nil {
			return err
		}
	}

	if v == "" && utils.Exists(shareNvidiaTargetPath) {
		// POST: Try to removing files with a name starting
		//       with `nvidia-application-profiles-`
		dirs, err := os.ReadDir(shareNvidiaTargetPath)
		for _, file := range dirs {
			if file.IsDir() {
				continue
			}

			if strings.HasPrefix(file.Name(), "nvidia-application-profiles-") {
				f := filepath.Join(shareNvidiaTargetPath,
					file.Name())

				err = b.removeFileIfExist(f)
				if err != nil {
					return err
				}
			}
		}
	}

	// Removing /usr/share/glvnd/egl_vendor.d/10_nvidia.json
	eglvendorTargetPath := "/usr/share/glvnd/egl_vendor.d"
	eglvendorTargetFile := filepath.Join(
		eglvendorTargetPath,
		"10_nvidia.json",
	)

	err = b.removeFileIfExist(eglvendorTargetFile)
	if err != nil {
		return err
	}

	// Removing /usr/share/X11/xorg.conf.d/nvidia-drm-outputclass.conf
	outputclassTargetPath := "/usr/share/X11/xorg.conf.d"
	outputclassTargetFile := filepath.Join(
		outputclassTargetPath,
		"nvidia-drm-outputclass.conf",
	)
	err = b.removeFileIfExist(outputclassTargetFile)
	if err != nil {
		return err
	}

	// Removing /usr/share/dbus-1/system.d/nvidia-dbus.conf
	dbusSystemTargetPath := "/usr/share/dbus-1/system.d"
	dbusTargetFile := filepath.Join(
		dbusSystemTargetPath, "nvidia-dbus.conf",
	)
	err = b.removeFileIfExist(dbusTargetFile)
	if err != nil {
		return err
	}

	// Removing /usr/share/man/man1/* files
	manTargetPath := "/usr/share/man/man1"
	for _, f := range manPages {
		manFile := filepath.Join(manTargetPath, f)
		err = b.removeFileIfExist(manFile)
		if err != nil {
			return err
		}
	}

	return nil
}

func (b *MacaroniBackend) purgePngfile() error {
	pixmapsDir := "/usr/share/pixmaps"
	pngFile := filepath.Join(pixmapsDir,
		"nvidia-settings.png",
	)

	err := b.removeFileIfExist(pngFile)
	if err != nil {
		return err
	}

	return nil
}

func (b *MacaroniBackend) purgeDesktopfile() error {
	appsDesktopDir := "/usr/share/applications"
	desktopFile := filepath.Join(appsDesktopDir,
		"nvidia-settings.desktop",
	)

	err := b.removeFileIfExist(desktopFile)
	if err != nil {
		return err
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

	err := b.removeFileIfExist(nvidiaFile)
	if err != nil {
		return err
	}

	err = b.removeFileIfExist(nvidiaSettingsFile)
	if err != nil {
		return err
	}

	err = b.removeFileIfExist(nvidiaTmpfilesd)
	if err != nil {
		return err
	}

	openCLDir := "/etc/OpenCL/vendors"
	linkOpenCLFile := filepath.Join(
		openCLDir, "nvidia.icd",
	)
	return b.removeFileIfExist(linkOpenCLFile)
}

func (b *MacaroniBackend) purgeNvidiaBins() error {
	for idx := range binariesBin {
		f := filepath.Join("/usr/bin/", binariesBin[idx])
		err := b.removeFileIfExist(f)
		if err != nil {
			return err
		}
	}

	return nil
}

func (b *MacaroniBackend) purgeNvidiaInitd() error {
	for idx := range initdscripts {
		f := filepath.Join("/etc/init.d/", initdscripts[idx])
		err := b.removeFileIfExist(f)
		if err != nil {
			return err
		}
	}

	return nil
}
