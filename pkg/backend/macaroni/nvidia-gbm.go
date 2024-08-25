/*
Copyright © 2024 Macaroni OS Linux
See AUTHORS and LICENSE for the license details and contributors.
*/
package macaroni

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/macaroni-os/gpu-configurator/pkg/logger"
	"github.com/macaroni-os/gpu-configurator/pkg/specs"
	"github.com/macaroni-os/macaronictl/pkg/utils"
)

const (
	gbmlibName = "nvidia-drm_gbm.so"
	// Some applications search for nvidia_gbm.so (for
	// example the electron applications).
	gbmlibNameShort = "nvidia_gbm.so"
)

func (b *MacaroniBackend) PurgeGBMLinks(setup *specs.System, libs []string) error {
	log := logger.GetDefaultLogger()
	if len(libs) == 0 {
		log.DebugC("No GBM Libraries to remove.")
	}

	gbmlibdir := b.GetGBMLibDir()

	for _, lib := range libs {
		f := filepath.Join(gbmlibdir, lib)

		if lib != gbmlibName && lib != gbmlibNameShort {
			log.DebugC("Ignoring library", lib)
			continue
		}

		syslib := setup.GetGBMLibrary(lib)
		if syslib != nil {
			if syslib.Disabled {
				err := b.removeFileIfExist(f + ".disabled")
				if err != nil {
					return err
				}
			} else {
				err := b.removeFileIfExist(f)
				if err != nil {
					return err
				}
			}

		} else {
			// POST: Try to force removing both file (with or without .disable)

			err := b.removeFileIfExist(f)
			if err != nil {
				return err
			}

			err = b.removeFileIfExist(f + ".disabled")
			if err != nil {
				return err
			}

		}

	}

	return nil
}

func (b *MacaroniBackend) ConfigureGBMLinks(setup *specs.System,
	nvidiaVersion string, enabled bool) error {
	log := logger.GetDefaultLogger()
	var err error

	gbmlibdir := b.GetGBMLibDir()
	// Create GBM library if doesn't exist
	if !utils.Exists(gbmlibdir) {
		err := os.MkdirAll(gbmlibdir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	libNvidia := setup.GetGBMLibrary(gbmlibName)
	libNvidiaShort := setup.GetGBMLibrary(gbmlibNameShort)

	nvidiaDriver := setup.GetNvidia().GetDriver(nvidiaVersion)
	if nvidiaVersion == "" {
		return fmt.Errorf("Nvidia driver %s not found",
			nvidiaVersion)
	}

	nvidiaLinkedFile := filepath.Join(
		nvidiaDriver.Path, "lib64", gbmlibName)
	nvidiaLinkFile := filepath.Join(gbmlibdir, gbmlibName)
	nvidiaLinkFileShort := filepath.Join(gbmlibdir, gbmlibNameShort)

	if libNvidia != nil && libNvidia.LinkedFile == nvidiaLinkedFile &&
		((enabled && !libNvidia.Disabled) || (!enabled && libNvidia.Disabled)) {
		log.InfoC(fmt.Sprintf(
			"Nvidia gbm library %s is already correctly configured.",
			gbmlibName))
		return nil
	}

	if libNvidia != nil {
		if libNvidia.Disabled {
			err = b.removeFileIfExist(nvidiaLinkFile + ".disabled")
		} else {
			err = b.removeFileIfExist(nvidiaLinkFile)
		}
		if err != nil {
			return err
		}
	}

	if libNvidiaShort != nil {
		if libNvidiaShort.Disabled {
			// NOTE: This is not created automatically because
			//       I remove it if gbmlib is disabled.
			err = b.removeFileIfExist(nvidiaLinkFileShort + ".disabled")
		} else {
			err = b.removeFileIfExist(nvidiaLinkFileShort)
		}
	}

	if enabled {

		log.DebugC(fmt.Sprintf("Creating link %s to %s...",
			nvidiaLinkFile, nvidiaLinkedFile))
		err = os.Symlink(nvidiaLinkedFile, nvidiaLinkFile)
		if err != nil {
			return fmt.Errorf(
				"error on create symlink on %s: %s",
				nvidiaLinkFile, err.Error())
		}

		log.DebugC(fmt.Sprintf("Creating link %s to %s...",
			nvidiaLinkFileShort, nvidiaLinkedFile))
		err = os.Symlink(nvidiaLinkedFile, nvidiaLinkFileShort)
		if err != nil {
			return fmt.Errorf(
				"error on create symlink on %s: %s",
				nvidiaLinkFileShort, err.Error())
		}

	} else {

		nvidiaLinkFile += ".disabled"
		log.DebugC(fmt.Sprintf("Creating link %s to %s...",
			nvidiaLinkFile, nvidiaLinkedFile))
		err = os.Symlink(nvidiaLinkedFile, nvidiaLinkFile)
		if err != nil {
			return fmt.Errorf(
				"error on create symlink on %s: %s",
				nvidiaLinkFile, err.Error())
		}

		// NOTE: Avoiding the creation of the short name
		//       link when disabled.
	}

	return nil
}
